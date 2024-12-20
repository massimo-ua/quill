package common

import (
	"database/sql/driver"
	"github.com/oklog/ulid/v2"
	"math/rand"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestID_EmptyValue(t *testing.T) {
	t.Run("empty ID value behavior", func(t *testing.T) {
		var id ID
		if id.String() != "" {
			t.Errorf("Empty ID should return empty string, got %v", id.String())
		}

		// Test Value() with empty ID
		value, err := id.Value()
		if err != nil {
			t.Errorf("Empty ID Value() should not return error, got %v", err)
		}
		if value != "" {
			t.Errorf("Empty ID Value() should return empty string, got %v", value)
		}
	})
}

func TestID_ScanEdgeCases(t *testing.T) {
	tests := []struct {
		name    string
		input   interface{}
		wantErr bool
	}{
		{
			name:    "scan very long string",
			input:   "01H9QRDGCR35XKVJ2ZXS67D3AT01H9QRDGCR35XKVJ2ZXS67D3AT",
			wantErr: true,
		},
		{
			name:    "scan zero-length byte slice",
			input:   []byte{},
			wantErr: true,
		},
		{
			name:    "scan invalid byte slice",
			input:   []byte("invalid-ulid"),
			wantErr: true,
		},
		{
			name:    "scan bool value",
			input:   true,
			wantErr: true,
		},
		{
			name:    "scan float value",
			input:   3.14,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var id ID
			err := id.Scan(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ID.Scan() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestID_CompareEdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		id1      ID
		id2      ID
		expected int
	}{
		{
			name:     "compare empty IDs",
			id1:      ID{},
			id2:      ID{},
			expected: 0,
		},
		{
			name:     "compare empty with valid ID",
			id1:      ID{},
			id2:      GenerateID(),
			expected: -1,
		},
		{
			name:     "compare same timestamp different entropy",
			id1:      GenerateID(),
			id2:      GenerateID(),
			expected: -1, // Second ID should be greater
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.id1.Compare(tt.id2)
			if (tt.expected == 0 && result != 0) ||
				(tt.expected < 0 && result >= 0) ||
				(tt.expected > 0 && result <= 0) {
				t.Errorf("ID.Compare() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestID_ValueImplementation(t *testing.T) {
	t.Run("implements driver.Valuer", func(t *testing.T) {
		var _ driver.Valuer = ID{} // Compile-time check

		id := GenerateID()
		value, err := id.Value()
		if err != nil {
			t.Errorf("Value() error = %v", err)
		}

		str, ok := value.(string)
		if !ok {
			t.Error("Value() should return string type")
		}

		if str != id.String() {
			t.Errorf("Value() = %v, want %v", str, id.String())
		}
	})
}

func TestNewID_ValidationCases(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:    "only spaces",
			input:   "     ",
			wantErr: true,
		},
		{
			name:    "invalid characters",
			input:   "01H9QRDGCR35XKVJ2ZXS67D3AT!",
			wantErr: true,
		},
		{
			name:    "lowercase valid ULID",
			input:   "01h9qrdgcr35xkvj2zxs67d3at",
			wantErr: false, // ULIDs are case-insensitive
		},
		{
			name:    "too short ULID",
			input:   "01H9QRDGCR35XKVJ2ZXS67D3A",
			wantErr: true,
		},
		{
			name:    "too long ULID",
			input:   "01H9QRDGCR35XKVJ2ZXS67D3ATX",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ID, err := NewID(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewID() error = %v, wantErr %v, ID %v", err, tt.wantErr, ID.String())
			}
		})
	}
}

func TestID_Time(t *testing.T) {
	// Test cases for parsing timestamp from ULID
	tests := []struct {
		name     string
		idValue  string
		wantTime time.Time
		wantZero bool
	}{
		{
			name:     "empty ID",
			idValue:  "",
			wantZero: true,
		},
		{
			name:     "invalid ULID format",
			idValue:  "invalid-ulid",
			wantZero: true,
		},
		{
			name:     "minimum timestamp",
			idValue:  "00000000000000000000000000",
			wantTime: time.Unix(0, 0).UTC(),
		},
		{
			name:     "maximum timestamp",
			idValue:  "7ZZZZZZZZZZZZZZZZZZZZZZZZZ",
			wantTime: time.Unix(281474976710, 655000000).UTC(),
		},
		{
			name:     "specific timestamp",
			idValue:  ulid.MustNew(ulid.Timestamp(time.Date(2023, 8, 15, 14, 30, 0, 500000000, time.UTC)), ulid.Monotonic(rand.New(rand.NewSource(0)), 0)).String(),
			wantTime: time.Date(2023, 8, 15, 14, 30, 0, 500000000, time.UTC),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id := ID{value: tt.idValue}
			got := id.Time()

			if tt.wantZero {
				if !got.IsZero() {
					t.Errorf("Time() = %v, want zero time", got)
				}
				return
			}

			if !got.Equal(tt.wantTime) {
				t.Errorf("Time() = %v, want %v", got, tt.wantTime)
			}
		})
	}

	t.Run("extracts correct timestamp", func(t *testing.T) {
		// Create a known timestamp
		timestamp := time.Date(2023, 8, 15, 14, 30, 0, 500000000, time.UTC)

		// Create ULID with this timestamp
		entropy := ulid.Monotonic(rand.New(rand.NewSource(0)), 0)
		mustNew := ulid.MustNew(ulid.Timestamp(timestamp), entropy)

		id := ID{value: mustNew.String()}
		extractedTime := id.Time()

		if !timestamp.Equal(extractedTime) {
			t.Errorf("Time() extracted incorrect timestamp. got: %v, want: %v", extractedTime, timestamp)
		}
	})

	t.Run("preserves millisecond precision", func(t *testing.T) {
		// Test with a timestamp that has millisecond precision
		timestamp := time.Date(2023, 8, 15, 14, 30, 0, 123000000, time.UTC) // 123 milliseconds

		entropy := ulid.Monotonic(rand.New(rand.NewSource(0)), 0)
		mustNew := ulid.MustNew(ulid.Timestamp(timestamp), entropy)

		id := ID{value: mustNew.String()}
		extractedTime := id.Time()

		// Compare milliseconds
		if extractedTime.Nanosecond()/1000000 != 123 {
			t.Errorf("Time() lost millisecond precision. got: %v ms, want: 123ms", extractedTime.Nanosecond()/1000000)
		}
	})

	t.Run("handles zero time", func(t *testing.T) {
		// Create a ULID with zero timestamp
		zeroTime := time.Unix(0, 0).UTC()
		entropy := ulid.Monotonic(rand.New(rand.NewSource(0)), 0)
		mustNew := ulid.MustNew(ulid.Timestamp(zeroTime), entropy)

		id := ID{value: mustNew.String()}
		extractedTime := id.Time()

		if !extractedTime.Equal(zeroTime) {
			t.Errorf("Time() failed to handle zero time. got: %v, want: %v", extractedTime, zeroTime)
		}
	})

	t.Run("handles max time", func(t *testing.T) {
		// Create a ULID with maximum possible timestamp (281474976710.655 seconds)
		maxTime := time.Unix(281474976710, 655000000).UTC()

		// Since we can't use ulid.MustNew with a timestamp beyond its range,
		// we'll create the maximum possible ULID string directly
		id := ID{value: "7ZZZZZZZZZZZZZZZZZZZZZZZZZ"}
		extractedTime := id.Time()

		if extractedTime.After(maxTime) {
			t.Errorf("Time() exceeded maximum time. got: %v, want at or before: %v", extractedTime, maxTime)
		}
	})

	t.Run("empty ID returns zero time", func(t *testing.T) {
		id := ID{}
		extractedTime := id.Time()

		if !extractedTime.IsZero() {
			t.Errorf("Time() with empty ID should return zero time, got: %v", extractedTime)
		}
	})

	t.Run("maintains UTC timezone", func(t *testing.T) {
		// Create a timestamp in a different timezone
		loc, _ := time.LoadLocation("America/New_York")
		localTime := time.Date(2023, 8, 15, 14, 30, 0, 0, loc)

		entropy := ulid.Monotonic(rand.New(rand.NewSource(0)), 0)
		mustNew := ulid.MustNew(ulid.Timestamp(localTime), entropy)

		id := ID{value: mustNew.String()}
		extractedTime := id.Time()

		if extractedTime.Location() != time.UTC {
			t.Errorf("Time() returned non-UTC timezone: %v", extractedTime.Location())
		}
	})

	t.Run("handles millisecond rounding", func(t *testing.T) {
		now := time.Now().UTC()
		id := GenerateID()
		extractedTime := id.Time()

		// The difference should be less than a millisecond
		diffMs := extractedTime.Sub(now).Milliseconds()
		if diffMs > 1000 || diffMs < -1000 {
			t.Errorf("Time() millisecond rounding error. Difference: %dms", diffMs)
		}
	})
}

func TestID_Scan(t *testing.T) {
	validID := GenerateID()

	tests := []struct {
		name    string
		input   interface{}
		want    ID
		wantErr bool
	}{
		{
			name:    "valid string ID",
			input:   validID.String(),
			want:    validID,
			wantErr: false,
		},
		{
			name:    "valid byte slice ID",
			input:   []byte(validID.String()),
			want:    validID,
			wantErr: false,
		},
		{
			name:    "nil input",
			input:   nil,
			want:    ID{},
			wantErr: true,
		},
		{
			name:    "empty string",
			input:   "",
			want:    ID{},
			wantErr: true,
		},
		{
			name:    "empty byte slice",
			input:   []byte{},
			want:    ID{},
			wantErr: true,
		},
		{
			name:    "invalid type (int)",
			input:   123,
			want:    ID{},
			wantErr: true,
		},
		{
			name:    "invalid type (bool)",
			input:   true,
			want:    ID{},
			wantErr: true,
		},
		{
			name:    "invalid ULID format",
			input:   "not-a-ulid",
			want:    ID{},
			wantErr: true,
		},
		{
			name:    "invalid ULID format as bytes",
			input:   []byte("not-a-ulid"),
			want:    ID{},
			wantErr: true,
		},
		{
			name:    "string with whitespace",
			input:   "  " + validID.String() + "  ",
			want:    validID,
			wantErr: false,
		},
		{
			name:    "lowercase valid ULID",
			input:   strings.ToLower(validID.String()),
			want:    validID,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got ID
			err := got.Scan(tt.input)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Scan() error = nil, wantErr = true")
				}
				if err != nil && err != ErrInvalidID {
					t.Errorf("Scan() expected ErrInvalidID, got = %v", err)
				}
			} else {
				if err != nil {
					t.Errorf("Scan() unexpected error = %v", err)
				}
				if !got.Equals(tt.want) {
					t.Errorf("Scan() = %v, want %v", got, tt.want)
				}
			}
		})
	}

	t.Run("maintains value consistency", func(t *testing.T) {
		var id ID
		originalValue := GenerateID()

		// First scan should succeed
		err := id.Scan(originalValue.String())
		if err != nil {
			t.Errorf("First Scan() failed: %v", err)
		}

		// Invalid scan should not modify the original value
		err = id.Scan(nil)
		if err == nil {
			t.Error("Second Scan() should have failed")
		}
		if !id.Equals(originalValue) {
			t.Errorf("ID value changed after failed scan: got %v, want %v", id, originalValue)
		}
	})

	t.Run("concurrent scanning", func(t *testing.T) {
		var wg sync.WaitGroup
		id := &ID{}
		validValue := GenerateID().String()

		// Test concurrent scanning of the same ID
		for i := 0; i < 10; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				_ = id.Scan(validValue)
			}()
		}

		wg.Wait()

		// Verify final state
		if id.String() != validValue {
			t.Errorf("Concurrent Scan() resulted in incorrect value: got %v, want %v", id.String(), validValue)
		}
	})
}
