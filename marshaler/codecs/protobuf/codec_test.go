package protobuf_test

import (
	"reflect"
	"testing"

	. "github.com/dogmatiq/enginekit/marshaler/codecs/protobuf"
	"github.com/dogmatiq/enginekit/protobuf/uuidpb"
)

func TestCodec_Query(t *testing.T) {
	var codec Codec

	t.Run("uses the protocol name as the portable type", func(t *testing.T) {
		rt := reflect.TypeFor[*uuidpb.UUID]()

		names := codec.PortableNames(
			[]reflect.Type{rt},
		)

		got := names[rt]
		want := "dogma.protobuf.UUID"

		if got != want {
			t.Errorf("unexpected portable name for %s: got %q, want %q", rt, got, want)
		}
	})

	t.Run("excludes non-protocol-buffers types", func(t *testing.T) {
		rt := reflect.TypeFor[int]()

		names := codec.PortableNames(
			[]reflect.Type{rt},
		)

		if len(names) != 0 {
			t.Errorf("expected no types, got %d", len(names))
		}
	})
}
