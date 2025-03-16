package xuid

import "github.com/google/uuid"

// ToUUID returns the string representation of the underlying UUID
// in standard UUID format (xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx).
// This is useful when you need the standard UUID representation
// rather than the XUID format.
func (x *XUID) ToUUID() string {
	return x.UUID.String()
}

// ParseUUID converts a standard UUID string into a XUID.
// It takes a UUID-formatted string and an optional prefix to assign to the XUID.
// If the prefix is empty, the XUID will have no type information.
//
// This is useful for converting existing UUIDs to XUIDs or
// for working with existing systems that use standard UUIDs.
//
// Returns a XUID pointer and any error encountered during parsing.
func ParseUUID(inputUuid, prefix string) (*XUID, error) {
	u, err := uuid.Parse(inputUuid)
	if err != nil {
		return nil, err
	}

	return FromUUID(u, prefix)
}

// MustParseUUID works like ParseUUID but panics if parsing fails.
// This is useful for constants or when the UUID string is known to be valid.
//
// Takes a UUID-formatted string and an optional prefix to assign to the XUID.
func MustParseUUID(inputUuid, prefix string) *XUID {
	return Must(FromUUID(uuid.MustParse(inputUuid), prefix))
}
