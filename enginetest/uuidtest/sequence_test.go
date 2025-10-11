package uuidtest_test

import (
	"slices"
	"testing"

	. "github.com/dogmatiq/enginekit/enginetest/uuidtest"
	"github.com/dogmatiq/enginekit/protobuf/uuidpb"
)

func TestSequence(t *testing.T) {
	seq := NewSequence()

	var uuids []*uuidpb.UUID

	for range 10 {
		id := seq.Next()

		if slices.ContainsFunc(
			uuids,
			func(x *uuidpb.UUID) bool {
				return x.Equal(id)
			},
		) {
			t.Fatalf("duplicate UUID generated: %s", id)
		}

		uuids = append(uuids, id)
	}

	for i, x := range uuids {
		if !seq.IsNth(x, uint64(i)) {
			t.Errorf(
				"expected UUID %s to be the %d'th in the sequence",
				x,
				i,
			)
		}
	}
}
