package openai

import (
	"errors"
	"strings"
)

var (
	ErrMissingAPIKey       = errors.New("OpenAI API key is required")
	ErrMissingModelName    = errors.New("model name is required")
	ErrInvalidTemperature  = errors.New("temperature must be between 0 and 2")
	ErrInvalidMaxTokens    = errors.New("max tokens must be greater than 0")
)

// Config contains OpenAI API configuration
type Config struct {
	// APIKey is the OpenAI API key
	APIKey string

	// BaseURL is the custom API endpoint (optional, uses OpenAI default if empty)
	BaseURL string

	// Model is the model to use (e.g., "gpt-4", "gpt-3.5-turbo")
	Model string

	// Temperature controls randomness (0-2, default: 0.7)
	Temperature float64

	// MaxTokens is the maximum number of tokens to generate (default: 1024)
	MaxTokens int

	// Organization is the OpenAI organization ID (optional)
	Organization string
}

// NewDefaultConfig creates a Config with default values
func NewDefaultConfig(apiKey, model string) *Config {
	return &Config{
		APIKey:      apiKey,
		Model:       model,
		Temperature: 0.7,
		MaxTokens:   1024,
	}
}

// Validate checks if the configuration is valid
func (c *Config) Validate() error {
	if strings.TrimSpace(c.APIKey) == "" {
		return ErrMissingAPIKey
	}

	if strings.TrimSpace(c.Model) == "" {
		return ErrMissingModelName
	}

	if c.Temperature < 0 || c.Temperature > 2 {
		return ErrInvalidTemperature
	}

	if c.MaxTokens <= 0 {
		return ErrInvalidMaxTokens
	}

	return nil
}