package xuid

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
)

// Scan sets the value of the XUID to the passed data, parsing it as needed.
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

// Value returns the xuid as a string value
func (x XUID) Value() (driver.Value, error) {
	return x.String(), nil
}
