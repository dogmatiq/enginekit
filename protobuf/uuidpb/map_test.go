package uuidpb_test

import (
	"fmt"
	"testing"

	"github.com/dogmatiq/dapper"
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

func TestMapKey_DapperString(t *testing.T) {
	t.Parallel()

	subject := (&UUID{
		Upper: 0xa967a8b93f9c4918,
		Lower: 0x9a4119577be5fec5,
	}).AsMapKey()

	actual := dapper.Format(subject)
	expect := "github.com/dogmatiq/enginekit/protobuf/uuidpb.MapKey [a967a8b9-3f9c-4918-9a41-19577be5fec5]"

	if actual != expect {
		t.Fatalf("got %q, want %q", actual, expect)
	}
}

func TestMapKey_Format(t *testing.T) {
	t.Parallel()

	subject := (&UUID{
		Upper: 0xa967a8b93f9c4918,
		Lower: 0x9a4119577be5fec5,
	}).AsMapKey()

	expect := subject.String()
	actual := fmt.Sprintf("%s", subject)

	if actual != expect {
		t.Fatalf("got %q, want %q", actual, expect)
	}
}
