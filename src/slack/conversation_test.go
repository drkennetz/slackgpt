package slackhandler

import (
	"bytes"
	"fmt"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
	"time"
)

func TestConversation_Get(t *testing.T) {
	c := newConversation()
	c.AddNew("test", "1")
	for i := 0; i < 10; i++ {
		go c.Get("test")
	}
}

func TestConversation_AddNew(t *testing.T) {
	c := newConversation()
	for i := 0; i < 10; i++ {
		go c.AddNew("test", "1")
	}
}

func TestConversation_AddTextToExisting(t *testing.T) {
	c := newConversation()
	c.AddNew("test", "1")
	for i := 0; i < 10; i++ {
		go c.AddTextToExisting("test", "2")
	}
}

func TestConversation_UpdateConversation(t *testing.T) {
	convo := newConversation()
	userChannel := "user"
	for i := 0; i < 10; i++ {
		tmp := fmt.Sprintf("%s%v", userChannel, i)
		convo.UpdateConversation(userChannel, tmp)
		if i < 8 {
			assert.Equal(t, convo.data[userChannel][0], "user0")
		} else {
			assert.Equal(t, convo.data[userChannel][0], fmt.Sprintf("%s%v", "user", i%7))
		}
	}
}

func TestConversation_ClearConversation_DataRaceOk(t *testing.T) {
	c := newConversation()
	c.AddNew("key1", "1")
	for i := 0; i < 10; i++ {
		go c.ClearConversation("key1")
	}
}

func TestConversation_ClearConversation_ClearsOk(t *testing.T) {
	var c = newConversation()
	c.data = map[string][]string{
		"key1": {"1", "2", "3", "4"},
		"key2": {"1", "2", "3", "4"},
	}
	var isCleared = c.ClearConversation("key1")
	var isNotCleared = c.ClearConversation("badKey")
	assert.Equal(t, 4, len(c.data["key2"]))
	assert.Equal(t, 1, len(c.data))
	assert.True(t, isCleared)
	assert.False(t, isNotCleared)
}

func TestConversation_LogConversationHistoryKvPairs(t *testing.T) {

	// Create a new conversation instance
	var c = newConversation()
	c.data = map[string][]string{"key1": {"value1", "value2"}, "key2": {"value3"}}

	// Create a new buffer to capture the log output
	var buf bytes.Buffer
	log.SetOutput(&buf)

	// Call the LogConversationHistoryKvPairs method with the conversation instance
	c.LogConversationHistoryKvPairs()
	//time now to mimic the timestamp on log.println
	now := time.Now()
	expectedTime := now.Format("2006/01/02 15:04:05")

	// Get the log output from the buffer
	output := buf.String()

	// Check that the log output contains the expected output
	expectedOutput := fmt.Sprintf("%s Key: key1, Value: [value1 value2], Length: 2\n%s Key: key2, Value: [value3], Length: 1\n", expectedTime, expectedTime)
	if output != expectedOutput {
		t.Errorf("Unexpected log output. Expected:\n%s\nActual:\n%s\n", expectedOutput, output)
	}

}
