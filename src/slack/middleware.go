package gptslack

import (
	"context"
	gogpt "github.com/sashabaranov/go-gpt3"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
	"go-slack-chat-gpt3/src/chatgpt"
	"log"
	"strings"
)

// TODO: implement middleware wrapper that handles everything through inner event data
// TODO: think of data structure that will handle:
//
//	base question from thread, answer from bot, question from user, answer from bot, question from user, answer from bot
//	this should store the last 4 conversation items (q+a) (roughly)
//	in case another user interleaves a question in while bot is responding
//	bot should be able to keep track of separate conversations even within a thread
//	start with dms because only 1 user
//	might use map[threadid][userid] = type struct conversation { chat []string }
//	func (c *conversation) Update(question, answer string) {
//	  if len(c.chat) < 7 { c.chat = append(c.chat, question); c.chat = append(c.chat, answer); }
//	  else { c.chat = c.chat[2:]; c.chat = append(c.chat, question); c.chat = append(c.chat, answer); }
//	}
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
	gptClient *gogpt.Client, ctx context.Context, logger *log.Logger, convo conversation) {
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

	userChannel := ev.User + ev.Channel
	convo.UpdateConversation(userChannel, ev.Text)
	gpt3Resp, err := chatgpt.GetStringResponse(gptClient, ctx, convo[userChannel])
	convo.UpdateConversation(userChannel, gpt3Resp)
	if err != nil {
		logger.Printf("Failed to get gpt3 response: %v\n", err)
		gpt3Resp = "I'm having some trouble communicating with our servers (my brain). Please try again in a little bit and hopefully the fuzz clears up."
	}
	_, _, err = client.Client.PostMessage(ev.Channel, slack.MsgOptionText(strings.Join([]string{"```", gpt3Resp, "```"}, ""), false))
	if err != nil {
		logger.Printf("failed posting message: %v", err)
		return
	}
}

func middlewareMessageEvent(evt *socketmode.Event, client *socketmode.Client, gptClient *gogpt.Client, ctx context.Context, logger *log.Logger, convo conversation) {
	logger.Println("Hello from Message middleware")
	eventsAPIEvent, ok := evt.Data.(slackevents.EventsAPIEvent)
	// only handle non-bot-id-events
	if !ok {
		logger.Printf("Ignored %+v\n", evt)
		return
	}

	client.Ack(*evt.Request)
	ev, ok := eventsAPIEvent.InnerEvent.Data.(*slackevents.MessageEvent)
	if !ok {
		logger.Printf("Ignored %+v\n", evt)
		return
	}
	if ev.BotID != "" {
		return
	}
	userChannel := ev.Username + ev.Channel
	convo.UpdateConversation(userChannel, ev.Text)
	gpt3Resp, err := chatgpt.GetStringResponse(gptClient, ctx, convo[userChannel])
	if err != nil {
		logger.Printf("Failed to get gpt3 response: %v\n", err)
		gpt3Resp = "I'm having some trouble communicating with our servers (my brain). Please try again in a little bit and hopefully the fuzz clears up."
	}
	convo.UpdateConversation(userChannel, gpt3Resp)
	logger.Println(convo)
	_, _, err = client.Client.PostMessage(ev.Channel, slack.MsgOptionText(strings.Join([]string{"```", gpt3Resp, "```"}, ""), false))
	if err != nil {
		logger.Printf("failed posting message: %v\n", err)
		return
	}
}
