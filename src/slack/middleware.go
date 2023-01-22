package gptslack

import (
	"context"
	"fmt"
	"github.com/PullRequestInc/go-gpt3"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
	"go-slack-chat-gpt3/src/chatgpt"
)

func middlewareConnecting(evt *socketmode.Event, client *socketmode.Client) {
	fmt.Println("Connecting to Slack with Socket Mode...")
}

func middlewareConnectionError(evt *socketmode.Event, client *socketmode.Client) {
	fmt.Println("Connection failed. Retrying later...")
}

func middlewareConnected(evt *socketmode.Event, client *socketmode.Client) {
	fmt.Println("Connected to Slack with Socket Mode.")
}

// we have to org this in such a way that this part does the chatGPT stuff
// but it needs the tokens from the environment
func middlewareAppMentionEvent(evt *socketmode.Event, client *socketmode.Client, gptClient gpt3.Client, ctx context.Context) {
	fmt.Println("Hello from AppMention middleware")
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
	fmt.Printf("we have been mentioned in %v\n", ev.Channel)
	question := ev.Text
	gpt3Resp, err := chatgpt.GetStringResponse(gptClient, ctx, question)
	if err != nil {
		fmt.Println("Failed to get gpt3 response", err)
		return
	}
	_, _, err = client.Client.PostMessage(ev.Channel, slack.MsgOptionText(gpt3Resp, false))
	if err != nil {
		fmt.Printf("failed posting message: %v", err)
		return
	}
}
