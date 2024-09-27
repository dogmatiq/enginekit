package stateless_test

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/dogmatiq/dogma"
	. "github.com/dogmatiq/enginekit/marshaler/codecs/stateless"
)

func TestCodec(t *testing.T) {
	var codec Codec

	t.Run("uses a short string as the portable name", func(t *testing.T) {
		rt := reflect.TypeOf(dogma.StatelessProcessRoot)

		names := codec.PortableNames(
			[]reflect.Type{rt},
		)

		got := names[rt]
		want := "process"

		if got != want {
			t.Errorf("unexpected portable name for %s: got %q, want %q", rt, got, want)
		}
	})

	t.Run("it uses the correct media-type", func(t *testing.T) {
		got := codec.BasicMediaType()
		want := "application/x-empty"

		if got != want {
			t.Errorf("unexpected media-type: got %v, want %v", got, want)
		}
	})

	t.Run("it marshals the value to an empty byte-slice", func(t *testing.T) {
		got, err := codec.Marshal(dogma.StatelessProcessRoot)
		if err != nil {
			t.Fatal(err)
		}

		want := []byte(nil)

		if !bytes.Equal(got, want) {
			t.Errorf("expected data: got %v, want %v", got, want)
		}
	})

	t.Run("it returns an error if passed any value other than dogma.StatelessProcessRoot", func(t *testing.T) {
		_, err := codec.Marshal(123)
		if err == nil {
			t.Fatalf("expected an error")
		}
	})

	t.Run("it does nothing on a successful unmarshal", func(t *testing.T) {
		before := dogma.StatelessProcessRoot

		err := codec.Unmarshal(nil, &dogma.StatelessProcessRoot)
		if err != nil {
			t.Fatal(err)
		}

		if !reflect.DeepEqual(dogma.StatelessProcessRoot, before) {
			t.Fatal("unexpected change to dogma.StatelessProcessRoot")
		}
	})

	t.Run("it returns an error when attempting to unmarshal a non-empty byte-slice", func(t *testing.T) {
		data := []byte(` `)
		err := codec.Unmarshal(data, &dogma.StatelessProcessRoot)
		if err == nil {
			t.Fatalf("expected an error")
		}

		got := err.Error()
		want := "expected empty data, got 1 byte(s)"

		if got != want {
			t.Errorf("unexpected error: got %q, want %q", got, want)
		}
	})

	t.Run("it returns an error if passed any value other than the address of dogma.StatelessProcessRoot", func(t *testing.T) {
		err := codec.Unmarshal(nil, 123)
		if err == nil {
			t.Fatalf("expected an error")
		}

		got := err.Error()
		want := "int is not a pointer to dogma.StatelessProcessRoot"

		if got != want {
			t.Errorf("unexpected error: got %q, want %q", got, want)
		}
	})
}
