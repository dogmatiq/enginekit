package maps_test

import (
	"cmp"
	"testing"

	. "github.com/dogmatiq/enginekit/collection/maps"
	"pgregory.net/rapid"
)

type reverseStringComparator struct{}

func (reverseStringComparator) Compare(a, b string) int {
	return -cmp.Compare(a, b)
}

func TestOrderedByComparator(t *testing.T) {
	cmp := reverseStringComparator{}

	testOrderedMap(
		t,
		func(pairs ...Pair[string, int]) OrderedByComparator[string, int, reverseStringComparator] {
			return NewOrderedByComparator(cmp, pairs...)
		},
		cmp.Compare,
		rapid.String(),
	)
}
