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
		Header: &Header{
			CausationId:   id,
			CorrelationId: id,
			Source: &Source{
				Site:        p.Site,
				Application: p.Application,
			},
		},
		Body: &Body{
			MessageId: id,
			Message: &Message{
				TypeId:      uuidpb.MustParse(mt.ID()),
				Description: m.MessageDescription(),
				Data:        data,
			},
		},
	}

	for _, opt := range options {
		opt.applyToEnvelope(env)
	}

	if env.Body.CreatedAt == nil {
		env.Body.CreatedAt = p.now()
	}

	if err := env.Validate(); err != nil {
		panic(err)
	}

	return env
}

// Unpack returns the message contained within an envelope.
func Unpack(env *Envelope) (dogma.Message, error) {
	message := env.GetBody().GetMessage()

	if err := message.validate(); err != nil {
		return nil, fmt.Errorf("invalid message: %w", err)
	}

	mt, ok := dogma.RegisteredMessageTypeByID(message.TypeId.AsString())
	if !ok {
		return nil, fmt.Errorf(
			"%s is not a registered message type ID",
			message.TypeId,
		)
	}

	m := mt.New()
	if err := m.UnmarshalBinary(message.Data); err != nil {
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

// PackOption is an option that alters the behavior of a [Packer.Pack]
// operation.
type PackOption interface {
	applyToEnvelope(*Envelope)
}

// SourcePackOption is a [PackOption] that modifies only the source
// information in an envelope header.
type SourcePackOption interface {
	PackOption
	applyToSource(*Source)
}

// BodyPackOption is a [PackOption] that modifies only envelope bodies.
type BodyPackOption interface {
	PackOption
	applyToBody(*Body)
}

type packOption func(*Envelope)

func (opt packOption) applyToEnvelope(env *Envelope) {
	opt(env)
}

type sourcePackOption func(*Source)

func (opt sourcePackOption) applyToEnvelope(env *Envelope) {
	opt(env.Header.Source)
}

func (opt sourcePackOption) applyToSource(source *Source) {
	opt(source)
}

type bodyPackOption func(*Body)

func (opt bodyPackOption) applyToEnvelope(env *Envelope) {
	opt(env.Body)
}

func (opt bodyPackOption) applyToBody(body *Body) {
	opt(body)
}

// WithCause sets env as the "cause" of the message being packed.
func WithCause(env *Envelope) PackOption {
	return packOption(
		func(packed *Envelope) {
			packed.Header.CausationId = env.GetBody().GetMessageId()
			packed.Header.CorrelationId = env.GetHeader().GetCorrelationId()
		},
	)
}

// WithHandler sets h as the identity of the handler that is the source of the
// message.
func WithHandler(h *identitypb.Identity) SourcePackOption {
	return sourcePackOption(
		func(source *Source) {
			source.Handler = h
		},
	)
}

// WithInstanceID sets the aggregate or process instance ID that is the
// source of the message.
func WithInstanceID(id string) SourcePackOption {
	return sourcePackOption(
		func(source *Source) {
			source.InstanceId = id
		},
	)
}

// WithCreatedAt sets the creation time of a message.
func WithCreatedAt(t time.Time) BodyPackOption {
	return bodyPackOption(
		func(body *Body) {
			body.CreatedAt = timestamppb.New(t)
		},
	)
}

// WithScheduledFor sets the scheduled time of a timeout message.
func WithScheduledFor(t time.Time) BodyPackOption {
	return bodyPackOption(
		func(body *Body) {
			body.ScheduledFor = timestamppb.New(t)
		},
	)
}

// WithIdempotencyKey sets the idempotency key of a command message.
func WithIdempotencyKey(key string) BodyPackOption {
	return bodyPackOption(
		func(body *Body) {
			body.IdempotencyKey = key
		},
	)
}
