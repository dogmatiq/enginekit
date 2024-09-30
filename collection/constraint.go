package collection

// Ordered is a constraint for types that can be ordered.
type Ordered[E any] interface {
	Compare(E) int
}
