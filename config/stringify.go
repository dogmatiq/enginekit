package config

import (
	"fmt"

	"github.com/dogmatiq/enginekit/optional"
)

func stringify[T any](
	label string,
	impl optional.Optional[Implementation[T]],
	identities []Identity,
) string {
	identifier := "?"

	if i, ok := impl.TryGet(); ok {
		identifier = i.TypeName
	} else {
		for _, id := range identities {
			if norm, err := Normalize(id); err == nil {
				identifier = norm.String()
				break
			}

			identifier = id.String()
		}
	}

	return fmt.Sprintf("%s(%s)", label, identifier)
}
