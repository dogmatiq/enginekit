package maps_test

import (
	"cmp"
	"testing"

	. "github.com/dogmatiq/enginekit/collections/maps"
	"pgregory.net/rapid"
)

func TestOrdered(t *testing.T) {
	testOrderedMap(
		t,
		NewOrdered[string, int],
		cmp.Compare[string],
		rapid.String(),
	)
}
