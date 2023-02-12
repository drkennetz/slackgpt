package gptslack

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/PullRequestInc/go-gpt3"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"testing"
)

func captureLog(f func()) string {
	var buf bytes.Buffer
	log.SetFlags(0)
	log.SetOutput(&buf)
	f()
	log.SetOutput(os.Stderr)
	return buf.String()
}

func TestMiddlewareConnecting(t *testing.T) {
	evt := &socketmode.Event{
		Type: socketmode.EventTypeConnecting,
	}
	client := &socketmode.Client{}
	output := captureLog(func() {
		middlewareConnecting(evt, client)
	})
	assert.Equal(t, "Connecting to Slack with Socket Mode...\n", output)
}

func TestMiddlewareConnectionError(t *testing.T) {
	evt := &socketmode.Event{
		Type: socketmode.EventTypeConnectionError,
	}
	client := &socketmode.Client{}
	output := captureLog(func() {
		middlewareConnectionError(evt, client)
	})
	assert.Equal(t, "Connection failed. Retrying later...\n", output)
}

func TestMiddlewareConnected(t *testing.T) {
	evt := &socketmode.Event{
		Type: socketmode.EventTypeConnected,
	}
	client := &socketmode.Client{}
	output := captureLog(func() {
		middlewareConnected(evt, client)
	})
	assert.Equal(t, "Connected to Slack with Socket Mode.\n", output)
}

func TestMiddlewareHello(t *testing.T) {
	evt := &socketmode.Event{
		Type: socketmode.EventTypeHello,
	}
	client := &socketmode.Client{}
	output := captureLog(func() {
		middlewareHello(evt, client)
	})
	assert.Equal(t, "Hello received from hello handler\n", output)
}

// Test1: Change the Data in socketmode.Event so that it isn't EventsAPIEvent
// Test2: Change the InnerEvent type so that it isn't AppMentionEvent
// See if we can mock the chat gpt response??
func TestMiddlewareAppMentionEvent(t *testing.T) {
	type payload struct {
		Text string `json:"text"`
	}
	x := payload{
		Text: "we made it",
	}
	send, err := json.Marshal(x)
	assert.NoError(t, err)
	type args struct {
		event *socketmode.Event
	}
	type expectedResult struct {
		e error
	}
	tests := []struct {
		name string
		arg  args
		want expectedResult
	}{
		{
			name: "nil data",
			arg: args{
				&socketmode.Event{
					Type: socketmode.EventTypeEventsAPI,
					Data: nil,
				},
			},
			want: expectedResult{errors.New("miss")},
		},
		{
			name: "wrong event type",
			arg: args{
				&socketmode.Event{
					Type: socketmode.EventTypeEventsAPI,
					Data: slackevents.EventsAPIEvent{
						Type: "event_callback",
						InnerEvent: slackevents.EventsAPIInnerEvent{
							Type: string(slackevents.Message),
						},
					},
					Request: &socketmode.Request{
						Type:           "test",
						NumConnections: 1,
						ConnectionInfo: socketmode.ConnectionInfo{"test-app"},
						Reason:         "test",
						EnvelopeID:     "1",
						Payload:        send,
					},
				},
			},
			want: expectedResult{
				errors.New("miss"),
			},
		},
		{
			name: "good event",
			arg: args{
				&socketmode.Event{
					Type: socketmode.EventTypeEventsAPI,
					Data: slackevents.EventsAPIEvent{
						Type: "event_callback",
						InnerEvent: slackevents.EventsAPIInnerEvent{
							Type: string(slackevents.AppMention),
							Data: &slackevents.AppMentionEvent{
								Type:    string(slackevents.AppMention),
								User:    "test",
								Text:    "Hello, test!",
								Channel: "app-dev",
							},
						},
					},
					Request: &socketmode.Request{
						Type:           "test",
						NumConnections: 1,
						ConnectionInfo: socketmode.ConnectionInfo{"test-app"},
						Reason:         "test",
						EnvelopeID:     "1",
						Payload:        send,
					},
				},
			},
			want: expectedResult{
				errors.New("miss"),
			},
		},
	}

	slackClient := slack.New("test")
	client := socketmode.New(slackClient)
	_, httpClient := fakeHttpClient()
	gptClient := gpt3.NewClient("test-key", gpt3.WithHTTPClient(httpClient))
	ctx := context.Background()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			middlewareAppMentionEvent(tt.arg.event, client, gptClient, ctx)
		})
	}
}

func TestMiddlewareMessageEvent(t *testing.T) {
	type payload struct {
		Text string `json:"text"`
	}
	x := payload{
		Text: "we made it",
	}
	send, err := json.Marshal(x)
	assert.NoError(t, err)
	type args struct {
		event *socketmode.Event
	}
	type expectedResult struct {
		e error
	}
	tests := []struct {
		name string
		arg  args
		want expectedResult
	}{
		{
			name: "nil data",
			arg: args{
				&socketmode.Event{
					Type: socketmode.EventTypeEventsAPI,
					Data: nil,
				},
			},
			want: expectedResult{errors.New("miss")},
		},
		{
			name: "wrong event type",
			arg: args{
				&socketmode.Event{
					Type: socketmode.EventTypeEventsAPI,
					Data: slackevents.EventsAPIEvent{
						Type: "event_callback",
						InnerEvent: slackevents.EventsAPIInnerEvent{
							Type: string(slackevents.AppMention),
						},
					},
					Request: &socketmode.Request{
						Type:           "test",
						NumConnections: 1,
						ConnectionInfo: socketmode.ConnectionInfo{"test-app"},
						Reason:         "test",
						EnvelopeID:     "1",
						Payload:        send,
					},
				},
			},
			want: expectedResult{
				errors.New("miss"),
			},
		},
		{
			name: "good event",
			arg: args{
				&socketmode.Event{
					Type: socketmode.EventTypeEventsAPI,
					Data: slackevents.EventsAPIEvent{
						Type: "event_callback",
						InnerEvent: slackevents.EventsAPIInnerEvent{
							Type: string(slackevents.Message),
							Data: &slackevents.MessageEvent{
								Type:    string(slackevents.Message),
								User:    "test",
								Text:    "Hello, test!",
								Channel: "app-dev",
							},
						},
					},
					Request: &socketmode.Request{
						Type:           "test",
						NumConnections: 1,
						ConnectionInfo: socketmode.ConnectionInfo{"test-app"},
						Reason:         "test",
						EnvelopeID:     "1",
						Payload:        send,
					},
				},
			},
			want: expectedResult{
				errors.New("miss"),
			},
		},
		{
			name: "bot event",
			arg: args{
				&socketmode.Event{
					Type: socketmode.EventTypeEventsAPI,
					Data: slackevents.EventsAPIEvent{
						Type: "event_callback",
						InnerEvent: slackevents.EventsAPIInnerEvent{
							Type: string(slackevents.Message),
							Data: &slackevents.MessageEvent{
								Type:    string(slackevents.Message),
								User:    "test",
								Text:    "Hello, test!",
								Channel: "app-dev",
								BotID:   "123",
							},
						},
					},
					Request: &socketmode.Request{
						Type:           "test",
						NumConnections: 1,
						ConnectionInfo: socketmode.ConnectionInfo{"test-app"},
						Reason:         "test",
						EnvelopeID:     "1",
						Payload:        send,
					},
				},
			},
			want: expectedResult{
				errors.New("miss"),
			},
		},
	}

	slackClient := slack.New("test")
	client := socketmode.New(slackClient)
	_, httpClient := fakeHttpClient()
	gptClient := gpt3.NewClient("test-key", gpt3.WithHTTPClient(httpClient))
	ctx := context.Background()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			middlewareMessageEvent(tt.arg.event, client, gptClient, ctx)
		})
	}
}
