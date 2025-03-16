[![GoDoc](https://godoc.org/github.com/KarpelesLab/xuid?status.svg)](https://godoc.org/github.com/KarpelesLab/xuid)

# XUID

XUID (eXtended Unique IDentifier) is a Go library that enhances UUIDs with a type prefix while maintaining compatibility with standard UUID bit size. XUIDs are encoded in base32 rather than base16, ensuring they never exceed 36 characters in length while providing better human readability.

## Features

- Fully compatible with standard UUIDs (same number of bits)
- Adds descriptive type prefixes (up to 5 characters)
- Encodes in base32 for shorter, more readable strings
- Case-insensitive matching
- Compatible with SQL databases via `sql.Scanner` and `driver.Valuer` interfaces
- JSON marshaling/unmarshaling support
- Easily convertible to and from standard UUIDs

## Installation

```bash
go get -u github.com/KarpelesLab/xuid
```

## Usage

```go
package main

import (
    "fmt"
    "github.com/KarpelesLab/xuid"
)

func main() {
    // Create a new random XUID with a prefix
    id := xuid.New("user")
    fmt.Println(id) // e.g., user-h4nu2n-zu3f-dmnn-kguv-6f643nei

    // Parse an existing XUID
    parsed, err := xuid.Parse("user-h4nu2n-zu3f-dmnn-kguv-6f643nei")
    if err != nil {
        panic(err)
    }
    
    // Convert XUID to UUID
    uuidStr := parsed.ToUUID()
    fmt.Println(uuidStr) // e.g., 3f1b4d37-34d9-46c6-b546-a57c5f736d22
    
    // Create a XUID from a UUID
    fromUUID, _ := xuid.ParseUUID(uuidStr, "user")
    fmt.Println(fromUUID) // user-h4nu2n-zu3f-dmnn-kguv-6f643nei
    
    // Create deterministic XUIDs based on keys
    keyID, _ := xuid.FromKeyPrefix("specific-resource-name", "res")
    fmt.Println(keyID) // Will always be the same for this key and prefix
}
```

## Type Prefixes

The type prefix (up to 5 characters) identifies what kind of object the ID represents, making it easier for both humans and automated systems to quickly identify the entity type without additional lookups.

Common prefix examples might include:
- `user` - User accounts
- `doc` - Documents
- `img` - Images
- `prod` - Products
- `ord` - Orders

## Format

XUIDs use the following format:
```
prefix-aaaaaa-aaaa-aaaa-aaaa-aaaaaaaa
```

Where:
- `prefix` is 1-5 characters identifying the entity type
- The remaining parts are the base32-encoded UUID with hyphens for readability

## Examples

* UUID `3f1b4d37-34d9-46c6-b546-a57c5f736d22` with type `shell` becomes `shell-h4nu2n-zu3f-dmnn-kguv-6f643nei`
* UUID `00000000-0000-0000-0000-000000000000` with type `null` becomes `null-aaaaaa-aaaa-aaaa-aaaa-aaaaaaaa`

## License

See [LICENSE](LICENSE) file.