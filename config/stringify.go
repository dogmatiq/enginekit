package config

import (
	"fmt"

	"github.com/dogmatiq/enginekit/optional"
)

func stringify(
	label string,
	typeName optional.Optional[string],
	identities []Identity,
) string {
	identifier := "?"

	if n, ok := typeName.TryGet(); ok {
		identifier = n
	} else {
		for _, id := range identities {
			identifier = id.String()
			if id.Validate() == nil {
				break
			}
		}
	}

	return fmt.Sprintf("%s(%s)", label, identifier)
}
