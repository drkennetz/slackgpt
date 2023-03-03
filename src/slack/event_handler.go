// Package slackhandler handles slack appMention events and responds with chat-gpt response
package slackhandler

import (
	"context"
	gogpt "github.com/sashabaranov/go-gpt3"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
	"go.uber.org/zap"
)

// EventHandler handles slack events
func EventHandler(appToken string, botToken string, gptClient *gogpt.Client, ctx context.Context, log *zap.SugaredLogger) error {

	desugared := zap.NewStdLog(log.Desugar())
	api := slack.New(
		botToken,
		slack.OptionDebug(true),
		slack.OptionLog(desugared),
		slack.OptionAppLevelToken(appToken),
	)
	client := socketmode.New(
		api,
		socketmode.OptionDebug(true),
		socketmode.OptionLog(desugared),
	)
	convo := newConversation()

	socketmodeHandler := socketmode.NewSocketmodeHandler(client)
	// should be a primary middleware handler, and these handle more granular events
	socketmodeHandler.Handle(socketmode.EventTypeConnecting, func(evt *socketmode.Event, client *socketmode.Client) {
		middlewareConnecting(evt, client, desugared)
	})
	socketmodeHandler.Handle(socketmode.EventTypeConnectionError, func(evt *socketmode.Event, client *socketmode.Client) {
		middlewareConnectionError(evt, client, desugared)
	})
	socketmodeHandler.Handle(socketmode.EventTypeConnected, func(evt *socketmode.Event, client *socketmode.Client) {
		middlewareConnected(evt, client, desugared)
	})
	socketmodeHandler.Handle(socketmode.EventTypeHello, func(evt *socketmode.Event, client *socketmode.Client) {
		middlewareHello(evt, client, desugared)
	})
	socketmodeHandler.HandleEvents(slackevents.AppMention, func(evt *socketmode.Event, client *socketmode.Client) {
		middlewareAppMentionEvent(evt, client, gptClient, ctx, desugared, convo)
	})
	socketmodeHandler.HandleEvents(slackevents.Message, func(evt *socketmode.Event, client *socketmode.Client) {
		middlewareMessageEvent(evt, client, gptClient, ctx, desugared, convo)
	})
	return socketmodeHandler.RunEventLoop()
}

// conversation stores user+channel conversations for up to 4 q+a before cycling out
type conversation map[string][]string

// newConversation creates a new conversation
func newConversation() conversation {
	convo := make(map[string][]string)
	return convo
}

// UpdateConversation stores records of 4 questions and answers for a given user channel combination
// to feed into the chatgpt API to enable conversations
func (c conversation) UpdateConversation(userChannel, chatText string) {
	// new userChannel combo
	if _, ok := c[userChannel]; !ok {
		c[userChannel] = append(c[userChannel], chatText)
		return
	}
	// this is around the maximum chat buffer chatgpt API can handle given 4096 tokens
	if len(c[userChannel]) < 8 {
		c[userChannel] = append(c[userChannel], chatText)
	} else {
		c[userChannel] = c[userChannel][1:]
		c[userChannel] = append(c[userChannel], chatText)
	}
}
