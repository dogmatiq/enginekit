package uuidpb_test

import (
	"testing"

	. "github.com/dogmatiq/enginekit/protobuf/uuidpb"
)

func TestSet(t *testing.T) {
	t.Parallel()

	id1 := Generate()
	id2 := Generate()

	s := Set{}

	if s.Has(id1) {
		t.Errorf("unexpected id1 in set")
	}
	if s.Has(id2) {
		t.Errorf("unexpected id1 in set")
	}

	s.Add(id1)
	s.Add(id2)

	if !s.Has(id1) {
		t.Errorf("expected id1 to be in set")
	}
	if !s.Has(id2) {
		t.Errorf("expected id1 to be in set")
	}

	s.Delete(id2)

	if !s.Has(id1) {
		t.Errorf("expected id1 to be in set")
	}
	if s.Has(id2) {
		t.Errorf("unexpected id1 in set")
	}
}
