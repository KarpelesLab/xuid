package xuid

import "github.com/google/uuid"

func (x *XUID) ToUUID() string {
	return x.UUID.String()
}

func ParseUUID(inputUuid, prefix string) (*XUID, error) {
	u, err := uuid.Parse(inputUuid)
	if err != nil {
		return nil, err
	}

	return FromUUID(u, prefix)
}

func MustParseUUID(inputUuid, prefix string) *XUID {
	return Must(FromUUID(uuid.MustParse(inputUuid), prefix))
}
