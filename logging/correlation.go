package logging

import (
	"io"
	"strings"

	"github.com/dogmatiq/enginekit/message"
	"github.com/dogmatiq/iago"
)

// WriteCorrelation writes a formatted representation of message correlation
// information.
//
// tr is the length to which IDs are truncated. If tr positive the tail of the
// ID is truncated. If negative, the head of the ID is truncated. If tr is zero
// no truncation is performed
func WriteCorrelation(w io.Writer, c message.Correlation, tr int) (n int,
	err error) {
	defer iago.Recover(&err)

	n += iago.MustWriteString(w, MessageIDIcon)
	n += iago.MustWriteByte(w, ' ')
	n += iago.MustWriteString(w, truncate(c.MessageID, tr))
	n += iago.MustWriteString(w, "  ")

	n += iago.MustWriteString(w, CausationIDIcon)
	n += iago.MustWriteByte(w, ' ')
	n += iago.MustWriteString(w, truncate(c.CausationID, tr))
	n += iago.MustWriteString(w, "  ")

	n += iago.MustWriteString(w, CorrelationIDIcon)
	n += iago.MustWriteByte(w, ' ')
	n += iago.MustWriteString(w, truncate(c.CorrelationID, tr))

	return
}

// FormatCorrelation returns a formatted representation of message correlation
// information.
//
// If larger than zero, tr is the length to which IDs are truncated.
func FormatCorrelation(c message.Correlation, tr int) string {
	var w strings.Builder
	iago.Must(WriteCorrelation(&w, c, tr))
	return w.String()
}
