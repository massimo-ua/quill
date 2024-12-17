package domain

import (
	"fmt"
	"strings"
	"time"
)

// MessageCategory represents different areas of project management
type MessageCategory string

const (
	Operations       MessageCategory = "operations"
	Development      MessageCategory = "development"
	Product          MessageCategory = "product"
	QualityAssurance MessageCategory = "quality-assurance"
	DataAnalysis     MessageCategory = "data-analysis"
)

func (c MessageCategory) IsValid() bool {
	switch c {
	case Operations, Development, Product, QualityAssurance, DataAnalysis:
		return true
	}
	return false
}

func (c MessageCategory) String() string {
	return string(c)
}

func ParseMessageCategory(s string) (MessageCategory, error) {
	category := MessageCategory(strings.ToLower(s))
	if !category.IsValid() {
		return "", fmt.Errorf("invalid message category: %s", s)
	}
	return category, nil
}

// MessageType represents different types of captured information
type MessageType string

const (
	IdeaType     MessageType = "idea"
	DecisionType MessageType = "decision"
	StatusType   MessageType = "status"
)

func (t MessageType) IsValid() bool {
	switch t {
	case IdeaType, DecisionType, StatusType:
		return true
	}
	return false
}

func (t MessageType) String() string {
	return string(t)
}

func ParseMessageType(s string) (MessageType, error) {
	msgType := MessageType(strings.ToLower(s))
	if !msgType.IsValid() {
		return "", fmt.Errorf("invalid message type: %s", s)
	}
	return msgType, nil
}

// MessageState represents the current state of a message
type MessageState string

const (
	StatePending   MessageState = "pending"
	StateProcessed MessageState = "processed"
	StateIgnored   MessageState = "ignored"
)

func (s MessageState) IsValid() bool {
	switch s {
	case StatePending, StateProcessed, StateIgnored:
		return true
	}
	return false
}

func (s MessageState) String() string {
	return string(s)
}

// Project represents a project configuration and metadata
type Project struct {
	ID          string          `json:"id"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Goals       []string        `json:"goals"`
	KPIs        []string        `json:"kpis"`
	Milestones  []string        `json:"milestones"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
	Settings    ProjectSettings `json:"settings"`
	ChannelID   string          `json:"channel_id"`
}

func (p *Project) Validate() error {
	if p.Name == "" {
		return fmt.Errorf("project name is required")
	}
	if p.ChannelID == "" {
		return fmt.Errorf("channel ID is required")
	}
	return p.Settings.Validate()
}

func (p *Project) UpdateTimestamps() {
	if p.CreatedAt.IsZero() {
		p.CreatedAt = time.Now()
	}
	p.UpdatedAt = time.Now()
}

// ProjectSettings contains configurable project settings
type ProjectSettings struct {
	AutoDetection       AutoDetectionConfig `json:"auto_detection"`
	DocumentationConfig DocumentationConfig `json:"documentation"`
}

func (s *ProjectSettings) Validate() error {
	if err := s.AutoDetection.Validate(); err != nil {
		return fmt.Errorf("auto detection config: %w", err)
	}
	if err := s.DocumentationConfig.Validate(); err != nil {
		return fmt.Errorf("documentation config: %w", err)
	}
	return nil
}

// AutoDetectionConfig contains settings for AI detection features
type AutoDetectionConfig struct {
	Enabled             bool              `json:"enabled"`
	ConfidenceThreshold float64           `json:"confidence_threshold"`
	EnabledTypes        []MessageType     `json:"enabled_types"`      // renamed from EnabledCategories
	EnabledCategories   []MessageCategory `json:"enabled_categories"` // renamed from EnabledDomains
}

func (c *AutoDetectionConfig) Validate() error {
	if c.ConfidenceThreshold < 0 || c.ConfidenceThreshold > 1 {
		return fmt.Errorf("confidence threshold must be between 0 and 1")
	}

	for _, msgType := range c.EnabledTypes {
		if !msgType.IsValid() {
			return fmt.Errorf("invalid message type: %s", msgType)
		}
	}

	for _, category := range c.EnabledCategories {
		if !category.IsValid() {
			return fmt.Errorf("invalid message category: %s", category)
		}
	}

	return nil
}

// DocumentationConfig contains settings for documentation generation
type DocumentationConfig struct {
	GitHubRepo   string `json:"github_repo"`
	GitHubBranch string `json:"github_branch"`
	BasePath     string `json:"base_path"`
}

func (c *DocumentationConfig) Validate() error {
	if c.GitHubRepo == "" {
		return fmt.Errorf("GitHub repository is required")
	}
	if c.GitHubBranch == "" {
		c.GitHubBranch = "main" // default branch
	}
	return nil
}

// Document represents a documentation entry
type Document struct {
	ID         string          `json:"id"`
	Type       MessageType     `json:"type"`
	Category   MessageCategory `json:"category"`
	Title      string          `json:"title"`
	Content    string          `json:"content"`
	CreatedAt  time.Time       `json:"created_at"`
	UpdatedAt  time.Time       `json:"updated_at"`
	AuthorID   string          `json:"author_id"`
	ThreadID   string          `json:"thread_id"`
	References []Reference     `json:"references"`
	Path       string          `json:"path"`
}

func (d *Document) Validate() error {
	if d.Title == "" {
		return fmt.Errorf("document title is required")
	}
	if !d.Type.IsValid() {
		return fmt.Errorf("invalid message type: %s", d.Type)
	}
	if !d.Category.IsValid() {
		return fmt.Errorf("invalid message category: %s", d.Category)
	}
	return nil
}

func (d *Document) UpdateTimestamps() {
	if d.CreatedAt.IsZero() {
		d.CreatedAt = time.Now()
	}
	d.UpdatedAt = time.Now()
}

// Helper method to add a reference
func (d *Document) AddReference(ref Reference) {
	d.References = append(d.References, ref)
}

// Reference represents a reference to another document
type Reference struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Path  string `json:"path"`
}

// Message represents a processed Slack message
type Message struct {
	ID        string          `json:"id"`
	Type      MessageType     `json:"type"`
	Category  MessageCategory `json:"category"`
	Content   string          `json:"content"`
	ChannelID string          `json:"channel_id"`
	UserID    string          `json:"user_id"`
	ThreadID  string          `json:"thread_id"`
	Timestamp string          `json:"timestamp"`
	State     MessageState    `json:"state"`
	Analysis  MessageAnalysis `json:"analysis"`
}

func (m *Message) Validate() error {
	if m.Content == "" {
		return fmt.Errorf("message content is required")
	}
	if m.ChannelID == "" {
		return fmt.Errorf("channel ID is required")
	}
	if !m.State.IsValid() {
		return fmt.Errorf("invalid message state: %s", m.State)
	}
	return nil
}

// IsThread returns true if the message is part of a thread
func (m *Message) IsThread() bool {
	return m.ThreadID != ""
}

// MessageAnalysis represents the AI analysis results
type MessageAnalysis struct {
	Type       MessageType     `json:"type"`     // idea, decision, status
	Category   MessageCategory `json:"category"` // operations, development, etc.
	Confidence float64         `json:"confidence"`
	Summary    string          `json:"summary"`
	Reasoning  string          `json:"reasoning"`
}

func (a *MessageAnalysis) Validate() error {
	if !a.Type.IsValid() {
		return fmt.Errorf("invalid message type: %s", a.Type)
	}
	if !a.Category.IsValid() {
		return fmt.Errorf("invalid message category: %s", a.Category)
	}
	if a.Confidence < 0 || a.Confidence > 1 {
		return fmt.Errorf("confidence must be between 0 and 1")
	}
	return nil
}

// NewProject Helper function to create a new project
func NewProject(name, description, channelID string) *Project {
	now := time.Now()
	return &Project{
		Name:        name,
		Description: description,
		ChannelID:   channelID,
		CreatedAt:   now,
		UpdatedAt:   now,
		Settings: ProjectSettings{
			AutoDetection: AutoDetectionConfig{
				Enabled:             true,
				ConfidenceThreshold: 0.8,
				EnabledTypes:        []MessageType{IdeaType, DecisionType, StatusType},
				EnabledCategories:   []MessageCategory{Operations, Development, Product, QualityAssurance, DataAnalysis},
			},
			DocumentationConfig: DocumentationConfig{
				GitHubBranch: "main", // Set default branch
				BasePath:     "docs", // Set default base path
			},
		},
	}
}

// NewDocumentFromMessage Helper function to create a new document from a message
func NewDocumentFromMessage(msg *Message) *Document {
	now := time.Now()
	return &Document{
		Type:      msg.Type,
		Category:  msg.Category,
		Content:   msg.Content,
		CreatedAt: now,
		UpdatedAt: now,
		AuthorID:  msg.UserID,
		ThreadID:  msg.ThreadID,
	}
}
