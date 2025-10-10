package envelopepb_test

import (
	"testing"
	"time"

	"github.com/dogmatiq/dogma"
	. "github.com/dogmatiq/enginekit/enginetest/stubs"
	. "github.com/dogmatiq/enginekit/internal/test"
	. "github.com/dogmatiq/enginekit/protobuf/envelopepb"
	"github.com/dogmatiq/enginekit/protobuf/identitypb"
	"github.com/dogmatiq/enginekit/protobuf/uuidpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestPacker_packAndUnpack(t *testing.T) {
	id := uuidpb.Generate()
	now := time.Now()

	packer := &Packer{
		Site: &identitypb.Identity{
			Name: "site",
			Key:  uuidpb.Generate(),
		},
		Application: &identitypb.Identity{
			Name: "app",
			Key:  uuidpb.Generate(),
		},
		GenerateID: func() *uuidpb.UUID {
			return id
		},
		Now: func() time.Time {
			return now
		},
	}

	got := packer.Pack(CommandA1)

	if err := got.Validate(); err != nil {
		t.Fatalf("packer produced an invalid envelope: %v", err)
	}

	want := &Envelope{
		MessageId:         id,
		CausationId:       id,
		CorrelationId:     id,
		SourceSite:        packer.Site,
		SourceApplication: packer.Application,
		SourceInstanceId:  "",
		CreatedAt:         timestamppb.New(now),
		Description:       `command(stubs.TypeA:A1, valid)`,
		TypeId:            uuidpb.MustParse(MessageTypeID[*CommandStub[TypeA]]()),
		Data:              []byte(`{"content":"A1"}`),
	}

	Expect(
		t,
		"unexpected envelope",
		got,
		want,
	)

	gotMessage, err := Unpack(want)
	if err != nil {
		t.Fatal(err)
	}

	Expect(
		t,
		"unexpected message",
		gotMessage,
		dogma.Message(CommandA1),
	)

	t.Run("func Pack()", func(t *testing.T) {
		t.Run("it panics if passed an unregistered message", func(t *testing.T) {
			ExpectPanic(
				t,
				"*envelopepb_test.T is not a registered message type",
				func() {
					type T struct{ dogma.Command }
					packer.Pack(&T{})
				},
			)
		})

		t.Run("it panics if the message cannot be marshaled", func(t *testing.T) {
			ExpectPanic(
				t,
				"unable to marshal *envelopepb_test.T: json: unsupported type: func()",
				func() {
					type T struct{ CommandStub[func()] }
					dogma.RegisterCommand[*T]("622003a4-01a5-4c77-8a4c-cb36b51994e7")
					packer.Pack(&T{})
				},
			)
		})

		t.Run("it panics if the envelope is invalid", func(t *testing.T) {
			before := packer.Site
			packer.Site = &identitypb.Identity{} // invalid
			defer func() { packer.Site = before }()

			ExpectPanic(
				t,
				"invalid source site (/00000000-0000-0000-0000-000000000000): invalid name: must be between 1 and 255 bytes",
				func() {
					packer.Pack(CommandA1)
				},
			)
		})
	})

	t.Run("func Unpack()", func(t *testing.T) {
		t.Run("it returns an error if the message type is not registered", func(t *testing.T) {
			env := &Envelope{
				TypeId: uuidpb.MustParse("f1816a71-3593-4771-8d8b-327650571288"),
				Data:   []byte(`{"content":"A1"}`),
			}

			want := "f1816a71-3593-4771-8d8b-327650571288 is not a registered message type ID"

			if _, err := Unpack(env); err == nil {
				t.Fatal("expected an error")
			} else if err.Error() != want {
				t.Fatalf("unexpected error: got %q, want %q", err, want)
			}
		})

		t.Run("it returns an error if the message cannot be unmarshaled", func(t *testing.T) {
			env := &Envelope{
				TypeId: uuidpb.MustParse(MessageTypeID[*CommandStub[TypeA]]()),
				Data:   []byte(`}`),
			}

			want := "unable to unmarshal *stubs.CommandStub[github.com/dogmatiq/enginekit/enginetest/stubs.TypeA]: invalid character '}' looking for beginning of value"

			if _, err := Unpack(env); err == nil {
				t.Fatal("expected an error")
			} else if err.Error() != want {
				t.Fatalf("unexpected error: got %q, want %q", err, want)
			}
		})
	})
}

func TestWithCause(t *testing.T) {
	packer := &Packer{
		Application: &identitypb.Identity{
			Name: "app",
			Key:  uuidpb.Generate(),
		},
	}

	root := packer.Pack(CommandA1)
	cause := packer.Pack(EventA1, WithCause(root))
	got := packer.Pack(CommandA2, WithCause(cause))

	if !got.CausationId.Equal(cause.MessageId) {
		t.Fatalf("unexpected causation ID: got %s, want %s", got.CausationId, cause.MessageId)
	}

	if !got.CorrelationId.Equal(cause.CorrelationId) {
		t.Fatalf("unexpected correlation ID: got %s, want %s", got.CorrelationId, root.MessageId)
	}
}

func TestWithHandler(t *testing.T) {
	packer := &Packer{
		Application: &identitypb.Identity{
			Name: "app",
			Key:  uuidpb.Generate(),
		},
	}

	handler := &identitypb.Identity{
		Name: "handler",
		Key:  uuidpb.Generate(),
	}

	got := packer.Pack(CommandA1, WithHandler(handler))

	if !got.SourceHandler.Equal(handler) {
		t.Fatalf("unexpected handler: got %s, want %s", got.SourceHandler, handler)
	}
}

func TestWithInstanceID(t *testing.T) {
	packer := &Packer{
		Application: &identitypb.Identity{
			Name: "app",
			Key:  uuidpb.Generate(),
		},
	}

	got := packer.Pack(
		CommandA1,
		WithInstanceID("instance"),

		// We cannot have an instance ID without saying which handler the
		// instance is managed by.
		WithHandler(&identitypb.Identity{
			Name: "handler",
			Key:  uuidpb.Generate(),
		}),
	)

	if got.SourceInstanceId != "instance" {
		t.Fatalf("unexpected instance ID: got %s, want instance", got.SourceInstanceId)
	}
}

func TestWithCreatedAt(t *testing.T) {
	packer := &Packer{
		Application: &identitypb.Identity{
			Name: "app",
			Key:  uuidpb.Generate(),
		},
		Now: func() time.Time {
			t.Fatal("unexpected call")
			return time.Time{}
		},
	}

	want := time.Now().Add(-time.Hour)

	got := packer.Pack(
		CommandA1,
		WithCreatedAt(want),
	).CreatedAt.AsTime()

	if !got.Equal(want) {
		t.Fatalf("unexpected creation time: got %s, want %s", got, want)
	}
}

func TestWithScheduledFor(t *testing.T) {
	packer := &Packer{
		Application: &identitypb.Identity{
			Name: "app",
			Key:  uuidpb.Generate(),
		},
	}

	want := time.Now().Add(time.Hour)

	got := packer.Pack(
		TimeoutA1,
		WithScheduledFor(want),

		// We cannot have a "scheduled-for" time without saying which handler
		// and process instance produced the timeout message.
		WithHandler(&identitypb.Identity{
			Name: "handler",
			Key:  uuidpb.Generate(),
		}),
		WithInstanceID("instance"),
	).ScheduledFor.AsTime()

	if !got.Equal(want) {
		t.Fatalf("unexpected scheduled time: got %s, want %s", got, want)
	}
}
