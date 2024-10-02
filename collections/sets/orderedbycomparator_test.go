package sets_test

import (
	"cmp"
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
		cmp.Compare,
		func(m string) bool { return len(m)%2 == 0 },
		rapid.String(),
	)
}
