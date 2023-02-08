package gptslack

import (
	"context"
	"github.com/PullRequestInc/go-gpt3"
	fakes "github.com/PullRequestInc/go-gpt3/go-gpt3fakes"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

func fakeHttpClient() (*fakes.FakeRoundTripper, *http.Client) {
	rt := &fakes.FakeRoundTripper{}
	return rt, &http.Client{
		Transport: rt,
	}
}

func TestEventHandlerNoHandles(t *testing.T) {
	appToken := "xapp-123"
	botTok := "xoxb-test"
	ctx := context.Background()
	_, httpClient := fakeHttpClient()
	client := gpt3.NewClient("test-key", gpt3.WithHTTPClient(httpClient))
	err := EventHandler(appToken, botTok, client, ctx)
	require.ErrorContains(t, err, "invalid_auth")
}
