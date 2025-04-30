package llm

import (
	"context"

	"github.com/sashabaranov/go-openai"
)

// OpenAIClientInterface defines the interface for OpenAI operations
type OpenAIClientInterface interface {
	CreateChatCompletion(ctx context.Context, req openai.ChatCompletionRequest) (openai.ChatCompletionResponse, error)
}

// Client represents an LLM client
type Client struct {
	model  string
	client OpenAIClientInterface
}

// NewClient creates a new LLM client
func NewClient(apiKey, model string) *Client {
	return &Client{
		model:  model,
		client: openai.NewClient(apiKey),
	}
}

// GenerateText generates text using the LLM
func (c *Client) GenerateText(ctx context.Context, prompt string) (string, error) {
	resp, err := c.client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: c.model,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
		},
	)
	if err != nil {
		return "", err
	}

	if len(resp.Choices) == 0 {
		return "", nil
	}

	return resp.Choices[0].Message.Content, nil
}
