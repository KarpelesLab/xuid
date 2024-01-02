package xuid

import "github.com/google/uuid"

var refNs = uuid.MustParse("d16b6139-8989-467f-a240-441df6734f45")

// New is a shorthand for Must(NewRandom(prefix)) and will always return a new
// xuid with the given prefix, or panic if for some reason the random generator
// does not work.
func New(prefix string) *XUID {
	return Must(NewRandom(prefix))
}

// Must can be used to transform any error thrown by the called method into
// a panic, and is useful when we know a call cannot fail, or that failure
// should be fatal.
func Must(x *XUID, err error) *XUID {
	if err != nil {
		panic(err)
	}
	return x
}

// FromUUID returns a new xuid with the passed uuid and chosen prefix
func FromUUID(u uuid.UUID, prefix string) (*XUID, error) {
	return &XUID{Prefix: prefix, UUID: u}, nil
}

// NewRandom returns a new random xuid with the given prefix set
func NewRandom(prefix string) (*XUID, error) {
	u, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	return FromUUID(u, prefix)
}

// FromKey returns a fixed ID based on a given key that can be guaranteed
// to always have the same ID as long as the key is the same value. This can
// be used for a Util/Ref object that have fixed IDs for a given key.
func FromKey(key string) (*XUID, error) {
	return FromUUID(uuid.NewSHA1(refNs, []byte(key)), "utref")
}

// FromKeyPrefix returns a fixed ID based on the key and prefix passed and
// is guaranteed to always return the same ID as long as the parameters are
// the same. This can be used for objects that need to keep the same ID no
// matter on which environment this runs.
func FromKeyPrefix(key, prefix string) (*XUID, error) {
	subRefNs := uuid.NewSHA1(refNs, []byte(prefix))
	return FromUUID(uuid.NewSHA1(subRefNs, []byte(key)), prefix)
}
