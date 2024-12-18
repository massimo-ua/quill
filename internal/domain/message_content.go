package domain

import (
	"encoding/json"
	"errors"
	"strings"
)

var (
	ErrEmptyContent  = errors.New("content cannot be empty")
	ErrInvalidJSON   = errors.New("invalid JSON format")
	MaxContentLength = 10000 // Maximum allowed content length
)

// MessageContent represents the content of a message
type MessageContent struct {
	text string
}

// NewMessageContent creates a new MessageContent instance
func NewMessageContent(text string) (*MessageContent, error) {
	text = strings.TrimSpace(text)
	if text == "" {
		return nil, ErrEmptyContent
	}

	if len(text) > MaxContentLength {
		text = text[:MaxContentLength]
	}

	return &MessageContent{
		text: text,
	}, nil
}

// MustNewMessageContent creates a new MessageContent instance and panics on error
func MustNewMessageContent(text string) *MessageContent {
	content, err := NewMessageContent(text)
	if err != nil {
		panic(err)
	}
	return content
}

// Text returns the message text content
func (mc *MessageContent) Text() string {
	return mc.text
}

// ContainsURL checks if the message content contains a URL
func (mc *MessageContent) ContainsURL() bool {
	return strings.Contains(mc.text, "http://") || strings.Contains(mc.text, "https://")
}

// WordCount returns the number of words in the message content
func (mc *MessageContent) WordCount() int {
	return len(strings.Fields(mc.text))
}

// ContainsTag checks if the message contains a specific tag (e.g., #idea, #decision)
func (mc *MessageContent) ContainsTag(tag string) bool {
	return strings.Contains(mc.text, "#"+tag)
}

// MarshalJSON implements the json.Marshaler interface
func (mc *MessageContent) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Text string `json:"text"`
	}{
		Text: mc.text,
	})
}

// UnmarshalJSON implements the json.Unmarshaler interface
func (mc *MessageContent) UnmarshalJSON(data []byte) error {
	var temp struct {
		Text string `json:"text"`
	}

	if err := json.Unmarshal(data, &temp); err != nil {
		return ErrInvalidJSON
	}

	content, err := NewMessageContent(temp.Text)
	if err != nil {
		return err
	}

	*mc = *content
	return nil
}
