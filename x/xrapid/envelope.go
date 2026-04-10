package xrapid

import (
	"github.com/dogmatiq/enginekit/protobuf/envelopepb"
	"github.com/dogmatiq/enginekit/protobuf/uuidpb"
	"pgregory.net/rapid"
)

// Envelope returns a generator of random [*envelopepb.Envelope] values.
//
// By design, the message type and data encoded within the envelope is not
// necessarily valid.
func Envelope() *rapid.Generator[*envelopepb.Envelope] {
	return rapid.Custom(
		func(t *rapid.T) *envelopepb.Envelope {
			env := &envelopepb.Envelope{
				Header: &envelopepb.Header{
					CausationId:   uuidpb.Generate(),
					CorrelationId: uuidpb.Generate(),
					Source: &envelopepb.Source{
						Site:        Nillable(Identity()).Draw(t, "source site"),
						Application: Identity().Draw(t, "source application"),
						Handler:     Nillable(Identity()).Draw(t, "source handler"),
					},
				},
				Body: &envelopepb.Body{
					MessageId: uuidpb.Generate(),
					CreatedAt: Timestamp().Draw(t, "created at"),
					Message: &envelopepb.Message{
						Description: rapid.StringN(1, -1, -1).Draw(t, "description"),
						TypeId:      uuidpb.Generate(),
						Data:        rapid.SliceOf(rapid.Byte()).Draw(t, "data"),
					},
					Extensions: &envelopepb.Extensions{
						Attributes: rapid.MapOf(rapid.String(), rapid.String()).Draw(t, "attributes"),
						Baggage:    rapid.MapOf(rapid.String(), rapid.String()).Draw(t, "baggage"),
					},
				},
			}

			if env.Header.Source.Handler != nil {
				env.Header.Source.InstanceId = rapid.String().Draw(t, "source instance id")

				if env.Header.Source.InstanceId != "" {
					env.Body.ScheduledFor = Nillable(Timestamp()).Draw(t, "scheduled for")
				}
			}

			if err := env.Validate(); err != nil {
				t.Fatalf("generated invalid envelope: %v", err)
			}

			return env
		},
	)
}
