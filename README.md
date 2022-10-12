[![GoDoc](https://godoc.org/github.com/KarpelesLab/xuid?status.svg)](https://godoc.org/github.com/KarpelesLab/xuid)

# XUID

XUID is an ID that contains the exact same number of bits as a UUID, adds a type prefix, and is encoded in a string that will not be larger than 36 characters.

This is achieved by using base32 instead of base16. Base32 still only uses one set of alphabetic characters, allowing caseless matches.

## Type

The type is a up to 5 characters long code that specifies what kind of object the ID is of. It is used to allow easy identification of an ID, either by humans or automated code.

## Examples

* UUID `3f1b4d37-34d9-46c6-b546-a57c5f736d22` with type `shell` becomes `shell-h4nu2n-zu3f-dmnn-kguv-6f643nei`
* UUID `00000000-0000-0000-0000-000000000000` with type `null` becomes `null-aaaaaa-aaaa-aaaa-aaaa-aaaaaaaa`
