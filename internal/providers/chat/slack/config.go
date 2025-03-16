package slack

// Config contains Slack configuration parameters
type Config struct {
	// BotToken is the Slack API token for the bot
	BotToken string

	// AppToken is the Slack App-level token for Socket Mode
	AppToken string

	// SigningSecret is used to verify incoming requests from Slack
	SigningSecret string

	// DebugMode enables detailed logging when true
	DebugMode bool
}

// NewConfig creates a new Slack configuration
func NewConfig(botToken, appToken, signingSecret string, debugMode bool) *Config {
	return &Config{
		BotToken:      botToken,
		AppToken:      appToken,
		SigningSecret: signingSecret,
		DebugMode:     debugMode,
	}
}