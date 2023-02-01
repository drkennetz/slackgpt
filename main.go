package main

import (
	"context"
	"fmt"
	"github.com/PullRequestInc/go-gpt3"
	"github.com/alexflint/go-arg"
	"github.com/spf13/viper"
	gptslack "go-slack-chat-gpt3/src/slack"
	"log"
	"os"
	"os/signal"
	"syscall"
)

type args struct {
	Config string `arg:"required,-c,--config" help:"config file with slack app+bot tokens, chat-gpt API token"`
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
	var arguments args
	arg.MustParse(&arguments)

	log.New(os.Stdout, "slack-gpt", log.Ldate|log.Ltime|log.Lshortfile)
	if err := run(arguments.Config); err != nil {
		os.Exit(1)
	}
}

func run(config string) error {
	log.SetOutput(os.Stdout)
	viper.SetConfigFile(config)
	viper.ReadInConfig()
	cgptApiKey := viper.GetString("CGPT_API_KEY")
	if cgptApiKey == "" {
		log.Fatalln("Missing chat-gpt API KEY")
	}
	slackAppToken := viper.GetString("SLACK_APP_TOKEN")
	if slackAppToken == "" {
		log.Fatalln("Missing slack app token")
	}
	slackBotToken := viper.GetString("SLACK_BOT_TOKEN")
	if slackBotToken == "" {
		log.Fatalln("Missing slack bot token")
	}
	log.Println("Config values parsed")
	ctx := context.Background()
	client := gpt3.NewClient(cgptApiKey)

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
		handlerErrors <- gptslack.EventHandler(slackAppToken, slackBotToken, client, ctx)
	}()

	// Blocking main and wiating for shutdown
	// This is a blocking select to handle errors - not shutdown
	select {
	case err := <-handlerErrors:
		return fmt.Errorf("handler error: %w", err)

	case sig := <-shutdown:

		log.Println("received shutdown signal, ", sig)
		// give outstanding requests a deadline for completion
		timeoutContext, cancel := context.WithTimeout(ctx, 10)
		defer cancel()

		log.Println("closing context", timeoutContext)
		// Asking listener to shutdown and shed load
		log.Println("Shutting down..")
	}
	return nil
}
