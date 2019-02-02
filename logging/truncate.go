package logging

import (
	"strings"
	"unicode/utf8"
)

// truncate trims a UTF-8 string to a specific length.
//
// tr is the length to which IDs are truncated. If tr positive the tail of the
// ID is truncated. If negative, the head of the ID is truncated. If tr is zero
// no truncation is performed
func truncate(s string, tr int) string {
	if tr == 0 {
		return s
	}

	abs := tr
	if tr < 0 {
		abs = -abs
	}

	runes := utf8.RuneCountInString(s)

	if runes == abs {
		return s
	}

	if runes < abs {
		return strings.Repeat(" ", abs-runes) + s
	}

	stop := tr
	if tr < 0 {
		stop = runes - abs
	}

	n := 0

	for i, r := range s {
		if i == stop {
			break
		}

		n += utf8.RuneLen(r)
	}

	if tr < 0 {
		return s[n:]
	}

	return s[:n]
}
