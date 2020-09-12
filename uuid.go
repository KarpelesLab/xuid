package xuid

import "encoding/hex"

func (x *XUID) ToUUID() string {
	// encode as hex format 00000000-0000-0000-0000-000000000000
	var dst [36]byte
	hex.Encode(dst[:8], x.UUID[:4])
	dst[8] = '-'
	hex.Encode(dst[9:13], x.UUID[4:6])
	dst[13] = '-'
	hex.Encode(dst[14:18], x.UUID[6:8])
	dst[18] = '-'
	hex.Encode(dst[19:23], x.UUID[8:10])
	dst[23] = '-'
	hex.Encode(dst[24:], x.UUID[10:])

	return string(dst[:])
}
