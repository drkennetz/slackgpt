package slackhandler

import (
	"context"
	"github.com/sashabaranov/go-openai"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"testing"
)

type mockHandler struct {
	mock.Mock
}

func (c *mockHandler) Handle(et socketmode.EventType, f socketmode.SocketmodeHandlerFunc) {}

func (c *mockHandler) HandleEvents(et slackevents.EventsAPIType, f socketmode.SocketmodeHandlerFunc) {
}

func (c *mockHandler) RunEventLoop() error {
	args := c.Called()
	return args.Error(0)
}

func initLogger(service string) (*log.Logger, error) {
	config := zap.NewProductionConfig()
	config.OutputPaths = []string{"stdout"}
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.DisableStacktrace = true
	config.InitialFields = map[string]any{
		"test": service,
	}
	log, err := config.Build()
	if err != nil {
		return nil, err
	}
	return zap.NewStdLog(log), nil
}

func TestEventHandlerArgs_NewSocketmodeHandler(t *testing.T) {
	logger, err := initLogger("test")
	assert.NoError(t, err)
	appToken := "xapp-123"
	botTok := "xoxb-test"
	ctx := context.Background()
	client := openai.NewClient("test-token")
	slackClient := slack.New(
		botTok,
		slack.OptionDebug(false),
		slack.OptionAppLevelToken(appToken),
		slack.OptionLog(logger),
	)
	socketmodeClient := socketmode.New(
		slackClient,
		socketmode.OptionDebug(false),
		socketmode.OptionLog(logger),
	)
	args := EventHandlerArgs{
		Logger:           logger,
		SlackClient:      slackClient,
		SocketModeClient: socketmodeClient,
		GPTClient:        client,
		Context:          ctx,
	}
	handler := args.NewSocketmodeHandler()
	assert.NotEmpty(t, handler)

	c := make(chan socketmode.Event)
	e := make(chan error)
	handler.Client.Events = c
	go func() {
		defer close(c)
		defer close(e)
		e <- EventHandler(args, handler)
	}()
	errNew := <-e

	assert.Error(t, errNew)

}
