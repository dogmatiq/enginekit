package config

import (
	"strings"

	"github.com/dogmatiq/enginekit/optional"
)

func stringify[T any](
	label string,
	ent Entity,
	impl optional.Optional[Implementation[T]],
) string {
	if !ent.IsExhaustive() {
		label = "partial " + label
	}

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

	if identifier == "" {
		return label
	}

	return label + ":" + identifier
}
