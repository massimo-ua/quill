package ollama

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/massimo-ua/quill/internal/domain"
	"strings"
)

// Provider implements the domain.AiAgentProvider interface using Ollama API
type Provider struct {
	client *Client
}

// NewProvider creates a new Ollama provider
func NewProvider(client *Client) *Provider {
	if client == nil {
		panic("client cannot be nil")
	}
	return &Provider{
		client: client,
	}
}

// AnalyzeMessage analyzes message content
func (p *Provider) AnalyzeMessage(ctx context.Context, content string) (*domain.MessageAnalysisResult, error) {
	if ctx == nil {
		return nil, fmt.Errorf("context cannot be nil")
	}
	
	if strings.TrimSpace(content) == "" {
		return nil, fmt.Errorf("content cannot be empty")
	}

	messages := []Message{
		{
			Role:    "system",
			Content: analyzeMessageSystemPrompt,
		},
		{
			Role:    "user",
			Content: content,
		},
	}

	response, err := p.client.GenerateChatCompletion(ctx, messages)
	if err != nil {
		return nil, fmt.Errorf("failed to generate chat completion: %w", err)
	}

	var analysis domain.Analysis
	if err := json.Unmarshal([]byte(response), &analysis); err != nil {
		// If JSON unmarshaling fails, try to extract structured data from the response
		analysis, err = parseAnalysisFromText(response)
		if err != nil {
			return nil, fmt.Errorf("failed to parse analysis: %w", err)
		}
	}

	// Validate message type and category
	msgType, err := domain.NewMessageType(string(analysis.Type))
	if err != nil {
		msgType = domain.MessageTypeUnknown
	}

	category, err := domain.NewCategory(string(analysis.Category))
	if err != nil {
		category = domain.CategoryUnknown
	}

	// Create analysis result
	result, err := domain.NewMessageAnalysisResult(
		msgType,
		category,
		nil, // References will be detected separately
		analysis.ConfidenceScore,
		analysis.SuggestedTags,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create message analysis result: %w", err)
	}

	return result, nil
}

// GenerateDocumentation generates documentation from a message
func (p *Provider) GenerateDocumentation(ctx context.Context, message string, metadata map[string]interface{}) (string, error) {
	if ctx == nil {
		return "", fmt.Errorf("context cannot be nil")
	}
	
	if strings.TrimSpace(message) == "" {
		return "", fmt.Errorf("message cannot be empty")
	}

	// Create a prompt that includes metadata
	prompt := generateDocumentationPrompt(message, metadata)

	messages := []Message{
		{
			Role:    "system",
			Content: generateDocumentationSystemPrompt,
		},
		{
			Role:    "user",
			Content: prompt,
		},
	}

	response, err := p.client.GenerateChatCompletion(ctx, messages)
	if err != nil {
		return "", fmt.Errorf("failed to generate chat completion: %w", err)
	}

	return response, nil
}

// CategorizeContent categorizes content
func (p *Provider) CategorizeContent(ctx context.Context, content string) (*domain.Category, error) {
	if ctx == nil {
		return nil, fmt.Errorf("context cannot be nil")
	}
	
	if strings.TrimSpace(content) == "" {
		return nil, fmt.Errorf("content cannot be empty")
	}

	messages := []Message{
		{
			Role:    "system",
			Content: categorizeContentSystemPrompt,
		},
		{
			Role:    "user",
			Content: content,
		},
	}

	response, err := p.client.GenerateChatCompletion(ctx, messages)
	if err != nil {
		return nil, fmt.Errorf("failed to generate chat completion: %w", err)
	}

	// Clean up the response
	categoryStr := strings.TrimSpace(response)
	categoryStr = strings.ToLower(categoryStr)
	
	// Remove any quotation marks
	categoryStr = strings.Trim(categoryStr, "\"'")

	// Create a category
	category, err := domain.NewCategory(categoryStr)
	if err != nil {
		return nil, fmt.Errorf("failed to create category: %w", err)
	}

	return &category, nil
}

// DetectReferences finds references in content
func (p *Provider) DetectReferences(ctx context.Context, content string) ([]*domain.Reference, error) {
	if ctx == nil {
		return nil, fmt.Errorf("context cannot be nil")
	}
	
	if strings.TrimSpace(content) == "" {
		return nil, fmt.Errorf("content cannot be empty")
	}

	messages := []Message{
		{
			Role:    "system",
			Content: detectReferencesSystemPrompt,
		},
		{
			Role:    "user",
			Content: content,
		},
	}

	response, err := p.client.GenerateChatCompletion(ctx, messages)
	if err != nil {
		return nil, fmt.Errorf("failed to generate chat completion: %w", err)
	}

	var referencesData []struct {
		Type  string `json:"type"`
		Value string `json:"value"`
	}

	// Try to parse as JSON
	if err := json.Unmarshal([]byte(response), &referencesData); err != nil {
		// If it fails, try to extract references from text
		return parseReferencesFromText(response)
	}

	// Convert to domain references
	var references []*domain.Reference
	for _, refData := range referencesData {
		refType := domain.ReferenceType(strings.ToLower(refData.Type))
		if !refType.IsValid() {
			continue
		}

		ref, err := domain.NewReference(refType, refData.Value)
		if err != nil {
			continue
		}
		references = append(references, ref)
	}

	return references, nil
}

// Generate documentation prompt with metadata
func generateDocumentationPrompt(message string, metadata map[string]interface{}) string {
	prompt := fmt.Sprintf("Generate comprehensive documentation from the following message:\n\n%s\n\n", message)
	
	if metadata != nil {
		prompt += "Additional context:\n"
		if msgType, ok := metadata["type"].(string); ok {
			prompt += fmt.Sprintf("- Type: %s\n", msgType)
		}
		if category, ok := metadata["category"].(string); ok {
			prompt += fmt.Sprintf("- Category: %s\n", category)
		}
		if timestamp, ok := metadata["created_at"].(string); ok {
			prompt += fmt.Sprintf("- Created: %s\n", timestamp)
		}
		if refs, ok := metadata["references"].([]*domain.Reference); ok && len(refs) > 0 {
			prompt += "- References:\n"
			for _, ref := range refs {
				prompt += fmt.Sprintf("  - %s: %s\n", ref.Type(), ref.Value())
			}
		}
	}
	
	prompt += "\nFormat the documentation in Markdown with proper sections, headings, and formatting."
	
	return prompt
}

// Helper function to parse analysis results from unstructured text
func parseAnalysisFromText(text string) (domain.Analysis, error) {
	var analysis domain.Analysis
	
	// Set defaults
	analysis.Type = domain.MessageTypeUnknown
	analysis.Category = domain.CategoryUnknown
	analysis.ConfidenceScore = 0.5
	
	lines := strings.Split(text, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(strings.ToLower(line), "type:") {
			typeStr := strings.TrimSpace(strings.TrimPrefix(strings.ToLower(line), "type:"))
			analysis.Type = domain.MessageType(typeStr)
		} else if strings.HasPrefix(strings.ToLower(line), "category:") {
			catStr := strings.TrimSpace(strings.TrimPrefix(strings.ToLower(line), "category:"))
			analysis.Category = domain.Category(catStr)
		} else if strings.HasPrefix(strings.ToLower(line), "confidence:") {
			scoreStr := strings.TrimSpace(strings.TrimPrefix(strings.ToLower(line), "confidence:"))
			fmt.Sscanf(scoreStr, "%f", &analysis.ConfidenceScore)
		} else if strings.HasPrefix(strings.ToLower(line), "tags:") {
			tagsStr := strings.TrimSpace(strings.TrimPrefix(strings.ToLower(line), "tags:"))
			tags := strings.Split(tagsStr, ",")
			for _, tag := range tags {
				tag = strings.TrimSpace(tag)
				if tag != "" {
					analysis.SuggestedTags = append(analysis.SuggestedTags, tag)
				}
			}
		}
	}
	
	return analysis, nil
}

// Helper function to parse references from unstructured text
func parseReferencesFromText(text string) ([]*domain.Reference, error) {
	var references []*domain.Reference
	
	lines := strings.Split(text, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		
		// Look for patterns like "message: MSG123" or "document: path/to/doc"
		if strings.Contains(line, ":") {
			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 {
				typeStr := strings.TrimSpace(strings.ToLower(parts[0]))
				valueStr := strings.TrimSpace(parts[1])
				
				var refType domain.ReferenceType
				if typeStr == "message" {
					refType = domain.ReferenceTypeMessage
				} else if typeStr == "document" {
					refType = domain.ReferenceTypeDocument
				} else {
					continue
				}
				
				ref, err := domain.NewReference(refType, valueStr)
				if err != nil {
					continue
				}
				references = append(references, ref)
			}
		}
	}
	
	return references, nil
}