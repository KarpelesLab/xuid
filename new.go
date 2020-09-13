package xuid

import "github.com/google/uuid"

func New(prefix string) *XUID {
	return Must(NewRandom(prefix))
}

func Must(x *XUID, err error) *XUID {
	if err != nil {
		panic(err)
	}
	return x
}

func FromUUID(u uuid.UUID, prefix string) (*XUID, error) {
	return &XUID{Prefix: prefix, UUID: u}, nil
}

func NewRandom(prefix string) (*XUID, error) {
	u, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	return FromUUID(u, prefix)
}
