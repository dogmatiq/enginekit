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
					e.SourceSite = nil
				}),
			},
			{
				"without source handler",
				newEnvelope(func(e *Envelope) {
					e.SourceHandler = nil
					e.SourceInstanceId = ""
					e.ScheduledFor = nil
				}),
			},
			{
				"without source instance ID",
				newEnvelope(func(e *Envelope) {
					e.SourceInstanceId = ""
				}),
			},
			{
				"without attributes",
				newEnvelope(func(e *Envelope) {
					e.Attributes = nil
				}),
			},
			{
				"without data",
				newEnvelope(func(e *Envelope) {
					e.Data = nil
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
				"invalid message ID (00000000-0000-0000-0000-000000000000): UUID must use version 4",
			},
			{
				"invalid message ID",
				newEnvelope(func(e *Envelope) {
					e.MessageId = &uuidpb.UUID{}
				}),
				"invalid message ID (00000000-0000-0000-0000-000000000000): UUID must use version 4",
			},
			{
				"invalid causation ID",
				newEnvelope(func(e *Envelope) {
					e.CausationId = &uuidpb.UUID{}
				}),
				"invalid causation ID (00000000-0000-0000-0000-000000000000): UUID must use version 4",
			},
			{
				"invalid correlation ID",
				newEnvelope(func(e *Envelope) {
					e.CorrelationId = &uuidpb.UUID{}
				}),
				"invalid correlation ID (00000000-0000-0000-0000-000000000000): UUID must use version 4",
			},
			{
				"invalid source site",
				newEnvelope(func(e *Envelope) {
					e.SourceSite = &identitypb.Identity{}
				}),
				"invalid source site (/00000000-0000-0000-0000-000000000000): invalid name: must be between 1 and 255 bytes",
			},
			{
				"invalid source application",
				newEnvelope(func(e *Envelope) {
					e.SourceApplication = &identitypb.Identity{}
				}),
				"invalid source application (/00000000-0000-0000-0000-000000000000): invalid name: must be between 1 and 255 bytes",
			},
			{
				"invalid source handler",
				newEnvelope(func(e *Envelope) {
					e.SourceHandler = &identitypb.Identity{}
				}),
				"invalid source handler (/00000000-0000-0000-0000-000000000000): invalid name: must be between 1 and 255 bytes",
			},
			{
				"source instance ID without source handler",
				newEnvelope(func(e *Envelope) {
					e.SourceHandler = nil
				}),
				"invalid source instance ID: must not be specified without a source handler",
			},
			{
				"scheduled-for time without source handler",
				newEnvelope(func(e *Envelope) {
					e.SourceHandler = nil
					e.SourceInstanceId = ""
				}),
				"invalid scheduled-for time: must not be specified without a source handler and instance ID",
			},
			{
				"invalid created-at time",
				newEnvelope(func(e *Envelope) {
					e.CreatedAt = nil
				}),
				"invalid created-at time: proto: invalid nil Timestamp",
			},
			{
				"invalid scheduled-for time",
				newEnvelope(func(e *Envelope) {
					e.ScheduledFor = &timestamppb.Timestamp{
						Seconds: math.MaxInt64,
						Nanos:   math.MaxInt32,
					}
				}),
				"invalid scheduled-for time: proto: timestamp (seconds:9223372036854775807 nanos:2147483647) after 9999-12-31",
			},
			{
				"invalid description",
				newEnvelope(func(e *Envelope) {
					e.Description = ""
				}),
				"invalid description: must not be empty",
			},
			{
				"without type ID",
				newEnvelope(func(e *Envelope) {
					e.TypeId = nil
				}),
				"invalid type ID (00000000-0000-0000-0000-000000000000): UUID must use version 4",
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
		MessageId:         uuidpb.Generate(),
		CausationId:       uuidpb.Generate(),
		CorrelationId:     uuidpb.Generate(),
		SourceSite:        identitypb.New("<site-name>", uuidpb.Generate()),
		SourceApplication: identitypb.New("<app-name>", uuidpb.Generate()),
		SourceHandler:     identitypb.New("<handler-name>", uuidpb.Generate()),
		SourceInstanceId:  "<instance>",
		CreatedAt:         timestamppb.Now(),
		ScheduledFor:      timestamppb.Now(),
		Description:       "<description>",
		TypeId:            uuidpb.Generate(),
		Data:              []byte("<data>"),
		Attributes:        map[string]string{"<key>": "<value>"},
	}

	for _, fn := range modifiers {
		fn(env)
	}

	return env
}
