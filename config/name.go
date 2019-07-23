package config

import "unicode"

// IsValidName returns true if n is a valid application or handler name.
//
// A valid name is a non-empty string consisting of Unicode printable
// characters, except whitespace.
func IsValidName(n string) bool {
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

// IsValidKey returns true if n is a valid application or handler key.
//
// A valid key is a non-empty string consisting of Unicode printable characters,
// except whitespace.
func IsValidKey(k string) bool {
	return IsValidName(k)
}
