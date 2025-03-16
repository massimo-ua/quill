# Slack Integration for Quill

This package implements the Slack integration for the Quill documentation bot.

## Setup

### 1. Create a Slack App

1. Go to [api.slack.com/apps](https://api.slack.com/apps) and create a new app.
2. Select "From scratch" and provide a name (e.g., "Quill") and workspace.

### 2. Configure Socket Mode

1. Navigate to "Socket Mode" in the sidebar and enable it.
2. Generate an app-level token with `connections:write` scope and save it.

### 3. Configure Event Subscriptions

1. Go to "Event Subscriptions" in the sidebar.
2. Enable events.
3. Subscribe to bot events:
   - `message.channels` - To receive channel messages
   - `message.groups` - For private channel messages
   - `message.im` - For direct messages

### 4. Configure Bot Permissions

1. Go to "OAuth & Permissions" in the sidebar.
2. Add the following scopes:
   - `channels:history` - To access channel messages
   - `chat:write` - To send messages
   - `groups:history` - To access private channel messages
   - `im:history` - To access direct messages
   - `users:read` - To access user information

### 5. Install the app to your workspace

1. Go to "Install App" in the sidebar.
2. Click "Install to Workspace" and authorize the app.

## Environment Variables

The Slack connector requires the following environment variables:

- `SLACK_BOT_TOKEN` - The OAuth token starting with `xoxb-`
- `SLACK_APP_TOKEN` - The app-level token starting with `xapp-`
- `SLACK_SIGNING_SECRET` - The signing secret from the "Basic Information" page

## Usage

```go
// Create config
config := slack.NewConfig(botToken, appToken, signingSecret, true)

// Create factory
factory := slack.NewFactory(config)

// Get chat provider
chatProvider, err := factory.CreateChatProvider()
if err != nil {
    log.Fatalf("Failed to create chat provider: %v", err)
}

// Start listening for messages
messageCh, err := chatProvider.ListenForMessages(context.Background())
if err != nil {
    log.Fatalf("Failed to start listening for messages: %v", err)
}

// Process messages
for msg := range messageCh {
    // Process the message
    fmt.Printf("Received message: %s\n", msg.Content().Text())
}
```