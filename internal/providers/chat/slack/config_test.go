package slack

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewConfig(t *testing.T) {
	tests := []struct {
		name          string
		botToken      string
		appToken      string
		signingSecret string
		debugMode     bool
		expectConfig  *Config
	}{
		{
			name:          "creates config with all values",
			botToken:      "xoxb-test-token",
			appToken:      "xapp-test-token",
			signingSecret: "test-signing-secret",
			debugMode:     true,
			expectConfig: &Config{
				BotToken:      "xoxb-test-token",
				AppToken:      "xapp-test-token",
				SigningSecret: "test-signing-secret",
				DebugMode:     true,
			},
		},
		{
			name:          "creates config with debug mode off",
			botToken:      "xoxb-test-token",
			appToken:      "xapp-test-token",
			signingSecret: "test-signing-secret",
			debugMode:     false,
			expectConfig: &Config{
				BotToken:      "xoxb-test-token",
				AppToken:      "xapp-test-token",
				SigningSecret: "test-signing-secret",
				DebugMode:     false,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			config := NewConfig(tc.botToken, tc.appToken, tc.signingSecret, tc.debugMode)
			assert.Equal(t, tc.expectConfig, config)
		})
	}
}