package ollama

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const (
	// DefaultTimeout is the default timeout for HTTP requests
	DefaultTimeout = 120 * time.Second
)

// Message represents a chat message
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// GenerateRequest represents an Ollama generate request
type GenerateRequest struct {
	Model       string    `json:"model"`
	Prompt      string    `json:"prompt,omitempty"`
	System      string    `json:"system,omitempty"`
	Messages    []Message `json:"messages,omitempty"`
	Stream      bool      `json:"stream,omitempty"`
	Temperature float64   `json:"temperature,omitempty"`
	MaxTokens   int       `json:"max_tokens,omitempty"`
}

// GenerateResponse represents an Ollama generate response
type GenerateResponse struct {
	Model     string `json:"model"`
	Response  string `json:"response"`
	Done      bool   `json:"done"`
	Context   []int  `json:"context,omitempty"`
	TotalDuration int64 `json:"total_duration,omitempty"`
	LoadDuration  int64 `json:"load_duration,omitempty"`
	PromptEvalCount   int  `json:"prompt_eval_count,omitempty"`
	PromptEvalDuration int64 `json:"prompt_eval_duration,omitempty"`
	EvalCount   int  `json:"eval_count,omitempty"`
	EvalDuration int64 `json:"eval_duration,omitempty"`
}

// Client represents an Ollama API client
type Client struct {
	config     *Config
	httpClient *http.Client
}

// NewClient creates a new Ollama API client
func NewClient(cfg *Config) (*Client, error) {
	if cfg == nil {
		return nil, fmt.Errorf("config cannot be nil")
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	httpClient := &http.Client{
		Timeout: DefaultTimeout,
	}

	return &Client{
		config:     cfg,
		httpClient: httpClient,
	}, nil
}

// GenerateCompletion sends a generate request to the Ollama API with prompt
func (c *Client) GenerateCompletion(ctx context.Context, prompt string) (string, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	endpoint := fmt.Sprintf("%s/api/generate", strings.TrimRight(c.config.ServerURL, "/"))

	request := GenerateRequest{
		Model:       c.config.Model,
		Prompt:      prompt,
		System:      c.config.SystemPrompt,
		Temperature: c.config.Temperature,
		MaxTokens:   c.config.MaxTokens,
	}

	return c.sendGenerateRequest(ctx, endpoint, request)
}

// GenerateChatCompletion sends a generate request to the Ollama API with messages
func (c *Client) GenerateChatCompletion(ctx context.Context, messages []Message) (string, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	endpoint := fmt.Sprintf("%s/api/chat", strings.TrimRight(c.config.ServerURL, "/"))

	request := GenerateRequest{
		Model:       c.config.Model,
		Messages:    messages,
		System:      c.config.SystemPrompt,
		Temperature: c.config.Temperature,
		MaxTokens:   c.config.MaxTokens,
	}

	return c.sendGenerateRequest(ctx, endpoint, request)
}

func (c *Client) sendGenerateRequest(ctx context.Context, endpoint string, request GenerateRequest) (string, error) {
	jsonData, err := json.Marshal(request)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, body)
	}

	var generateResponse GenerateResponse
	if err := json.Unmarshal(body, &generateResponse); err != nil {
		return "", fmt.Errorf("failed to parse response: %w", err)
	}

	return generateResponse.Response, nil
}