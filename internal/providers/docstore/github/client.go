package github

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
	"time"
)

const (
	// GitHubAPIBaseURL is the base URL for the GitHub API
	GitHubAPIBaseURL = "https://api.github.com"
	// DefaultTimeout is the default HTTP client timeout
	DefaultTimeout = 30 * time.Second
)

// Client represents a GitHub API client
type Client struct {
	config     *Config
	httpClient *http.Client
	apiBaseURL string
}

// NewClient creates a new GitHub API client
func NewClient(cfg *Config) (*Client, error) {
	if cfg == nil {
		return nil, fmt.Errorf("config cannot be nil")
	}
	
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	httpClient := &http.Client{
		Timeout: DefaultTimeout,
	}

	return &Client{
		config:     cfg,
		httpClient: httpClient,
		apiBaseURL: GitHubAPIBaseURL,
	}, nil
}

// GetContent retrieves the content of a file from GitHub
func (c *Client) GetContent(ctx context.Context, path string) (*GitHubContent, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	fullPath := c.buildContentPath(path)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fullPath, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	c.addAuthHeader(req)
	q := req.URL.Query()
	q.Add("ref", c.config.Branch)
	req.URL.RawQuery = q.Encode()

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("file not found: %s", path)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, body)
	}

	var content GitHubContent
	if err := json.Unmarshal(body, &content); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &content, nil
}

// CreateContent creates a new file in GitHub
func (c *Client) CreateContent(ctx context.Context, path string, content []byte, message string) (*GitHubCommitResponse, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	file := GitHubFile{
		Path:    path,
		Content: base64.StdEncoding.EncodeToString(content),
		Message: message,
		Branch:  c.config.Branch,
		Committer: &GitHubCommitter{
			Name:  c.config.CommitterName,
			Email: c.config.CommitterEmail,
		},
	}

	return c.putContent(ctx, path, file)
}

// UpdateContent updates an existing file in GitHub
func (c *Client) UpdateContent(ctx context.Context, path string, content []byte, message string) (*GitHubCommitResponse, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	// Get the current file to get its SHA
	existingContent, err := c.GetContent(ctx, path)
	if err != nil {
		return nil, fmt.Errorf("failed to get existing content: %w", err)
	}

	file := GitHubFile{
		Path:    path,
		Content: base64.StdEncoding.EncodeToString(content),
		Message: message,
		Branch:  c.config.Branch,
		SHA:     existingContent.SHA,
		Committer: &GitHubCommitter{
			Name:  c.config.CommitterName,
			Email: c.config.CommitterEmail,
		},
	}

	return c.putContent(ctx, path, file)
}

// DeleteContent deletes a file from GitHub
func (c *Client) DeleteContent(ctx context.Context, path string, message string) (*GitHubCommitResponse, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	// Get the current file to get its SHA
	existingContent, err := c.GetContent(ctx, path)
	if err != nil {
		return nil, fmt.Errorf("failed to get existing content: %w", err)
	}

	file := GitHubFile{
		Path:    path,
		Message: message,
		Branch:  c.config.Branch,
		SHA:     existingContent.SHA,
		Committer: &GitHubCommitter{
			Name:  c.config.CommitterName,
			Email: c.config.CommitterEmail,
		},
	}

	fullPath := c.buildContentPath(path)
	jsonData, err := json.Marshal(file)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, fullPath, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	c.addAuthHeader(req)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, body)
	}

	var commitResponse GitHubCommitResponse
	if err := json.Unmarshal(body, &commitResponse); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &commitResponse, nil
}

// ListContents lists files and directories in a path
func (c *Client) ListContents(ctx context.Context, path string) ([]GitHubContentListItem, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	fullPath := c.buildContentPath(path)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fullPath, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	c.addAuthHeader(req)
	q := req.URL.Query()
	q.Add("ref", c.config.Branch)
	req.URL.RawQuery = q.Encode()

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode == http.StatusNotFound {
		return []GitHubContentListItem{}, nil
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, body)
	}

	var items []GitHubContentListItem
	if err := json.Unmarshal(body, &items); err != nil {
		// May be a file, not a directory
		var item GitHubContentListItem
		if err := json.Unmarshal(body, &item); err != nil {
			return nil, fmt.Errorf("failed to parse response: %w", err)
		}
		return []GitHubContentListItem{item}, nil
	}

	return items, nil
}

// CreateDirectory creates an empty directory by creating a .gitkeep file
func (c *Client) CreateDirectory(ctx context.Context, path string) error {
	if ctx == nil {
		ctx = context.Background()
	}

	// GitHub doesn't natively support empty directories, so we create a .gitkeep file
	gitkeepPath := filepath.Join(path, ".gitkeep")
	_, err := c.CreateContent(ctx, gitkeepPath, []byte{}, "Create directory "+path)
	if err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	return nil
}

// putContent puts content to GitHub (either create or update)
func (c *Client) putContent(ctx context.Context, path string, file GitHubFile) (*GitHubCommitResponse, error) {
	fullPath := c.buildContentPath(path)
	jsonData, err := json.Marshal(file)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, fullPath, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	c.addAuthHeader(req)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, body)
	}

	var commitResponse GitHubCommitResponse
	if err := json.Unmarshal(body, &commitResponse); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &commitResponse, nil
}

// ensureDirectoryExists ensures that a directory exists in the repository
func (c *Client) ensureDirectoryExists(ctx context.Context, path string) error {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) == 0 {
		return nil
	}

	current := ""
	for _, part := range parts {
		if part == "" {
			continue
		}

		current = filepath.Join(current, part)
		items, err := c.ListContents(ctx, current)
		if err != nil || len(items) == 0 {
			// Directory doesn't exist, create it
			if err := c.CreateDirectory(ctx, current); err != nil {
				return fmt.Errorf("failed to create directory %s: %w", current, err)
			}
		}
	}

	return nil
}

// buildContentPath builds the full URL path for content operations
func (c *Client) buildContentPath(path string) string {
	contentPath := strings.TrimPrefix(path, "/")
	if c.config.BasePath != "" {
		contentPath = filepath.Join(c.config.BasePath, contentPath)
	}
	return fmt.Sprintf("%s/repos/%s/%s/contents/%s", c.apiBaseURL, c.config.Owner, c.config.Repo, url.PathEscape(contentPath))
}

// addAuthHeader adds the Authorization header to the request
func (c *Client) addAuthHeader(req *http.Request) {
	req.Header.Set("Authorization", "Bearer "+c.config.Token)
	req.Header.Set("Accept", "application/vnd.github.v3+json")
}