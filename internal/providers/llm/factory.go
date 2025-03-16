package llm

import (
	"fmt"
	"github.com/massimo-ua/quill/internal/domain/ports"
	"github.com/massimo-ua/quill/internal/providers/llm/ollama"
	"github.com/massimo-ua/quill/internal/providers/llm/openai"
)

// ProviderType represents the type of LLM provider
type ProviderType string

const (
	// ProviderTypeOpenAI represents the OpenAI provider
	ProviderTypeOpenAI ProviderType = "openai"
	// ProviderTypeOllama represents the Ollama provider
	ProviderTypeOllama ProviderType = "ollama"
)

// Config contains configuration for creating an LLM provider
type Config struct {
	// Type of provider (openai or ollama)
	Type ProviderType

	// OpenAI-specific configuration
	OpenAI *openai.Config

	// Ollama-specific configuration
	Ollama *ollama.Config
}

// NewLLMProvider creates a new AiAgentProvider based on the specified provider type
func NewLLMProvider(cfg *Config) (ports.AiAgentProvider, error) {
	if cfg == nil {
		return nil, fmt.Errorf("config cannot be nil")
	}

	switch cfg.Type {
	case ProviderTypeOpenAI:
		if cfg.OpenAI == nil {
			return nil, fmt.Errorf("OpenAI config cannot be nil for OpenAI provider")
		}
		return openai.NewOpenAIProvider(cfg.OpenAI)
	case ProviderTypeOllama:
		if cfg.Ollama == nil {
			return nil, fmt.Errorf("Ollama config cannot be nil for Ollama provider")
		}
		return ollama.NewOllamaProvider(cfg.Ollama)
	default:
		return nil, fmt.Errorf("unsupported provider type: %s", cfg.Type)
	}
}