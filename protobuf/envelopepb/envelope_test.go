package envelopepb_test

import (
	"math"
	"strings"
	"testing"

	"github.com/dogmatiq/enginekit/protobuf/envelopepb"
	. "github.com/dogmatiq/enginekit/protobuf/envelopepb"
	identitypb "github.com/dogmatiq/enginekit/protobuf/identitypb"
	uuidpb "github.com/dogmatiq/enginekit/protobuf/uuidpb"
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
					e.Header.Source.Site = nil
				}),
			},
			{
				"without source handler",
				newEnvelope(func(e *Envelope) {
					e.Header.Source.Handler = nil
					e.Header.Source.InstanceId = ""
					e.Body.ScheduledFor = nil
				}),
			},
			{
				"without source instance ID",
				newEnvelope(func(e *Envelope) {
					e.Header.Source.InstanceId = ""
					e.Body.ScheduledFor = nil
				}),
			},
			{
				"without attributes",
				newEnvelope(func(e *Envelope) {
					e.Body.Extensions.Attributes = nil
				}),
			},
			{
				"without data",
				newEnvelope(func(e *Envelope) {
					e.Body.Message.Data = nil
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
					e.Body.MessageId = &uuidpb.UUID{}
				}),
				"invalid body: invalid message ID (00000000-0000-0000-0000-000000000000): UUID must use version 4 or 5",
			},
			{
				"invalid causation ID",
				newEnvelope(func(e *Envelope) {
					e.Header.CausationId = &uuidpb.UUID{}
				}),
				"invalid header: invalid causation ID (00000000-0000-0000-0000-000000000000): UUID must use version 4 or 5",
			},
			{
				"invalid correlation ID",
				newEnvelope(func(e *Envelope) {
					e.Header.CorrelationId = &uuidpb.UUID{}
				}),
				"invalid header: invalid correlation ID (00000000-0000-0000-0000-000000000000): UUID must use version 4 or 5",
			},
			{
				"invalid source site",
				newEnvelope(func(e *Envelope) {
					e.Header.Source.Site = &identitypb.Identity{}
				}),
				"invalid header: invalid source: invalid site (/00000000-0000-0000-0000-000000000000): invalid name: must be between 1 and 255 bytes",
			},
			{
				"invalid source application",
				newEnvelope(func(e *Envelope) {
					e.Header.Source.Application = &identitypb.Identity{}
				}),
				"invalid header: invalid source: invalid application (/00000000-0000-0000-0000-000000000000): invalid name: must be between 1 and 255 bytes",
			},
			{
				"invalid source handler",
				newEnvelope(func(e *Envelope) {
					e.Header.Source.Handler = &identitypb.Identity{}
				}),
				"invalid header: invalid source: invalid handler (/00000000-0000-0000-0000-000000000000): invalid name: must be between 1 and 255 bytes",
			},
			{
				"source instance ID without source handler",
				newEnvelope(func(e *Envelope) {
					e.Header.Source.Handler = nil
				}),
				"invalid header: invalid source: invalid instance ID: must not be specified without a handler",
			},
			{
				"scheduled-for time without source handler",
				newEnvelope(func(e *Envelope) {
					e.Header.Source.Handler = nil
					e.Header.Source.InstanceId = ""
				}),
				"invalid body: invalid scheduled-for time: must not be specified without a source handler and instance ID",
			},
			{
				"invalid created-at time",
				newEnvelope(func(e *Envelope) {
					e.Body.CreatedAt = nil
				}),
				"invalid body: invalid created-at time: proto: invalid nil Timestamp",
			},
			{
				"invalid scheduled-for time",
				newEnvelope(func(e *Envelope) {
					e.Body.ScheduledFor = &timestamppb.Timestamp{
						Seconds: math.MaxInt64,
						Nanos:   math.MaxInt32,
					}
				}),
				"invalid body: invalid scheduled-for time: proto: timestamp (seconds:9223372036854775807 nanos:2147483647) after 9999-12-31",
			},
			{
				"invalid description",
				newEnvelope(func(e *Envelope) {
					e.Body.Message.Description = ""
				}),
				"invalid body: invalid message: invalid description: must not be empty",
			},
			{
				"without type ID",
				newEnvelope(func(e *Envelope) {
					e.Body.Message.TypeId = nil
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

func newEnvelope(modifiers ...func(*Envelope)) *envelopepb.Envelope {
	env := &Envelope{
		Header: &Header{
			CausationId:   uuidpb.Generate(),
			CorrelationId: uuidpb.Generate(),
			Source: &Source{
				Site:        identitypb.New("<site-name>", uuidpb.Generate()),
				Application: identitypb.New("<app-name>", uuidpb.Generate()),
				Handler:     identitypb.New("<handler-name>", uuidpb.Generate()),
				InstanceId:  "<instance>",
			},
		},
		Body: &Body{
			MessageId:    uuidpb.Generate(),
			CreatedAt:    timestamppb.Now(),
			ScheduledFor: timestamppb.Now(),
			Message: &Message{
				Description: "<description>",
				TypeId:      uuidpb.Generate(),
				Data:        []byte("<data>"),
			},
			Extensions: &Extensions{
				Attributes: map[string]string{"<attr-key>": "<attr-value>"},
				Baggage:    map[string]string{"<baggage-key": "<baggage-value>"},
			},
		},
	}

	for _, fn := range modifiers {
		fn(env)
	}

	return env
}
