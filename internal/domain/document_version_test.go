package domain

import (
	"testing"
	"time"
)

func TestNewDocumentVersion(t *testing.T) {
	tests := []struct {
		name      string
		version   uint
		timestamp time.Time
		wantErr   bool
	}{
		{
			name:      "valid version",
			version:   1,
			timestamp: time.Now(),
			wantErr:   false,
		},
		{
			name:      "zero version is valid",
			version:   0,
			timestamp: time.Now(),
			wantErr:   false,
		},
		{
			name:      "zero timestamp is invalid",
			version:   1,
			timestamp: time.Time{},
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewDocumentVersion(tt.version, tt.timestamp)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewDocumentVersion() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Error("NewDocumentVersion() returned nil but expected a value")
			}
			if !tt.wantErr && got.Version() != tt.version {
				t.Errorf("NewDocumentVersion() version = %v, want %v", got.Version(), tt.version)
			}
		})
	}
}

func TestNewDefaultDocumentVersion(t *testing.T) {
	dv := NewDefaultDocumentVersion()
	if dv == nil {
		t.Fatal("NewDefaultDocumentVersion() returned nil")
	}

	if dv.Version() != 1 {
		t.Errorf("NewDefaultDocumentVersion() version = %v, want 1", dv.Version())
	}

	if dv.Timestamp().IsZero() {
		t.Error("NewDefaultDocumentVersion() timestamp is zero")
	}
}

func TestDocumentVersion_String(t *testing.T) {
	timestamp := time.Now()
	tests := []struct {
		name    string
		version uint
		want    string
	}{
		{
			name:    "version 1",
			version: 1,
			want:    "v1",
		},
		{
			name:    "version 0",
			version: 0,
			want:    "v0",
		},
		{
			name:    "version 999",
			version: 999,
			want:    "v999",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dv, _ := NewDocumentVersion(tt.version, timestamp)
			if got := dv.String(); got != tt.want {
				t.Errorf("DocumentVersion.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDocumentVersion_IsNewer(t *testing.T) {
	now := time.Now()
	later := now.Add(time.Hour)

	tests := []struct {
		name    string
		current *DocumentVersion
		other   *DocumentVersion
		want    bool
		wantErr bool
	}{
		{
			name:    "higher version is newer",
			current: &DocumentVersion{version: 2, timestamp: now},
			other:   &DocumentVersion{version: 1, timestamp: now},
			want:    true,
			wantErr: false,
		},
		{
			name:    "lower version is older",
			current: &DocumentVersion{version: 1, timestamp: now},
			other:   &DocumentVersion{version: 2, timestamp: now},
			want:    false,
			wantErr: false,
		},
		{
			name:    "same version, later timestamp is newer",
			current: &DocumentVersion{version: 1, timestamp: later},
			other:   &DocumentVersion{version: 1, timestamp: now},
			want:    true,
			wantErr: false,
		},
		{
			name:    "same version, earlier timestamp is older",
			current: &DocumentVersion{version: 1, timestamp: now},
			other:   &DocumentVersion{version: 1, timestamp: later},
			want:    false,
			wantErr: false,
		},
		{
			name:    "nil current version",
			current: nil,
			other:   &DocumentVersion{version: 1, timestamp: now},
			want:    false,
			wantErr: true,
		},
		{
			name:    "nil other version",
			current: &DocumentVersion{version: 1, timestamp: now},
			other:   nil,
			want:    false,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.current.IsNewer(tt.other)
			if (err != nil) != tt.wantErr {
				t.Errorf("DocumentVersion.IsNewer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("DocumentVersion.IsNewer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDocumentVersion_Equals(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name    string
		current *DocumentVersion
		other   *DocumentVersion
		want    bool
	}{
		{
			name:    "same version and timestamp",
			current: &DocumentVersion{version: 1, timestamp: now},
			other:   &DocumentVersion{version: 1, timestamp: now},
			want:    true,
		},
		{
			name:    "different version",
			current: &DocumentVersion{version: 1, timestamp: now},
			other:   &DocumentVersion{version: 2, timestamp: now},
			want:    false,
		},
		{
			name:    "different timestamp",
			current: &DocumentVersion{version: 1, timestamp: now},
			other:   &DocumentVersion{version: 1, timestamp: now.Add(time.Hour)},
			want:    false,
		},
		{
			name:    "nil current",
			current: nil,
			other:   &DocumentVersion{version: 1, timestamp: now},
			want:    false,
		},
		{
			name:    "nil other",
			current: &DocumentVersion{version: 1, timestamp: now},
			other:   nil,
			want:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.current.Equals(tt.other); got != tt.want {
				t.Errorf("DocumentVersion.Equals() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDocumentVersion_Increment(t *testing.T) {
	initial, _ := NewDocumentVersion(1, time.Now())
	time.Sleep(time.Millisecond) // Ensure timestamp difference

	incremented := initial.Increment()

	if incremented == nil {
		t.Fatal("Increment() returned nil")
	}

	if incremented.Version() != initial.Version()+1 {
		t.Errorf("Increment() version = %v, want %v", incremented.Version(), initial.Version()+1)
	}

	if !incremented.Timestamp().After(initial.Timestamp()) {
		t.Error("Increment() timestamp should be later than original")
	}

	// Verify immutability
	if initial.Version() != 1 {
		t.Error("Original version should remain unchanged")
	}
}
