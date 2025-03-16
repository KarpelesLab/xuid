package xuid

import "github.com/google/uuid"

// refNs is a reference namespace UUID used for generating deterministic XUIDs
// from keys using the SHA-1 hash algorithm
var refNs = uuid.MustParse("d16b6139-8989-467f-a240-441df6734f45")

// New creates a new random XUID with the given prefix.
// It's a shorthand for Must(NewRandom(prefix)) and will panic if the random 
// generator fails for some reason.
//
// This is the most common method for creating new XUIDs.
func New(prefix string) *XUID {
	return Must(NewRandom(prefix))
}

// Must is a generic helper that transforms any error from the called function
// into a panic. It's useful for operations that shouldn't fail under normal
// circumstances, or when failure should be fatal.
//
// If err is nil, it returns the value x. Otherwise, it panics with the error.
func Must[T any](x T, err error) T {
	if err != nil {
		panic(err)
	}
	return x
}

// FromUUID creates a new XUID from an existing UUID and a prefix.
// This is useful when you want to convert a standard UUID to a XUID
// with type information.
func FromUUID(u uuid.UUID, prefix string) (*XUID, error) {
	return &XUID{Prefix: prefix, UUID: u}, nil
}

// NewRandom generates a new random XUID with the given prefix.
// It uses the underlying uuid.NewRandom() function to generate
// a version 4 (random) UUID and assigns the provided prefix.
func NewRandom(prefix string) (*XUID, error) {
	u, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	return FromUUID(u, prefix)
}

// FromKey generates a deterministic XUID based on the provided key string.
// It always produces the same XUID for the same key value.
//
// The generated XUID always has the prefix "utref" (utility reference).
// This is useful for creating reference objects with predictable IDs.
func FromKey(key string) (*XUID, error) {
	return FromUUID(uuid.NewSHA1(refNs, []byte(key)), "utref")
}

// FromKeyPrefix generates a deterministic XUID based on both a key and a prefix.
// It always produces the same XUID for the same combination of key and prefix.
//
// This is useful for objects that need consistent IDs across different environments.
// The prefix is used both for generating the UUID and as the type prefix in the XUID.
func FromKeyPrefix(key, prefix string) (*XUID, error) {
	// Create a unique namespace for this prefix
	subRefNs := uuid.NewSHA1(refNs, []byte(prefix))
	// Generate a SHA-1 UUID using the key in the prefix-specific namespace
	return FromUUID(uuid.NewSHA1(subRefNs, []byte(key)), prefix)
}
