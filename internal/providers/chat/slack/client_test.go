package slack

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewClient(t *testing.T) {
	tests := []struct {
		name        string
		config      *Config
		expectError bool
	}{
		{
			name: "creates client with valid config",
			config: &Config{
				BotToken:      "xoxb-test-token",
				AppToken:      "xapp-test-token",
				SigningSecret: "test-signing-secret",
				DebugMode:     true,
			},
			expectError: false,
		},
		{
			name:        "returns error with nil config",
			config:      nil,
			expectError: true,
		},
		{
			name: "returns error with missing bot token",
			config: &Config{
				BotToken:      "",
				AppToken:      "xapp-test-token",
				SigningSecret: "test-signing-secret",
			},
			expectError: true,
		},
		{
			name: "returns error with missing app token",
			config: &Config{
				BotToken:      "xoxb-test-token",
				AppToken:      "",
				SigningSecret: "test-signing-secret",
			},
			expectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			client, err := NewClient(tc.config)
			if tc.expectError {
				assert.Error(t, err)
				assert.Nil(t, client)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, client)
				assert.Equal(t, tc.config, client.config)
				assert.NotNil(t, client.api)
				assert.NotNil(t, client.socket)
				assert.NotNil(t, client.messageCh)
				assert.NotNil(t, client.threadMap)
			}
		})
	}
}