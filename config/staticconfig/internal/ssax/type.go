package ssax

import (
	"fmt"
	"go/types"

	"github.com/dogmatiq/enginekit/internal/typename"
	"github.com/dogmatiq/enginekit/optional"
	"golang.org/x/tools/go/ssa"
)

// Implements reports whether t implements i, regardless of whether it
// uses pointer or non-pointer method receivers.
//
// If ok is true, r is the receiver type that implements the interface, which
// may be either t or *t.
func Implements(t *ssa.Type, i *types.Interface) (r types.Type, ok bool) {
	r = t.Type()

	if IsAbstract(r) {
		return nil, false
	}

	if types.Implements(r, i) {
		return r, true
	}

	r = types.NewPointer(r)
	if types.Implements(r, i) {
		return r, true
	}

	return nil, false
}

// IsAbstract returns true if t is abstract, either because it refers to an
// interface or because it is a generic type that has not been instantiated.
func IsAbstract(t types.Type) bool {
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

// Package returns the package in which the elemental type of t is declared.
func Package(t types.Type) *types.Package {
	switch t := t.(type) {
	case *types.Named:
		return t.Obj().Pkg()
	case *types.Alias:
		return t.Obj().Pkg()
	case *types.Pointer:
		return Package(t.Elem())
	default:
		panic(fmt.Sprintf("cannot determine package for anonymous or built-in type %v", t))
	}
}

// ConcreteType returns the concrete type of v, if it can be determined at
// compile-time.
func ConcreteType(v ssa.Value) optional.Optional[types.Type] {
	t := v.Type()

	if !IsAbstract(t) {
		return optional.Some(t)
	}

	switch v := v.(type) {
	case *ssa.Alloc:
	case *ssa.BinOp:
	case *ssa.Builtin:
	case *ssa.Call:
	case *ssa.ChangeInterface:
	case *ssa.ChangeType:
	case *ssa.Const:
		// We made it past the IsAbstract() check so we know this is a constant
		// nil value for an interface, and hence no type information is present.
		return optional.None[types.Type]()
	case *ssa.Convert:
	case *ssa.Extract:
	case *ssa.Field:
	case *ssa.FieldAddr:
	case *ssa.FreeVar:
	case *ssa.Function:
	case *ssa.Global:
	case *ssa.Index:
	case *ssa.IndexAddr:
	case *ssa.Lookup:
	case *ssa.MakeChan:
	case *ssa.MakeClosure:
	case *ssa.MakeInterface:
		return ConcreteType(v.X)
	case *ssa.MakeMap:
	case *ssa.MakeSlice:
	case *ssa.MultiConvert:
	case *ssa.Next:
	case *ssa.Parameter:
	case *ssa.Phi:
	case *ssa.Slice:
	case *ssa.SliceToArrayPointer:
	case *ssa.TypeAssert:
	case *ssa.UnOp:
		_ = v

	case *ssa.Range, *ssa.Select:
		// These types implement ssa.Value, but they can not actually be used as
		// expressions in Go.
		return optional.None[types.Type]()
	}

	panic(fmt.Sprintf("unhandled %T of type %s", v, typename.OfStatic(t)))
}
