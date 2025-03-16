package openai

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		config  *Config
		wantErr bool
	}{
		{
			name: "valid config",
			config: &Config{
				APIKey:      "sk-test123",
				Model:       "gpt-4",
				Temperature: 0.7,
				MaxTokens:   1024,
			},
			wantErr: false,
		},
		{
			name: "missing API key",
			config: &Config{
				Model:       "gpt-4",
				Temperature: 0.7,
				MaxTokens:   1024,
			},
			wantErr: true,
		},
		{
			name: "missing model",
			config: &Config{
				APIKey:      "sk-test123",
				Temperature: 0.7,
				MaxTokens:   1024,
			},
			wantErr: true,
		},
		{
			name: "temperature too low",
			config: &Config{
				APIKey:      "sk-test123",
				Model:       "gpt-4",
				Temperature: -0.1,
				MaxTokens:   1024,
			},
			wantErr: true,
		},
		{
			name: "temperature too high",
			config: &Config{
				APIKey:      "sk-test123",
				Model:       "gpt-4",
				Temperature: 2.1,
				MaxTokens:   1024,
			},
			wantErr: true,
		},
		{
			name: "max tokens too low",
			config: &Config{
				APIKey:      "sk-test123",
				Model:       "gpt-4",
				Temperature: 0.7,
				MaxTokens:   0,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestNewDefaultConfig(t *testing.T) {
	apiKey := "sk-test123"
	model := "gpt-4"
	
	config := NewDefaultConfig(apiKey, model)
	
	assert.Equal(t, apiKey, config.APIKey)
	assert.Equal(t, model, config.Model)
	assert.Equal(t, 0.7, config.Temperature)
	assert.Equal(t, 1024, config.MaxTokens)
}