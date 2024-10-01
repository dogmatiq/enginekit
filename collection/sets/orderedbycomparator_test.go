package sets_test

import (
	"cmp"
	"testing"

	. "github.com/dogmatiq/enginekit/collection/sets"
	"pgregory.net/rapid"
)

type reverseStringComparator struct{}

func (reverseStringComparator) Compare(a, b string) int {
	return -cmp.Compare(a, b)
}

func TestOrderedByComparator(t *testing.T) {
	cmp := reverseStringComparator{}

	testOrderedSet(
		t,
		func(members ...string) OrderedByComparator[string, reverseStringComparator] {
			return NewOrderedByComparator(cmp, members...)
		},
		cmp.Compare,
		func(m string) bool { return len(m)%2 == 0 },
		rapid.String(),
	)
}
