package envelopepb_test

import (
	"math"
	"strings"
	"testing"

	. "github.com/dogmatiq/enginekit/internal/test"
	"github.com/dogmatiq/enginekit/protobuf/envelopepb"
	. "github.com/dogmatiq/enginekit/protobuf/envelopepb"
	identitypb "github.com/dogmatiq/enginekit/protobuf/identitypb"
	uuidpb "github.com/dogmatiq/enginekit/protobuf/uuidpb"
	anypb "google.golang.org/protobuf/types/known/anypb"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
)

func TestEnvelope_Validate(t *testing.T) {
	t.Parallel()

	t.Run("when the envelope is valid", func(t *testing.T) {
		t.Parallel()

		cases := []struct {
			Desc    string
			Subject *Envelope
		}{
			{
				"complete",
				newEnvelope(),
			},
			{
				"without source site",
				newEnvelope(func(e *Envelope) {
					e.GetHeader().GetSource().SetSite(nil)
				}),
			},
			{
				"without source handler",
				newEnvelope(func(e *Envelope) {
					e.GetHeader().GetSource().SetHandler(nil)
					e.GetHeader().GetSource().SetInstanceId("")
					e.GetBody().SetScheduledFor(nil)
				}),
			},
			{
				"without source instance ID",
				newEnvelope(func(e *Envelope) {
					e.GetHeader().GetSource().SetInstanceId("")
					e.GetBody().SetScheduledFor(nil)
				}),
			},
			{
				"without extensions",
				newEnvelope(func(e *Envelope) {
					e.GetBody().SetExtensions(nil)
				}),
			},
			{
				"with empty extension payload",
				newEnvelope(func(e *Envelope) {
					e.GetBody().SetExtensions([]*anypb.Any{{
						TypeUrl: "type.googleapis.com/example.Extension",
					}})
				}),
			},
			{
				"with empty baggage payload",
				newEnvelope(func(e *Envelope) {
					e.GetBody().SetBaggage([]*anypb.Any{{
						TypeUrl: "type.googleapis.com/example.Baggage",
					}})
				}),
			},
			{
				"without data",
				newEnvelope(func(e *Envelope) {
					e.GetBody().GetMessage().SetData(nil)
				}),
			},
		}

		for _, c := range cases {
			t.Run(c.Desc, func(t *testing.T) {
				t.Parallel()

				if err := c.Subject.Validate(); err != nil {
					t.Fatal(err)
				}
			})
		}
	})

	t.Run("when the envelope is invalid", func(t *testing.T) {
		t.Parallel()

		cases := []struct {
			Desc    string
			Subject *Envelope
			Expect  string
		}{
			{
				"empty",
				&Envelope{},
				"invalid header: must not be nil",
			},
			{
				"invalid message ID",
				newEnvelope(func(e *Envelope) {
					e.GetBody().SetMessageId(&uuidpb.UUID{})
				}),
				"invalid body: invalid message ID (00000000-0000-0000-0000-000000000000): UUID must use version 4 or 5",
			},
			{
				"invalid causation ID",
				newEnvelope(func(e *Envelope) {
					e.GetHeader().SetCausationId(&uuidpb.UUID{})
				}),
				"invalid header: invalid causation ID (00000000-0000-0000-0000-000000000000): UUID must use version 4 or 5",
			},
			{
				"invalid correlation ID",
				newEnvelope(func(e *Envelope) {
					e.GetHeader().SetCorrelationId(&uuidpb.UUID{})
				}),
				"invalid header: invalid correlation ID (00000000-0000-0000-0000-000000000000): UUID must use version 4 or 5",
			},
			{
				"invalid source site",
				newEnvelope(func(e *Envelope) {
					e.GetHeader().GetSource().SetSite(&identitypb.Identity{})
				}),
				"invalid header: invalid source: invalid site (00000000-0000-0000-0000-000000000000 ?): invalid name: must be between 1 and 255 bytes",
			},
			{
				"invalid source application",
				newEnvelope(func(e *Envelope) {
					e.GetHeader().GetSource().SetApplication(&identitypb.Identity{})
				}),
				"invalid header: invalid source: invalid application (00000000-0000-0000-0000-000000000000 ?): invalid name: must be between 1 and 255 bytes",
			},
			{
				"invalid source handler",
				newEnvelope(func(e *Envelope) {
					e.GetHeader().GetSource().SetHandler(&identitypb.Identity{})
				}),
				"invalid header: invalid source: invalid handler (00000000-0000-0000-0000-000000000000 ?): invalid name: must be between 1 and 255 bytes",
			},
			{
				"source instance ID without source handler",
				newEnvelope(func(e *Envelope) {
					e.GetHeader().GetSource().SetHandler(nil)
				}),
				"invalid header: invalid source: invalid instance ID: must not be specified without a handler",
			},
			{
				"scheduled-for time without source handler",
				newEnvelope(func(e *Envelope) {
					e.GetHeader().GetSource().SetHandler(nil)
					e.GetHeader().GetSource().SetInstanceId("")
				}),
				"invalid body: invalid scheduled-for time: must not be specified without a source handler and instance ID",
			},
			{
				"invalid created-at time",
				newEnvelope(func(e *Envelope) {
					e.GetBody().SetCreatedAt(nil)
				}),
				"invalid body: invalid created-at time: proto: invalid nil Timestamp",
			},
			{
				"invalid scheduled-for time",
				newEnvelope(func(e *Envelope) {
					e.GetBody().SetScheduledFor(&timestamppb.Timestamp{
						Seconds: math.MaxInt64,
						Nanos:   math.MaxInt32,
					})
				}),
				"invalid body: invalid scheduled-for time: proto: timestamp (seconds:9223372036854775807 nanos:2147483647) after 9999-12-31",
			},
			{
				"invalid description",
				newEnvelope(func(e *Envelope) {
					e.GetBody().GetMessage().SetDescription("")
				}),
				"invalid body: invalid message: invalid description: must not be empty",
			},
			{
				"empty header extension",
				newEnvelope(func(e *Envelope) {
					e.GetHeader().SetExtensions([]*anypb.Any{{}})
				}),
				"invalid header: invalid extensions at index 0: type URL must not be empty",
			},
			{
				"empty header baggage",
				newEnvelope(func(e *Envelope) {
					e.GetHeader().SetBaggage([]*anypb.Any{nil})
				}),
				"invalid header: invalid baggage at index 0: type URL must not be empty",
			},
			{
				"empty body extension",
				newEnvelope(func(e *Envelope) {
					e.GetBody().SetExtensions([]*anypb.Any{{}})
				}),
				"invalid body: invalid extensions at index 0: type URL must not be empty",
			},
			{
				"empty body baggage",
				newEnvelope(func(e *Envelope) {
					e.GetBody().SetBaggage([]*anypb.Any{{}})
				}),
				"invalid body: invalid baggage at index 0: type URL must not be empty",
			},
			{
				"without type ID",
				newEnvelope(func(e *Envelope) {
					e.GetBody().GetMessage().SetTypeId(nil)
				}),
				"invalid body: invalid message: invalid type ID (00000000-0000-0000-0000-000000000000): UUID must use version 4 or 5",
			},
		}

		for _, c := range cases {
			t.Run(c.Desc, func(t *testing.T) {
				t.Parallel()

				err := c.Subject.Validate()
				if err == nil {
					t.Fatal("expected an error")
				}

				// defeat the protobuf's random injection of of whitespace into
				// error messages.
				actual := strings.ReplaceAll(err.Error(), "\u00a0", " ")
				actual = strings.ReplaceAll(actual, "  ", " ")

				if actual != c.Expect {
					t.Fatalf("got %q, want %q", err, c.Expect)
				}
			})
		}
	})
}

func TestWireCompatibility(t *testing.T) {
	t.Parallel()

	cases := []struct {
		Desc             string
		HeaderExtensions []*anypb.Any
		HeaderBaggage    []*anypb.Any
		BodyExtensions   []*anypb.Any
		BodyBaggage      []*anypb.Any
		WantExtensions   []*anypb.Any
		WantBaggage      []*anypb.Any
	}{
		{
			Desc: "header only",
			HeaderExtensions: []*anypb.Any{
				{
					TypeUrl: "type.googleapis.com/example.ExtensionA",
					Value:   []byte("header-a"),
				},
			},
			HeaderBaggage: []*anypb.Any{
				{
					TypeUrl: "type.googleapis.com/example.BaggageA",
					Value:   []byte("header-a"),
				},
			},
			WantExtensions: []*anypb.Any{
				{
					TypeUrl: "type.googleapis.com/example.ExtensionA",
					Value:   []byte("header-a"),
				},
			},
			WantBaggage: []*anypb.Any{
				{
					TypeUrl: "type.googleapis.com/example.BaggageA",
					Value:   []byte("header-a"),
				},
			},
		},
		{
			Desc: "body only",
			BodyExtensions: []*anypb.Any{
				{
					TypeUrl: "type.googleapis.com/example.ExtensionA",
					Value:   []byte("body-a"),
				},
			},
			BodyBaggage: []*anypb.Any{
				{
					TypeUrl: "type.googleapis.com/example.BaggageA",
					Value:   []byte("body-a"),
				},
			},
			WantExtensions: []*anypb.Any{
				{
					TypeUrl: "type.googleapis.com/example.ExtensionA",
					Value:   []byte("body-a"),
				},
			},
			WantBaggage: []*anypb.Any{
				{
					TypeUrl: "type.googleapis.com/example.BaggageA",
					Value:   []byte("body-a"),
				},
			},
		},
		{
			Desc: "split with overrides",
			HeaderExtensions: []*anypb.Any{
				{
					TypeUrl: "type.googleapis.com/example.ExtensionA",
					Value:   []byte("header-a"),
				},
				{
					TypeUrl: "type.googleapis.com/example.ExtensionB",
					Value:   []byte("header-b"),
				},
			},
			HeaderBaggage: []*anypb.Any{
				{
					TypeUrl: "type.googleapis.com/example.BaggageA",
					Value:   []byte("header-a"),
				},
				{
					TypeUrl: "type.googleapis.com/example.BaggageB",
					Value:   []byte("header-b"),
				},
			},
			BodyExtensions: []*anypb.Any{
				{
					TypeUrl: "type.googleapis.com/example.ExtensionB",
					Value:   []byte("body-b"),
				},
				{
					TypeUrl: "type.googleapis.com/example.ExtensionC",
					Value:   []byte("body-c"),
				},
			},
			BodyBaggage: []*anypb.Any{
				{
					TypeUrl: "type.googleapis.com/example.BaggageB",
					Value:   []byte("body-b"),
				},
				{
					TypeUrl: "type.googleapis.com/example.BaggageC",
					Value:   []byte("body-c"),
				},
			},
			WantExtensions: []*anypb.Any{
				{
					TypeUrl: "type.googleapis.com/example.ExtensionA",
					Value:   []byte("header-a"),
				},
				{
					TypeUrl: "type.googleapis.com/example.ExtensionB",
					Value:   []byte("body-b"),
				},
				{
					TypeUrl: "type.googleapis.com/example.ExtensionC",
					Value:   []byte("body-c"),
				},
			},
			WantBaggage: []*anypb.Any{
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
		},
	}

	for _, c := range cases {
		t.Run(c.Desc, func(t *testing.T) {
			t.Parallel()

			env := newEnvelope(func(e *Envelope) {
				e.GetHeader().SetExtensions(c.HeaderExtensions)
				e.GetHeader().SetBaggage(c.HeaderBaggage)
				e.GetBody().SetExtensions(c.BodyExtensions)
				e.GetBody().SetBaggage(c.BodyBaggage)
			})

			expectEnvelopeWireCompatibleExtensions(
				t,
				env,
				c.WantExtensions,
				c.WantBaggage,
			)

			multi := NewMultiEnvelopeBuilder().
				WithHeader(env.GetHeader()).
				WithBodies([]*Body{env.GetBody()}).
				Build()

			expectSingleBodyMultiEnvelopeWireCompatibleExtensions(
				t,
				multi,
				c.WantExtensions,
				c.WantBaggage,
			)
		})
	}
}

func expectEnvelopeWireCompatibleExtensions(
	t *testing.T,
	env *Envelope,
	wantExtensions, wantBaggage []*anypb.Any,
) {
	t.Helper()

	data, err := env.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}

	var got MultiEnvelope
	if err := got.UnmarshalBinary(data); err != nil {
		t.Fatal(err)
	}

	expectSingleBodyMultiEnvelopeEffectiveExtensions(
		t,
		&got,
		wantExtensions,
		wantBaggage,
	)
}

func expectSingleBodyMultiEnvelopeWireCompatibleExtensions(
	t *testing.T,
	env *MultiEnvelope,
	wantExtensions, wantBaggage []*anypb.Any,
) {
	t.Helper()

	data, err := env.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}

	var got Envelope
	if err := got.UnmarshalBinary(data); err != nil {
		t.Fatal(err)
	}

	expectEffectiveExtensionsForTest(
		t,
		got.GetHeader().GetExtensions(),
		got.GetBody().GetExtensions(),
		got.GetHeader().GetBaggage(),
		got.GetBody().GetBaggage(),
		wantExtensions,
		wantBaggage,
	)
}

func expectSingleBodyMultiEnvelopeEffectiveExtensions(
	t *testing.T,
	env *MultiEnvelope,
	wantExtensions, wantBaggage []*anypb.Any,
) {
	t.Helper()

	bodies := env.GetBodies()
	if len(bodies) != 1 {
		t.Fatalf("unexpected body count: got %d, want 1", len(bodies))
	}

	expectEffectiveExtensionsForTest(
		t,
		env.GetHeader().GetExtensions(),
		bodies[0].GetExtensions(),
		env.GetHeader().GetBaggage(),
		bodies[0].GetBaggage(),
		wantExtensions,
		wantBaggage,
	)
}

func expectEffectiveExtensionsForTest(
	t *testing.T,
	headerExtensions,
	bodyExtensions,
	headerBaggage,
	bodyBaggage,
	wantExtensions,
	wantBaggage []*anypb.Any,
) {
	t.Helper()

	Expect(
		t,
		"unexpected effective extensions",
		mergedAnyValuesByTypeURLForTest(
			headerExtensions,
			bodyExtensions,
		),
		wantExtensions,
	)

	Expect(
		t,
		"unexpected effective baggage",
		mergedAnyValuesByTypeURLForTest(
			headerBaggage,
			bodyBaggage,
		),
		wantBaggage,
	)
}

func mergedAnyValuesByTypeURLForTest(header, body []*anypb.Any) []*anypb.Any {
	if len(header) == 0 {
		return body
	}

	if len(body) == 0 {
		return header
	}

	values := make([]*anypb.Any, 0, len(header)+len(body))

	for _, v := range header {
		if containsAnyValueWithTypeURLForTest(body, v.GetTypeUrl()) {
			continue
		}

		values = append(values, v)
	}

	for _, v := range body {
		values = append(values, v)
	}

	return values
}

func containsAnyValueWithTypeURLForTest(values []*anypb.Any, typeURL string) bool {
	for _, v := range values {
		if v.GetTypeUrl() == typeURL {
			return true
		}
	}

	return false
}

func newEnvelope(modifiers ...func(*Envelope)) *envelopepb.Envelope {
	env := NewEnvelopeBuilder().
		WithHeader(
			NewHeaderBuilder().
				WithCausationId(uuidpb.Generate()).
				WithCorrelationId(uuidpb.Generate()).
				WithSource(NewSourceBuilder().
					WithSite(identitypb.New("<site-name>", uuidpb.Generate())).
					WithApplication(identitypb.New("<app-name>", uuidpb.Generate())).
					WithHandler(identitypb.New("<handler-name>", uuidpb.Generate())).
					WithInstanceId("<instance>").
					Build()).
				Build(),
		).
		WithBody(
			NewBodyBuilder().
				WithMessageId(uuidpb.Generate()).
				WithCreatedAt(timestamppb.Now()).
				WithScheduledFor(timestamppb.Now()).
				WithMessage(NewMessageBuilder().
					WithDescription("<description>").
					WithTypeId(uuidpb.Generate()).
					WithData([]byte("<data>")).
					Build(),
				).
				WithExtensions([]*anypb.Any{
					{
						TypeUrl: "type.googleapis.com/example.envelope.v1.Extension",
						Value:   []byte("<extension-value>"),
					},
				}).
				WithBaggage([]*anypb.Any{
					{
						TypeUrl: "type.googleapis.com/example.envelope.v1.Baggage",
						Value:   []byte("<baggage-value>"),
					},
				}).
				Build()).
		Build()

	for _, fn := range modifiers {
		fn(env)
	}

	return env
}
