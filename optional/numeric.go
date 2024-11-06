package optional

type numeric interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 |
		~float32 | ~float64 |
		~complex64 | ~complex128
}

// Sum returns the sum of all of the given values. If are of the values are
// none, then the result is also none.
func Sum[T numeric](values ...Optional[T]) Optional[T] {
	var sum T

	for _, value := range values {
		v, ok := value.TryGet()
		if !ok {
			return None[T]()
		}

		sum += v
	}

	return Some(sum)
}
