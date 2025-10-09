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

	gotMessage, err := packer.Unpack(want)
	if err != nil {
		t.Fatal(err)
	}

	Expect(
		t,
		"unexpected message",
		gotMessage,
		dogma.Message(CommandA1),
	)
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
