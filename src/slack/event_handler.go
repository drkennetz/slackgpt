// Package gptslack handles slack appMention events and responds with chat-gpt response
package gptslack

import (
	"context"
	"github.com/PullRequestInc/go-gpt3"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
	"log"
	"os"
	"strings"
)

// EventHandler handles slack events
func EventHandler(appToken string, botToken string, gptClient gpt3.Client, ctx context.Context) error {
	if appToken == "" {
		panic("need an app token to listen to events")
	}

	if !strings.HasPrefix(appToken, "xapp-") {
		panic("slack app tokens start with xapp- but the one passed does not")
	}

	if botToken == "" {
		panic("need a bot token to interact with workspace")
	}

	if !strings.HasPrefix(botToken, "xoxb-") {
		panic("slack bot tokens start with xoxb- but the one passed does not.")
	}

	api := slack.New(
		botToken,
		slack.OptionDebug(true),
		slack.OptionLog(log.New(os.Stdout, "api: ", log.Lshortfile|log.LstdFlags)),
		slack.OptionAppLevelToken(appToken),
	)

	client := socketmode.New(
		api,
		socketmode.OptionDebug(true),
		socketmode.OptionLog(log.New(os.Stdout, "socketmode: ", log.Lshortfile|log.LstdFlags)),
	)
	socketmodeHandler := socketmode.NewSocketmodeHandler(client)
	socketmodeHandler.Handle(socketmode.EventTypeConnecting, middlewareConnecting)
	socketmodeHandler.Handle(socketmode.EventTypeConnectionError, middlewareConnectionError)
	socketmodeHandler.Handle(socketmode.EventTypeConnected, middlewareConnected)
	socketmodeHandler.Handle(socketmode.EventTypeHello, middlewareHello)
	socketmodeHandler.HandleEvents(slackevents.AppMention, func(evt *socketmode.Event, client *socketmode.Client) {
		middlewareAppMentionEvent(evt, client, gptClient, ctx)
	})

	return socketmodeHandler.RunEventLoop()
}
