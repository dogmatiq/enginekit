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

// PartialError indicates that a [Component] has only partial
// configuration.
//
// See [Fidelity] for more information.
type PartialError struct{}

func (e PartialError) Error() string {
	return "could not evaluate entire configuration"
}

// SpeculativeError indicates that a [Component]'s inclusion in the
// configuration is subject to some condition that could not be evaluated at the
// time the configuration was built.
//
// See [Fidelity] for more information.
type SpeculativeError struct{}

func (e SpeculativeError) Error() string {
	return "conditions for the component's inclusion in the configuration could not be evaluated"
}

// UnresolvedError indicates that a [Component] is contains values that could
// not be resolved at the time the configuration was built.
//
// See [Fidelity] for more information.
type UnresolvedError struct{}

func (e UnresolvedError) Error() string {
	return "configuration includes values that could not be evaluated"
}
