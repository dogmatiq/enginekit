package optional_test

import (
	"strconv"
	"testing"

	. "github.com/dogmatiq/enginekit/optional"
)

func TestSome(t *testing.T) {
	v := Some(42)

	if !v.IsPresent() {
		t.Fatal("expected value to be present")
	}

	if x := v.Get(); x != 42 {
		t.Fatalf("unexpected value: got %v, want 42", x)
	}

	if x, ok := v.TryGet(); !ok {
		t.Fatal("expected value to be present")
	} else if x != 42 {
		t.Fatalf("unexpected value: got %v, want 42", x)
	}
}

func TestNone(t *testing.T) {
	v := None[int]()

	if v.IsPresent() {
		t.Fatal("expected value to be absent")
	}

	func() {
		defer func() {
			if x := recover(); x == nil {
				t.Fatal("expected a panic")
			}
		}()

		v.Get()
	}()

	if x, ok := v.TryGet(); ok {
		t.Fatal("expected value to be absent")
	} else if x != 0 {
		t.Fatalf("unexpected value: got %v, want 0", x)
	}
}

func TestTransform(t *testing.T) {
	v := Some(42)
	u := Transform(v, strconv.Itoa)

	if !u.IsPresent() {
		t.Fatal("expected transformed value to be present")
	}

	if x := u.Get(); x != "42" {
		t.Fatalf("unexpected transformed value: got %v, want 'x42'", x)
	}

	v = None[int]()
	u = Transform(v, strconv.Itoa)

	if u.IsPresent() {
		t.Fatal("expected transformed value to be absent")
	}
}
