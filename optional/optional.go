package optional

// Optional represents an optional value of type T.
type Optional[T any] struct {
	value T
	ok    bool
}

// Some returns an Optional[T] that contains the given value.
func Some[T any](v T) Optional[T] {
	return Optional[T]{v, true}
}

// None returns an Optional[T] that does not contain a value.
func None[T any]() Optional[T] {
	return Optional[T]{}
}

// IsPresent returns true if the optional value is present.
func (o Optional[T]) IsPresent() bool {
	return o.ok
}

// Value returns the optional value, or panics if it is not present.
func (o Optional[T]) Value() T {
	if o.ok {
		return o.value
	}
	panic("value is not present")
}

// TryValue returns the optional value and a boolean indicating whether it is
// present.
func (o Optional[T]) TryValue() (T, bool) {
	return o.value, o.ok
}

// Transform applies a transformation to v, returning a new optional value that
// contains the result of the function.
//
// If v's value is not present, the returned optional value will also not
// contain a value.
func Transform[T, U any](
	v Optional[T],
	fn func(T) U,
) Optional[U] {
	if v.ok {
		return Some(fn(v.value))
	}
	return None[U]()
}
