package slackhandler

import "sync"

// conversation stores user+channel conversations in a concurrency safe way
type conversation struct {
	sync.Mutex
	data map[string][]string
}

// newConversation creates a new conversation
func newConversation() *conversation {
	convo := make(map[string][]string)
	return &conversation{
		data: convo,
	}
}

// UpdateConversation stores records of 4 questions and answers for a given user channel combination
// to feed into the chat-gpt API to enable conversations
func (c *conversation) UpdateConversation(key string, chatText string) {
	// All locking is handled in other methods

	// new userChannel+(thread) combo, guard
	if _, ok := c.Get(key); !ok {
		c.AddNew(key, chatText)
		return
	}

	c.AddTextToExisting(key, chatText)
}

// AddNew is used to add a new value to conversation map
func (c *conversation) AddNew(key, value string) {
	c.Lock()
	defer c.Unlock()
	c.data[key] = append(c.data[key], value)
}

// AddTextToExisting updates an existing conversation
func (c *conversation) AddTextToExisting(key, value string) {
	c.Lock()
	defer c.Unlock()

	// this is around the maximum chat buffer chat-gpt API can handle given 4096 token
	// number may need tweaking if found to be too large
	if len(c.data[key]) < 8 {
		c.data[key] = append(c.data[key], value)
	} else {
		// slice off only the first message to preserve context
		c.data[key] = c.data[key][1:]
		c.data[key] = append(c.data[key], value)
	}
}

// Get safely retrieves a value from a map
func (c *conversation) Get(key string) ([]string, bool) {
	c.Lock()
	defer c.Unlock()
	value, ok := c.data[key]
	return value, ok
}
