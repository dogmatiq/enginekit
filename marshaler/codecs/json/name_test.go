package json

import (
	"reflect"
	"testing"
)

type (
	noTypeParameters                 struct{}
	singleTypeParameter[A any]       struct{}
	multipleTypeParameters[A, B any] struct{}
)

func TestPortableName(t *testing.T) {
	cases := []struct {
		Type reflect.Type
		Want string
	}{
		{reflect.TypeFor[*noTypeParameters](), "noTypeParameters"},
		{reflect.TypeFor[singleTypeParameter[noTypeParameters]](), "singleTypeParameter[noTypeParameters]"},
		{reflect.TypeFor[singleTypeParameter[int]](), "singleTypeParameter[int]"},
		{reflect.TypeFor[multipleTypeParameters[noTypeParameters, int]](), "multipleTypeParameters[noTypeParameters,int]"},
	}

	for _, c := range cases {
		got, ok := portableName(c.Type)
		if !ok {
			t.Fatal("expected portable name to be valid")
		}

		if got != c.Want {
			t.Errorf("unexpected portable name for %s: got %q, want %q", c.Type, got, c.Want)
		}
	}

	_, ok := portableName(reflect.TypeOf(struct{}{}))
	if ok {
		t.Fatal("expected portable name to be invalid for anonymous type")
	}
}
