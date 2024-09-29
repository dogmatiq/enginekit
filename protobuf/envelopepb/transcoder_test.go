package envelopepb_test

import (
	reflect "reflect"
	"testing"

	"github.com/dogmatiq/enginekit/internal/test"
	"github.com/dogmatiq/enginekit/marshaler"
	"github.com/dogmatiq/enginekit/marshaler/codecs/protobuf"
	. "github.com/dogmatiq/enginekit/protobuf/envelopepb"
	"github.com/dogmatiq/enginekit/protobuf/envelopepb/internal/stubs"
	"google.golang.org/protobuf/proto"
)

func TestTranscoder(t *testing.T) {
	m, err := marshaler.New(
		[]reflect.Type{
			reflect.TypeFor[*stubs.ProtoMessage](),
		},
		[]marshaler.Codec{
			protobuf.DefaultJSONCodec,
			protobuf.DefaultTextCodec,
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	transcoder := &Transcoder{
		MediaTypes: map[string][]string{
			`dogmatiq.enginekit.protobuf.envelopepb.stubs.ProtoMessage`: {
				`application/vnd.google.protobuf+json; type=different`,
				`application/vnd.google.protobuf+json; type=different; extra=true`,
				`application/vnd.google.protobuf+json; no-type=true`,
				`application/vnd.google.protobuf+json; type=dogmatiq.enginekit.protobuf.envelopepb.stubs.ProtoMessage`,
				`application/vnd.google.protobuf; type=dogmatiq.enginekit.protobuf.envelopepb.stubs.ProtoMessage`,
			},
		},
		Marshaler: m,
	}

	t.Run("it does not transcode envelopes that are supported by the recipient unchanged", func(t *testing.T) {
		want := &Envelope{
			PortableName: `dogmatiq.enginekit.protobuf.envelopepb.stubs.ProtoMessage`,
			// note non-canonical capitalization & spacing in media-type
			MediaType: `application/vnd.Google.protobuf+json;TYPE="dogmatiq.enginekit.protobuf.envelopepb.stubs.ProtoMessage"`,
			Data:      []byte(`{"value":"A1"}`),
		}

		got, ok, err := transcoder.Transcode(want)
		if err != nil {
			t.Fatal(err)
		}

		if got != want {
			t.Fatalf("unexpected envelope: got %p, want %p", got, want)
		}

		if !ok {
			t.Error("expected ok to be true")
		}
	})

	t.Run("it transcodes using the recipients preferred encoding if necessary", func(t *testing.T) {
		original := &Envelope{
			PortableName: `dogmatiq.enginekit.protobuf.envelopepb.stubs.ProtoMessage`,
			MediaType:    `text/vnd.google.protobuf; type=dogmatiq.enginekit.protobuf.envelopepb.stubs.ProtoMessage`,
			Data:         []byte(`value: "A1"`),
		}
		snapshot := proto.Clone(original).(*Envelope)

		got, ok, err := transcoder.Transcode(original)
		if err != nil {
			t.Fatal(err)
		}

		want := &Envelope{
			PortableName: `dogmatiq.enginekit.protobuf.envelopepb.stubs.ProtoMessage`,
			MediaType:    `application/vnd.google.protobuf+json; type=dogmatiq.enginekit.protobuf.envelopepb.stubs.ProtoMessage`,
			Data:         []byte(`{"value":"A1"}`),
		}

		test.Expect(
			t,
			"unexpected envelope",
			got,
			want,
		)

		if !ok {
			t.Error("expected ok to be true")
		}

		test.Expect(
			t,
			"original envelope was modified",
			original,
			snapshot,
		)
	})

	t.Run("it returns an error if the recipient does not support any encodings", func(t *testing.T) {
		_, ok, err := transcoder.Transcode(
			&Envelope{
				PortableName: `Unrecognized`,
				MediaType:    `text/plain`,
			},
		)
		if err != nil {
			t.Fatal(err)
		}

		if ok {
			t.Error("expected ok to be false")
		}
	})

	t.Run("it returns an error if the marshaler does not support any of the encodings supported by the recipient", func(t *testing.T) {
		transcoder := &Transcoder{
			MediaTypes: map[string][]string{
				`dogmatiq.enginekit.protobuf.envelopepb.stubs.ProtoMessage`: {
					`application/vnd.google.protobuf; type=different`,
					`application/vnd.google.protobuf; type=different; extra=true`,
					`application/vnd.google.protobuf`,
				},
			},
			Marshaler: m,
		}

		_, ok, err := transcoder.Transcode(
			&Envelope{
				PortableName: `dogmatiq.enginekit.protobuf.envelopepb.stubs.ProtoMessage`,
				MediaType:    `application/vnd.google.protobuf+json; type=dogmatiq.enginekit.protobuf.envelopepb.stubs.ProtoMessage`,
				Data:         []byte(`{"value":"A1"}`),
			},
		)
		if err != nil {
			t.Fatal(err)
		}

		if ok {
			t.Error("expected ok to be false")
		}
	})
}
