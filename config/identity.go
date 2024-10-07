package config

import (
	"fmt"
	"unicode"

	"github.com/dogmatiq/enginekit/protobuf/uuidpb"
)

// Identity represents the (potentially invalid) identity of an entity.
type Identity struct {
	Name string
	Key  string
}

// Errors returns a list of errors that were encountered while building the
// identity.
func (i Identity) Errors() []error {
	var errors []error

	if !isValidIdentityName(i.Name) {
		errors = append(
			errors,
			fmt.Errorf("invalid identity name %q: names must be non-empty, printable UTF-8 strings with no whitespace", i.Name),
		)
	}

	if _, err := uuidpb.Parse(i.Key); err != nil {
		errors = append(
			errors,
			fmt.Errorf("invalid identity key %q: keys must be RFC 4122/9562 UUIDs: %w", i.Key, err),
		)
	}

	return errors
}

// isValidIdentityName returns true if n is a valid application or handler name.
//
// A valid name is a non-empty string consisting of Unicode printable
// characters, except whitespace.
func isValidIdentityName(n string) bool {
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
