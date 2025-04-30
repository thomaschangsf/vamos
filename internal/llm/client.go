package llm

import (
	"context"
	"fmt"
	"strings"

	"github.com/sashabaranov/go-openai"
)

// Client represents an LLM client
type Client struct {
	client *openai.Client
	model  string
}

// NewClient creates a new LLM client
func NewClient(apiKey, model string) *Client {
	return &Client{
		client: openai.NewClient(apiKey),
		model:  model,
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
		return "", fmt.Errorf("failed to generate text: %v", err)
	}

	return strings.TrimSpace(resp.Choices[0].Message.Content), nil
}
