package envelopepb_test

import (
	reflect "reflect"
	"testing"

	. "github.com/dogmatiq/enginekit/enginetest/stubs"
	. "github.com/dogmatiq/enginekit/internal/stubs"
	. "github.com/dogmatiq/enginekit/internal/test"
	"github.com/dogmatiq/enginekit/marshaler"
	"github.com/dogmatiq/enginekit/marshaler/codecs/json"
	"github.com/dogmatiq/enginekit/marshaler/codecs/protobuf"
	. "github.com/dogmatiq/enginekit/protobuf/envelopepb"
	"google.golang.org/protobuf/proto"
)

func TestTranscoder(t *testing.T) {
	m, err := marshaler.New(
		[]reflect.Type{
			reflect.TypeFor[*ProtoStubA](),
			reflect.TypeOf(CommandA1),
		},
		[]marshaler.Codec{
			protobuf.DefaultJSONCodec,
			protobuf.DefaultTextCodec,
			json.DefaultCodec,
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	transcoder := &Transcoder{
		MediaTypes: map[reflect.Type][]string{
			reflect.TypeFor[*ProtoStubA](): {
				`application/vnd.google.protobuf+json; type=different`,
				`application/vnd.google.protobuf+json; type=different; extra=true`,
				`application/vnd.google.protobuf+json; no-type=true`,
				`application/vnd.google.protobuf+json; type=dogmatiq.enginekit.stubs.ProtoStubA`,
				`application/vnd.google.protobuf; type=dogmatiq.enginekit.stubs.ProtoStubA`,
			},
		},
		Marshaler: m,
	}

	t.Run("it does not transcode envelopes that are supported by the recipient unchanged", func(t *testing.T) {
		want := &Envelope{
			// note non-canonical capitalization & spacing in media-type
			MediaType: `application/vnd.Google.protobuf+json;TYPE="dogmatiq.enginekit.stubs.ProtoStubA"`,
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
			MediaType: `text/vnd.google.protobuf; type=dogmatiq.enginekit.stubs.ProtoStubA`,
			Data:      []byte(`value: "A1"`),
		}
		snapshot := proto.Clone(original).(*Envelope)

		got, ok, err := transcoder.Transcode(original)
		if err != nil {
			t.Fatal(err)
		}

		want := &Envelope{
			MediaType: `application/vnd.google.protobuf+json; type=dogmatiq.enginekit.stubs.ProtoStubA`,
			Data:      []byte(`{"value":"A1"}`),
		}

		Expect(
			t,
			"unexpected envelope",
			got,
			want,
		)

		if !ok {
			t.Error("expected ok to be true")
		}

		Expect(
			t,
			"original envelope was modified",
			original,
			snapshot,
		)
	})

	t.Run("it returns false if the recipient does not support any encodings", func(t *testing.T) {
		_, ok, err := transcoder.Transcode(
			&Envelope{
				MediaType: `text/plain; type="CommandStub[TypeA]"`,
			},
		)
		if err != nil {
			t.Fatal(err)
		}

		if ok {
			t.Error("expected ok to be false")
		}
	})

	t.Run("it returns false if the marshaler does not support any of the encodings supported by the recipient", func(t *testing.T) {
		transcoder := &Transcoder{
			MediaTypes: map[reflect.Type][]string{
				reflect.TypeFor[*ProtoStubA](): {
					`application/vnd.google.protobuf; type=different`,
					`application/vnd.google.protobuf; type=different; extra=true`,
					`application/vnd.google.protobuf; no-type=true`,
				},
			},
			Marshaler: m,
		}

		_, ok, err := transcoder.Transcode(
			&Envelope{
				MediaType: `application/vnd.google.protobuf+json; type=dogmatiq.enginekit.stubs.ProtoStubA`,
				Data:      []byte(`{"value":"A1"}`),
			},
		)
		if err != nil {
			t.Fatal(err)
		}

		if ok {
			t.Error("expected ok to be false")
		}
	})

	t.Run("it returns an error if the marshaler does not support the original encoding", func(t *testing.T) {
		_, _, err := transcoder.Transcode(
			&Envelope{
				MediaType: `application/unsupported; type=irrelevant`,
			},
		)
		if err == nil {
			t.Fatal("expected an error")
		}

		got := err.Error()
		want := `the portable type name 'irrelevant' is not recognized`

		if got != want {
			t.Errorf("unexpected error: got %q, want %q", got, want)
		}
	})

	t.Run("it returns an error if the original encoding does not have a type", func(t *testing.T) {
		_, _, err := transcoder.Transcode(
			&Envelope{
				MediaType: `application/unsupported`,
			},
		)
		if err == nil {
			t.Fatal("expected an error")
		}

		got := err.Error()
		want := `the media-type does not specify a 'type' parameter`

		if got != want {
			t.Errorf("unexpected error: got %q, want %q", got, want)
		}
	})
}
