package common

import (
	"database/sql/driver"
	"errors"
	"github.com/oklog/ulid/v2"
	"math/rand"
	"strings"
	"time"
)

// ID represents a unique identifier for domain entities
type ID struct {
	value string
}

var (
	// ErrInvalidID indicates that the ID is invalid
	ErrInvalidID = errors.New("invalid ID format")

	// entropy source for ULID generation
	entropy = ulid.Monotonic(rand.New(rand.NewSource(time.Now().UnixNano())), 0)
)

// NewID creates a new ID instance
func NewID(value string) (ID, error) {
	value = strings.ToUpper(strings.TrimSpace(value))
	if value == "" {
		return ID{}, ErrInvalidID
	}

	// Validate ULID format
	if _, err := ulid.Parse(value); err != nil {
		return ID{}, ErrInvalidID
	}

	return ID{value: value}, nil
}

// GenerateID creates a new unique ID
func GenerateID() ID {
	id := ulid.MustNew(ulid.Timestamp(time.Now()), entropy)
	return ID{value: id.String()}
}

// String returns the string representation of the ID
func (id ID) String() string {
	return id.value
}

// Equals checks if two IDs are equal
func (id ID) Equals(other ID) bool {
	return id.value == other.value
}

// Compare returns an integer comparing two IDs lexicographically.
// The result will be 0 if id==other, -1 if id < other, and +1 if id > other.
func (id ID) Compare(other ID) int {
	return strings.Compare(id.value, other.value)
}

// Time returns the timestamp encoded in the ID
func (id ID) Time() time.Time {
	if id.value == "" {
		return time.Time{} // Return zero time for empty ID
	}

	uid, err := ulid.Parse(id.value)
	if err != nil {
		return time.Time{} // Return zero time for invalid ULID
	}

	return time.Unix(int64(uid.Time())/1000, (int64(uid.Time())%1000)*1000000).UTC()
}

// Value implements the driver.Valuer interface for database operations
func (id ID) Value() (driver.Value, error) {
	return id.value, nil
}

// Scan implements the sql.Scanner interface for database operations
func (id *ID) Scan(value interface{}) error {
	if value == nil {
		return ErrInvalidID
	}

	str, ok := value.(string)
	if !ok {
		bytes, ok := value.([]byte)
		if !ok {
			return ErrInvalidID
		}
		str = string(bytes)
	}

	parsedID, err := NewID(str)
	if err != nil {
		return err
	}

	*id = parsedID
	return nil
}
