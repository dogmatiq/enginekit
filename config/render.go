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
	ent Entity,
	impl optional.Optional[Source[T]],
) string {
	var w strings.Builder

	f := ent.Fidelity()

	if !f.IsExhaustive {
		w.WriteString("partial ")
	}

	if f.IsSpeculative {
		w.WriteString("speculative ")
	}

	if f.IsUnresolved || f.HasSpeculativeSubcomponents || f.HasMutuallyExclusiveSubcomponents {
		w.WriteString("non-deterministic ")
	}

	w.WriteString(label)

	identifier := ""

	if i, ok := impl.TryGet(); ok {
		identifier = strings.TrimPrefix(i.TypeName, "*")
	} else {
		for _, id := range ent.identities() {
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
