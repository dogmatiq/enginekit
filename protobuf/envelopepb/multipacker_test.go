package envelopepb_test

import (
	"testing"
	"time"

	. "github.com/dogmatiq/enginekit/enginetest/stubs"
	. "github.com/dogmatiq/enginekit/internal/test"
	. "github.com/dogmatiq/enginekit/protobuf/envelopepb"
	"github.com/dogmatiq/enginekit/protobuf/identitypb"
	"github.com/dogmatiq/enginekit/protobuf/uuidpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestPacker_CausedBy(t *testing.T) {
	t.Run("it panics if the cause is nil", func(t *testing.T) {
		packer := &Packer{
			Application: &identitypb.Identity{
				Name: "app",
				Key:  uuidpb.Generate(),
			},
		}

		ExpectPanic(
			t,
			"cause must not be nil",
			func() {
				packer.CausedBy(nil)
			},
		)
	})

	t.Run("it panics if the cause envelope is invalid", func(t *testing.T) {
		packer := &Packer{
			Application: &identitypb.Identity{
				Name: "app",
				Key:  uuidpb.Generate(),
			},
		}

		ExpectPanic(
			t,
			"invalid cause envelope: invalid header: must not be nil",
			func() {
				packer.CausedBy(&Envelope{})
			},
		)
	})

	t.Run("it snapshots the packer state", func(t *testing.T) {
		causeID := uuidpb.Generate()
		childID := uuidpb.Generate()
		now := time.Now()

		ids := []*uuidpb.UUID{causeID, childID}
		site := &identitypb.Identity{
			Name: "site",
			Key:  uuidpb.Generate(),
		}
		application := &identitypb.Identity{
			Name: "app",
			Key:  uuidpb.Generate(),
		}

		packer := &Packer{
			Site:        site,
			Application: application,
			GenerateID: func() *uuidpb.UUID {
				id := ids[0]
				ids = ids[1:]
				return id
			},
			Now: func() time.Time {
				return now
			},
		}

		cause := packer.Pack(CommandA1)
		mp := packer.CausedBy(cause)

		packer.Site = &identitypb.Identity{Name: "other-site", Key: uuidpb.Generate()}
		packer.Application = &identitypb.Identity{Name: "other-app", Key: uuidpb.Generate()}
		packer.GenerateID = func() *uuidpb.UUID {
			return uuidpb.Generate()
		}
		packer.Now = func() time.Time {
			return now.Add(24 * time.Hour)
		}

		mp.Pack(CommandA2)

		got, ok := mp.Seal()
		if !ok {
			t.Fatal("expected a multi-envelope")
		}

		if got == nil {
			t.Fatal("expected a non-nil multi-envelope")
		}

		if !got.Header.Source.Site.Equal(site) {
			t.Fatalf("unexpected site: got %s, want %s", got.Header.Source.Site, site)
		}

		if !got.Header.Source.Application.Equal(application) {
			t.Fatalf("unexpected application: got %s, want %s", got.Header.Source.Application, application)
		}

		if !got.Bodies[0].MessageId.Equal(childID) {
			t.Fatalf("unexpected child message ID: got %s, want %s", got.Bodies[0].MessageId, childID)
		}

		if !got.Bodies[0].CreatedAt.AsTime().Equal(now) {
			t.Fatalf("unexpected created-at: got %s, want %s", got.Bodies[0].CreatedAt.AsTime(), now)
		}
	})

	t.Run("it applies source options", func(t *testing.T) {
		packer := &Packer{
			Application: &identitypb.Identity{
				Name: "app",
				Key:  uuidpb.Generate(),
			},
		}

		handlerA := &identitypb.Identity{
			Name: "handler-a",
			Key:  uuidpb.Generate(),
		}
		handlerB := &identitypb.Identity{
			Name: "handler-b",
			Key:  uuidpb.Generate(),
		}

		cause := packer.Pack(CommandA1)

		left := packer.CausedBy(
			cause,
			WithHandler(handlerA),
			WithInstanceID("instance-a"),
		)
		left.Pack(TimeoutA1, WithScheduledFor(time.Now().Add(time.Hour)))

		right := packer.CausedBy(
			cause,
			WithHandler(handlerB),
			WithInstanceID("instance-b"),
		)
		right.Pack(TimeoutA2, WithScheduledFor(time.Now().Add(2*time.Hour)))

		leftEnv, ok := left.Seal()
		if !ok {
			t.Fatal("expected left multi-envelope")
		}

		rightEnv, ok := right.Seal()
		if !ok {
			t.Fatal("expected right multi-envelope")
		}

		if !leftEnv.Header.Source.Handler.Equal(handlerA) {
			t.Fatalf("unexpected left handler: got %s, want %s", leftEnv.Header.Source.Handler, handlerA)
		}

		if leftEnv.Header.Source.InstanceId != "instance-a" {
			t.Fatalf("unexpected left instance ID: got %q, want %q", leftEnv.Header.Source.InstanceId, "instance-a")
		}

		if !rightEnv.Header.Source.Handler.Equal(handlerB) {
			t.Fatalf("unexpected right handler: got %s, want %s", rightEnv.Header.Source.Handler, handlerB)
		}

		if rightEnv.Header.Source.InstanceId != "instance-b" {
			t.Fatalf("unexpected right instance ID: got %q, want %q", rightEnv.Header.Source.InstanceId, "instance-b")
		}
	})
}

func TestMultiPacker_Pack(t *testing.T) {
	t.Run("it panics if the body is invalid", func(t *testing.T) {
		packer := &Packer{
			Application: &identitypb.Identity{
				Name: "app",
				Key:  uuidpb.Generate(),
			},
		}

		cause := packer.Pack(CommandA1)
		mp := packer.CausedBy(cause)

		ExpectPanic(
			t,
			"invalid body: invalid scheduled-for time: must not be specified without a source handler and instance ID",
			func() {
				mp.Pack(TimeoutA1, WithScheduledFor(time.Now().Add(time.Hour)))
			},
		)
	})
}

func TestMultiPacker_Seal(t *testing.T) {
	t.Run("it returns the packed messages in insertion order", func(t *testing.T) {
		id0 := uuidpb.Generate()
		id1 := uuidpb.Generate()
		id2 := uuidpb.Generate()
		now := time.Now()

		ids := []*uuidpb.UUID{id0, id1, id2}
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
				id := ids[0]
				ids = ids[1:]
				return id
			},
			Now: func() time.Time {
				return now
			},
		}

		cause := packer.Pack(CommandA1)

		mp := packer.CausedBy(cause)
		mp.Pack(CommandA1, WithIdempotencyKey("key-a1"))
		mp.Pack(CommandA2)

		got, ok := mp.Seal()
		if !ok {
			t.Fatal("expected a multi-envelope")
		}

		if got == nil {
			t.Fatal("expected a non-nil multi-envelope")
		}

		if err := got.Validate(); err != nil {
			t.Fatalf("multi-packer produced an invalid envelope: %v", err)
		}

		want := &MultiEnvelope{
			Header: &Header{
				CausationId:   cause.Body.MessageId,
				CorrelationId: cause.Header.CorrelationId,
				Source: &Source{
					Site:        packer.Site,
					Application: packer.Application,
				},
			},
			Bodies: []*Body{
				{
					MessageId:      id1,
					CreatedAt:      timestamppb.New(now),
					IdempotencyKey: "key-a1",
					Message: &Message{
						Description: `command(stubs.TypeA:A1, valid)`,
						TypeId:      uuidpb.MustParse(MessageTypeID[*CommandStub[TypeA]]()),
						Data:        []byte(`{"content":"A1"}`),
					},
				},
				{
					MessageId: id2,
					CreatedAt: timestamppb.New(now),
					Message: &Message{
						Description: `command(stubs.TypeA:A2, valid)`,
						TypeId:      uuidpb.MustParse(MessageTypeID[*CommandStub[TypeA]]()),
						Data:        []byte(`{"content":"A2"}`),
					},
				},
			},
		}

		Expect(
			t,
			"unexpected multi-envelope",
			got,
			want,
		)
	})

	t.Run("it returns no envelope when empty", func(t *testing.T) {
		packer := &Packer{
			Application: &identitypb.Identity{
				Name: "app",
				Key:  uuidpb.Generate(),
			},
		}

		cause := packer.Pack(CommandA1)
		got, ok := packer.CausedBy(cause).Seal()

		if ok {
			t.Fatal("expected no multi-envelope")
		}

		if got != nil {
			t.Fatalf("unexpected multi-envelope: got %#v, want nil", got)
		}
	})

	t.Run("it panics after sealing", func(t *testing.T) {
		packer := &Packer{
			Application: &identitypb.Identity{
				Name: "app",
				Key:  uuidpb.Generate(),
			},
		}

		cause := packer.Pack(CommandA1)

		t.Run("after empty seal", func(t *testing.T) {
			mp := packer.CausedBy(cause)
			_, _ = mp.Seal()

			ExpectPanic(
				t,
				"already sealed",
				func() {
					mp.Pack(CommandA1)
				},
			)

			ExpectPanic(
				t,
				"already sealed",
				func() {
					mp.Seal()
				},
			)
		})

		t.Run("after non-empty seal", func(t *testing.T) {
			mp := packer.CausedBy(cause)
			mp.Pack(CommandA1)
			_, _ = mp.Seal()

			ExpectPanic(
				t,
				"already sealed",
				func() {
					mp.Pack(CommandA2)
				},
			)

			ExpectPanic(
				t,
				"already sealed",
				func() {
					mp.Seal()
				},
			)
		})
	})
}
