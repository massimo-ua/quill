package domain

import (
	"testing"
)

func TestNewMessageType(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    MessageType
		wantErr bool
	}{
		{
			name:    "valid idea type",
			input:   "idea",
			want:    MessageTypeIdea,
			wantErr: false,
		},
		{
			name:    "valid decision type",
			input:   "decision",
			want:    MessageTypeDecision,
			wantErr: false,
		},
		{
			name:    "valid status type",
			input:   "status",
			want:    MessageTypeStatus,
			wantErr: false,
		},
		{
			name:    "valid type with spaces",
			input:   "  idea  ",
			want:    MessageTypeIdea,
			wantErr: false,
		},
		{
			name:    "valid type with mixed case",
			input:   "StAtUs",
			want:    MessageTypeStatus,
			wantErr: false,
		},
		{
			name:    "invalid type",
			input:   "invalid",
			want:    MessageTypeUnknown,
			wantErr: true,
		},
		{
			name:    "empty string",
			input:   "",
			want:    MessageTypeUnknown,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewMessageType(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewMessageType() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("NewMessageType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMustNewMessageType(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		want      MessageType
		wantPanic bool
	}{
		{
			name:      "valid idea type",
			input:     "idea",
			want:      MessageTypeIdea,
			wantPanic: false,
		},
		{
			name:      "invalid type",
			input:     "invalid",
			wantPanic: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				r := recover()
				if (r != nil) != tt.wantPanic {
					t.Errorf("MustNewMessageType() panic = %v, wantPanic %v", r, tt.wantPanic)
				}
			}()

			if got := MustNewMessageType(tt.input); !tt.wantPanic && got != tt.want {
				t.Errorf("MustNewMessageType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMessageType_IsValid(t *testing.T) {
	tests := []struct {
		name string
		mt   MessageType
		want bool
	}{
		{
			name: "valid idea type",
			mt:   MessageTypeIdea,
			want: true,
		},
		{
			name: "valid decision type",
			mt:   MessageTypeDecision,
			want: true,
		},
		{
			name: "valid status type",
			mt:   MessageTypeStatus,
			want: true,
		},
		{
			name: "invalid type",
			mt:   MessageType("invalid"),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.mt.IsValid(); got != tt.want {
				t.Errorf("MessageType.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMessageType_TypeChecks(t *testing.T) {
	tests := []struct {
		name       string
		mt         MessageType
		isIdea     bool
		isStatus   bool
		isDecision bool
		isUnknown  bool
	}{
		{
			name:       "idea type",
			mt:         MessageTypeIdea,
			isIdea:     true,
			isStatus:   false,
			isDecision: false,
			isUnknown:  false,
		},
		{
			name:       "status type",
			mt:         MessageTypeStatus,
			isIdea:     false,
			isStatus:   true,
			isDecision: false,
			isUnknown:  false,
		},
		{
			name:       "decision type",
			mt:         MessageTypeDecision,
			isIdea:     false,
			isStatus:   false,
			isDecision: true,
			isUnknown:  false,
		},
		{
			name:       "unknown type",
			mt:         MessageTypeUnknown,
			isIdea:     false,
			isStatus:   false,
			isDecision: false,
			isUnknown:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.mt.IsIdea(); got != tt.isIdea {
				t.Errorf("MessageType.IsIdea() = %v, want %v", got, tt.isIdea)
			}
			if got := tt.mt.IsStatus(); got != tt.isStatus {
				t.Errorf("MessageType.IsStatus() = %v, want %v", got, tt.isStatus)
			}
			if got := tt.mt.IsDecision(); got != tt.isDecision {
				t.Errorf("MessageType.IsDecision() = %v, want %v", got, tt.isDecision)
			}
			if got := tt.mt.IsUnknown(); got != tt.isUnknown {
				t.Errorf("MessageType.IsUnknown() = %v, want %v", got, tt.isUnknown)
			}
		})
	}
}

func TestMessageType_String(t *testing.T) {
	tests := []struct {
		name string
		mt   MessageType
		want string
	}{
		{
			name: "idea type",
			mt:   MessageTypeIdea,
			want: "idea",
		},
		{
			name: "decision type",
			mt:   MessageTypeDecision,
			want: "decision",
		},
		{
			name: "status type",
			mt:   MessageTypeStatus,
			want: "status",
		},
		{
			name: "unknown type",
			mt:   MessageTypeUnknown,
			want: "unknown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.mt.String(); got != tt.want {
				t.Errorf("MessageType.String() = %v, want %v", got, tt.want)
			}
		})
	}
}
