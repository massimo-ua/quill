package slack

import (
	"github.com/slack-go/slack"
)

// MessageData represents additional data for Slack messages
type MessageData struct {
	// SlackChannelID is the Slack channel ID
	SlackChannelID string

	// SlackThreadTS is the Slack thread timestamp
	SlackThreadTS string

	// SlackMessageTS is the Slack message timestamp
	SlackMessageTS string

	// SlackUserID is the Slack user ID
	SlackUserID string
}

// InteractionResponse represents a response to a Slack interaction
type InteractionResponse struct {
	// ResponseType determines if the response is visible to all users or only to the user who triggered it
	ResponseType string

	// Text is the message text
	Text string

	// Blocks are the message blocks to send
	Blocks []slack.Block

	// ReplaceOriginal determines if the response should replace the original message
	ReplaceOriginal bool
}

// NewInteractionResponse creates a new interaction response
func NewInteractionResponse(responseType string, text string, blocks []slack.Block, replaceOriginal bool) *InteractionResponse {
	return &InteractionResponse{
		ResponseType:    responseType,
		Text:            text,
		Blocks:          blocks,
		ReplaceOriginal: replaceOriginal,
	}
}

// MessageCategory represents a message category button
type MessageCategory struct {
	// Value is the category value
	Value string

	// Text is the button text
	Text string

	// Description is the button description
	Description string
}

// GetCategoryBlocks returns blocks for category selection
func GetCategoryBlocks() []slack.Block {
	// Define categories
	categories := []MessageCategory{
		{Value: "development", Text: "Development", Description: "Code, architecture, technical implementation"},
		{Value: "product", Text: "Product", Description: "Features, requirements, UX decisions"},
		{Value: "operations", Text: "Operations", Description: "Infrastructure, deployment, monitoring"},
		{Value: "qa", Text: "QA", Description: "Testing, quality assurance, validation"},
		{Value: "data", Text: "Data", Description: "Analytics, metrics, insights"},
		{Value: "other", Text: "Other", Description: "Other information"},
	}

	// Create header section
	headerText := slack.NewTextBlockObject(slack.MarkdownType, "Please categorize this message:", false, false)
	headerSection := slack.NewSectionBlock(headerText, nil, nil)

	// Create buttons for each category
	var buttons []slack.BlockElement
	for _, category := range categories {
		buttonText := slack.NewTextBlockObject(slack.PlainTextType, category.Text, false, false)
		button := slack.NewButtonBlockElement(
			"category_"+category.Value,
			category.Value,
			buttonText,
		)
		buttons = append(buttons, button)
	}

	// Create the action block with the buttons
	actionBlock := slack.NewActionBlock("category_selection", buttons...)

	// Return the blocks
	return []slack.Block{headerSection, actionBlock}
}