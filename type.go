package xuid

import "github.com/google/uuid"

// XUID represents an extended UUID with a type prefix.
// It consists of a standard UUID and an optional prefix string
// that identifies the type of entity this ID represents.
//
// The prefix is limited to 5 characters maximum for string representation.
// When encoded to string, XUIDs use base32 encoding rather than base16 (hex)
// to produce more compact strings while maintaining the same information.
type XUID struct {
	// Prefix is a string of up to 5 characters identifying the type of entity
	Prefix string
	
	// UUID is the standard universally unique identifier that provides
	// the uniqueness guarantee
	UUID uuid.UUID
}
