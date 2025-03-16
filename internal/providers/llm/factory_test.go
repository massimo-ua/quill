package llm

import (
	"testing"

	"github.com/massimo-ua/quill/internal/providers/llm/ollama"
	"github.com/massimo-ua/quill/internal/providers/llm/openai"
	"github.com/stretchr/testify/assert"
)

func TestNewLLMProvider(t *testing.T) {
	tests := []struct {
		name    string
		config  *Config
		wantErr bool
	}{
		{
			name: "valid OpenAI config",
			config: &Config{
				Type: ProviderTypeOpenAI,
				OpenAI: &openai.Config{
					APIKey:      "test-key",
					Model:       "gpt-4",
					Temperature: 0.7,
					MaxTokens:   1024,
				},
			},
			wantErr: false,
		},
		{
			name: "valid Ollama config",
			config: &Config{
				Type: ProviderTypeOllama,
				Ollama: &ollama.Config{
					ServerURL:   "http://localhost:11434",
					Model:       "llama2",
					Temperature: 0.7,
					MaxTokens:   1024,
				},
			},
			wantErr: false,
		},
		{
			name:    "nil config",
			config:  nil,
			wantErr: true,
		},
		{
			name: "missing OpenAI config",
			config: &Config{
				Type: ProviderTypeOpenAI,
			},
			wantErr: true,
		},
		{
			name: "missing Ollama config",
			config: &Config{
				Type: ProviderTypeOllama,
			},
			wantErr: true,
		},
		{
			name: "invalid provider type",
			config: &Config{
				Type: "invalid",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			provider, err := NewLLMProvider(tt.config)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, provider)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, provider)
			}
		})
	}
}