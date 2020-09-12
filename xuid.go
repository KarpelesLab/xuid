package xuid

import (
	"encoding/base32"
	"strings"
)

var b32enc = base32.StdEncoding.WithPadding(base32.NoPadding)

func (x *XUID) String() string {
	var dst [26]byte
	// convert to base32
	b32enc.Encode(dst[:], x.UUID[:])

	var final [36]byte
	b := final[:]

	// prefx-aaaaaa-aaaa-aaaa-aaaa-aaaaaaaa
	// 5-6-4-4-4-8

	pfxLn := copy(b[:5], x.Prefix)
	b = b[pfxLn:]
	b[0] = '-'
	b = b[1:]

	copy(b[:6], dst[:6]) // 6
	b[7] = '-'
	copy(b[7:11], dst[6:10]) // 4
	b[11] = '-'
	copy(b[12:16], dst[10:14]) // 4
	b[16] = '-'
	copy(b[17:21], dst[14:18]) // 4
	b[21] = '-'
	copy(b[22:], dst[18:]) // 8

	return strings.ToLower(string(final[:31+pfxLn]))
}
