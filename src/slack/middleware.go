package gptslack

import (
	"context"
	"fmt"
	"github.com/PullRequestInc/go-gpt3"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
	"go-slack-chat-gpt3/src/chatgpt"
	"log"
	"strings"
)

func middlewareConnecting(evt *socketmode.Event, client *socketmode.Client) {
	log.Println("Connecting to Slack with Socket Mode...")
}

func middlewareConnectionError(evt *socketmode.Event, client *socketmode.Client) {
	log.Println("Connection failed. Retrying later...")
}

func middlewareConnected(evt *socketmode.Event, client *socketmode.Client) {
	log.Println("Connected to Slack with Socket Mode.")
}

func middlewareHello(evt *socketmode.Event, client *socketmode.Client) {
	log.Println("Hello received from hello handler")
}

// we have to org this in such a way that this part does the chatGPT stuff
// but it needs the tokens from the environment
func middlewareAppMentionEvent(evt *socketmode.Event, client *socketmode.Client, gptClient gpt3.Client, ctx context.Context) {
	log.Println("Hello from AppMention middleware")
	eventsAPIEvent, ok := evt.Data.(slackevents.EventsAPIEvent)
	if !ok {
		fmt.Printf("Ignored %+v\n", evt)
		return
	}
	client.Ack(*evt.Request)
	ev, ok := eventsAPIEvent.InnerEvent.Data.(*slackevents.AppMentionEvent)
	if !ok {
		fmt.Printf("Ignored %+v\n", ev)
		return
	}
	log.Printf("we have been mentioned in %v\n", ev.Channel)
	question := ev.Text
	gpt3Resp, err := chatgpt.GetStringResponse(gptClient, ctx, question)
	if err != nil {
		fmt.Println("Failed to get gpt3 response", err)
		return
	}
	_, _, err = client.Client.PostMessage(ev.Channel, slack.MsgOptionText(strings.Join([]string{"```", gpt3Resp, "```"}, ""), false))
	if err != nil {
		fmt.Printf("failed posting message: %v", err)
		return
	}
}

func middlewareMessageEvent(evt *socketmode.Event, client *socketmode.Client, gptClient gpt3.Client, ctx context.Context) {
	log.Println("Hello from Message middleware")
	eventsAPIEvent, ok := evt.Data.(slackevents.EventsAPIEvent)
	// only handle non-bot-id-events
	if !ok {
		fmt.Printf("Ignored %+v\n", evt)
		return
	}

	client.Ack(*evt.Request)
	ev, ok := eventsAPIEvent.InnerEvent.Data.(*slackevents.MessageEvent)
	if !ok {
		fmt.Printf("Ignored %+v\n", evt)
		return
	}
	if ev.BotID != "" {
		fmt.Println("The bot will answer itself.")
		return
	}
	question := ev.Text
	gpt3Resp, err := chatgpt.GetStringResponse(gptClient, ctx, question)
	if err != nil {
		fmt.Printf("Failed to get gpt3 response: %v\n", err)
		gpt3Resp = "I'm having some trouble communicating with our servers (my brain). Please try again in a little bit and hopefully the fuzz clears up."
	}

	_, _, err = client.Client.PostMessage(ev.Channel, slack.MsgOptionText(strings.Join([]string{"```", gpt3Resp, "```"}, ""), false))
	if err != nil {
		fmt.Printf("failed posting message: %v\n", err)
		return
	}
}
