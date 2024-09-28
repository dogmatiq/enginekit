package protobuf_test

import (
	"encoding/json"
	"reflect"
	"testing"

	. "github.com/dogmatiq/enginekit/marshaler/codecs/protobuf"
	"github.com/dogmatiq/enginekit/protobuf/uuidpb"
	"google.golang.org/protobuf/proto"
)

func TestCodec_json(t *testing.T) {
	t.Run("it uses the correct media-type", func(t *testing.T) {
		got := DefaultJSONCodec.BasicMediaType()
		want := "application/vnd.google.protobuf+json"

		if got != want {
			t.Errorf("unexpected media-type: got %v, want %v", got, want)
		}
	})

	t.Run("when marshaling", func(t *testing.T) {
		t.Run("it marshals the value to JSON format", func(t *testing.T) {
			id := uuidpb.MustParse("c3d830ff-bdb1-4042-8c6c-08f6b7739f3e")
			data, err := DefaultJSONCodec.Marshal(id)
			if err != nil {
				t.Fatal(err)
			}

			got := map[string]any{}

			err = json.Unmarshal(data, &got)
			if err != nil {
				t.Fatal(err)
			}

			want := map[string]any{
				"upper": "14112083307322753090",
				"lower": "10118472318527446846",
			}

			if !reflect.DeepEqual(got, want) {
				t.Errorf("unexpected value: got %v, want %v", got, want)
			}
		})

		t.Run("it returns an error if the value is not a protocol buffers message", func(t *testing.T) {
			_, err := DefaultJSONCodec.Marshal(42)
			if err == nil {
				t.Fatal("expected an error")
			}

			got := err.Error()
			want := "'int' is not a protocol buffers message"

			if got != want {
				t.Errorf("unexpected error: got %q, want %q", got, want)
			}
		})
	})

	t.Run("when unmarshaling", func(t *testing.T) {
		t.Run("it unmarshals the data", func(t *testing.T) {
			data := []byte(`{"upper":"14112083307322753090", "lower":"10118472318527446846"}`)

			got := &uuidpb.UUID{}
			err := DefaultJSONCodec.Unmarshal(data, got)
			if err != nil {
				t.Fatal(err)
			}

			want := uuidpb.MustParse("c3d830ff-bdb1-4042-8c6c-08f6b7739f3e")

			if !proto.Equal(got, want) {
				t.Errorf("unexpected value: got %v, want %v", got, want)
			}
		})

		t.Run("it returns an error if the type is not a protocol buffers message", func(t *testing.T) {
			var v int
			err := DefaultJSONCodec.Unmarshal(nil, &v)
			if err == nil {
				t.Fatal("expected an error")
			}

			got := err.Error()
			want := "'*int' is not a protocol buffers message"

			if got != want {
				t.Errorf("unexpected error: got %q, want %q", got, want)
			}
		})
	})
}
