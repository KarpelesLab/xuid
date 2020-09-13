package xuid

import (
	"encoding/base32"
	"strings"

	"github.com/google/uuid"
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

	var pfxLn int
	if x.Prefix != "" {
		pfxLn = copy(b[:5], x.Prefix)
		b = b[pfxLn:]
		b[0] = '-'
		b = b[1:]
	}

	copy(b[:6], dst[:6]) // 6
	b[6] = '-'
	copy(b[7:11], dst[6:10]) // 4
	b[11] = '-'
	copy(b[12:16], dst[10:14]) // 4
	b[16] = '-'
	copy(b[17:21], dst[14:18]) // 4
	b[21] = '-'
	copy(b[22:], dst[18:]) // 8

	return strings.ToLower(string(final[:31+pfxLn]))
}

func Parse(s string) (*XUID, error) {
	// parse can handle many type of formats, and will fallback to uuid parsing if failing
	// a xuid length can be 30bytes (no prefix), or 32~36 (prefix len 1 to 5)
	l := len(s)
	var pfx, v string

	switch l {
	case 30:
		v = s
	case 32, 33, 34, 35, 36:
		pfxLn := l - 31
		if s[pfxLn] != '-' {
			// fallback
			return ParseUUID(s, "")
		}
		pfx = s[:pfxLn]
		v = s[pfxLn+1:]
	default:
		// fallback
		return ParseUUID(s, "")
	}

	// v has a len of 30, and should be in format h4nu2n-zu3f-dmnn-kguv-6f643nei
	// positions of -: 6, 11, 16, 21
	if v[6] != '-' || v[11] != '-' || v[16] != '-' || v[21] != '-' {
		return ParseUUID(s, "")
	}

	parts := []string{
		v[0:6],
		v[7:11],
		v[12:16],
		v[17:21],
		v[22:],
	}
	var data uuid.UUID
	_, err := b32enc.Decode(data[:], []byte(strings.ToUpper(strings.Join(parts, ""))))
	if err != nil {
		return nil, err
	}

	return &XUID{Prefix: pfx, UUID: data}, nil
}

func MustParse(s string) *XUID {
	return Must(Parse(s))
}
