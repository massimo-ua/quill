package domain

import (
	"errors"
	"github.com/massimo-ua/quill/internal/domain/common"
	"strings"
	"time"
)

var (
	ErrInvalidUsername = errors.New("invalid username")
	ErrInvalidEmail    = errors.New("invalid email")
)

// User represents a user in the system
type User struct {
	id        common.ID
	username  string
	email     string
	roles     []string
	createdAt time.Time
	updatedAt time.Time
}

// NewUser creates a new User instance
func NewUser(username, email string) (*User, error) {
	if strings.TrimSpace(username) == "" {
		return nil, ErrInvalidUsername
	}
	if strings.TrimSpace(email) == "" {
		return nil, ErrInvalidEmail
	}

	now := time.Now()
	return &User{
		id:        common.GenerateID(),
		username:  username,
		email:     email,
		roles:     make([]string, 0),
		createdAt: now,
		updatedAt: now,
	}, nil
}

// ID returns user identifier
func (u *User) ID() common.ID {
	return u.id
}

// Username returns username
func (u *User) Username() string {
	return u.username
}

// Email returns user email
func (u *User) Email() string {
	return u.email
}

// Roles returns user roles
func (u *User) Roles() []string {
	roles := make([]string, len(u.roles))
	copy(roles, u.roles)
	return roles
}

// AddRole adds a role to the user
func (u *User) AddRole(role string) {
	if role = strings.TrimSpace(role); role != "" {
		u.roles = append(u.roles, role)
		u.updatedAt = time.Now()
	}
}

// HasRole checks if user has specific role
func (u *User) HasRole(role string) bool {
	for _, r := range u.roles {
		if r == role {
			return true
		}
	}
	return false
}

// CreatedAt returns creation timestamp
func (u *User) CreatedAt() time.Time {
	return u.createdAt
}

// UpdatedAt returns last update timestamp
func (u *User) UpdatedAt() time.Time {
	return u.updatedAt
}
