package config

import (
	"fmt"
	"strings"

	"github.com/dogmatiq/enginekit/optional"
)

// renderList returns a human-readable list of items.
func renderList[T any](items []T) string {
	var s string

	for i, item := range items {
		if i == len(items)-1 {
			s += " and "
		} else if i > 0 {
			s += ", "
		}
		s += fmt.Sprint(item)
	}

	return s
}

func renderEntity[T any](
	label string,
	e Entity,
	s optional.Optional[Source[T]],
) string {
	w := &strings.Builder{}

	writeComponentPrefix(w, label, e)

	identifier := ""

	if i, ok := s.TryGet(); ok {
		identifier = strings.TrimPrefix(i.TypeName, "*")
	} else {
		for _, id := range e.identitiesAsConfigured() {
			if norm, err := Normalize(id); err == nil {
				identifier = norm.String()
				break
			}
			identifier = id.String()
		}
	}

	if identifier != "" {
		w.WriteByte(':')
		w.WriteString(identifier)
	}

	return w.String()
}

func writeComponentPrefix(
	w *strings.Builder,
	label string,
	c Component,
) string {
	f := c.Fidelity()

	if f.IsPartial {
		w.WriteString("partial ")
	}

	if f.IsSpeculative {
		w.WriteString("speculative ")
	}

	if f.IsUnresolved {
		w.WriteString("unresolved ")
	}

	w.WriteString(label)

	return w.String()
}
