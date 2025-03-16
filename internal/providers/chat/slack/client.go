package slack

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"

	"github.com/massimo-ua/quill/internal/domain"
	"github.com/massimo-ua/quill/internal/domain/common"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/socketmode"
)

var (
	ErrInvalidConfig    = errors.New("invalid slack configuration")
	ErrConnectionFailed = errors.New("failed to connect to slack")
)

// Client implements the ChatAccessProvider interface for Slack
type Client struct {
	config     *Config
	api        *slack.Client
	socket     *socketmode.Client
	messageCh  chan *domain.Message
	threadMap  map[string]common.ID // Maps Slack thread TS to our ThreadID
	threadLock sync.RWMutex
}

// NewClient creates a new Slack client
func NewClient(config *Config) (*Client, error) {
	if config == nil || config.BotToken == "" || config.AppToken == "" {
		return nil, ErrInvalidConfig
	}

	api := slack.New(
		config.BotToken,
		slack.OptionAppLevelToken(config.AppToken),
		slack.OptionDebug(config.DebugMode),
	)

	socketClient := socketmode.New(
		api,
		socketmode.OptionDebug(config.DebugMode),
	)

	return &Client{
		config:    config,
		api:       api,
		socket:    socketClient,
		messageCh: make(chan *domain.Message, 100),
		threadMap: make(map[string]common.ID),
	}, nil
}

// SendMessage sends a message to a Slack channel
func (c *Client) SendMessage(ctx context.Context, channelID, content string) error {
	_, _, err := c.api.PostMessageContext(ctx, channelID, slack.MsgOptionText(content, false))
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}
	return nil
}

// ReplyToMessage replies to a specific Slack message
func (c *Client) ReplyToMessage(ctx context.Context, messageID, content string) error {
	_, _, err := c.api.PostMessageContext(
		ctx,
		messageID, // This should be the channel ID
		slack.MsgOptionText(content, false),
		slack.MsgOptionTS(messageID), // Thread TS
	)
	if err != nil {
		return fmt.Errorf("failed to reply to message: %w", err)
	}
	return nil
}

// ListenForMessages starts listening for Slack messages
func (c *Client) ListenForMessages(ctx context.Context) (<-chan *domain.Message, error) {
	// Start the socket mode client in a separate goroutine
	go func() {
		err := c.socket.Run()
		if err != nil {
			log.Printf("Error running socket mode client: %v", err)
		}
	}()

	// Handle events in a separate goroutine
	go c.handleEvents(ctx)

	return c.messageCh, nil
}

// HandleInteraction processes interactive components from Slack
func (c *Client) HandleInteraction(ctx context.Context, interaction interface{}) error {
	// Process different types of interactions based on the interface type
	switch i := interaction.(type) {
	case *slack.InteractionCallback:
		return c.processInteractionCallback(ctx, i)
	default:
		return fmt.Errorf("unsupported interaction type: %T", interaction)
	}
}

// handleEvents processes incoming Slack events
func (c *Client) handleEvents(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case event := <-c.socket.Events:
			// Handle different event types
			switch event.Type {
			case socketmode.EventTypeEventsAPI:
				eventsAPIEvent, ok := event.Data.(EventsAPIEvent)
				if !ok {
					log.Printf("Error: could not type cast to EventsAPIEvent")
					continue
				}

				// Acknowledge the event
				c.socket.Ack(*event.Request)

				// Process the inner event
				switch innerEvent := eventsAPIEvent.Data.(type) {
				case *slack.MessageEvent:
					c.processMessageEvent(eventsAPIEvent.EventID, innerEvent)
				}

			case socketmode.EventTypeInteractive:
				interaction, ok := event.Data.(slack.InteractionCallback)
				if !ok {
					log.Printf("Error: could not type cast to InteractionCallback")
					continue
				}

				// Acknowledge the interaction
				c.socket.Ack(*event.Request)

				// Process the interaction in a separate goroutine
				go func() {
					if err := c.processInteractionCallback(ctx, &interaction); err != nil {
						log.Printf("Error processing interaction: %v", err)
					}
				}()
			}
		}
	}
}

// processMessageEvent converts a Slack message to our domain Message
func (c *Client) processMessageEvent(eventID string, msg *slack.MessageEvent) {
	// Skip bot messages
	if msg.BotID != "" || msg.User == "" {
		return
	}

	// Get user info
	userInfo, err := c.api.GetUserInfo(msg.User)
	if err != nil {
		log.Printf("Error fetching user info: %v", err)
		return
	}

	// Handle thread mapping
	var threadID common.ID
	c.threadLock.RLock()
	if msg.ThreadTimestamp != "" {
		// This is part of a thread
		if id, exists := c.threadMap[msg.ThreadTimestamp]; exists {
			threadID = id
		} else {
			// The parent message's thread ID wasn't mapped yet
			threadID = common.GenerateID()
			c.threadLock.RUnlock()
			c.threadLock.Lock()
			c.threadMap[msg.ThreadTimestamp] = threadID
			c.threadLock.Unlock()
			return
		}
	} else {
		// This is a new message, not in a thread
		threadID = common.GenerateID()
		c.threadLock.RUnlock()
		c.threadLock.Lock()
		c.threadMap[msg.Timestamp] = threadID
		c.threadLock.Unlock()
	}

	// Create message content
	messageContent, err := domain.NewMessageContent(msg.Text)
	if err != nil {
		log.Printf("Error creating message content: %v", err)
		return
	}

	// For now, use default type and category - these will be determined later by AI analysis
	defaultType := domain.MessageTypeInformation
	defaultCategory := domain.CategoryOther

	// Create domain message
	domainMsg, err := domain.NewMessage(
		threadID,
		userInfo.Name,
		messageContent,
		defaultType,
		defaultCategory,
		nil, // no references initially
	)

	if err != nil {
		log.Printf("Error creating domain message: %v", err)
		return
	}

	// Send to message channel for processing
	c.messageCh <- domainMsg
}

// processInteractionCallback handles Slack interaction callbacks
func (c *Client) processInteractionCallback(ctx context.Context, interaction *slack.InteractionCallback) error {
	// Process different interaction types
	switch interaction.Type {
	case slack.InteractionTypeBlockActions:
		// Handle block actions
		for _, action := range interaction.ActionCallback.BlockActions {
			log.Printf("Received block action: %s with value %s", action.ActionID, action.Value)
			// Here we would process specific actions based on action.ActionID
		}
	case slack.InteractionTypeViewSubmission:
		// Handle modal submissions
		log.Printf("Received view submission: %s", interaction.View.ID)
		// Process form submission
	}

	return nil
}