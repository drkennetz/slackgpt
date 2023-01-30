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
	responseText := strings.Join([]string{"Hello, here's what I found from chat-gpt:", "```", gpt3Resp, "```"}, " ")
	_, _, err = client.Client.PostMessage(ev.Channel, slack.MsgOptionText(responseText, false))
	if err != nil {
		fmt.Printf("failed posting message: %v", err)
		return
	}
}
