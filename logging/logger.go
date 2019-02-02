package logging

import (
	"strings"

	"github.com/dogmatiq/enginekit/message"
)

// Logger logs about Dogma messages in a standardized format.
type Logger struct {
	Log             func(string)
	FormatMessageID func(string) string
}

// LogGeneric prints a log message consisting of message correlation
// information, icons and human-readable text.
func (l *Logger) LogGeneric(
	c message.Correlation,
	icons []string,
	text ...string,
) {
	w := &strings.Builder{}

	writeCorrelation(w, c, l.FormatMessageID)

	if len(icons) > 0 {
		w.WriteString("  ")

		for i, icon := range icons {
			if i > 0 {
				w.WriteByte(' ')
			}

			if icon == "" {
				w.WriteByte(' ')
			} else {
				w.WriteString(icon)
			}
		}
	}

	if len(text) > 0 {
		w.WriteString("  ")
		i := 0

		for _, t := range text {
			if t == "" {
				continue
			}

			if i > 0 {
				w.WriteByte(' ')
				w.WriteString(SeparatorIcon)
				w.WriteByte(' ')
			}

			w.WriteString(t)
			i++
		}
	}

	l.Log(w.String())
}
