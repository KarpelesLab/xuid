package xuid

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
)

// Scan implements the sql.Scanner interface for XUID.
// It allows XUIDs to be scanned directly from database query results.
//
// The scan supports the following types:
// - string: Parses the string as a XUID
// - sql.RawBytes: Converts to string and parses as a XUID
//
// This functionality enables XUIDs to be used directly with database/sql
// operations without manual type conversion.
func (x *XUID) Scan(value any) error {
	switch v := value.(type) {
	case string:
		nv, err := Parse(v)
		if err != nil {
			return err
		}
		x.Prefix = nv.Prefix
		x.UUID = nv.UUID
		return nil
	case sql.RawBytes:
		nv, err := Parse(string(v))
		if err != nil {
			return err
		}
		x.Prefix = nv.Prefix
		x.UUID = nv.UUID
		return nil
	default:
		return fmt.Errorf("Scan type %T unsupported to store into XUID", v)
	}
}

// Value implements the driver.Valuer interface for XUID.
// It returns the string representation of the XUID, which can be
// directly stored in database fields.
//
// This allows XUIDs to be used directly in database operations
// without manual type conversion.
func (x XUID) Value() (driver.Value, error) {
	return x.String(), nil
}
