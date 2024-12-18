package domain

import (
	"strings"
	"testing"
)

func TestNewReference(t *testing.T) {
	tests := []struct {
		name      string
		refType   ReferenceType
		value     string
		wantError bool
	}{
		{
			name:      "valid message reference",
			refType:   ReferenceTypeMessage,
			value:     "msg_123",
			wantError: false,
		},
		{
			name:      "valid document reference",
			refType:   ReferenceTypeDocument,
			value:     "docs/decisions/001.md",
			wantError: false,
		},
		{
			name:      "invalid reference type",
			refType:   ReferenceType("invalid"),
			value:     "test",
			wantError: true,
		},
		{
			name:      "empty value",
			refType:   ReferenceTypeMessage,
			value:     "",
			wantError: true,
		},
		{
			name:      "whitespace value",
			refType:   ReferenceTypeMessage,
			value:     "   ",
			wantError: true,
		},
		{
			name:      "trims whitespace from valid value",
			refType:   ReferenceTypeMessage,
			value:     "  msg_123  ",
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ref, err := NewReference(tt.refType, tt.value)
			if tt.wantError {
				if err == nil {
					t.Error("expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if ref.Type() != tt.refType {
					t.Errorf("expected type %v, got %v", tt.refType, ref.Type())
				}
				if ref.Value() != strings.TrimSpace(tt.value) {
					t.Errorf("expected value %v, got %v", strings.TrimSpace(tt.value), ref.Value())
				}
			}
		})
	}
}

func TestNewMessageReference(t *testing.T) {
	tests := []struct {
		name      string
		messageID string
		wantError bool
	}{
		{
			name:      "valid message ID",
			messageID: "msg_123",
			wantError: false,
		},
		{
			name:      "empty message ID",
			messageID: "",
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ref, err := NewMessageReference(tt.messageID)
			if tt.wantError {
				if err == nil {
					t.Error("expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if !ref.Type().IsMessage() {
					t.Error("expected message reference type")
				}
				if ref.Value() != tt.messageID {
					t.Errorf("expected value %v, got %v", tt.messageID, ref.Value())
				}
			}
		})
	}
}

func TestNewDocumentReference(t *testing.T) {
	tests := []struct {
		name         string
		documentPath string
		wantError    bool
	}{
		{
			name:         "valid document path",
			documentPath: "docs/decisions/001.md",
			wantError:    false,
		},
		{
			name:         "empty document path",
			documentPath: "",
			wantError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ref, err := NewDocumentReference(tt.documentPath)
			if tt.wantError {
				if err == nil {
					t.Error("expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if !ref.Type().IsDocument() {
					t.Error("expected document reference type")
				}
				if ref.Value() != tt.documentPath {
					t.Errorf("expected value %v, got %v", tt.documentPath, ref.Value())
				}
			}
		})
	}
}

func TestMustNewReference(t *testing.T) {
	tests := []struct {
		name      string
		refType   ReferenceType
		value     string
		wantPanic bool
	}{
		{
			name:      "valid reference",
			refType:   ReferenceTypeMessage,
			value:     "msg_123",
			wantPanic: false,
		},
		{
			name:      "invalid reference",
			refType:   ReferenceType("invalid"),
			value:     "test",
			wantPanic: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				r := recover()
				if tt.wantPanic && r == nil {
					t.Error("expected panic, got none")
				}
				if !tt.wantPanic && r != nil {
					t.Errorf("unexpected panic: %v", r)
				}
			}()

			ref := MustNewReference(tt.refType, tt.value)
			if tt.wantPanic {
				t.Error("expected panic, got none")
			} else {
				if ref.Type() != tt.refType {
					t.Errorf("expected type %v, got %v", tt.refType, ref.Type())
				}
				if ref.Value() != tt.value {
					t.Errorf("expected value %v, got %v", tt.value, ref.Value())
				}
			}
		})
	}
}

func TestReference_String(t *testing.T) {
	tests := []struct {
		name     string
		ref      *Reference
		expected string
	}{
		{
			name:     "message reference",
			ref:      MustNewReference(ReferenceTypeMessage, "msg_123"),
			expected: "message:msg_123",
		},
		{
			name:     "document reference",
			ref:      MustNewReference(ReferenceTypeDocument, "docs/decisions/001.md"),
			expected: "document:docs/decisions/001.md",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.ref.String(); got != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, got)
			}
		})
	}
}

func TestReference_Equals(t *testing.T) {
	ref1 := MustNewReference(ReferenceTypeMessage, "msg_123")
	ref2 := MustNewReference(ReferenceTypeMessage, "msg_123")
	ref3 := MustNewReference(ReferenceTypeMessage, "msg_456")
	ref4 := MustNewReference(ReferenceTypeDocument, "msg_123")

	tests := []struct {
		name     string
		ref1     *Reference
		ref2     *Reference
		expected bool
	}{
		{
			name:     "same reference",
			ref1:     ref1,
			ref2:     ref1,
			expected: true,
		},
		{
			name:     "equal references",
			ref1:     ref1,
			ref2:     ref2,
			expected: true,
		},
		{
			name:     "different values",
			ref1:     ref1,
			ref2:     ref3,
			expected: false,
		},
		{
			name:     "different types",
			ref1:     ref1,
			ref2:     ref4,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.ref1.Equals(*tt.ref2); got != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, got)
			}
		})
	}
}

func TestParseReference(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantType  ReferenceType
		wantValue string
		wantError bool
	}{
		{
			name:      "valid message reference",
			input:     "message:msg_123",
			wantType:  ReferenceTypeMessage,
			wantValue: "msg_123",
			wantError: false,
		},
		{
			name:      "valid document reference",
			input:     "document:docs/decisions/001.md",
			wantType:  ReferenceTypeDocument,
			wantValue: "docs/decisions/001.md",
			wantError: false,
		},
		{
			name:      "invalid format",
			input:     "invalid_format",
			wantError: true,
		},
		{
			name:      "invalid reference type",
			input:     "invalid:value",
			wantError: true,
		},
		{
			name:      "empty value",
			input:     "message:",
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ref, err := ParseReference(tt.input)
			if tt.wantError {
				if err == nil {
					t.Error("expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if ref.Type() != tt.wantType {
					t.Errorf("expected type %v, got %v", tt.wantType, ref.Type())
				}
				if ref.Value() != tt.wantValue {
					t.Errorf("expected value %v, got %v", tt.wantValue, ref.Value())
				}
			}
		})
	}
}

func TestReferenceType_Methods(t *testing.T) {
	tests := []struct {
		name string
		rt   ReferenceType
		want struct {
			isValid    bool
			isMessage  bool
			isDocument bool
			isUnknown  bool
			string     string
		}
	}{
		{
			name: "message type",
			rt:   ReferenceTypeMessage,
			want: struct {
				isValid    bool
				isMessage  bool
				isDocument bool
				isUnknown  bool
				string     string
			}{
				isValid:    true,
				isMessage:  true,
				isDocument: false,
				isUnknown:  false,
				string:     "message",
			},
		},
		{
			name: "document type",
			rt:   ReferenceTypeDocument,
			want: struct {
				isValid    bool
				isMessage  bool
				isDocument bool
				isUnknown  bool
				string     string
			}{
				isValid:    true,
				isMessage:  false,
				isDocument: true,
				isUnknown:  false,
				string:     "document",
			},
		},
		{
			name: "unknown type",
			rt:   ReferenceTypeUnknown,
			want: struct {
				isValid    bool
				isMessage  bool
				isDocument bool
				isUnknown  bool
				string     string
			}{
				isValid:    true,
				isMessage:  false,
				isDocument: false,
				isUnknown:  true,
				string:     "unknown",
			},
		},
		{
			name: "invalid type",
			rt:   ReferenceType("invalid"),
			want: struct {
				isValid    bool
				isMessage  bool
				isDocument bool
				isUnknown  bool
				string     string
			}{
				isValid:    false,
				isMessage:  false,
				isDocument: false,
				isUnknown:  false,
				string:     "invalid",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.rt.IsValid(); got != tt.want.isValid {
				t.Errorf("IsValid() = %v, want %v", got, tt.want.isValid)
			}
			if got := tt.rt.IsMessage(); got != tt.want.isMessage {
				t.Errorf("IsMessage() = %v, want %v", got, tt.want.isMessage)
			}
			if got := tt.rt.IsDocument(); got != tt.want.isDocument {
				t.Errorf("IsDocument() = %v, want %v", got, tt.want.isDocument)
			}
			if got := tt.rt.IsUnknown(); got != tt.want.isUnknown {
				t.Errorf("IsUnknown() = %v, want %v", got, tt.want.isUnknown)
			}
			if got := tt.rt.String(); got != tt.want.string {
				t.Errorf("String() = %v, want %v", got, tt.want.string)
			}
		})
	}
}
