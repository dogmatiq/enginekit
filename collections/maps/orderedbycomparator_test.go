package maps_test

import (
	"cmp"
	"testing"

	. "github.com/dogmatiq/enginekit/collections/maps"
	"pgregory.net/rapid"
)

type reverseStringComparator struct{}

func (c *reverseStringComparator) Compare(x, y string) int {
	if c == nil {
		panic("comparator value was not propagated")
	}
	return -cmp.Compare(x, y)
}

func TestOrderedByComparator(t *testing.T) {
	cmp := &reverseStringComparator{}

	testOrderedMap(
		t,
		func(pairs ...Pair[string, int]) *OrderedByComparator[string, int, *reverseStringComparator] {
			return NewOrderedByComparator(cmp, pairs...)
		},
		cmp.Compare,
		rapid.String(),
	)
}
