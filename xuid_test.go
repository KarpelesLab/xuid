package xuid

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"github.com/google/uuid"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		wantPrefix  string
		wantUUID    string
		expectError bool
	}{
		{
			name:       "Empty string prefix",
			input:      "aaaaaa-aaaa-aaaa-aaaa-aaaaaaaa",
			wantPrefix: "",
			wantUUID:   "00000000-0000-0000-0000-000000000000",
		},
		{
			name:       "Null UUID with prefix",
			input:      "null-aaaaaa-aaaa-aaaa-aaaa-aaaaaaaa",
			wantPrefix: "null",
			wantUUID:   "00000000-0000-0000-0000-000000000000",
		},
		{
			name:       "Shell example from README",
			input:      "shell-h4nu2n-zu3f-dmnn-kguv-6f643nei",
			wantPrefix: "shell",
			wantUUID:   "3f1b4d37-34d9-46c6-b546-a57c5f736d22",
		},
		// Use a known-good test example instead of making one up
		{
			name:       "New random example",
			input:      "test-h4nu2n-zu3f-dmnn-kguv-6f643nei",
			wantPrefix: "test",
			wantUUID:   "3f1b4d37-34d9-46c6-b546-a57c5f736d22",
		},
		{
			name:       "Long but valid input, UUID fallback",
			input:      "3f1b4d37-34d9-46c6-b546-a57c5f736d22",
			wantPrefix: "",
			wantUUID:   "3f1b4d37-34d9-46c6-b546-a57c5f736d22",
		},
		{
			name:        "Invalid format",
			input:       "invalid-format",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.input)
			if tt.expectError {
				if err == nil {
					t.Errorf("Parse() expected error for %q, got nil", tt.input)
				}
				return
			}
			if err != nil {
				t.Errorf("Parse() error = %v", err)
				return
			}
			if got.Prefix != tt.wantPrefix {
				t.Errorf("Parse() prefix = %q, want %q", got.Prefix, tt.wantPrefix)
			}
			if got.ToUUID() != tt.wantUUID {
				t.Errorf("Parse() UUID = %q, want %q", got.ToUUID(), tt.wantUUID)
			}
		})
	}
}

func TestMustParse(t *testing.T) {
	t.Run("Valid input", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("MustParse() panicked unexpectedly: %v", r)
			}
		}()

		uid := MustParse("shell-h4nu2n-zu3f-dmnn-kguv-6f643nei")
		if uid.ToUUID() != "3f1b4d37-34d9-46c6-b546-a57c5f736d22" {
			t.Errorf("MustParse() UUID = %q, want %q", uid.ToUUID(), "3f1b4d37-34d9-46c6-b546-a57c5f736d22")
		}
	})

	t.Run("Invalid input", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("MustParse() did not panic for invalid input")
			}
		}()

		_ = MustParse("invalid-format")
	})
}

func TestParsePrefix(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		prefix      string
		expectError bool
		wantErrType error
	}{
		{
			name:        "Matching prefix",
			input:       "shell-h4nu2n-zu3f-dmnn-kguv-6f643nei",
			prefix:      "shell",
			expectError: false,
		},
		{
			name:        "Non-matching prefix",
			input:       "shell-h4nu2n-zu3f-dmnn-kguv-6f643nei",
			prefix:      "admin",
			expectError: true,
			wantErrType: ErrBadPrefix,
		},
		{
			name:        "Invalid input",
			input:       "invalid-format",
			prefix:      "user",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParsePrefix(tt.input, tt.prefix)
			if tt.expectError {
				if err == nil {
					t.Errorf("ParsePrefix() expected error for %q with prefix %q, got nil", tt.input, tt.prefix)
				}
				if tt.wantErrType != nil && !errors.Is(err, tt.wantErrType) {
					t.Errorf("ParsePrefix() expected error type %v, got %v", tt.wantErrType, err)
				}
				return
			}
			if err != nil {
				t.Errorf("ParsePrefix() error = %v", err)
				return
			}
			if got.Prefix != tt.prefix {
				t.Errorf("ParsePrefix() prefix = %q, want %q", got.Prefix, tt.prefix)
			}
		})
	}
}

func TestNew(t *testing.T) {
	prefix := "test"
	xuid := New(prefix)

	if xuid.Prefix != prefix {
		t.Errorf("New() prefix = %q, want %q", xuid.Prefix, prefix)
	}

	// UUID should be valid
	if len(xuid.UUID) != 16 {
		t.Errorf("New() UUID length = %d, want 16", len(xuid.UUID))
	}

	// String representation should start with prefix
	strRep := xuid.String()
	expectedPrefix := prefix + "-"
	if strRep[:len(expectedPrefix)] != expectedPrefix {
		t.Errorf("New().String() prefix = %q, want %q", strRep[:len(expectedPrefix)], expectedPrefix)
	}
}

func TestNewRandom(t *testing.T) {
	prefix := "test"
	xuid, err := NewRandom(prefix)
	if err != nil {
		t.Errorf("NewRandom() error = %v", err)
		return
	}

	if xuid.Prefix != prefix {
		t.Errorf("NewRandom() prefix = %q, want %q", xuid.Prefix, prefix)
	}

	// UUID should be valid
	if len(xuid.UUID) != 16 {
		t.Errorf("NewRandom() UUID length = %d, want 16", len(xuid.UUID))
	}

	// Generate another one and make sure they're different
	xuid2, err := NewRandom(prefix)
	if err != nil {
		t.Errorf("NewRandom() second call error = %v", err)
		return
	}

	if xuid.UUID == xuid2.UUID {
		t.Errorf("NewRandom() generated the same UUID twice")
	}
}

func TestFromUUID(t *testing.T) {
	u := uuid.MustParse("3f1b4d37-34d9-46c6-b546-a57c5f736d22")
	prefix := "shell"

	xuid, err := FromUUID(u, prefix)
	if err != nil {
		t.Errorf("FromUUID() error = %v", err)
		return
	}

	if xuid.Prefix != prefix {
		t.Errorf("FromUUID() prefix = %q, want %q", xuid.Prefix, prefix)
	}

	if xuid.UUID.String() != u.String() {
		t.Errorf("FromUUID() UUID = %q, want %q", xuid.UUID.String(), u.String())
	}

	// Check string representation matches expected format
	expected := "shell-h4nu2n-zu3f-dmnn-kguv-6f643nei"
	if xuid.String() != expected {
		t.Errorf("FromUUID().String() = %q, want %q", xuid.String(), expected)
	}
}

func TestParseUUID(t *testing.T) {
	uuidStr := "3f1b4d37-34d9-46c6-b546-a57c5f736d22"
	prefix := "shell"

	xuid, err := ParseUUID(uuidStr, prefix)
	if err != nil {
		t.Errorf("ParseUUID() error = %v", err)
		return
	}

	if xuid.Prefix != prefix {
		t.Errorf("ParseUUID() prefix = %q, want %q", xuid.Prefix, prefix)
	}

	if xuid.UUID.String() != uuidStr {
		t.Errorf("ParseUUID() UUID = %q, want %q", xuid.UUID.String(), uuidStr)
	}

	// Test invalid UUID
	_, err = ParseUUID("invalid-uuid", prefix)
	if err == nil {
		t.Errorf("ParseUUID() with invalid UUID did not return an error")
	}
}

func TestEquality(t *testing.T) {
	u1 := MustParse("shell-h4nu2n-zu3f-dmnn-kguv-6f643nei")
	u2 := MustParse("shell-h4nu2n-zu3f-dmnn-kguv-6f643nei")
	u3 := MustParse("user-h4nu2n-zu3f-dmnn-kguv-6f643nei")  // Same UUID, different prefix
	u4 := MustParse("shell-aaaaaa-aaaa-aaaa-aaaa-aaaaaaaa") // Different UUID, same prefix

	if !u1.Equals(*u2) {
		t.Errorf("Equals() for identical XUIDs returned false")
	}

	if u1.Equals(*u3) {
		t.Errorf("Equals() for XUIDs with different prefixes returned true")
	}

	if u1.Equals(*u4) {
		t.Errorf("Equals() for XUIDs with different UUIDs returned true")
	}
}

func TestFromKey(t *testing.T) {
	key := "test-key"

	// FromKey should produce consistent results
	x1, err := FromKey(key)
	if err != nil {
		t.Errorf("FromKey() error = %v", err)
		return
	}

	x2, err := FromKey(key)
	if err != nil {
		t.Errorf("FromKey() error = %v", err)
		return
	}

	if !x1.Equals(*x2) {
		t.Errorf("FromKey() with same key produced different XUIDs")
	}

	// Check prefix is set correctly
	if x1.Prefix != "utref" {
		t.Errorf("FromKey() prefix = %q, want %q", x1.Prefix, "utref")
	}

	// Different keys should produce different XUIDs
	x3, _ := FromKey("different-key")
	if x1.Equals(*x3) {
		t.Errorf("FromKey() with different keys produced same XUID")
	}
}

func TestFromKeyPrefix(t *testing.T) {
	key := "test-key"
	prefix := "test"

	// FromKeyPrefix should produce consistent results
	x1, err := FromKeyPrefix(key, prefix)
	if err != nil {
		t.Errorf("FromKeyPrefix() error = %v", err)
		return
	}

	x2, err := FromKeyPrefix(key, prefix)
	if err != nil {
		t.Errorf("FromKeyPrefix() error = %v", err)
		return
	}

	if !x1.Equals(*x2) {
		t.Errorf("FromKeyPrefix() with same key/prefix produced different XUIDs")
	}

	// Check prefix is set correctly
	if x1.Prefix != prefix {
		t.Errorf("FromKeyPrefix() prefix = %q, want %q", x1.Prefix, prefix)
	}

	// Different keys should produce different XUIDs
	x3, _ := FromKeyPrefix("different-key", prefix)
	if x1.Equals(*x3) {
		t.Errorf("FromKeyPrefix() with different keys produced same XUID")
	}

	// Different prefixes should produce different XUIDs
	x4, _ := FromKeyPrefix(key, "diff")
	if x1.Equals(*x4) {
		t.Errorf("FromKeyPrefix() with different prefixes produced same XUID")
	}
}

func TestJSON(t *testing.T) {
	type TestStruct struct {
		ID XUID `json:"id"`
	}

	// Test marshaling
	t.Run("Marshal", func(t *testing.T) {
		original := TestStruct{
			ID: *MustParse("shell-h4nu2n-zu3f-dmnn-kguv-6f643nei"),
		}

		jsonData, err := json.Marshal(original)
		if err != nil {
			t.Errorf("json.Marshal() error = %v", err)
			return
		}

		expected := `{"id":"shell-h4nu2n-zu3f-dmnn-kguv-6f643nei"}`
		if string(jsonData) != expected {
			t.Errorf("json.Marshal() = %q, want %q", string(jsonData), expected)
		}
	})

	// Test unmarshaling
	t.Run("Unmarshal", func(t *testing.T) {
		jsonData := []byte(`{"id":"shell-h4nu2n-zu3f-dmnn-kguv-6f643nei"}`)

		var decoded TestStruct
		err := json.Unmarshal(jsonData, &decoded)
		if err != nil {
			t.Errorf("json.Unmarshal() error = %v", err)
			return
		}

		expected := "shell-h4nu2n-zu3f-dmnn-kguv-6f643nei"
		if decoded.ID.String() != expected {
			t.Errorf("Unmarshaled ID = %q, want %q", decoded.ID.String(), expected)
		}

		expectedUUID := "3f1b4d37-34d9-46c6-b546-a57c5f736d22"
		if decoded.ID.ToUUID() != expectedUUID {
			t.Errorf("Unmarshaled UUID = %q, want %q", decoded.ID.ToUUID(), expectedUUID)
		}
	})

	// Test unmarshaling with invalid data
	t.Run("Unmarshal invalid", func(t *testing.T) {
		jsonData := []byte(`{"id":"invalid-format"}`)

		var decoded TestStruct
		err := json.Unmarshal(jsonData, &decoded)
		if err == nil {
			t.Errorf("json.Unmarshal() with invalid data did not return an error")
		}
	})
}

func TestScan(t *testing.T) {
	t.Run("Scan from string", func(t *testing.T) {
		var xuid XUID
		err := xuid.Scan("shell-h4nu2n-zu3f-dmnn-kguv-6f643nei")
		if err != nil {
			t.Errorf("Scan() from string error = %v", err)
			return
		}

		if xuid.Prefix != "shell" {
			t.Errorf("Scan() prefix = %q, want %q", xuid.Prefix, "shell")
		}

		if xuid.ToUUID() != "3f1b4d37-34d9-46c6-b546-a57c5f736d22" {
			t.Errorf("Scan() UUID = %q, want %q", xuid.ToUUID(), "3f1b4d37-34d9-46c6-b546-a57c5f736d22")
		}
	})

	t.Run("Scan from sql.RawBytes", func(t *testing.T) {
		var xuid XUID
		rawBytes := sql.RawBytes("shell-h4nu2n-zu3f-dmnn-kguv-6f643nei")
		err := xuid.Scan(rawBytes)
		if err != nil {
			t.Errorf("Scan() from sql.RawBytes error = %v", err)
			return
		}

		if xuid.Prefix != "shell" {
			t.Errorf("Scan() prefix = %q, want %q", xuid.Prefix, "shell")
		}

		if xuid.ToUUID() != "3f1b4d37-34d9-46c6-b546-a57c5f736d22" {
			t.Errorf("Scan() UUID = %q, want %q", xuid.ToUUID(), "3f1b4d37-34d9-46c6-b546-a57c5f736d22")
		}
	})

	t.Run("Scan from string with invalid XUID", func(t *testing.T) {
		var xuid XUID
		err := xuid.Scan("not-a-valid-xuid")
		if err == nil {
			t.Errorf("Scan() from invalid XUID string did not return an error")
		}
	})

	t.Run("Scan from invalid type", func(t *testing.T) {
		var xuid XUID
		err := xuid.Scan(123) // not a supported type
		if err == nil {
			t.Errorf("Scan() from invalid type did not return an error")
		}
	})
}

func TestValue(t *testing.T) {
	xuid := MustParse("shell-h4nu2n-zu3f-dmnn-kguv-6f643nei")

	val, err := xuid.Value()
	if err != nil {
		t.Errorf("Value() error = %v", err)
		return
	}

	// Check that Value returns a string
	str, ok := val.(string)
	if !ok {
		t.Errorf("Value() did not return a string, got %T", val)
		return
	}

	// Check the string value
	expected := "shell-h4nu2n-zu3f-dmnn-kguv-6f643nei"
	if str != expected {
		t.Errorf("Value() = %q, want %q", str, expected)
	}

	// Check value implements driver.Valuer
	var _ driver.Valuer = xuid
}

func TestMustParseUUID(t *testing.T) {
	t.Run("Valid UUID", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("MustParseUUID() panicked unexpectedly: %v", r)
			}
		}()

		uid := MustParseUUID("3f1b4d37-34d9-46c6-b546-a57c5f736d22", "shell")
		if uid.ToUUID() != "3f1b4d37-34d9-46c6-b546-a57c5f736d22" {
			t.Errorf("MustParseUUID() returned incorrect UUID, got %s", uid.ToUUID())
		}
		if uid.Prefix != "shell" {
			t.Errorf("MustParseUUID() returned incorrect prefix, got %s, want %s", uid.Prefix, "shell")
		}
		if uid.String() != "shell-h4nu2n-zu3f-dmnn-kguv-6f643nei" {
			t.Errorf("MustParseUUID() returned incorrect string, got %s", uid.String())
		}
	})

	t.Run("Invalid UUID", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("MustParseUUID() did not panic for invalid UUID")
			}
		}()

		_ = MustParseUUID("invalid-uuid", "test")
	})
}

func TestStringWithoutPrefix(t *testing.T) {
	// Test XUID string representation with an empty prefix
	// Create a known UUID that we can test
	u := uuid.MustParse("3f1b4d37-34d9-46c6-b546-a57c5f736d22")
	
	// Create a XUID without a prefix
	xuid, err := FromUUID(u, "")
	if err != nil {
		t.Errorf("FromUUID() error = %v", err)
		return
	}
	
	// Get the string representation
	result := xuid.String()
	
	// The string representation without a prefix should match our example from the README
	// but without the prefix and dash
	if !strings.HasPrefix(result, "h4nu2n-") {
		t.Errorf("String() with empty prefix does not start with h4nu2n-, got %s", result)
	}
	
	// Verify that we can convert it back to a UUID
	uuidStr := xuid.ToUUID()
	if uuidStr != "3f1b4d37-34d9-46c6-b546-a57c5f736d22" {
		t.Errorf("ToUUID() = %q, want %q", uuidStr, "3f1b4d37-34d9-46c6-b546-a57c5f736d22")
	}
	
	// Verify the prefix is empty
	if xuid.Prefix != "" {
		t.Errorf("Prefix = %q, want empty string", xuid.Prefix)
	}
}