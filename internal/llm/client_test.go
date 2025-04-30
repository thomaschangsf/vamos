package llm

import (
	"context"
	"errors"
	"testing"

	"github.com/sashabaranov/go-openai"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockOpenAIClient implements the OpenAIClientInterface
type MockOpenAIClient struct {
	mock.Mock
}

func (m *MockOpenAIClient) CreateChatCompletion(ctx context.Context, req openai.ChatCompletionRequest) (openai.ChatCompletionResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return openai.ChatCompletionResponse{}, args.Error(1)
	}
	return args.Get(0).(openai.ChatCompletionResponse), args.Error(1)
}

func TestNewClient(t *testing.T) {
	client := NewClient("test-api-key", "gpt-3.5-turbo")
	assert.NotNil(t, client)
	assert.Equal(t, "gpt-3.5-turbo", client.model)
	assert.NotNil(t, client.client)
}

func TestGenerateText(t *testing.T) {
	tests := []struct {
		name          string
		prompt        string
		response      string
		expectedError error
	}{
		{
			name:          "Success",
			prompt:        "What is the capital of France?",
			response:      "The capital of France is Paris.",
			expectedError: nil,
		},
		{
			name:          "Error from API",
			prompt:        "Invalid prompt",
			response:      "",
			expectedError: errors.New("API error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockOpenAI := new(MockOpenAIClient)
			client := &Client{
				model:  "gpt-3.5-turbo",
				client: mockOpenAI,
			}

			response := openai.ChatCompletionResponse{
				Choices: []openai.ChatCompletionChoice{
					{
						Message: openai.ChatCompletionMessage{
							Content: tt.response,
						},
					},
				},
			}

			mockOpenAI.On("CreateChatCompletion", mock.Anything, mock.Anything).Return(response, tt.expectedError)

			result, err := client.GenerateText(context.Background(), tt.prompt)
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.response, result)
			}

			mockOpenAI.AssertExpectations(t)
		})
	}
}
