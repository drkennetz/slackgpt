package main

import (
	"context"
	"fmt"
	"github.com/PullRequestInc/go-gpt3"
	"github.com/alexflint/go-arg"
	configs "go-slack-chat-gpt3/config"
	gptslack "go-slack-chat-gpt3/src/slack"
	"go.uber.org/automaxprocs/maxprocs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"os/signal"
	"runtime"
	"syscall"
)

type args struct {
	Config string `arg:"required,-c,--config" help:"config file with slack app+bot tokens, chat-gpt API token"`
	Type   string `arg:"-t, --type" default:"" help:"the config type [json, toml, yaml, hcl, ini, env, properties]; if not passed, inferred from file ext"`
}

func (args) Version() string {
	return "VERSION: development\n"
}

func (args) Description() string {
	return "This program is a slack bot that sends mentions to chat-gpt and responds with chat-gpt result\n"
}

func (args) Epilogue() string {
	return "for more information, visit https://github.com/drkennetz/go-slack-chat-gpt3"
}

func main() {
	// Perform the startup and shutdown sequence
	log, err := initLogger("SLACKGPT-BOT")
	if err != nil {
		fmt.Println("Error constructing logger:", err)
		os.Exit(1)
	}
	defer log.Sync()

	var arguments args
	arg.MustParse(&arguments)

	log.Infow("startup", "version", arguments.Version())
	if err := run(arguments.Config, arguments.Type, log); err != nil {
		os.Exit(1)
	}
}

func run(config, cfgType string, log *zap.SugaredLogger) error {
	// ========================
	// GOMAXPROCS

	// set the correct number of threads for the service
	// based on either machine or quotas in kub
	if _, err := maxprocs.Set(); err != nil {
		return fmt.Errorf("maxprocs: %w", err)
	}
	log.Infow("startup", "GOMAXPROCS", runtime.GOMAXPROCS(0))
	cfgParts, err := configs.ParseConfigFromPath(config, cfgType)
	cfg, err := configs.LoadConfig(cfgParts)
	if err != nil {
		return err
	}
	ctx := context.Background()

	log.Infow("startup", "status", "gpt3 client started")
	client := gpt3.NewClient(cfg.ChatGPTKey)

	// make a channel to listen for an interrupt or term signal from the os
	// use a buffered channel because the signal package requires it
	shutdown := make(chan os.Signal, 1)
	// Should I capture more?
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	// our event handler will have a  buffer of 1, sends happen before receives, so this
	// goroutine will return before server shuts down.
	// In the future, certain errors may trigger a shutdown, but not right now
	handlerErrors := make(chan error, 1)

	// Start the service listening for events
	go func() {
		log.Infow("startup", "status", "slack event handler started")
		// refactor this to take in a struct
		handlerErrors <- gptslack.EventHandler(cfg.SlackAppToken, cfg.SlackBotToken, client, ctx, log)
	}()

	// Blocking main and waiting for shutdown
	// This is a blocking select to handle errors - not shutdown
	select {
	case err := <-handlerErrors:
		return fmt.Errorf("handler error: %w", err)

	case sig := <-shutdown:
		log.Infow("shutdown", "status", "shutdown started", "signal", sig)
		defer log.Infow("shutdown", "status", "shutdown complete", "signal", sig)
		// give outstanding requests a deadline for completion
		_, cancel := context.WithTimeout(ctx, 10)
		defer cancel()
	}
	return nil
}

func initLogger(service string) (*zap.SugaredLogger, error) {
	config := zap.NewProductionConfig()
	config.OutputPaths = []string{"stdout"}
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.DisableStacktrace = true
	config.InitialFields = map[string]any{
		"service": service,
	}
	log, err := config.Build()
	if err != nil {
		return nil, err
	}
	return log.Sugar(), nil
}
