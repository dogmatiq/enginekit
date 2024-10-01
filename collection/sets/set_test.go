package sets_test

import (
	"testing"

	. "github.com/dogmatiq/enginekit/collection/sets"
	"pgregory.net/rapid"
)

func TestSet(t *testing.T) {
	testSet(
		t,
		New[string],
		func(a, b string) bool { return a == b },
		func(m string) bool { return len(m)%2 == 0 },
		rapid.String(),
	)
}
