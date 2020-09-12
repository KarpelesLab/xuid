package xuid

type XUID struct {
	Prefix string
	UUID   [16]byte
}
