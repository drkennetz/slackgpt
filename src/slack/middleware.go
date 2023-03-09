package slackhandler

import (
	"context"
	"github.com/sashabaranov/go-openai"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
	"go-slack-chat-gpt3/src/chatgpt"
	"log"
	"strings"
)

// TODO: implement middleware wrapper that handles everything through inner event data

func middlewareConnecting(evt *socketmode.Event, client *socketmode.Client, logger *log.Logger) {
	logger.Println("Connecting")
}

func middlewareConnectionError(evt *socketmode.Event, client *socketmode.Client, logger *log.Logger) {
	logger.Println("Connection failed. Retrying later...")
}

func middlewareConnected(evt *socketmode.Event, client *socketmode.Client, logger *log.Logger) {
	logger.Println("Connected to Slack with Socket Mode.")
}

func middlewareHello(evt *socketmode.Event, client *socketmode.Client, logger *log.Logger) {
	logger.Println("Hello received from hello handler")
}

// we have to org this in such a way that this part does the chatGPT stuff
// but it needs the tokens from the environment
func middlewareAppMentionEvent(evt *socketmode.Event, client *socketmode.Client,
	gptClient *openai.Client, ctx context.Context, logger *log.Logger, convo *conversation) {
	logger.Println("Hello from AppMention middleware")
	eventsAPIEvent, ok := evt.Data.(slackevents.EventsAPIEvent)
	if !ok {
		logger.Printf("Ignored %+v\n", evt)
		return
	}

	client.Ack(*evt.Request)
	ev, ok := eventsAPIEvent.InnerEvent.Data.(*slackevents.AppMentionEvent)
	if !ok {
		logger.Printf("Ignored %+v\n", ev)
		return
	}
	logger.Printf("we have been mentioned in %v\n", ev.Channel)
	logger.Println(ev)
	if ev.ThreadTimeStamp == "" {
		ev.ThreadTimeStamp = ev.TimeStamp
	}
	// found a unique way to identify a thread
	timestampUserChannel := ev.ThreadTimeStamp + ev.User + ev.Channel
	log.Printf("timestamp: %v\n", ev.TimeStamp)
	log.Printf("thread_timestamp: %v\n", ev.ThreadTimeStamp)
	convo.UpdateConversation(timestampUserChannel, ev.Text)
	gpt3Resp, err := chatgpt.GetStringResponse(gptClient, ctx, convo.data[timestampUserChannel])
	convo.UpdateConversation(timestampUserChannel, gpt3Resp)
	if err != nil {
		logger.Printf("Failed to get gpt3 response: %v\n", err)
		gpt3Resp = "I'm having some trouble communicating with our servers (my brain). Please try again in a little bit and hopefully the fuzz clears up."
	}
	_, _, err = client.Client.PostMessage(ev.Channel,
		slack.MsgOptionText(strings.Join([]string{"```", gpt3Resp, "```"}, ""), false),
		slack.MsgOptionTS(ev.ThreadTimeStamp))
	if err != nil {
		logger.Printf("failed posting message: %v", err)
		return
	}
}

func middlewareMessageEvent(evt *socketmode.Event, client *socketmode.Client, gptClient *openai.Client, ctx context.Context, logger *log.Logger, convo *conversation) {
	logger.Println("Hello from Message middleware")
	eventsAPIEvent, ok := evt.Data.(slackevents.EventsAPIEvent)
	// only handle non-bot-id-events
	if !ok {
		logger.Printf("Ignored %+v\n", evt)
		return
	}

	client.Ack(*evt.Request)
	ev, ok := eventsAPIEvent.InnerEvent.Data.(*slackevents.MessageEvent)
	logger.Println(ev)
	if !ok {
		logger.Printf("Ignored %+v\n", evt)
		return
	}
	if ev.BotID != "" {
		return
	}
	userChannel := ev.Username + ev.Channel
	convo.UpdateConversation(userChannel, ev.Text)
	gpt3Resp, err := chatgpt.GetStringResponse(gptClient, ctx, convo.data[userChannel])
	if err != nil {
		logger.Printf("Failed to get gpt3 response: %v\n", err)
		gpt3Resp = "I'm having some trouble communicating with our servers (my brain). Please try again in a little bit and hopefully the fuzz clears up."
	}
	convo.UpdateConversation(userChannel, gpt3Resp)
	_, _, err = client.Client.PostMessage(ev.Channel, slack.MsgOptionText(strings.Join([]string{"```", gpt3Resp, "```"}, ""), false))
	if err != nil {
		logger.Printf("failed posting message: %v\n", err)
		return
	}
}
