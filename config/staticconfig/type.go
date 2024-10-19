package staticconfig

import "go/types"

func isAbstract(t types.Type) bool {
	if types.IsInterface(t) {
		return true
	}

	// Check if the type is a generic type that has not been instantiated
	// (meaning that it has no concrete values for its type parameters).
	switch t := t.(type) {
	case *types.Named:
		return t.Origin() == t && t.TypeParams().Len() != 0
	case *types.Alias:
		return t.Origin() == t && t.TypeParams().Len() != 0
	}

	return false
}
