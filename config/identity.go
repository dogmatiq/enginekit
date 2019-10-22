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

// Validate returns an error if i is not a valid identity.
func (i Identity) Validate() error {
	if !isValidName(i.Name) {
		return fmt.Errorf(
			"invalid name %#v, names must be non-empty, printable UTF-8 strings with no whitespace",
			i.Name,
		)
	}

	if !isValidName(i.Key) {
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

// isValidName returns true if n is a valid application or handler name or key.
//
// A valid name/key is a non-empty string consisting of Unicode printable
// characters, except whitespace.
func isValidName(n string) bool {
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
