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

	if identifier == "" {
		return label
	}

	return label + ":" + identifier
}
