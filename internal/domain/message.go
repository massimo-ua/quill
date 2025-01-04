package domain

import (
	"errors"
	"time"

	"github.com/massimo-ua/quill/internal/domain/common"
)

var (
	ErrInvalidSender = errors.New("invalid sender")
	ErrNoContent     = errors.New("message must have content")
	ErrNoType        = errors.New("message must have type")
)

// Message represents a chat message in the system
type Message struct {
	id          common.ID
	threadID    common.ID
	sender      string
	content     *MessageContent
	messageType MessageType
	category    Category
	references  []*Reference
	timestamp   time.Time
}

// NewMessage creates a new Message instance
func NewMessage(
	threadID common.ID,
	sender string,
	content *MessageContent,
	messageType MessageType,
	category Category,
	references []*Reference,
) (*Message, error) {
	if sender == "" {
		return nil, ErrInvalidSender
	}
	if content == nil {
		return nil, ErrNoContent
	}
	if !messageType.IsValid() {
		return nil, ErrNoType
	}

	return &Message{
		id:          common.GenerateID(),
		threadID:    threadID,
		sender:      sender,
		content:     content,
		messageType: messageType,
		category:    category,
		references:  references,
		timestamp:   time.Now(),
	}, nil
}

// ID returns the message's identifier
func (m *Message) ID() common.ID {
	return m.id
}

// ThreadID returns the message's thread identifier
func (m *Message) ThreadID() common.ID {
	return m.threadID
}

// Sender returns the message sender
func (m *Message) Sender() string {
	return m.sender
}

// Content returns the message content
func (m *Message) Content() *MessageContent {
	return m.content
}

// Type returns the message type
func (m *Message) Type() MessageType {
	return m.messageType
}

// Category returns the message category
func (m *Message) Category() Category {
	return m.category
}

// References returns the message references
func (m *Message) References() []*Reference {
	refs := make([]*Reference, len(m.references))
	copy(refs, m.references)
	return refs
}

// Timestamp returns the message timestamp
func (m *Message) Timestamp() time.Time {
	return m.timestamp
}

// AddReference adds a new reference to the message
func (m *Message) AddReference(ref *Reference) {
	if ref != nil {
		m.references = append(m.references, ref)
	}
}

// HasReferences checks if the message has any references
func (m *Message) HasReferences() bool {
	return len(m.references) > 0
}

// UpdateCategory updates the message category
func (m *Message) UpdateCategory(category Category) {
	if category.IsValid() {
		m.category = category
	}
}
