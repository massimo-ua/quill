package domain

import (
	"errors"
)

var (
	ErrInvalidConfidenceScore = errors.New("confidence score must be between 0 and 1")
	ErrInvalidAnalysisResult  = errors.New("invalid analysis result")
)

// Analysis represents the raw AI analysis result
type Analysis struct {
	Type            MessageType
	Category        Category
	ConfidenceScore float64
	SuggestedTags   []string
}

// MessageAnalysisResult represents the complete analysis of a message
type MessageAnalysisResult struct {
	messageType     MessageType
	category        Category
	references      []*Reference
	confidenceScore float64
	suggestedTags   []string
}

// NewMessageAnalysisResult creates a new MessageAnalysisResult instance
func NewMessageAnalysisResult(
	msgType MessageType,
	category Category,
	references []*Reference,
	confidence float64,
	tags []string,
) (*MessageAnalysisResult, error) {
	if confidence < 0 || confidence > 1 {
		return nil, ErrInvalidConfidenceScore
	}

	if !msgType.IsValid() || !category.IsValid() {
		return nil, ErrInvalidAnalysisResult
	}

	return &MessageAnalysisResult{
		messageType:     msgType,
		category:        category,
		references:      references,
		confidenceScore: confidence,
		suggestedTags:   tags,
	}, nil
}

// MessageType returns the detected message type
func (r *MessageAnalysisResult) MessageType() MessageType {
	return r.messageType
}

// Category returns the detected category
func (r *MessageAnalysisResult) Category() Category {
	return r.category
}

// References returns the detected references
func (r *MessageAnalysisResult) References() []*Reference {
	refs := make([]*Reference, len(r.references))
	copy(refs, r.references)
	return refs
}

// ConfidenceScore returns the confidence score of the analysis
func (r *MessageAnalysisResult) ConfidenceScore() float64 {
	return r.confidenceScore
}

// SuggestedTags returns the suggested tags
func (r *MessageAnalysisResult) SuggestedTags() []string {
	tags := make([]string, len(r.suggestedTags))
	copy(tags, r.suggestedTags)
	return tags
}

// IsHighConfidence checks if the analysis has high confidence (>= 0.8)
func (r *MessageAnalysisResult) IsHighConfidence() bool {
	return r.confidenceScore >= 0.8
}

// HasReferences checks if any references were detected
func (r *MessageAnalysisResult) HasReferences() bool {
	return len(r.references) > 0
}

// HasSuggestedTags checks if any tags were suggested
func (r *MessageAnalysisResult) HasSuggestedTags() bool {
	return len(r.suggestedTags) > 0
}
