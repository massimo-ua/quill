package domain

import (
	"errors"
	"github.com/massimo-ua/quill/internal/domain/common"
	"time"
)

var (
	ErrEmptyThreadTitle = errors.New("thread title cannot be empty")
	ErrInvalidMessages  = errors.New("invalid messages list")
)

// Thread represents a conversation thread
type Thread struct {
	id        common.ID
	title     string
	messages  []*Message
	createdAt time.Time
	updatedAt time.Time
}

// NewThread creates a new Thread instance
func NewThread(title string) (*Thread, error) {
	if title == "" {
		return nil, ErrEmptyThreadTitle
	}

	now := time.Now()
	return &Thread{
		id:        common.GenerateID(),
		title:     title,
		messages:  make([]*Message, 0),
		createdAt: now,
		updatedAt: now,
	}, nil
}

// ID returns thread identifier
func (t *Thread) ID() common.ID {
	return t.id
}

// Title returns thread title
func (t *Thread) Title() string {
	return t.title
}

// Messages returns thread messages
func (t *Thread) Messages() []*Message {
	msgs := make([]*Message, len(t.messages))
	copy(msgs, t.messages)
	return msgs
}

// AddMessage adds a message to the thread
func (t *Thread) AddMessage(msg *Message) error {
	if msg == nil {
		return ErrInvalidMessages
	}

	t.messages = append(t.messages, msg)
	t.updatedAt = time.Now()
	return nil
}

// LastMessage returns the most recent message
func (t *Thread) LastMessage() *Message {
	if len(t.messages) == 0 {
		return nil
	}
	return t.messages[len(t.messages)-1]
}

// CreatedAt returns thread creation time
func (t *Thread) CreatedAt() time.Time {
	return t.createdAt
}

// UpdatedAt returns last update time
func (t *Thread) UpdatedAt() time.Time {
	return t.updatedAt
}

// MessageCount returns total messages in thread
func (t *Thread) MessageCount() int {
	return len(t.messages)
}
