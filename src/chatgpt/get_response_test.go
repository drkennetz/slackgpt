package chatgpt

import (
	"context"
	"errors"
	"github.com/PullRequestInc/go-gpt3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

type MockClient struct {
	mock.Mock
}

func (c *MockClient) Engines(ctx context.Context) (*gpt3.EnginesResponse, error) {
	args := c.Called(ctx)
	return args.Get(0).(*gpt3.EnginesResponse), args.Error(1)
}

func (c *MockClient) Engine(ctx context.Context, engine string) (*gpt3.EngineObject, error) {
	args := c.Called(ctx)
	return args.Get(0).(*gpt3.EngineObject), args.Error(1)
}

func (c *MockClient) Completion(ctx context.Context, request gpt3.CompletionRequest) (*gpt3.CompletionResponse, error) {
	args := c.Called(ctx, request)
	return args.Get(0).(*gpt3.CompletionResponse), args.Error(1)
}

func (c *MockClient) CompletionStream(ctx context.Context, request gpt3.CompletionRequest, onData func(*gpt3.CompletionResponse)) error {
	args := c.Called(ctx, request, onData)
	return args.Error(0)
}

func (c *MockClient) CompletionWithEngine(ctx context.Context, engine string, request gpt3.CompletionRequest) (*gpt3.CompletionResponse, error) {
	args := c.Called(ctx, engine, request)
	return args.Get(0).(*gpt3.CompletionResponse), args.Error(1)
}

func (c *MockClient) CompletionStreamWithEngine(ctx context.Context, engine string, request gpt3.CompletionRequest, onData func(*gpt3.CompletionResponse)) error {
	args := c.Called(ctx, engine, request, onData)
	return args.Error(0)
}

func (c *MockClient) Edits(ctx context.Context, request gpt3.EditsRequest) (*gpt3.EditsResponse, error) {
	args := c.Called(ctx, request)
	return args.Get(0).(*gpt3.EditsResponse), args.Error(1)
}

func (c *MockClient) Search(ctx context.Context, request gpt3.SearchRequest) (*gpt3.SearchResponse, error) {
	args := c.Called(ctx, request)
	return args.Get(0).(*gpt3.SearchResponse), args.Error(1)
}

func (c *MockClient) SearchWithEngine(ctx context.Context, engine string, request gpt3.SearchRequest) (*gpt3.SearchResponse, error) {
	args := c.Called(ctx, engine, request)
	return args.Get(0).(*gpt3.SearchResponse), args.Error(1)
}

func (c *MockClient) Embeddings(ctx context.Context, request gpt3.EmbeddingsRequest) (*gpt3.EmbeddingsResponse, error) {
	args := c.Called(ctx, request)
	return args.Get(0).(*gpt3.EmbeddingsResponse), args.Error(1)
}

func TestGetStringResponse(t *testing.T) {
	mockClient := MockClient{}
	ctx := context.Background()
	// define test cases
	testCases := []struct {
		name        string
		question    string
		expected    gpt3.CompletionResponse
		expectedErr error
	}{
		{
			name:     "returns response for valid question",
			question: "What is the meaning of life?",
			expected: gpt3.CompletionResponse{
				Choices: []gpt3.CompletionResponseChoice{
					{Text: "42"},
				},
			},
			expectedErr: nil,
		},
		{
			name:     "returns error for invalid question",
			question: "",
			expected: gpt3.CompletionResponse{
				Choices: []gpt3.CompletionResponseChoice{
					{Text: ""},
				},
			},
			expectedErr: ErrorEmptyPrompt,
		},
		{
			name:     "simulates an error from the api call",
			question: "This Forces Fake Error",
			expected: gpt3.CompletionResponse{
				Choices: []gpt3.CompletionResponseChoice{
					{Text: ""},
				},
			},
			expectedErr: errors.New("Simulated err"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// setup our mock client to return a response or error based on the test
			if tc.expectedErr == nil {
				mockClient.On("CompletionWithEngine", ctx, gpt3.TextDavinci003Engine, gpt3.CompletionRequest{
					Prompt:      []string{tc.question},
					MaxTokens:   gpt3.IntPtr(3000),
					Temperature: gpt3.Float32Ptr(0),
				}).Return(&tc.expected, nil)
			} else if tc.question == "" {
				mockClient.On("CompletionWithEngine", ctx, gpt3.TextDavinci003Engine, gpt3.CompletionRequest{
					Prompt:      []string{tc.question},
					MaxTokens:   gpt3.IntPtr(3000),
					Temperature: gpt3.Float32Ptr(0),
				}).Return(&tc.expected, tc.expectedErr)
			} else if tc.question == "This Forces Fake Error" {
				mockClient.On("CompletionWithEngine", ctx, gpt3.TextDavinci003Engine, gpt3.CompletionRequest{
					Prompt:      []string{tc.question},
					MaxTokens:   gpt3.IntPtr(3000),
					Temperature: gpt3.Float32Ptr(0),
				}).Return(&tc.expected, tc.expectedErr)
			}

			response, err := GetStringResponse(&mockClient, ctx, tc.question)
			assert.Equal(t, tc.expected.Choices[0].Text, response)
			if tc.expectedErr == nil {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.expectedErr.Error())
			}

			// assert that the mock client's CompletionWithEngine method was called with the expected arguments
			if tc.question != "" {
				mockClient.AssertCalled(t, "CompletionWithEngine", ctx, gpt3.TextDavinci003Engine, gpt3.CompletionRequest{
					Prompt:      []string{tc.question},
					MaxTokens:   gpt3.IntPtr(3000),
					Temperature: gpt3.Float32Ptr(0),
				})
			}

		})
	}
}
