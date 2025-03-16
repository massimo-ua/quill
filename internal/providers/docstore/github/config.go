package github

import (
	"errors"
	"strings"
)

// Config contains GitHub repository configuration
type Config struct {
	// Token is the GitHub personal access token
	Token string
	// Owner is the repository owner (user or organization)
	Owner string
	// Repo is the repository name
	Repo string
	// Branch is the branch to use (default: main)
	Branch string
	// BasePath is the base path in the repository to store documents
	BasePath string
	// Committer information
	CommitterName  string
	CommitterEmail string
}

var (
	ErrMissingToken         = errors.New("GitHub token is required")
	ErrMissingOwner         = errors.New("repository owner is required")
	ErrMissingRepo          = errors.New("repository name is required")
	ErrMissingCommitterName = errors.New("committer name is required")
	ErrMissingCommitterEmail = errors.New("committer email is required")
)

// Validate checks if the configuration is valid
func (c *Config) Validate() error {
	if strings.TrimSpace(c.Token) == "" {
		return ErrMissingToken
	}
	if strings.TrimSpace(c.Owner) == "" {
		return ErrMissingOwner
	}
	if strings.TrimSpace(c.Repo) == "" {
		return ErrMissingRepo
	}
	if strings.TrimSpace(c.CommitterName) == "" {
		return ErrMissingCommitterName
	}
	if strings.TrimSpace(c.CommitterEmail) == "" {
		return ErrMissingCommitterEmail
	}

	// Set default branch if not specified
	if strings.TrimSpace(c.Branch) == "" {
		c.Branch = "main"
	}

	// Clean up base path
	c.BasePath = strings.Trim(strings.TrimSpace(c.BasePath), "/")

	return nil
}