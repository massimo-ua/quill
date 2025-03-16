package ollama

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
				ServerURL:   "http://localhost:11434",
				Model:       "llama2",
				Temperature: 0.7,
				MaxTokens:   1024,
			},
			wantErr: false,
		},
		{
			name: "missing server URL",
			config: &Config{
				Model:       "llama2",
				Temperature: 0.7,
				MaxTokens:   1024,
			},
			wantErr: true,
		},
		{
			name: "missing model",
			config: &Config{
				ServerURL:   "http://localhost:11434",
				Temperature: 0.7,
				MaxTokens:   1024,
			},
			wantErr: true,
		},
		{
			name: "temperature too low",
			config: &Config{
				ServerURL:   "http://localhost:11434",
				Model:       "llama2",
				Temperature: -0.1,
				MaxTokens:   1024,
			},
			wantErr: true,
		},
		{
			name: "temperature too high",
			config: &Config{
				ServerURL:   "http://localhost:11434",
				Model:       "llama2",
				Temperature: 2.1,
				MaxTokens:   1024,
			},
			wantErr: true,
		},
		{
			name: "max tokens too low",
			config: &Config{
				ServerURL:   "http://localhost:11434",
				Model:       "llama2",
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
	serverURL := "http://localhost:11434"
	model := "llama2"
	
	config := NewDefaultConfig(serverURL, model)
	
	assert.Equal(t, serverURL, config.ServerURL)
	assert.Equal(t, model, config.Model)
	assert.Equal(t, 0.7, config.Temperature)
	assert.Equal(t, 1024, config.MaxTokens)
}