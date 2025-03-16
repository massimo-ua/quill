package github

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		config  *Config
		wantErr bool
	}{
		{
			name: "valid config",
			config: &Config{
				Token:          "token123",
				Owner:          "owner",
				Repo:           "repo",
				Branch:         "main",
				BasePath:       "docs",
				CommitterName:  "Test User",
				CommitterEmail: "test@example.com",
			},
			wantErr: false,
		},
		{
			name: "missing token",
			config: &Config{
				Owner:          "owner",
				Repo:           "repo",
				Branch:         "main",
				BasePath:       "docs",
				CommitterName:  "Test User",
				CommitterEmail: "test@example.com",
			},
			wantErr: true,
		},
		{
			name: "missing owner",
			config: &Config{
				Token:          "token123",
				Repo:           "repo",
				Branch:         "main",
				BasePath:       "docs",
				CommitterName:  "Test User",
				CommitterEmail: "test@example.com",
			},
			wantErr: true,
		},
		{
			name: "missing repo",
			config: &Config{
				Token:          "token123",
				Owner:          "owner",
				Branch:         "main",
				BasePath:       "docs",
				CommitterName:  "Test User",
				CommitterEmail: "test@example.com",
			},
			wantErr: true,
		},
		{
			name: "missing committer name",
			config: &Config{
				Token:          "token123",
				Owner:          "owner",
				Repo:           "repo",
				Branch:         "main",
				BasePath:       "docs",
				CommitterEmail: "test@example.com",
			},
			wantErr: true,
		},
		{
			name: "missing committer email",
			config: &Config{
				Token:         "token123",
				Owner:         "owner",
				Repo:          "repo",
				Branch:        "main",
				BasePath:      "docs",
				CommitterName: "Test User",
			},
			wantErr: true,
		},
		{
			name: "empty branch uses default",
			config: &Config{
				Token:          "token123",
				Owner:          "owner",
				Repo:           "repo",
				CommitterName:  "Test User",
				CommitterEmail: "test@example.com",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				if tt.config.Branch == "" {
					assert.Equal(t, "main", tt.config.Branch)
				}
			}
		})
	}
}