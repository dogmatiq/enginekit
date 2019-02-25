package config

import "unicode"

// IsValidName returns true if n is a valid application or handler name.
func IsValidName(n string) bool {
	if len(n) == 0 {
		return false
	}

	for _, r := range n {
		if !unicode.IsPrint(r) {
			return false
		}
	}

	return true
}
