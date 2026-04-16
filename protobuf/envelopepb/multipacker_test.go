package envelopepb_test

import (
	"testing"
	"time"

	. "github.com/dogmatiq/enginekit/enginetest/stubs"
	. "github.com/dogmatiq/enginekit/internal/test"
	. "github.com/dogmatiq/enginekit/protobuf/envelopepb"
	"github.com/dogmatiq/enginekit/protobuf/identitypb"
	"github.com/dogmatiq/enginekit/protobuf/uuidpb"
	anypb "google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func TestPacker_CausedBy(t *testing.T) {
	t.Run("it panics if the handler is nil", func(t *testing.T) {
		packer := &Packer{
			Application: &identitypb.Identity{
				Name: "app",
				Key:  uuidpb.Generate(),
			},
		}

		cause := packer.PackCommand(CommandA1)

		ExpectPanic(
			t,
			"handler must not be nil",
			func() {
				packer.PackEffects(cause, nil)
			},
		)
	})

	t.Run("it panics if the cause is nil", func(t *testing.T) {
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

		ExpectPanic(
			t,
			"cause must not be nil",
			func() {
				packer.PackEffects(nil, handler)
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

		handler := &identitypb.Identity{
			Name: "handler",
			Key:  uuidpb.Generate(),
		}

		ExpectPanic(
			t,
			"invalid cause envelope: invalid header: must not be nil",
			func() {
				packer.PackEffects(&Envelope{}, handler)
			},
		)
	})

	t.Run("it snapshots the packer state", func(t *testing.T) {
		causeID := uuidpb.Generate()
		packedID := uuidpb.Generate()
		now := time.Now()

		ids := []*uuidpb.UUID{causeID, packedID}
		site := &identitypb.Identity{
			Name: "site",
			Key:  uuidpb.Generate(),
		}
		application := &identitypb.Identity{
			Name: "app",
			Key:  uuidpb.Generate(),
		}
		handler := &identitypb.Identity{
			Name: "handler",
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

		cause := packer.PackCommand(CommandA1)
		ep := packer.PackEffects(cause, handler)

		packer.Site = &identitypb.Identity{Name: "other-site", Key: uuidpb.Generate()}
		packer.Application = &identitypb.Identity{Name: "other-app", Key: uuidpb.Generate()}
		packer.GenerateID = func() *uuidpb.UUID {
			return uuidpb.Generate()
		}
		packer.Now = func() time.Time {
			return now.Add(24 * time.Hour)
		}

		ep.PackCommand(CommandA2)

		got, ok := ep.Seal()
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

		if !got.Header.Source.Handler.Equal(handler) {
			t.Fatalf("unexpected handler: got %s, want %s", got.Header.Source.Handler, handler)
		}

		if !got.Bodies[0].MessageId.Equal(packedID) {
			t.Fatalf("unexpected packed message ID: got %s, want %s", got.Bodies[0].MessageId, packedID)
		}

		if !got.Bodies[0].CreatedAt.AsTime().Equal(now) {
			t.Fatalf("unexpected created-at: got %s, want %s", got.Bodies[0].CreatedAt.AsTime(), now)
		}
	})

	t.Run("it applies effect packer options to the shared source", func(t *testing.T) {
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

		cause := packer.PackCommand(CommandA1)

		left := packer.PackEffects(
			cause,
			handlerA,
			WithInstanceID("instance-a"),
		)
		left.PackTimeout(TimeoutA1, WithScheduledFor(time.Now().Add(time.Hour)))

		right := packer.PackEffects(
			cause,
			handlerB,
			WithInstanceID("instance-b"),
		)
		right.PackTimeout(TimeoutA2, WithScheduledFor(time.Now().Add(2*time.Hour)))

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

		leftScheduledFor := leftEnv.Bodies[0].ScheduledFor.AsTime()
		if leftScheduledFor.IsZero() {
			t.Fatal("expected left scheduled-for time")
		}

		if !rightEnv.Header.Source.Handler.Equal(handlerB) {
			t.Fatalf("unexpected right handler: got %s, want %s", rightEnv.Header.Source.Handler, handlerB)
		}

		if rightEnv.Header.Source.InstanceId != "instance-b" {
			t.Fatalf("unexpected right instance ID: got %q, want %q", rightEnv.Header.Source.InstanceId, "instance-b")
		}

		rightScheduledFor := rightEnv.Bodies[0].ScheduledFor.AsTime()
		if rightScheduledFor.IsZero() {
			t.Fatal("expected right scheduled-for time")
		}
	})

	t.Run("it applies extension and baggage options to the header", func(t *testing.T) {
		packer := &Packer{
			Application: &identitypb.Identity{
				Name: "app",
				Key:  uuidpb.Generate(),
			},
		}

		extension := wrapperspb.String("extension")
		wantExtension, err := anypb.New(extension)
		if err != nil {
			t.Fatal(err)
		}

		baggage := wrapperspb.String("baggage")
		wantBaggage, err := anypb.New(baggage)
		if err != nil {
			t.Fatal(err)
		}

		handler := &identitypb.Identity{
			Name: "handler",
			Key:  uuidpb.Generate(),
		}

		cause := packer.PackCommand(CommandA1)
		ep := packer.PackEffects(
			cause,
			handler,
			WithExtension(extension),
			WithBaggage(baggage),
		)
		ep.PackCommand(CommandA2)

		got, ok := ep.Seal()
		if !ok {
			t.Fatal("expected a multi-envelope")
		}

		Expect(
			t,
			"unexpected header extensions",
			got.Header.Extensions,
			[]*anypb.Any{wantExtension},
		)

		Expect(
			t,
			"unexpected header baggage",
			got.Header.Baggage,
			[]*anypb.Any{wantBaggage},
		)

		if got.Bodies[0].Extensions != nil || got.Bodies[0].Baggage != nil {
			t.Fatalf(
				"unexpected body extension state: got extensions %#v, baggage %#v, want nil",
				got.Bodies[0].Extensions,
				got.Bodies[0].Baggage,
			)
		}
	})

	t.Run("it keeps only the last baggage value for a type URL when propagating from the cause", func(t *testing.T) {
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

		cause := packer.PackCommand(CommandA1)
		cause.Header.Baggage = []*anypb.Any{
			{
				TypeUrl: "type.googleapis.com/example.BaggageA",
				Value:   []byte("header-a-1"),
			},
			{
				TypeUrl: "type.googleapis.com/example.BaggageA",
				Value:   []byte("header-a-2"),
			},
			{
				TypeUrl: "type.googleapis.com/example.BaggageB",
				Value:   []byte("header-b"),
			},
		}

		cause.Body.Baggage = []*anypb.Any{
			{
				TypeUrl: "type.googleapis.com/example.BaggageB",
				Value:   []byte("body-b-1"),
			},
			{
				TypeUrl: "type.googleapis.com/example.BaggageB",
				Value:   []byte("body-b-2"),
			},
			{
				TypeUrl: "type.googleapis.com/example.BaggageC",
				Value:   []byte("body-c"),
			},
		}

		ep := packer.PackEffects(cause, handler)
		ep.PackCommand(CommandA2)

		got, ok := ep.Seal()
		if !ok {
			t.Fatal("expected a multi-envelope")
		}

		Expect(
			t,
			"unexpected propagated baggage",
			got.Header.Baggage,
			[]*anypb.Any{
				{
					TypeUrl: "type.googleapis.com/example.BaggageA",
					Value:   []byte("header-a-2"),
				},
				{
					TypeUrl: "type.googleapis.com/example.BaggageB",
					Value:   []byte("body-b-2"),
				},
				{
					TypeUrl: "type.googleapis.com/example.BaggageC",
					Value:   []byte("body-c"),
				},
			},
		)
	})

	t.Run("it propagates baggage but not extensions from the cause", func(t *testing.T) {
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

		cause := packer.PackCommand(CommandA1)
		cause.Header.Extensions = []*anypb.Any{
			{
				TypeUrl: "type.googleapis.com/example.ExtensionA",
				Value:   []byte("header-a"),
			},
			{
				TypeUrl: "type.googleapis.com/example.ExtensionB",
				Value:   []byte("header-b"),
			},
		}

		cause.Header.Baggage = []*anypb.Any{
			{
				TypeUrl: "type.googleapis.com/example.BaggageA",
				Value:   []byte("header-a"),
			},
			{
				TypeUrl: "type.googleapis.com/example.BaggageB",
				Value:   []byte("header-b"),
			},
		}

		cause.Body.Extensions = []*anypb.Any{
			{
				TypeUrl: "type.googleapis.com/example.ExtensionB",
				Value:   []byte("body-b"),
			},
			{
				TypeUrl: "type.googleapis.com/example.ExtensionC",
				Value:   []byte("body-c"),
			},
		}

		cause.Body.Baggage = []*anypb.Any{
			{
				TypeUrl: "type.googleapis.com/example.BaggageB",
				Value:   []byte("body-b"),
			},
			{
				TypeUrl: "type.googleapis.com/example.BaggageC",
				Value:   []byte("body-c"),
			},
		}

		ep := packer.PackEffects(cause, handler)
		ep.PackCommand(CommandA2)

		got, ok := ep.Seal()
		if !ok {
			t.Fatal("expected a multi-envelope")
		}

		expectSingleBodyMultiEnvelopeWireCompatibleExtensions(
			t,
			got,
			nil,
			[]*anypb.Any{
				{
					TypeUrl: "type.googleapis.com/example.BaggageA",
					Value:   []byte("header-a"),
				},
				{
					TypeUrl: "type.googleapis.com/example.BaggageB",
					Value:   []byte("body-b"),
				},
				{
					TypeUrl: "type.googleapis.com/example.BaggageC",
					Value:   []byte("body-c"),
				},
			},
		)
	})
}

func TestEffectPacker_Pack(t *testing.T) {
	t.Run("it panics if the body is invalid", func(t *testing.T) {
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

		cause := packer.PackCommand(CommandA1)
		ep := packer.PackEffects(cause, handler)

		ExpectPanic(
			t,
			"invalid body: invalid scheduled-for time: must not be specified without a source handler and instance ID",
			func() {
				ep.PackTimeout(TimeoutA1, WithScheduledFor(time.Now().Add(time.Hour)))
			},
		)
	})

	t.Run("it applies baggage options to the body", func(t *testing.T) {
		packer := &Packer{
			Application: &identitypb.Identity{
				Name: "app",
				Key:  uuidpb.Generate(),
			},
		}

		baggage := wrapperspb.String("baggage")
		want, err := anypb.New(baggage)
		if err != nil {
			t.Fatal(err)
		}

		handler := &identitypb.Identity{
			Name: "handler",
			Key:  uuidpb.Generate(),
		}

		cause := packer.PackCommand(CommandA1)
		ep := packer.PackEffects(cause, handler)
		ep.PackCommand(CommandA2, WithBaggage(baggage))

		got, ok := ep.Seal()
		if !ok {
			t.Fatal("expected a multi-envelope")
		}

		Expect(
			t,
			"unexpected baggage values",
			got.Bodies[0].Baggage,
			[]*anypb.Any{want},
		)
	})
}

func TestEffectPacker_Seal(t *testing.T) {
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

		handler := &identitypb.Identity{
			Name: "handler",
			Key:  uuidpb.Generate(),
		}

		cause := packer.PackCommand(CommandA1)

		ep := packer.PackEffects(cause, handler)
		ep.PackCommand(CommandA1)
		ep.PackCommand(CommandA2)

		got, ok := ep.Seal()
		if !ok {
			t.Fatal("expected a multi-envelope")
		}

		if got == nil {
			t.Fatal("expected a non-nil multi-envelope")
		}

		if err := got.Validate(); err != nil {
			t.Fatalf("effect packer produced an invalid envelope: %v", err)
		}

		want := &MultiEnvelope{
			Header: &Header{
				CausationId:   cause.Body.MessageId,
				CorrelationId: cause.Header.CorrelationId,
				Source: &Source{
					Site:        packer.Site,
					Application: packer.Application,
					Handler:     handler,
				},
			},
			Bodies: []*Body{
				{
					MessageId: id1,
					CreatedAt: timestamppb.New(now),
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

		handler := &identitypb.Identity{
			Name: "handler",
			Key:  uuidpb.Generate(),
		}

		cause := packer.PackCommand(CommandA1)
		got, ok := packer.PackEffects(cause, handler).Seal()

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

		handler := &identitypb.Identity{
			Name: "handler",
			Key:  uuidpb.Generate(),
		}

		cause := packer.PackCommand(CommandA1)

		t.Run("after empty seal", func(t *testing.T) {
			ep := packer.PackEffects(cause, handler)
			_, _ = ep.Seal()

			ExpectPanic(
				t,
				"already sealed",
				func() {
					ep.PackCommand(CommandA1)
				},
			)

			ExpectPanic(
				t,
				"already sealed",
				func() {
					ep.Seal()
				},
			)
		})

		t.Run("after non-empty seal", func(t *testing.T) {
			ep := packer.PackEffects(cause, handler)
			ep.PackCommand(CommandA1)
			_, _ = ep.Seal()

			ExpectPanic(
				t,
				"already sealed",
				func() {
					ep.PackCommand(CommandA2)
				},
			)

			ExpectPanic(
				t,
				"already sealed",
				func() {
					ep.Seal()
				},
			)
		})
	})
}
