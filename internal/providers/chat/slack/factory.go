package slack

import (
	"github.com/massimo-ua/quill/internal/domain/ports"
)

// Factory creates and configures Slack clients
type Factory struct {
	config *Config
}

// NewFactory creates a new Slack client factory
func NewFactory(config *Config) *Factory {
	return &Factory{
		config: config,
	}
}

// CreateChatProvider creates a new Slack client that implements the ChatAccessProvider interface
func (f *Factory) CreateChatProvider() (ports.ChatAccessProvider, error) {
	return NewClient(f.config)
}