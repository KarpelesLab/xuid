// Package xuid provides an extended UUID implementation that adds a type prefix
// while maintaining compatibility with standard UUIDs.
//
// XUIDs consist of a standard UUID with an optional type prefix (up to 5 characters),
// encoded in base32 rather than base16 for better readability and shorter string representation.
// They are designed to be used as identifiers with built-in type information.
package xuid

import (
	"encoding/base32"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

// b32enc is the base32 encoder used for XUID string representation
// It uses standard encoding without padding to create compact representations
var b32enc = base32.StdEncoding.WithPadding(base32.NoPadding)

// String formats the XUID as a string with the format:
// prefix-aaaaaa-aaaa-aaaa-aaaa-aaaaaaaa
// where 'prefix' is the type prefix (up to 5 characters) and the remaining parts
// are the base32-encoded UUID with hyphens for readability.
// The result is always lowercase and contains no padding characters.
func (x XUID) String() string {
	var dst [26]byte
	// convert to base32
	b32enc.Encode(dst[:], x.UUID[:])

	var final [36]byte
	b := final[:]

	// Format: prefx-aaaaaa-aaaa-aaaa-aaaa-aaaaaaaa
	// Lengths: 5-6-4-4-4-8

	var pfxLn int
	if x.Prefix != "" {
		// Copy prefix (up to 5 chars) and add separator
		pfxLn = copy(b[:5], x.Prefix)
		b = b[pfxLn:]
		b[0] = '-'
		b = b[1:]
	}

	// Format the base32 encoded UUID with hyphens in the same positions as a regular UUID
	copy(b[:6], dst[:6])      // First 6 chars
	b[6] = '-'
	copy(b[7:11], dst[6:10])  // Next 4 chars
	b[11] = '-'
	copy(b[12:16], dst[10:14]) // Next 4 chars
	b[16] = '-'
	copy(b[17:21], dst[14:18]) // Next 4 chars
	b[21] = '-'
	copy(b[22:], dst[18:])     // Final 8 chars

	// Return the string representation, ensuring it's lowercase
	return strings.ToLower(string(final[:31+pfxLn]))
}

// Equals compares two XUIDs and returns true if they are equal, 
// meaning they have the same prefix and UUID.
func (x XUID) Equals(y XUID) bool {
	return x == y
}

// Parse parses a string representation and returns the resulting XUID.
// It can handle XUID formatted strings as well as standard UUIDs.
// For XUIDs, it supports both prefixed and non-prefixed formats.
//
// If the input string doesn't conform to XUID format, Parse will attempt
// to interpret it as a standard UUID and assign an empty prefix.
//
// Returns the parsed XUID and any error encountered.
func Parse(s string) (*XUID, error) {
	// XUID length can be:
	// - 30 bytes (no prefix, base32 with hyphens)
	// - 32-36 bytes (prefix length 1-5 + hyphen + 30 byte base32 representation)
	l := len(s)
	var pfx, v string

	switch l {
	case 30: // No prefix
		v = s
	case 32, 33, 34, 35, 36: // With prefix (length 1-5)
		pfxLn := l - 31
		if s[pfxLn] != '-' {
			// If there's no hyphen after the prefix, fallback to UUID parsing
			return ParseUUID(s, "")
		}
		pfx = s[:pfxLn]
		v = s[pfxLn+1:]
	default:
		// Invalid length for XUID, try parsing as UUID
		return ParseUUID(s, "")
	}

	// Validate the XUID format with hyphens at correct positions
	// A valid XUID body (v) has format: aaaaaa-aaaa-aaaa-aaaa-aaaaaaaa
	if v[6] != '-' || v[11] != '-' || v[16] != '-' || v[21] != '-' {
		return ParseUUID(s, "")
	}

	// Extract the parts without hyphens for base32 decoding
	parts := []string{
		v[0:6],
		v[7:11],
		v[12:16],
		v[17:21],
		v[22:],
	}
	var data uuid.UUID
	// Decode the base32 representation back to UUID bytes
	_, err := b32enc.Decode(data[:], []byte(strings.ToUpper(strings.Join(parts, ""))))
	if err != nil {
		return nil, err
	}

	return &XUID{Prefix: pfx, UUID: data}, nil
}

// MustParse parses a string into a XUID, panicking if parsing fails.
// This is useful for constants or when the XUID string is known to be valid.
func MustParse(s string) *XUID {
	return Must(Parse(s))
}

// ParsePrefix parses a XUID string and verifies that the prefix matches the provided one.
// Returns a parsed XUID if successful, or an error if either:
// - The input cannot be parsed as a XUID
// - The prefix doesn't match the expected value
//
// This is useful for type-safe XUID parsing where the caller expects a specific entity type.
func ParsePrefix(s, prefix string) (*XUID, error) {
	v, err := Parse(s)
	if err != nil {
		return nil, err
	}
	if v.Prefix != prefix {
		return nil, fmt.Errorf("%w, expected prefix %s", ErrBadPrefix, prefix)
	}
	return v, nil
}
