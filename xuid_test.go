package xuid

import "testing"

func TestXUID(t *testing.T) {
	uid := MustParse("null-aaaaaa-aaaa-aaaa-aaaa-aaaaaaaa")

	if uid.ToUUID() != "00000000-0000-0000-0000-000000000000" {
		t.Errorf("invalid uuid, got %s", uid.ToUUID())
	}
}
