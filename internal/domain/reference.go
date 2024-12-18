package domain

import (
	"errors"
	"fmt"
	"strings"
)

// ReferenceType represents the type of reference in the system
type ReferenceType string

const (
	// ReferenceTypeMessage represents a reference to a message
	ReferenceTypeMessage ReferenceType = "message"
	// ReferenceTypeDocument represents a reference to a document
	ReferenceTypeDocument ReferenceType = "document"
	// ReferenceTypeUnknown represents an unrecognized reference type
	ReferenceTypeUnknown ReferenceType = "unknown"
)

var (
	ErrInvalidReference    = errors.New("invalid reference")
	ErrEmptyReferenceValue = errors.New("empty reference value")

	// validReferenceTypes contains all valid reference types for validation
	validReferenceTypes = map[ReferenceType]bool{
		ReferenceTypeMessage:  true,
		ReferenceTypeDocument: true,
		ReferenceTypeUnknown:  true,
	}
)

// Reference is a value object that represents a reference to another entity in the system
type Reference struct {
	refType ReferenceType
	value   string // either message ID or document path
}

// NewReference creates a new Reference instance
func NewReference(refType ReferenceType, value string) (*Reference, error) {
	if !refType.IsValid() {
		return nil, fmt.Errorf("%w: invalid reference type", ErrInvalidReference)
	}

	if strings.TrimSpace(value) == "" {
		return nil, ErrEmptyReferenceValue
	}

	return &Reference{
		refType: refType,
		value:   strings.TrimSpace(value),
	}, nil
}

// MustNewReference creates a new Reference instance
// It panics if the reference is invalid
func MustNewReference(refType ReferenceType, value string) *Reference {
	ref, err := NewReference(refType, value)
	if err != nil {
		panic(err)
	}
	return ref
}

// NewMessageReference creates a new Reference to a message
func NewMessageReference(messageID string) (*Reference, error) {
	return NewReference(ReferenceTypeMessage, messageID)
}

// NewDocumentReference creates a new Reference to a document
func NewDocumentReference(documentPath string) (*Reference, error) {
	return NewReference(ReferenceTypeDocument, documentPath)
}

// Type returns the type of the reference
func (r Reference) Type() ReferenceType {
	return r.refType
}

// Value returns the value of the reference (message ID or document path)
func (r Reference) Value() string {
	return r.value
}

// String returns a string representation of the reference
func (r Reference) String() string {
	return fmt.Sprintf("%s:%s", r.refType, r.value)
}

// Equals checks if two references are equal
func (r Reference) Equals(other Reference) bool {
	return r.refType == other.refType && r.value == other.value
}

// ParseReference Parse creates a Reference from a string representation
func ParseReference(ref string) (*Reference, error) {
	parts := strings.SplitN(ref, ":", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("%w: invalid format", ErrInvalidReference)
	}

	refType := ReferenceType(strings.ToLower(strings.TrimSpace(parts[0])))
	if !refType.IsValid() {
		return nil, fmt.Errorf("%w: invalid reference type", ErrInvalidReference)
	}

	return NewReference(refType, parts[1])
}

// IsValid checks if the ReferenceType is valid
func (rt ReferenceType) IsValid() bool {
	return validReferenceTypes[rt]
}

// String returns the string representation of the ReferenceType
func (rt ReferenceType) String() string {
	return string(rt)
}

// IsMessage checks if the ReferenceType is a message reference
func (rt ReferenceType) IsMessage() bool {
	return rt == ReferenceTypeMessage
}

// IsDocument checks if the ReferenceType is a document reference
func (rt ReferenceType) IsDocument() bool {
	return rt == ReferenceTypeDocument
}

// IsUnknown checks if the ReferenceType is unknown
func (rt ReferenceType) IsUnknown() bool {
	return rt == ReferenceTypeUnknown
}
