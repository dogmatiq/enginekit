package logging

import (
	"io"
	"strings"

	"github.com/dogmatiq/enginekit/message"
	"github.com/dogmatiq/iago/must"
)

// FormatCorrelation returns a formatted representation of message correlation
// information.
//
// f is a function that formats each message ID. If it is nil, no formatting is
// performed. f must accept empty message IDs.
func FormatCorrelation(
	c message.Correlation,
	f func(string) string,
) string {
	w := &strings.Builder{}
	writeCorrelation(w, c, f)
	return w.String()
}

// FormatCorrelation returns a formatted representation of message correlation
// information.
//
// f is a function that formats each message ID. If it is nil, no formatting is
// performed. f must accept empty message IDs.
func writeCorrelation(
	w io.Writer,
	c message.Correlation,
	f func(string) string,
) {
	if f == nil {
		f = func(s string) string {
			if s == "" {
				return "-"
			}

			return s
		}
	}

	must.WriteString(w, MessageIDIcon)
	must.WriteByte(w, ' ')
	must.WriteString(w, f(c.MessageID))
	must.WriteString(w, "  ")

	must.WriteString(w, CausationIDIcon)
	must.WriteByte(w, ' ')
	must.WriteString(w, f(c.CausationID))
	must.WriteString(w, "  ")

	must.WriteString(w, CorrelationIDIcon)
	must.WriteByte(w, ' ')
	must.WriteString(w, f(c.CorrelationID))
}
