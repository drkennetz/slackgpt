package gptslack

import (
	"context"
	"github.com/PullRequestInc/go-gpt3"
	fakes "github.com/PullRequestInc/go-gpt3/go-gpt3fakes"
	"github.com/magiconair/properties/assert"
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

func TestEventHandlerNoAppTok(t *testing.T) {
	appToken := ""
	botTok := ""
	ctx := context.Background()
	_, httpClient := fakeHttpClient()
	client := gpt3.NewClient("test-key", gpt3.WithHTTPClient(httpClient))
	assert.Panic(t, func() { EventHandler(appToken, botTok, client, ctx) }, "need an app token to listen to events")
}

func TestEventHandlerBadAppTok(t *testing.T) {
	appToken := "test"
	botTok := ""
	ctx := context.Background()
	_, httpClient := fakeHttpClient()
	client := gpt3.NewClient("test-key", gpt3.WithHTTPClient(httpClient))
	assert.Panic(t, func() { EventHandler(appToken, botTok, client, ctx) }, "slack app tokens start with xapp- but the one passed does not")
}

func TestEventHandlerNoBotTok(t *testing.T) {
	appToken := "xapp-123"
	botTok := ""
	ctx := context.Background()
	_, httpClient := fakeHttpClient()
	client := gpt3.NewClient("test-key", gpt3.WithHTTPClient(httpClient))
	assert.Panic(t, func() { EventHandler(appToken, botTok, client, ctx) }, "need a bot token to interact with workspace")
}

func TestEventHandlerBadBotTok(t *testing.T) {
	appToken := "xapp-123"
	botTok := "test"
	ctx := context.Background()
	_, httpClient := fakeHttpClient()
	client := gpt3.NewClient("test-key", gpt3.WithHTTPClient(httpClient))
	assert.Panic(t, func() { EventHandler(appToken, botTok, client, ctx) }, "slack bot tokens start with xoxb- but the one passed does not.")
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
