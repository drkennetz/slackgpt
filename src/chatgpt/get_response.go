package chatgpt

import (
	"context"
	"errors"
	gogpt "github.com/sashabaranov/go-gpt3"
	"strings"
)

var ErrorEmptyPrompt error = errors.New("Error empty prompt")

func GetStringResponse(client *gogpt.Client, ctx context.Context, chat []string) (string, error) {
	if len(chat) == 0 {
		return "", ErrorEmptyPrompt
	}

	req := gogpt.CompletionRequest{
		Model:       gogpt.GPT3TextDavinci003,
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
