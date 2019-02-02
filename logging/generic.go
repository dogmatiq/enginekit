package logging

import (
	"io"
	"strings"

	"github.com/dogmatiq/enginekit/message"
	"github.com/dogmatiq/iago"
)

// Write writes a generic log message to w.
//
// tr is the length to which IDs are truncated. If tr positive the tail of the
// ID is truncated. If negative, the head of the ID is truncated. If tr is zero
// no truncation is performed
//
// icons is a set of "status icons" to include after the correlation information.
// text is a set of human-readable messages.
//
// The output does not include a trailing newline.
func Write(
	w io.Writer,
	c message.Correlation,
	tr int,
	icons []string,
	text ...string,
) (n int, err error) {
	defer iago.Recover(&err)

	iago.Must(WriteCorrelation(w, c, tr))

	if len(icons) > 0 {
		n += iago.MustWriteString(w, "  ")

		for i, icon := range icons {
			if i > 0 {
				n += iago.MustWriteByte(w, ' ')
			}

			if icon == "" {
				n += iago.MustWriteByte(w, ' ')
			} else {
				n += iago.MustWriteString(w, icon)
			}
		}
	}

	if len(text) > 0 {
		n += iago.MustWriteString(w, "  ")
		i := 0

		for _, t := range text {
			if t == "" {
				continue
			}

			if i > 0 {
				n += iago.MustWriteByte(w, ' ')
				n += iago.MustWriteString(w, SeparatorIcon)
				n += iago.MustWriteByte(w, ' ')
			}

			n += iago.MustWriteString(w, t)
			i++
		}
	}

	return
}

// Format returns a generic log message.
//
// tr is the length to which IDs are truncated. If tr positive the tail of the
// ID is truncated. If negative, the head of the ID is truncated. If tr is zero
// no truncation is performed
//
// icons is a set of "status icons" to include after the correlation information.
// text is a set of human-readable messages.
//
// The output does not include a trailing newline.
func Format(
	c message.Correlation,
	tr int,
	icons []string,
	text ...string,
) string {
	var w strings.Builder

	iago.Must(Write(
		&w,
		c,
		tr,
		icons,
		text...,
	))

	return w.String()
}
