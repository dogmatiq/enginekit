package config

import "strings"

// ComponentError indicates that a [Component] is invalid.
type ComponentError struct {
	Component Component
	Causes    []error
}

func (e ComponentError) Error() string {
	var w strings.Builder

	w.WriteString(e.Component.String())
	w.WriteString(" is invalid")

	if len(e.Causes) == 1 {
		w.WriteString(": ")
		w.WriteString(e.Causes[0].Error())
	} else if len(e.Causes) > 1 {
		w.WriteString(":")

		for _, cause := range e.Causes {
			lines := strings.Split(cause.Error(), "\n")

			for i, line := range lines {
				w.WriteByte('\n')
				if i == 0 {
					w.WriteString("- ")
				} else {
					w.WriteString("  ")
				}
				w.WriteString(line)
			}
		}
	}

	return w.String()
}

func (e ComponentError) Unwrap() []error {
	return e.Causes
}
