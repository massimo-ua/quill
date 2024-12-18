package domain

import (
	"errors"
	"fmt"
	"time"
)

// DocumentVersion represents a specific version of a document
type DocumentVersion struct {
	version   uint
	timestamp time.Time
}

var (
	ErrInvalidTimestamp    = errors.New("invalid timestamp")
	ErrVersionComparison   = errors.New("cannot compare versions with nil values")
	defaultDocumentVersion = DocumentVersion{
		version:   1,
		timestamp: time.Now(),
	}
)

// NewDocumentVersion creates a new DocumentVersion instance
func NewDocumentVersion(version uint, timestamp time.Time) (*DocumentVersion, error) {
	if timestamp.IsZero() {
		return nil, ErrInvalidTimestamp
	}

	return &DocumentVersion{
		version:   version,
		timestamp: timestamp,
	}, nil
}

// NewDefaultDocumentVersion creates a new DocumentVersion with default values (version 1)
func NewDefaultDocumentVersion() *DocumentVersion {
	return &DocumentVersion{
		version:   defaultDocumentVersion.version,
		timestamp: time.Now(),
	}
}

// String returns the string representation of the DocumentVersion
func (dv *DocumentVersion) String() string {
	return fmt.Sprintf("v%d", dv.version)
}

// Version returns the version number
func (dv *DocumentVersion) Version() uint {
	return dv.version
}

// Timestamp returns the version's timestamp
func (dv *DocumentVersion) Timestamp() time.Time {
	return dv.timestamp
}

// IsNewer checks if this version is newer than the provided version
func (dv *DocumentVersion) IsNewer(other *DocumentVersion) (bool, error) {
	if dv == nil || other == nil {
		return false, ErrVersionComparison
	}

	if dv.version != other.version {
		return dv.version > other.version, nil
	}

	return dv.timestamp.After(other.timestamp), nil
}

// Equals checks if two versions are equal
func (dv *DocumentVersion) Equals(other *DocumentVersion) bool {
	if dv == nil || other == nil {
		return false
	}

	return dv.version == other.version &&
		dv.timestamp.Equal(other.timestamp)
}

// Increment returns a new DocumentVersion with incremented version number
func (dv *DocumentVersion) Increment() *DocumentVersion {
	return &DocumentVersion{
		version:   dv.version + 1,
		timestamp: time.Now(),
	}
}
