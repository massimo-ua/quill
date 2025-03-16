package domain

import (
	"errors"
	"strings"
)

// Category represents the category of message in the system
type Category string

const (
	// CategoryOperations represents operations-related content
	CategoryOperations Category = "operations"
	// CategoryDevelopment represents development-related content
	CategoryDevelopment Category = "development"
	// CategoryProduct represents product-related content
	CategoryProduct Category = "product"
	// CategoryQualityAssurance represents QA-related content
	CategoryQualityAssurance Category = "quality_assurance"
	// CategoryDataAnalysis represents data analysis-related content
	CategoryDataAnalysis Category = "data_analysis"
	// CategoryOther represents other content
	CategoryOther Category = "other"
	// CategoryUnknown represents an unrecognized category
	CategoryUnknown Category = "unknown"
)

var (
	ErrInvalidCategory = errors.New("invalid category")

	// validCategories contains all valid categories for validation
	validCategories = map[Category]bool{
		CategoryOperations:       true,
		CategoryDevelopment:      true,
		CategoryProduct:          true,
		CategoryQualityAssurance: true,
		CategoryDataAnalysis:     true,
		CategoryOther:            true,
		CategoryUnknown:          true,
	}
)

// NewCategory creates a new Category instance from a string
func NewCategory(c string) (Category, error) {
	category := Category(strings.ToLower(strings.TrimSpace(c)))
	if !category.IsValid() {
		return CategoryUnknown, ErrInvalidCategory
	}
	return category, nil
}

// MustNewCategory creates a new Category instance from a string
// It panics if the category is invalid
func MustNewCategory(c string) Category {
	category, err := NewCategory(c)
	if err != nil {
		panic(err)
	}
	return category
}

// String returns the string representation of the Category
func (c Category) String() string {
	return string(c)
}

// IsValid checks if the Category is valid
func (c Category) IsValid() bool {
	return validCategories[c]
}

// IsOperations checks if the Category is operations
func (c Category) IsOperations() bool {
	return c == CategoryOperations
}

// IsDevelopment checks if the Category is development
func (c Category) IsDevelopment() bool {
	return c == CategoryDevelopment
}

// IsProduct checks if the Category is product
func (c Category) IsProduct() bool {
	return c == CategoryProduct
}

// IsQualityAssurance checks if the Category is quality assurance
func (c Category) IsQualityAssurance() bool {
	return c == CategoryQualityAssurance
}

// IsDataAnalysis checks if the Category is data analysis
func (c Category) IsDataAnalysis() bool {
	return c == CategoryDataAnalysis
}

// IsOther checks if the Category is other
func (c Category) IsOther() bool {
	return c == CategoryOther
}

// IsUnknown checks if the Category is unknown
func (c Category) IsUnknown() bool {
	return c == CategoryUnknown
}
