package config

import (
	"errors"
	"fmt"
	"strconv"
	"unicode"

	"github.com/dogmatiq/enginekit/protobuf/uuidpb"
)

// Identity represents the (potentially invalid) identity of an entity.
type Identity struct {
	Name string
	Key  string
}

// Validate returns an error if the identity is invalid.
func (i Identity) Validate(options ...ValidationOption) error {
	return i.validate(newValidationOptions(options))
}

func (i Identity) validate(validationOptions) error {
	var problems error

	if !isValidIdentityName(i.Name) {
		problems = errors.Join(problems, InvalidIdentityNameError{i.Name})
	}

	if _, err := uuidpb.Parse(i.Key); err != nil {
		problems = errors.Join(problems, InvalidIdentityKeyError{i.Key, err})
	}

	return problems
}

// Normalize returns a normalized copy of the identity, or an error if the
// identity is invalid.
func (i Identity) Normalize(options ...ValidationOption) (Identity, error) {
	return i.normalize(newValidationOptions(options))
}

func (i Identity) normalize(opts validationOptions) (Identity, error) {
	if err := i.validate(opts); err != nil {
		return Identity{}, err
	}

	return Identity{
		Name: i.Name,
		Key:  uuidpb.MustParse(i.Key).AsString(),
	}, nil
}

func (i Identity) String() string {
	name := "?"
	if i.Name != "" {
		if isValidIdentityName(i.Name) {
			name = i.Name
		} else {
			name = strconv.Quote(i.Name)
		}
	}

	key := "?"
	if i.Key != "" {
		if normalized, err := uuidpb.Parse(i.Key); err == nil {
			key = normalized.AsString()
		} else {
			key = strconv.Quote(i.Key)
		}
	}

	return name + "/" + key
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

// InvalidIdentityNameError is an error that occurs when an identity name is
// invalid.
type InvalidIdentityNameError struct {
	InvalidName string
}

func (e InvalidIdentityNameError) Error() string {
	return fmt.Sprintf("invalid identity name (%q): names must be non-empty, printable UTF-8 strings with no whitespace", e.InvalidName)
}

// InvalidIdentityKeyError is an error that occurs when an identity key is
// invalid.
type InvalidIdentityKeyError struct {
	InvalidKey string
	ParseError error
}

func (e InvalidIdentityKeyError) Unwrap() error {
	return e.ParseError
}

func (e InvalidIdentityKeyError) Error() string {
	return fmt.Sprintf("invalid identity key (%q): keys must be RFC 4122/9562 UUIDs: %s", e.InvalidKey, e.ParseError)
}
