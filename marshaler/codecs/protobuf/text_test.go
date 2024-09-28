package protobuf_test

import (
	"regexp"
	"testing"

	. "github.com/dogmatiq/enginekit/marshaler/codecs/protobuf"
	"github.com/dogmatiq/enginekit/protobuf/uuidpb"
	"google.golang.org/protobuf/proto"
)

func TestCodec_text(t *testing.T) {
	t.Run("it uses the correct media-type", func(t *testing.T) {
		got := DefaultTextCodec.BasicMediaType()
		want := "text/vnd.google.protobuf"

		if got != want {
			t.Errorf("unexpected media-type: got %v, want %v", got, want)
		}
	})

	t.Run("when marshaling", func(t *testing.T) {
		t.Run("it marshals the value to text format", func(t *testing.T) {
			id := uuidpb.MustParse("c3d830ff-bdb1-4042-8c6c-08f6b7739f3e")
			got, err := DefaultTextCodec.Marshal(id)
			if err != nil {
				t.Fatal(err)
			}

			// Note that we need to use a regex to match an arbitrary amount of
			// whitespace in between the key and value as a result of the behavior
			// described in https://github.com/golang/protobuf/issues/1121.
			p := regexp.MustCompile(`upper:\s*14112083307322753090\nlower:\s*10118472318527446846\n`)

			if !p.Match(got) {
				t.Errorf("unexpected data: %q", string(got))
			}
		})

		t.Run("it returns an error if the value is not a protocol buffers message", func(t *testing.T) {
			_, err := DefaultTextCodec.Marshal(42)
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
			data := []byte("upper: 14112083307322753090\nlower: 10118472318527446846\n")

			got := &uuidpb.UUID{}
			err := DefaultTextCodec.Unmarshal(data, got)
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
			err := DefaultTextCodec.Unmarshal(nil, &v)
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
