// Package gptslack handles slack events
package gptslack

import (
	"context"
	"fmt"
	"github.com/PullRequestInc/go-gpt3"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
	"log"
	"os"
	"strings"
)

const (
	noAppTokenExit      = 1
	wrongAppTokenPrefix = 2
	noBotTokenExit      = 3
	wrongBotTokenPrefix = 4
)

// EventHandler handles slack events
func EventHandler(appToken string, botToken string, gptClient gpt3.Client, ctx context.Context) {
	if appToken == "" {
		fmt.Println("need an app token to listen to events")
		os.Exit(noAppTokenExit)
	}

	if !strings.HasPrefix(appToken, "xapp-") {
		fmt.Println("slack app tokens start with xapp- but the one passed does not. Exiting")
		os.Exit(wrongAppTokenPrefix)
	}

	if botToken == "" {
		fmt.Println("need a bot token to interact with workspace")
		os.Exit(noBotTokenExit)
	}

	if !strings.HasPrefix(botToken, "xoxb-") {
		fmt.Println("slack bot tokens start with xoxb- but the one passed does not.")
		os.Exit(wrongBotTokenPrefix)
	}

	api := slack.New(
		botToken,
		slack.OptionDebug(true),
		slack.OptionLog(log.New(os.Stdout, "api: ", log.Lshortfile|log.LstdFlags)),
		slack.OptionAppLevelToken(appToken),
	)
	channelID, timestampe, err := api.PostMessage(
		"#slack-chat-gpt-app",
		slack.MsgOptionText("Hello from your bot", false),
	)
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}

	fmt.Printf("Message successfully sent to channel %s at %s", channelID, timestampe)
	client := socketmode.New(
		api,
		socketmode.OptionDebug(true),
		socketmode.OptionLog(log.New(os.Stdout, "socketmode: ", log.Lshortfile|log.LstdFlags)),
	)
	socketmodeHandler := socketmode.NewSocketmodeHandler(client)
	socketmodeHandler.Handle(socketmode.EventTypeConnecting, middlewareConnecting)
	socketmodeHandler.Handle(socketmode.EventTypeConnectionError, middlewareConnectionError)
	socketmodeHandler.Handle(socketmode.EventTypeConnected, middlewareConnected)
	socketmodeHandler.Handle(socketmode.EventTypeHello, func(evt *socketmode.Event, client *socketmode.Client) {
		fmt.Println("Hello received from hello handler")
	})
	socketmodeHandler.HandleEvents(slackevents.AppMention, func(evt *socketmode.Event, client *socketmode.Client) {
		middlewareAppMentionEvent(evt, client, gptClient, ctx)
	})
	//socketmodeHandler.Handle(socketmode.EventTypeEventsAPI, func(evt *socketmode.Event, client *socketmode.Client) {
	//	middlewareEventsAPI(evt, client, gptClient, ctx)
	//})
	socketmodeHandler.RunEventLoop()
}
