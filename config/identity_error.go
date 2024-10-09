package config

import "fmt"

// InvalidIdentityNameError indicates that the "name" component of an [Identity]
// is invalid.
type InvalidIdentityNameError struct {
	InvalidName string
}

func (e InvalidIdentityNameError) Error() string {
	return fmt.Sprintf("invalid identity name (%q): names must be non-empty, printable UTF-8 strings with no whitespace", e.InvalidName)
}

// InvalidIdentityKeyError indicates that the "key" component of an [Identity]
// is invalid.
type InvalidIdentityKeyError struct {
	InvalidKey string
	Cause      error
}

func (e InvalidIdentityKeyError) Error() string {
	return fmt.Sprintf("invalid identity key (%q): keys must be RFC 4122/9562 UUIDs: %s", e.InvalidKey, e.Cause)
}

func (e InvalidIdentityKeyError) Unwrap() error {
	return e.Cause
}
