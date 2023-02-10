package chatgpt

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/PullRequestInc/go-gpt3"
	fakes "github.com/PullRequestInc/go-gpt3/go-gpt3fakes"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"testing"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 net/http.RoundTripper

func fakeHttpClient() (*fakes.FakeRoundTripper, *http.Client) {
	rt := &fakes.FakeRoundTripper{}

	return rt, &http.Client{
		Transport: rt,
	}
}

func TestGetStringResponse(t *testing.T) {
	ctx := context.Background()
	rt, httpClient := fakeHttpClient()

	client := gpt3.NewClient("test-key", gpt3.WithHTTPClient(httpClient))
	completionResponse := &gpt3.CompletionResponse{
		ID:      "ABC",
		Object:  "list",
		Created: 123456789,
		Model:   gpt3.TextDavinci003Engine,
		Choices: []gpt3.CompletionResponseChoice{
			gpt3.CompletionResponseChoice{
				Text:         "why?",
				FinishReason: "stop",
			},
		},
	}
	data, err := json.Marshal(completionResponse)
	assert.NoError(t, err)
	mockResponse := &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(bytes.NewBuffer(data)),
	}
	rt.RoundTripReturns(mockResponse, nil)

	text, err := GetStringResponse(client, ctx, "what?")
	fmt.Println(text)
	fmt.Println(err)
}
