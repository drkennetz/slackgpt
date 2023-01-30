package chatgpt

import (
	"context"
	"fmt"
	"github.com/PullRequestInc/go-gpt3"
	fakes "github.com/PullRequestInc/go-gpt3/go-gpt3fakes"
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
	_, httpClient := fakeHttpClient()
	client := gpt3.NewClient("test-key", gpt3.WithHTTPClient(httpClient))
	text, err := GetStringResponse(client, ctx, "what?")
	fmt.Println(text)
	fmt.Println(err)
}
