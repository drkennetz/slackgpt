package slackhandler

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
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
