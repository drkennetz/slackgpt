package slackhandler

import (
	"context"
	"fmt"
	gogpt "github.com/sashabaranov/go-gpt3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"testing"
)

func initLogger(service string) (*zap.SugaredLogger, error) {
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
	return log.Sugar(), nil
}

func TestEventHandlerNoHandles(t *testing.T) {
	logger, err := initLogger("test")
	assert.NoError(t, err)
	appToken := "xapp-123"
	botTok := "xoxb-test"
	ctx := context.Background()
	client := gogpt.NewClient("test-token")
	err = EventHandler(appToken, botTok, client, ctx, logger)
	require.ErrorContains(t, err, "invalid_auth")
}

func TestConversation_UpdateConversation(t *testing.T) {
	convo := newConversation()
	userChannel := "user"
	for i := 0; i < 10; i++ {
		tmp := fmt.Sprintf("%s%v", userChannel, i)
		convo.UpdateConversation(userChannel, tmp)
		if i < 8 {
			assert.Equal(t, convo[userChannel][0], "user0")
		} else {
			assert.Equal(t, convo[userChannel][0], fmt.Sprintf("%s%v", "user", i%7))
		}
	}
}
