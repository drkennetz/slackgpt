// Package slackhandler handles slack appMention events and responds with chat-gpt response
package slackhandler

import (
	"context"
	"github.com/sashabaranov/go-openai"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
	"log"
)

type EventHandlerArgs struct {
	Logger           *log.Logger
	SlackClient      *slack.Client
	SocketModeClient *socketmode.Client
	GPTClient        *openai.Client
	Context          context.Context
}

// NewSocketmodeHandler returns a new instance of a socketmode.SocketmodeHandler
func (e *EventHandlerArgs) NewSocketmodeHandler() *socketmode.SocketmodeHandler {
	return socketmode.NewSocketmodeHandler(e.SocketModeClient)
}

// EventHandler handles slack events
func EventHandler(args EventHandlerArgs, handler *socketmode.SocketmodeHandler) error {

	convo := newConversation()

	// should be a primary middleware handler, and these handle more granular events
	handler.Handle(socketmode.EventTypeConnecting, func(evt *socketmode.Event, client *socketmode.Client) {
		middlewareConnecting(evt, client, args.Logger)
	})
	handler.Handle(socketmode.EventTypeConnectionError, func(evt *socketmode.Event, client *socketmode.Client) {
		middlewareConnectionError(evt, client, args.Logger)
	})
	handler.Handle(socketmode.EventTypeConnected, func(evt *socketmode.Event, client *socketmode.Client) {
		middlewareConnected(evt, client, args.Logger)
	})
	handler.Handle(socketmode.EventTypeHello, func(evt *socketmode.Event, client *socketmode.Client) {
		middlewareHello(evt, client, args.Logger)
	})

	handler.HandleEvents(slackevents.AppMention, func(evt *socketmode.Event, client *socketmode.Client) {
		middlewareAppMentionEvent(evt, client, args.GPTClient, args.Context, args.Logger, convo)
	})
	handler.HandleEvents(slackevents.Message, func(evt *socketmode.Event, client *socketmode.Client) {
		middlewareMessageEvent(evt, client, args.GPTClient, args.Context, args.Logger, convo)
	})
	return handler.RunEventLoop()
}
