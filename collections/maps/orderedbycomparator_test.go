package maps_test

import (
	"cmp"
	"iter"
	"testing"

	. "github.com/dogmatiq/enginekit/collections/maps"
	"pgregory.net/rapid"
)

type reverseStringComparator struct{}

func (c *reverseStringComparator) Compare(x, y string) int {
	return -cmp.Compare(x, y)
}

func TestOrderedByComparator(t *testing.T) {
	cmp := &reverseStringComparator{}

	testOrderedMap(
		t,
		func(pairs ...Pair[string, int]) *OrderedByComparator[string, int, *reverseStringComparator] {
			return NewOrderedByComparator(cmp, pairs...)
		},
		func(pairs iter.Seq2[string, int]) *OrderedByComparator[string, int, *reverseStringComparator] {
			return NewOrderedByComparatorFromSeq(cmp, pairs)
		},
		cmp.Compare,
		rapid.String(),
	)
}
