package chatgpt

import (
	"context"
	"errors"
	"github.com/PullRequestInc/go-gpt3"
	"strings"
)

var ErrorEmptyPrompt error = errors.New("Error empty prompt")

func GetStringResponse(client gpt3.Client, ctx context.Context, question string) (string, error) {
	if question == "" {
		return "", ErrorEmptyPrompt
	}

	prompt := []string{question}
	req := gpt3.CompletionRequest{
		Prompt:      prompt,
		MaxTokens:   gpt3.IntPtr(3000),
		Temperature: gpt3.Float32Ptr(0),
	}
	resp, err := client.CompletionWithEngine(ctx, gpt3.TextDavinci003Engine, req)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(resp.Choices[0].Text), nil
}
