package sets_test

import (
	"cmp"
	"iter"
	"testing"

	. "github.com/dogmatiq/enginekit/collections/sets"
	"pgregory.net/rapid"
)

type reverseStringComparator struct{}

func (reverseStringComparator) Compare(x, y string) int {
	return -cmp.Compare(x, y)
}

func TestOrderedByComparator(t *testing.T) {
	cmp := reverseStringComparator{}

	testOrderedSet(
		t,
		func(members ...string) *OrderedByComparator[string, reverseStringComparator] {
			return NewOrderedByComparator(cmp, members...)
		},
		func(seq iter.Seq[string]) *OrderedByComparator[string, reverseStringComparator] {
			return NewOrderedByComparatorFromSeq(cmp, seq)
		},
		func(seq iter.Seq2[string, any]) *OrderedByComparator[string, reverseStringComparator] {
			return NewOrderedByComparatorFromKeys(cmp, seq)
		},
		func(seq iter.Seq2[any, string]) *OrderedByComparator[string, reverseStringComparator] {
			return NewOrderedByComparatorFromValues(cmp, seq)
		},
		cmp.Compare,
		func(m string) bool { return len(m)%2 == 0 },
		rapid.String(),
	)
}
