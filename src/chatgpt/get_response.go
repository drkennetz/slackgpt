package chatgpt

import (
	"context"
	"errors"
	"github.com/sashabaranov/go-openai"
	"strings"
)

// GPTClient implements CreateCompletion from gogpt.Client for testing and future methods
type GPTClient interface {
	CreateCompletion(ctx context.Context, req openai.CompletionRequest) (response openai.CompletionResponse, err error)
}

// ErrorEmptyPrompt implements an Error raised by passing an empty prompt
var ErrorEmptyPrompt error = errors.New("Error empty prompt")

// GetStringResponse sends a completion request to the GPT-3 API to generate a response
// for a given conversation using the specified GPT-3 model. The function takes in a GPT-3
// client, a context, and a slice of strings representing the conversation.
//
// If the length of the conversation slice is 0, an error called ErrorEmptyPrompt is returned.
//
// The function returns the generated response text from the GPT-3 API as a string, with any leading
// or trailing spaces removed using strings.TrimSpace().
//
// Parameters:
// - client: a GPT-3 client object used to make API requests
// - ctx: a context object used to handle timeouts and cancellations
// - chat: a slice of strings representing the conversation
//
// Returns:
// - a string containing the generated response from the GPT-3 API
// - an error, if any
func GetStringResponse(client GPTClient, ctx context.Context, chat []string) (string, error) {
	if len(chat) == 0 {
		return "", ErrorEmptyPrompt
	}

	req := openai.CompletionRequest{
		Model:       openai.GPT3TextDavinci003,
		Prompt:      strings.Join(chat, " "),
		MaxTokens:   2000,
		Temperature: 0,
	}
	resp, err := client.CreateCompletion(ctx, req)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(resp.Choices[0].Text), nil
}
