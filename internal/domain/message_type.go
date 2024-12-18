package domain

import (
	"errors"
	"strings"
)

// MessageType represents the type of message in the system
type MessageType string

const (
	// MessageTypeIdea Idea represents a new idea or concept
	MessageTypeIdea MessageType = "idea"
	// MessageTypeDecision Decision represents a decision made
	MessageTypeDecision MessageType = "decision"
	// MessageTypeStatus Status represents a status update
	MessageTypeStatus MessageType = "status"
	// MessageTypeUnknown Unknown represents an unrecognized message type
	MessageTypeUnknown MessageType = "unknown"
)

var (
	ErrInvalidMessageType = errors.New("invalid message type")

	// validMessageTypes contains all valid message types for validation
	validMessageTypes = map[MessageType]bool{
		MessageTypeIdea:     true,
		MessageTypeDecision: true,
		MessageTypeStatus:   true,
		MessageTypeUnknown:  true,
	}
)

// NewMessageType creates a new MessageType instance from a string
func NewMessageType(t string) (MessageType, error) {
	mt := MessageType(strings.ToLower(strings.TrimSpace(t)))
	if !mt.IsValid() {
		return MessageTypeUnknown, ErrInvalidMessageType
	}
	return mt, nil
}

// MustNewMessageType creates a new MessageType instance from a string
// It panics if the message type is invalid
func MustNewMessageType(t string) MessageType {
	mt, err := NewMessageType(t)
	if err != nil {
		panic(err)
	}
	return mt
}

// String returns the string representation of the MessageType
func (mt MessageType) String() string {
	return string(mt)
}

// IsValid checks if the MessageType is valid
func (mt MessageType) IsValid() bool {
	return validMessageTypes[mt]
}

// IsIdea checks if the MessageType is an idea
func (mt MessageType) IsIdea() bool {
	return mt == MessageTypeIdea
}

// IsDecision checks if the MessageType is a decision
func (mt MessageType) IsDecision() bool {
	return mt == MessageTypeDecision
}

// IsStatus checks if the MessageType is a status update
func (mt MessageType) IsStatus() bool {
	return mt == MessageTypeStatus
}

// IsUnknown checks if the MessageType is unknown
func (mt MessageType) IsUnknown() bool {
	return mt == MessageTypeUnknown
}
