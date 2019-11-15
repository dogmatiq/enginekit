package config

import (
	"fmt"
	"unicode"
)

// Identity contains the application-defined identity of an application or its
// handlers.
type Identity struct {
	Name string
	Key  string
}

// NewIdentity returns a new identity.
func NewIdentity(n, k string) (Identity, error) {
	i := Identity{n, k}
	return i, i.Validate()
}

// MustNewIdentity returns a new identity, or panics if the given name or key are invalid.
func MustNewIdentity(n, k string) Identity {
	i, err := NewIdentity(n, k)
	if err != nil {
		panic(err)
	}

	return i
}

// IsZero returns true if the identity is the zero-value.
func (i Identity) IsZero() bool {
	return i.Name == "" && i.Key == ""
}

// Validate returns an error if i is not a valid identity.
func (i Identity) Validate() error {
	if !isValid(i.Name) {
		return fmt.Errorf(
			"invalid name %#v, names must be non-empty, printable UTF-8 strings with no whitespace",
			i.Name,
		)
	}

	if !isValid(i.Key) {
		return fmt.Errorf(
			"invalid key %#v, keys must be non-empty, printable UTF-8 strings with no whitespace",
			i.Key,
		)
	}

	return nil
}

func (i Identity) String() string {
	return fmt.Sprintf("%s (%s)", i.Name, i.Key)
}

// isValid returns true if n is a valid application or handler name or key.
//
// A valid name/key is a non-empty string consisting of Unicode printable
// characters, except whitespace.
func isValid(n string) bool {
	if len(n) == 0 {
		return false
	}

	for _, r := range n {
		if unicode.IsSpace(r) || !unicode.IsPrint(r) {
			return false
		}
	}

	return true
}
