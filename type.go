package xuid

import "github.com/google/uuid"

type XUID struct {
	Prefix string
	UUID   uuid.UUID
}
