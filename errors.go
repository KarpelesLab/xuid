package xuid

import "errors"

// Package errors that can be returned by the XUID library functions
var (
	// ErrBadPrefix is returned when a XUID's prefix doesn't match the expected value
	// This is typically used by ParsePrefix to validate that an ID belongs to a specific entity type
	ErrBadPrefix = errors.New("xuid: bad prefix")
)
