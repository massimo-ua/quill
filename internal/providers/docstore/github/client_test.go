package github

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewClient(t *testing.T) {
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
			name:    "nil config",
			config:  nil,
			wantErr: true,
		},
		{
			name: "invalid config",
			config: &Config{
				// Missing required fields
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := NewClient(tt.config)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, client)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, client)
				assert.Equal(t, tt.config, client.config)
				assert.Equal(t, GitHubAPIBaseURL, client.apiBaseURL)
				assert.NotNil(t, client.httpClient)
			}
		})
	}
}

func TestClient_buildContentPath(t *testing.T) {
	tests := []struct {
		name     string
		config   *Config
		path     string
		expected string
	}{
		{
			name: "with base path",
			config: &Config{
				Owner:    "owner",
				Repo:     "repo",
				BasePath: "docs",
			},
			path:     "folder/file.md",
			expected: "https://api.github.com/repos/owner/repo/contents/docs%2Ffolder%2Ffile.md",
		},
		{
			name: "without base path",
			config: &Config{
				Owner: "owner",
				Repo:  "repo",
			},
			path:     "folder/file.md",
			expected: "https://api.github.com/repos/owner/repo/contents/folder%2Ffile.md",
		},
		{
			name: "with leading slash in path",
			config: &Config{
				Owner: "owner",
				Repo:  "repo",
			},
			path:     "/folder/file.md",
			expected: "https://api.github.com/repos/owner/repo/contents/folder%2Ffile.md",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &Client{
				config:     tt.config,
				apiBaseURL: GitHubAPIBaseURL,
			}
			path := client.buildContentPath(tt.path)
			assert.Equal(t, tt.expected, path)
		})
	}
}