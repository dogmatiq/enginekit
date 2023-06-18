package uuidpb_test

import (
	"testing"

	. "github.com/dogmatiq/enginekit/protobuf/uuidpb"
)

func TestMap(t *testing.T) {
	t.Parallel()

	k1 := Generate()
	k2 := Generate()

	m := Map[string]{}

	expect := ""
	if v := m.Get(k1); v != expect {
		t.Errorf("unexpected value for k1: got %q, want %q", v, expect)
	}
	if _, ok := m.TryGet(k1); ok {
		t.Errorf("expected k1 to be absent")
	}

	m.Set(k1, "foo")
	m.Set(k2, "bar")

	expect = "foo"
	if v := m.Get(k1); v != expect {
		t.Errorf("unexpected value for k1: got %q, want %q", v, expect)
	}
	if v, ok := m.TryGet(k1); !ok {
		t.Errorf("expected k1 to be present")
	} else if v != expect {
		t.Errorf("unexpected value for k1: got %q, want %q", v, expect)
	}

	expect = "bar"
	if v := m.Get(k2); v != expect {
		t.Errorf("unexpected value for k1: got %q, want %q", v, expect)
	}
	if v, ok := m.TryGet(k2); !ok {
		t.Errorf("expected k2 to be present")
	} else if v != expect {
		t.Errorf("unexpected value for k2: got %q, want %q", v, expect)
	}

	m.Delete(k1)

	expect = ""
	if v := m.Get(k1); v != expect {
		t.Errorf("unexpected value for k1: got %q, want %q", v, expect)
	}
	if _, ok := m.TryGet(k1); ok {
		t.Errorf("expected k1 to be absent")
	}

	expect = "bar"
	for k, v := range m {
		id := k.AsUUID()

		if !id.Equal(k2) {
			t.Errorf("unexpected key: got %q, want %q", id, k2)
		}
		if v != expect {
			t.Errorf("unexpected value: got %q, want %q", v, expect)
		}
	}
}
