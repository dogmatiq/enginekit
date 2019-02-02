package logging

import (
	"io"
	"strings"

	"github.com/dogmatiq/enginekit/message"
	"github.com/dogmatiq/iago"
)

// WriteMessage writes a log message about a Dogma message to w.
//
// tr is the length to which IDs are truncated. If tr positive the tail of the
// ID is truncated. If negative, the head of the ID is truncated. If tr is zero
// no truncation is performed
//
// text is a set of human-readable messages.
//
// The output does not include a trailing newline.
func WriteMessage(
	w io.Writer,
	md message.MetaData,
	tr int,
	isRetry bool,
	err error,
	text []string,
) (int, error) {
	icons := []string{
		DirectionIcon(md.Direction, err != nil),
	}

	before := []string{
		md.Type.String() + md.Role.Marker(),
	}

	if err != nil {
		icons = append(icons, ErrorIcon)
		before = append(before, err.Error())
	} else if isRetry {
		icons = append(icons, RetryIcon)
	} else {
		icons = append(icons, "")
	}

	return Write(
		w,
		md.Correlation,
		tr,
		icons,
		append(before, text...),
	)
}

// FormatMessage returns a log message about a Dogma message to w.
//
// tr is the length to which IDs are truncated. If tr positive the tail of the
// ID is truncated. If negative, the head of the ID is truncated. If tr is zero
// no truncation is performed
//
// text is a set of human-readable messages.
//
// The output does not include a trailing newline.
func FormatMessage(
	md message.MetaData,
	tr int,
	isRetry bool,
	err error,
	text []string,
) string {
	var w strings.Builder

	iago.Must(WriteMessage(
		&w,
		md,
		tr,
		isRetry,
		err,
		text,
	))

	return w.String()
}
