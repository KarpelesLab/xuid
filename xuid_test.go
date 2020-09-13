package xuid

import "testing"

func TestXUID(t *testing.T) {
	uid := MustParse("null-aaaaaa-aaaa-aaaa-aaaa-aaaaaaaa")

	if uid.ToUUID() != "00000000-0000-0000-0000-000000000000" {
		t.Errorf("invalid uuid, got %s", uid.ToUUID())
	}

	uid = MustParse("shell-h4nu2n-zu3f-dmnn-kguv-6f643nei")

	if uid.ToUUID() != "3f1b4d37-34d9-46c6-b546-a57c5f736d22" {
		t.Errorf("invalid uuid, got %s", uid.ToUUID())
	}

	if uid.String() != "shell-h4nu2n-zu3f-dmnn-kguv-6f643nei" {
		t.Errorf("invalid xuid, got %s instead of shell-h4nu2n-zu3f-dmnn-kguv-6f643nei", uid.String())
	}
}
