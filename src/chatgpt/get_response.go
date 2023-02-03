package chatgpt

import (
	"context"
	"github.com/PullRequestInc/go-gpt3"
)

func GetStringResponse(client gpt3.Client, ctx context.Context, question string) (string, error) {
	resp, err := client.CompletionWithEngine(ctx, gpt3.TextDavinci003Engine, gpt3.CompletionRequest{
		Prompt: []string{
			question,
		},
		MaxTokens:   gpt3.IntPtr(3000),
		Temperature: gpt3.Float32Ptr(0),
	})
	if err != nil {
		return "", err
	}
	response := resp.Choices[0].Text
	return response[2:], nil
}
