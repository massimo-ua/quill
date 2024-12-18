package domain

import (
	"testing"
)

func TestNewCategory(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		want        Category
		wantErr     bool
		errExpected error
	}{
		{
			name:    "valid operations category",
			input:   "operations",
			want:    CategoryOperations,
			wantErr: false,
		},
		{
			name:    "valid development category",
			input:   "development",
			want:    CategoryDevelopment,
			wantErr: false,
		},
		{
			name:    "valid product category",
			input:   "product",
			want:    CategoryProduct,
			wantErr: false,
		},
		{
			name:    "valid quality_assurance category",
			input:   "quality_assurance",
			want:    CategoryQualityAssurance,
			wantErr: false,
		},
		{
			name:    "valid data_analysis category",
			input:   "data_analysis",
			want:    CategoryDataAnalysis,
			wantErr: false,
		},
		{
			name:    "uppercase input",
			input:   "DEVELOPMENT",
			want:    CategoryDevelopment,
			wantErr: false,
		},
		{
			name:    "mixed case input",
			input:   "DaTa_AnAlYsIs",
			want:    CategoryDataAnalysis,
			wantErr: false,
		},
		{
			name:    "input with spaces",
			input:   "  development  ",
			want:    CategoryDevelopment,
			wantErr: false,
		},
		{
			name:        "invalid category",
			input:       "invalid",
			want:        CategoryUnknown,
			wantErr:     true,
			errExpected: ErrInvalidCategory,
		},
		{
			name:        "empty string",
			input:       "",
			want:        CategoryUnknown,
			wantErr:     true,
			errExpected: ErrInvalidCategory,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewCategory(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Errorf("NewCategory() error = nil, wantErr %v", tt.wantErr)
					return
				}
				if err != tt.errExpected {
					t.Errorf("NewCategory() error = %v, wantErr %v", err, tt.errExpected)
					return
				}
			} else if err != nil {
				t.Errorf("NewCategory() unexpected error = %v", err)
				return
			}

			if got != tt.want {
				t.Errorf("NewCategory() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMustNewCategory(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		want      Category
		wantPanic bool
	}{
		{
			name:      "valid category",
			input:     "development",
			want:      CategoryDevelopment,
			wantPanic: false,
		},
		{
			name:      "invalid category",
			input:     "invalid",
			wantPanic: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				r := recover()
				if (r != nil) != tt.wantPanic {
					t.Errorf("MustNewCategory() panic = %v, wantPanic %v", r, tt.wantPanic)
				}
			}()

			got := MustNewCategory(tt.input)
			if !tt.wantPanic && got != tt.want {
				t.Errorf("MustNewCategory() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCategory_String(t *testing.T) {
	tests := []struct {
		name     string
		category Category
		want     string
	}{
		{
			name:     "operations category",
			category: CategoryOperations,
			want:     "operations",
		},
		{
			name:     "development category",
			category: CategoryDevelopment,
			want:     "development",
		},
		{
			name:     "product category",
			category: CategoryProduct,
			want:     "product",
		},
		{
			name:     "quality_assurance category",
			category: CategoryQualityAssurance,
			want:     "quality_assurance",
		},
		{
			name:     "data_analysis category",
			category: CategoryDataAnalysis,
			want:     "data_analysis",
		},
		{
			name:     "unknown category",
			category: CategoryUnknown,
			want:     "unknown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.category.String(); got != tt.want {
				t.Errorf("Category.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCategory_IsValid(t *testing.T) {
	tests := []struct {
		name     string
		category Category
		want     bool
	}{
		{
			name:     "valid operations category",
			category: CategoryOperations,
			want:     true,
		},
		{
			name:     "valid development category",
			category: CategoryDevelopment,
			want:     true,
		},
		{
			name:     "valid product category",
			category: CategoryProduct,
			want:     true,
		},
		{
			name:     "valid quality_assurance category",
			category: CategoryQualityAssurance,
			want:     true,
		},
		{
			name:     "valid data_analysis category",
			category: CategoryDataAnalysis,
			want:     true,
		},
		{
			name:     "valid unknown category",
			category: CategoryUnknown,
			want:     true,
		},
		{
			name:     "invalid category",
			category: "invalid",
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.category.IsValid(); got != tt.want {
				t.Errorf("Category.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCategory_TypeCheckers(t *testing.T) {
	tests := []struct {
		name     string
		category Category
		checks   map[string]bool
	}{
		{
			name:     "operations category",
			category: CategoryOperations,
			checks: map[string]bool{
				"IsOperations":       true,
				"IsDevelopment":      false,
				"IsProduct":          false,
				"IsQualityAssurance": false,
				"IsDataAnalysis":     false,
				"IsUnknown":          false,
			},
		},
		{
			name:     "development category",
			category: CategoryDevelopment,
			checks: map[string]bool{
				"IsOperations":       false,
				"IsDevelopment":      true,
				"IsProduct":          false,
				"IsQualityAssurance": false,
				"IsDataAnalysis":     false,
				"IsUnknown":          false,
			},
		},
		{
			name:     "unknown category",
			category: CategoryUnknown,
			checks: map[string]bool{
				"IsOperations":       false,
				"IsDevelopment":      false,
				"IsProduct":          false,
				"IsQualityAssurance": false,
				"IsDataAnalysis":     false,
				"IsUnknown":          true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.category.IsOperations(); got != tt.checks["IsOperations"] {
				t.Errorf("Category.IsOperations() = %v, want %v", got, tt.checks["IsOperations"])
			}
			if got := tt.category.IsDevelopment(); got != tt.checks["IsDevelopment"] {
				t.Errorf("Category.IsDevelopment() = %v, want %v", got, tt.checks["IsDevelopment"])
			}
			if got := tt.category.IsProduct(); got != tt.checks["IsProduct"] {
				t.Errorf("Category.IsProduct() = %v, want %v", got, tt.checks["IsProduct"])
			}
			if got := tt.category.IsQualityAssurance(); got != tt.checks["IsQualityAssurance"] {
				t.Errorf("Category.IsQualityAssurance() = %v, want %v", got, tt.checks["IsQualityAssurance"])
			}
			if got := tt.category.IsDataAnalysis(); got != tt.checks["IsDataAnalysis"] {
				t.Errorf("Category.IsDataAnalysis() = %v, want %v", got, tt.checks["IsDataAnalysis"])
			}
			if got := tt.category.IsUnknown(); got != tt.checks["IsUnknown"] {
				t.Errorf("Category.IsUnknown() = %v, want %v", got, tt.checks["IsUnknown"])
			}
		})
	}
}
