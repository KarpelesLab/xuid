package xuid

import "encoding/json"

// UnmarshalJSON implements the json.Unmarshaler interface for XUID.
// It allows XUIDs to be directly unmarshaled from JSON string representations.
//
// The JSON value should be a string in valid XUID format or convertible
// to a valid XUID format through the Parse function.
func (x *XUID) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	nv, err := Parse(s)
	if err != nil {
		return err
	}
	x.Prefix = nv.Prefix
	x.UUID = nv.UUID
	return nil
}

// MarshalJSON implements the json.Marshaler interface for XUID.
// It converts the XUID to its string representation (prefix-aaaaaa-aaaa-aaaa-aaaa-aaaaaaaa)
// and then marshals that string to JSON.
//
// This makes XUIDs appear as strings in JSON output, which is more readable
// and works better with other systems that expect string IDs.
func (x XUID) MarshalJSON() ([]byte, error) {
	return json.Marshal(x.String())
}
