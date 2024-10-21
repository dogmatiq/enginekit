package staticconfig

import (
	"fmt"
	"go/types"
)

// isAbstract returns true if t is abstract, either because it refers to an
// interface or because it is a generic type that has not been instantiated.
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

// analyzeType analyzes a type that was discovered within a package.
//
// THe currently implementation only looks for [dogma.Application]
// implementations; handler implementations are ignored unless they are actually
// used within an application.
func analyzeType(ctx *context, t types.Type) {
	if isAbstract(t) {
		// We're only interested in concrete types; otherwise there's nothing to
		// analyze!
		return
	}

	// The sequence of the if-blocks below is important as a type
	// implements an interface only if the methods in the interface's
	// method set have non-pointer receivers. Hence the implementation
	// check for the "raw" (non-pointer) type is made first.
	//
	// A pointer to the type, on the other hand, implements the
	// interface regardless of whether pointer receivers are used or
	// not.

	if types.Implements(t, ctx.Dogma.Application) {
		analyzeApplicationType(ctx, t)
		return
	}

	p := types.NewPointer(t)
	if types.Implements(p, ctx.Dogma.Application) {
		analyzeApplicationType(ctx, p)
		return
	}
}

// packageOf returns the package in which t is declared.
//
// It panics if t is not a named type or a pointer to a named type.
func packageOf(t types.Type) *types.Package {
	switch t := t.(type) {
	case *types.Named:
		return t.Obj().Pkg()
	case *types.Alias:
		return t.Obj().Pkg()
	case *types.Pointer:
		return packageOf(t.Elem())
	default:
		panic(fmt.Sprintf("cannot determine package for anonymous or built-in type %v", t))
	}
}
