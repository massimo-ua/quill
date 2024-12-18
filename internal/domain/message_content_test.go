package domain

import (
	"errors"
	"strings"
	"testing"
)

func TestNewMessageContent(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectError bool
		errorType   error
	}{
		{
			name:        "valid content",
			input:       "This is a valid message",
			expectError: false,
		},
		{
			name:        "empty string",
			input:       "",
			expectError: true,
			errorType:   ErrEmptyContent,
		},
		{
			name:        "only whitespace",
			input:       "   \t\n",
			expectError: true,
			errorType:   ErrEmptyContent,
		},
		{
			name:        "trims whitespace",
			input:       "  valid message  ",
			expectError: false,
		},
		{
			name:        "very long content",
			input:       strings.Repeat("a", MaxContentLength+100),
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			content, err := NewMessageContent(tt.input)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
				if err != tt.errorType {
					t.Errorf("expected error %v but got %v", tt.errorType, err)
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if content == nil {
				t.Error("expected content to not be nil")
				return
			}

			actualText := content.Text()
			expectedText := strings.TrimSpace(tt.input)
			if len(tt.input) > MaxContentLength {
				expectedText = expectedText[:MaxContentLength]
			}

			if actualText != expectedText {
				t.Errorf("expected text %q but got %q", expectedText, actualText)
			}
		})
	}
}

func TestMustNewMessageContent(t *testing.T) {
	t.Run("valid content", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("unexpected panic: %v", r)
			}
		}()

		content := MustNewMessageContent("valid message")
		if content == nil {
			t.Error("expected content to not be nil")
		}
	})

	t.Run("invalid content should panic", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic but got none")
			}
		}()

		MustNewMessageContent("")
	})
}

func TestMessageContent_ContainsURL(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "contains http URL",
			input:    "Check this link http://example.com",
			expected: true,
		},
		{
			name:     "contains https URL",
			input:    "Check this link https://example.com",
			expected: true,
		},
		{
			name:     "no URL",
			input:    "Plain text without URL",
			expected: false,
		},
		{
			name:     "partial URL",
			input:    "Invalid example.com URL",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			content := MustNewMessageContent(tt.input)
			if content.ContainsURL() != tt.expected {
				t.Errorf("expected ContainsURL() to be %v", tt.expected)
			}
		})
	}
}

func TestMessageContent_WordCount(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int
	}{
		{
			name:     "empty content",
			input:    "test",
			expected: 1,
		},
		{
			name:     "multiple words",
			input:    "this is a test message",
			expected: 5,
		},
		{
			name:     "multiple spaces",
			input:    "this   is   a   test",
			expected: 4,
		},
		{
			name:     "with newlines",
			input:    "this\nis\na\ntest",
			expected: 4,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			content := MustNewMessageContent(tt.input)
			if count := content.WordCount(); count != tt.expected {
				t.Errorf("expected WordCount() to be %d but got %d", tt.expected, count)
			}
		})
	}
}

func TestMessageContent_ContainsTag(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		tag      string
		expected bool
	}{
		{
			name:     "contains idea tag",
			content:  "This is an #idea for improvement",
			tag:      "idea",
			expected: true,
		},
		{
			name:     "contains decision tag",
			content:  "We made a #decision today",
			tag:      "decision",
			expected: true,
		},
		{
			name:     "no tag",
			content:  "Regular message without tags",
			tag:      "idea",
			expected: false,
		},
		{
			name:     "partial tag match",
			content:  "This is my ideas",
			tag:      "idea",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			content := MustNewMessageContent(tt.content)
			if result := content.ContainsTag(tt.tag); result != tt.expected {
				t.Errorf("expected ContainsTag(%q) to be %v but got %v", tt.tag, tt.expected, result)
			}
		})
	}
}

func TestMessageContent_JSON(t *testing.T) {
	t.Run("marshal and unmarshal", func(t *testing.T) {
		original := MustNewMessageContent("test message")

		// Marshal
		data, err := original.MarshalJSON()
		if err != nil {
			t.Fatalf("failed to marshal: %v", err)
		}

		// Unmarshal
		var decoded MessageContent
		err = decoded.UnmarshalJSON(data)
		if err != nil {
			t.Fatalf("failed to unmarshal: %v", err)
		}

		// Compare
		if decoded.Text() != original.Text() {
			t.Errorf("expected text %q but got %q", original.Text(), decoded.Text())
		}
	})

	t.Run("empty content", func(t *testing.T) {
		var content MessageContent
		err := content.UnmarshalJSON([]byte(`{"text": ""}`))
		if err == nil {
			t.Error("expected error for empty content but got none")
		}
		if !errors.Is(err, ErrEmptyContent) {
			t.Errorf("expected ErrEmptyContent but got %v", err)
		}
	})

	t.Run("malformed json", func(t *testing.T) {
		var content MessageContent
		err := content.UnmarshalJSON([]byte(`{malformed`))
		if err == nil {
			t.Error("expected error for malformed JSON but got none")
		}
		if !errors.Is(err, ErrInvalidJSON) {
			t.Errorf("expected ErrInvalidJSON but got %v", err)
		}
	})

	t.Run("invalid type in json", func(t *testing.T) {
		var content MessageContent
		err := content.UnmarshalJSON([]byte(`{"text": 123}`))
		if err == nil {
			t.Error("expected error for invalid JSON type but got none")
		}
		if !errors.Is(err, ErrInvalidJSON) {
			t.Errorf("expected ErrInvalidJSON but got %v", err)
		}
	})
}
