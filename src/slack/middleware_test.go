package gptslack

import (
	"bytes"
	"github.com/magiconair/properties/assert"
	"github.com/slack-go/slack/socketmode"
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
