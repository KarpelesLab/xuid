package xuid

import "encoding/json"

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

func (x XUID) MarshalJSON() ([]byte, error) {
	return json.Marshal(x.String())
}
