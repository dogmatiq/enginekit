package fmtbackport

import (
	"fmt"
	"strconv"
	"unicode/utf8"
)

// FormatString returns a string representing the fully qualified formatting
// directive captured by the State, followed by the argument verb.
//
// It is a copy of the fmt.FormatString() function in Go v1.20, and may be
// removed once Go v1.19 support is dropped.
func FormatString(state fmt.State, verb rune) string {
	var tmp [16]byte // Use a local buffer.
	b := append(tmp[:0], '%')
	for _, c := range " +-#0" { // All known flags
		if state.Flag(int(c)) { // The argument is an int for historical reasons.
			b = append(b, byte(c))
		}
	}
	if w, ok := state.Width(); ok {
		b = strconv.AppendInt(b, int64(w), 10)
	}
	if p, ok := state.Precision(); ok {
		b = append(b, '.')
		b = strconv.AppendInt(b, int64(p), 10)
	}
	b = utf8.AppendRune(b, verb)
	return string(b)
}
