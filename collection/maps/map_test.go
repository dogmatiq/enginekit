package maps_test

import (
	"testing"

	. "github.com/dogmatiq/enginekit/collection/maps"
	"pgregory.net/rapid"
)

func TestMap(t *testing.T) {
	testMap(
		t,
		New[string, int],
		func(a, b string) bool { return a == b },
		rapid.String(),
	)
}
