package json_test

import (
	"reflect"
	"testing"

	. "github.com/dogmatiq/enginekit/marshaler/codecs/json"
)

func TestCodec(t *testing.T) {
	var codec Codec

	type Value struct {
		Value string `json:"value"`
	}

	t.Run("uses the unqualified type-name as the portable type", func(t *testing.T) {
		rt := reflect.TypeFor[Value]()

		names := codec.PortableNames(
			[]reflect.Type{rt},
		)

		got := names[rt]
		want := "Value"

		if got != want {
			t.Errorf("unexpected portable name for %s: got %q, want %q", rt, got, want)
		}
	})

	t.Run("it uses the correct media-type", func(t *testing.T) {
		got := codec.BasicMediaType()
		want := "application/json"

		if got != want {
			t.Errorf("unexpected media-type: got %v, want %v", got, want)
		}
	})

	t.Run("it marshals the value to JSON format", func(t *testing.T) {
		got, err := codec.Marshal(Value{"<value>"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		want := `{"value":"\u003cvalue\u003e"}`

		if string(got) != want {
			t.Errorf("expected data: got %s, want %s", got, want)
		}
	})

	t.Run("it unmarshals the data", func(t *testing.T) {
		data := []byte(`{"value":"\u003cvalue\u003e"}`)

		var got Value
		err := codec.Unmarshal(data, &got)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		want := Value{"<value>"}

		if got != want {
			t.Errorf("unexpected value: got %v, want %v", got, want)
		}
	})
}
