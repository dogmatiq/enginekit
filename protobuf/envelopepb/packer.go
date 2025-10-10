package envelopepb

import (
	"fmt"
	"time"

	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/enginekit/protobuf/identitypb"
	"github.com/dogmatiq/enginekit/protobuf/uuidpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// A Packer puts messages into envelopes.
type Packer struct {
	// Site is the (optional) identity of the site that the source application
	// is running within.
	//
	// The site is used to disambiguate between messages from different
	// installations of the same application.
	Site *identitypb.Identity

	// Application is the identity of the application that is the source of the
	// messages.
	Application *identitypb.Identity

	// GenerateID is a function used to generate new message IDs.
	//
	// If it is nil, a random UUID is generated.
	GenerateID func() *uuidpb.UUID

	// Now is a function used to get the current time. If it is nil, time.Now()
	// is used.
	Now func() time.Time
}

// Pack returns an envelope containing the given message.
func (p *Packer) Pack(m dogma.Message, options ...PackOption) *Envelope {
	mt, ok := dogma.RegisteredMessageTypeOf(m)
	if !ok {
		panic(fmt.Sprintf(
			"%T is not a registered message type",
			m,
		))
	}

	data, err := m.MarshalBinary()
	if err != nil {
		panic(fmt.Sprintf(
			"unable to marshal %T: %s",
			m,
			err,
		))
	}

	id := p.generateID()

	env := &Envelope{
		MessageId:         id,
		CorrelationId:     id,
		CausationId:       id,
		SourceSite:        p.Site,
		SourceApplication: p.Application,
		Description:       m.MessageDescription(),
		TypeId:            uuidpb.MustParse(mt.ID()),
		Data:              data,
	}

	for _, opt := range options {
		opt(env)
	}

	if env.CreatedAt == nil {
		env.CreatedAt = p.now()
	}

	if err := env.Validate(); err != nil {
		panic(err)
	}

	return env
}

// Unpack returns the message contained within an envelope.
func (p *Packer) Unpack(env *Envelope) (dogma.Message, error) {
	mt, ok := dogma.RegisteredMessageTypeByID(env.TypeId.AsString())
	if !ok {
		return nil, fmt.Errorf(
			"%s is not a registered message type ID",
			env.TypeId,
		)
	}

	m := mt.New()
	if err := m.UnmarshalBinary(env.Data); err != nil {
		return nil, fmt.Errorf("unable to unmarshal %T: %w", m, err)
	}

	return m, nil
}

// now returns the current time.
func (p *Packer) now() *timestamppb.Timestamp {
	if p.Now == nil {
		return timestamppb.Now()
	}

	return timestamppb.New(p.Now())
}

// generateID generates a new message ID.
func (p *Packer) generateID() *uuidpb.UUID {
	if p.GenerateID != nil {
		return p.GenerateID()
	}

	return uuidpb.Generate()
}

// PackOption is an option that alters the behavior of a Pack operation.
type PackOption func(*Envelope)

// WithCause sets env as the "cause" of the message being packed.
func WithCause(env *Envelope) PackOption {
	return func(e *Envelope) {
		e.CausationId = env.MessageId
		e.CorrelationId = env.CorrelationId
	}
}

// WithHandler sets h as the identity of the handler that is the source of the
// message.
func WithHandler(h *identitypb.Identity) PackOption {
	return func(e *Envelope) {
		e.SourceHandler = h
	}
}

// WithInstanceID sets the aggregate or process instance ID that is the
// source of the message.
func WithInstanceID(id string) PackOption {
	return func(e *Envelope) {
		e.SourceInstanceId = id
	}
}

// WithCreatedAt sets the creation time of a message.
func WithCreatedAt(t time.Time) PackOption {
	return func(e *Envelope) {
		e.CreatedAt = timestamppb.New(t)
	}
}

// WithScheduledFor sets the scheduled time of a timeout message.
func WithScheduledFor(t time.Time) PackOption {
	return func(e *Envelope) {
		e.ScheduledFor = timestamppb.New(t)
	}
}
