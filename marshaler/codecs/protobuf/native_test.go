package protobuf_test

import (
	"bytes"
	"testing"

	. "github.com/dogmatiq/enginekit/marshaler/codecs/protobuf"
	"github.com/dogmatiq/enginekit/protobuf/uuidpb"
	"google.golang.org/protobuf/proto"
)

func TestCodec_native(t *testing.T) {
	t.Run("it uses the correct media-type", func(t *testing.T) {
		got := DefaultNativeCodec.BasicMediaType()
		want := "application/vnd.google.protobuf"

		if got != want {
			t.Errorf("unexpected media-type: got %v, want %v", got, want)
		}
	})

	t.Run("when marshaling", func(t *testing.T) {
		t.Run("it marshals the value to text format", func(t *testing.T) {
			id := uuidpb.MustParse("c3d830ff-bdb1-4042-8c6c-08f6b7739f3e")
			got, err := DefaultNativeCodec.Marshal(id)
			if err != nil {
				t.Fatal(err)
			}

			want, err := proto.Marshal(id)
			if err != nil {
				t.Fatal(err)
			}

			if !bytes.Equal(got, want) {
				t.Errorf("unexpected data: got %q, want %q", string(got), string(want))
			}
		})

		t.Run("it returns an error if the value is not a protocol buffers message", func(t *testing.T) {
			_, err := DefaultNativeCodec.Marshal(42)
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
			want := uuidpb.MustParse("c3d830ff-bdb1-4042-8c6c-08f6b7739f3e")

			data, err := proto.Marshal(want)
			if err != nil {
				t.Fatal(err)
			}

			got := &uuidpb.UUID{}
			err = DefaultNativeCodec.Unmarshal(data, got)
			if err != nil {
				t.Fatal(err)
			}

			if !proto.Equal(got, want) {
				t.Errorf("unexpected value: got %v, want %v", got, want)
			}
		})

		t.Run("it returns an error if the type is not a protocol buffers message", func(t *testing.T) {
			var v int
			err := DefaultNativeCodec.Unmarshal(nil, &v)
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
