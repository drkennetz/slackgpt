package chatgpt

import (
	"context"
	"errors"
	openai "github.com/sashabaranov/go-openai"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

type MockClient struct {
	mock.Mock
}

func (c *MockClient) CreateCompletion(ctx context.Context, req openai.CompletionRequest) (openai.CompletionResponse, error) {
	args := c.Called(ctx, req)
	return args.Get(0).(openai.CompletionResponse), args.Error(1)
}

func TestGetStringResponse(t *testing.T) {
	mockClient := &MockClient{}
	ctx := context.Background()
	// define test cases
	testCases := []struct {
		name        string
		question    string
		expected    openai.CompletionResponse
		expectedErr error
	}{
		{
			name:     "returns response for valid question",
			question: "What is the meaning of life?",
			expected: openai.CompletionResponse{
				Choices: []openai.CompletionChoice{
					{Text: "42"},
				},
			},
			expectedErr: nil,
		},
		{
			name:     "returns error for invalid question",
			question: "",
			expected: openai.CompletionResponse{
				Choices: []openai.CompletionChoice{
					{Text: ""},
				},
			},
			expectedErr: ErrorEmptyPrompt,
		},
		{
			name:     "simulates an error from the api call",
			question: "This Forces Fake Error",
			expected: openai.CompletionResponse{
				Choices: []openai.CompletionChoice{
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
				mockClient.On("CreateCompletion", ctx, openai.CompletionRequest{
					Model:       openai.GPT3TextDavinci003,
					Prompt:      tc.question,
					MaxTokens:   2000,
					Temperature: 0,
				}).Return(tc.expected, nil)
			} else if tc.question == "" {
				mockClient.On("CreateCompletion", ctx, openai.CompletionRequest{
					Model:       openai.GPT3TextDavinci003,
					Prompt:      tc.question,
					MaxTokens:   2000,
					Temperature: 0,
				}).Return(tc.expected, tc.expectedErr)
			} else if tc.question == "This Forces Fake Error" {
				mockClient.On("CreateCompletion", ctx, openai.CompletionRequest{
					Model:       openai.GPT3TextDavinci003,
					Prompt:      tc.question,
					MaxTokens:   2000,
					Temperature: 0,
				}).Return(tc.expected, tc.expectedErr)
			}

			response, err := GetStringResponse(mockClient, ctx, []string{tc.question})
			if tc.question != "" {
				assert.Equal(t, tc.expected.Choices[0].Text, response)
				if tc.expectedErr != nil {
					assert.EqualError(t, err, tc.expectedErr.Error())
				} else {
					assert.NoError(t, err)
					assert.Equal(t, tc.expected.Choices[0].Text, "42")
				}
			} else {
				response, err = GetStringResponse(mockClient, ctx, []string{})
				assert.EqualError(t, err, tc.expectedErr.Error())
				assert.Equal(t, tc.expected.Choices[0].Text, "")
			}

			// assert that the mock client's CompletionWithEngine method was called with the expected arguments
			if tc.question != "" {
				mockClient.AssertCalled(t, "CreateCompletion", ctx, openai.CompletionRequest{
					Model:       openai.GPT3TextDavinci003,
					Prompt:      tc.question,
					MaxTokens:   2000,
					Temperature: 0,
				})
			}

		})
	}
}
