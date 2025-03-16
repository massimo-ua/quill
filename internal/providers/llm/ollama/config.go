package ollama

import (
	"errors"
	"strings"
)

var (
	ErrMissingServerURL  = errors.New("Ollama server URL is required")
	ErrMissingModelName  = errors.New("model name is required")
	ErrInvalidTemperature = errors.New("temperature must be between 0 and 2")
	ErrInvalidMaxTokens   = errors.New("max tokens must be greater than 0")
)

// Config contains Ollama API configuration
type Config struct {
	// ServerURL is the Ollama API endpoint (e.g., "http://localhost:11434")
	ServerURL string

	// Model is the model to use (e.g., "llama2", "mistral", "gemma")
	Model string

	// Temperature controls randomness (0-2, default: 0.7)
	Temperature float64

	// MaxTokens is the maximum number of tokens to generate (default: 1024)
	MaxTokens int

	// SystemPrompt is the default system prompt to use (optional)
	SystemPrompt string
}

// NewDefaultConfig creates a Config with default values
func NewDefaultConfig(serverURL, model string) *Config {
	return &Config{
		ServerURL:   serverURL,
		Model:       model,
		Temperature: 0.7,
		MaxTokens:   1024,
	}
}

// Validate checks if the configuration is valid
func (c *Config) Validate() error {
	if strings.TrimSpace(c.ServerURL) == "" {
		return ErrMissingServerURL
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