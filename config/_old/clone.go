package config

import "slices"

// clonable is an interface/constraint for types that can be cloned.
type clonable interface {
	// clone returns a deep copy of the receiver. The return value must be of
	// the same type as the receiver.
	clone() any
}

func clone[T clonable](c T) T {
	return c.clone().(T)
}

func cloneInPlace[T clonable](c *T) {
	*c = clone(*c)
}

func clonePointee[
	P interface {
		clonable
		*T
	},
	T any,
](c T) T {
	return *clone[P](&c)
}

func clonePointeeInPlace[
	P interface {
		clonable
		*T
	},
	T any,
](c P) {
	*c = clonePointee[P](*c)
}

func cloneSlice[T clonable](values []T) []T {
	clones := slices.Clone(values)

	for i, c := range values {
		clones[i] = clone(c)
	}

	return clones
}

func cloneSliceInPlace[T clonable](values *[]T) {
	*values = cloneSlice(*values)
}
