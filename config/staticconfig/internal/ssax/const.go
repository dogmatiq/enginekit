package ssax

import (
	"go/constant"
	"math"

	"github.com/dogmatiq/enginekit/optional"
	"golang.org/x/tools/go/ssa"
)

// Const returns the singlar constant value of v if possible.
func Const(v ssa.Value) optional.Optional[constant.Value] {
	return optional.TryTransform(
		StaticValue(v),
		func(v ssa.Value) (constant.Value, bool) {
			if c, ok := v.(*ssa.Const); ok {
				return c.Value, true
			}
			return nil, false
		},
	)
}

// AsString returns the singular constant string value of v if possible.
func AsString(v ssa.Value) optional.Optional[string] {
	return constAs(constant.StringVal, v)
}

// AsBool returns the singular constant boolean value of v if possible.
func AsBool(v ssa.Value) optional.Optional[bool] {
	return constAs(constant.BoolVal, v)
}

// AsInt returns the singular constant integer value of v if possible.
func AsInt(v ssa.Value) optional.Optional[int] {
	return optional.TryTransform(
		Const(v),
		func(c constant.Value) (_ int, ok bool) {
			i, ok := constant.Int64Val(c)
			return int(i), ok && i >= math.MinInt && i <= math.MaxInt
		},
	)
}

// constAsX returns the constant value of v, converted to type T by fn.
func constAs[T any](
	fn func(constant.Value) T,
	v ssa.Value,
) optional.Optional[T] {
	return optional.TryTransform(
		Const(v),
		func(c constant.Value) (_ T, ok bool) {
			defer func() {
				if recover() != nil {
					// ignore panics about type conversion
					ok = false
				}
			}()

			return fn(c), true
		},
	)
}
