package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/JuanVilla424/qwen-pet/internal/config"
)

// Message represents a chat message.
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// Usage tracks token consumption.
type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type chatRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	Temperature float64   `json:"temperature,omitempty"`
	MaxTokens   int       `json:"max_tokens,omitempty"`
}

type chatResponse struct {
	Choices []struct {
		Message Message `json:"message"`
	} `json:"choices"`
	Usage Usage `json:"usage"`
}

// Client is an OpenRouter-compatible HTTP client.
type Client struct {
	httpClient *http.Client
	baseURL    string
	apiKey     string
	model      string
	temp       float64
	maxTokens  int
}

// NewClient creates an OpenRouter API client.
func NewClient(cfg config.AICfg, apiKey string) *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: time.Duration(cfg.TimeoutSeconds) * time.Second,
		},
		baseURL:   cfg.BaseURL,
		apiKey:    apiKey,
		model:     cfg.Model,
		temp:      cfg.Temperature,
		maxTokens: cfg.MaxTokens,
	}
}

// ChatCompletion sends a chat completion request and returns the response text.
func (c *Client) ChatCompletion(ctx context.Context, messages []Message) (string, Usage, error) {
	body := chatRequest{
		Model:       c.model,
		Messages:    messages,
		Temperature: c.temp,
		MaxTokens:   c.maxTokens,
	}

	jsonData, err := json.Marshal(body)
	if err != nil {
		return "", Usage{}, fmt.Errorf("marshal request: %w", err)
	}

	text, usage, err := c.doRequest(ctx, jsonData)
	if err != nil {
		slog.Warn("first attempt failed, retrying", "error", err)
		time.Sleep(2 * time.Second)
		return c.doRequest(ctx, jsonData)
	}
	return text, usage, nil
}

func (c *Client) doRequest(ctx context.Context, jsonData []byte) (string, Usage, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL, bytes.NewReader(jsonData))
	if err != nil {
		return "", Usage{}, fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", Usage{}, fmt.Errorf("http request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", Usage{}, fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", Usage{}, fmt.Errorf("API error (HTTP %d): %s", resp.StatusCode, string(respBody))
	}

	var chatResp chatResponse
	if err := json.Unmarshal(respBody, &chatResp); err != nil {
		return "", Usage{}, fmt.Errorf("parse response: %w", err)
	}

	if len(chatResp.Choices) == 0 {
		return "", Usage{}, fmt.Errorf("no choices in response")
	}

	return chatResp.Choices[0].Message.Content, chatResp.Usage, nil
}
