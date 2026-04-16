package envelopepb

import (
	"fmt"
	"time"

	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/enginekit/protobuf/identitypb"
	"github.com/dogmatiq/enginekit/protobuf/uuidpb"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
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

// PackCommand returns an envelope containing the given command.
func (p *Packer) PackCommand(m dogma.Command, options ...PackCommandOption) *Envelope {
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
		opt.applyPackCommandOption(env)
	}

	if env.Body.CreatedAt == nil {
		env.Body.CreatedAt = p.now()
	}

	// For a "single envelope" we normalize extensions and baggage into the body.
	env.Body.Extensions = flattenAnyValues(env.Header.Extensions, env.Body.Extensions)
	env.Body.Baggage = flattenAnyValues(env.Header.Baggage, env.Body.Baggage)
	env.Header.Extensions = nil
	env.Header.Baggage = nil

	if err := env.Validate(); err != nil {
		panic(err)
	}

	return env
}

// Unpack returns the message contained within an envelope.
// TODO: Make `UnpackCommand`, etc.
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

// PackCommandOption is an option that modifies the behavior of
// [Packer.PackCommand].
type PackCommandOption interface {
	applyPackCommandOption(*Envelope)
}

// PackEffectsOption is an option that modifies the behavior of
// [Packer.PackEffects].
type PackEffectsOption interface {
	applyPackEffectsOption(*Header)
}

// PackEffectCommandOption is an option that modifies the behavior of
// [EffectPacker.PackCommand].
type PackEffectCommandOption interface {
	applyPackEffectCommandOption(*Body)
}

// PackEffectEventOption is an option that modifies the behavior of
// [EffectPacker.PackEvent].
type PackEffectEventOption interface {
	applyPackEffectEventOption(*Body)
}

// PackEffectTimeoutOption is an option that modifies the behavior of
// [EffectPacker.PackTimeout].
type PackEffectTimeoutOption interface {
	applyPackEffectTimeoutOption(*Body)
}

type (
	packEffectsOption           func(*Header)
	packCommandOptionFunc       func(*Body)
	packEffectTimeoutOptionFunc func(*Body)
	universalOption             struct {
		applyToBodyFunc   func(*Body)
		applyToHeaderFunc func(*Header)
	}
)

func (o packEffectsOption) applyPackEffectsOption(header *Header)             { o(header) }
func (o packCommandOptionFunc) applyPackCommandOption(env *Envelope)          { o(env.Body) }
func (o packEffectTimeoutOptionFunc) applyPackEffectTimeoutOption(body *Body) { o(body) }
func (o universalOption) applyPackCommandOption(env *Envelope)                { o.applyToBodyFunc(env.Body) }
func (o universalOption) applyPackEffectCommandOption(body *Body)             { o.applyToBodyFunc(body) }
func (o universalOption) applyPackEffectEventOption(body *Body)               { o.applyToBodyFunc(body) }
func (o universalOption) applyPackEffectTimeoutOption(body *Body)             { o.applyToBodyFunc(body) }
func (o universalOption) applyPackEffectsOption(header *Header)               { o.applyToHeaderFunc(header) }

// WithIdempotencyKey sets the idempotency key of a command packed via
// [Packer.PackCommand].
func WithIdempotencyKey(key string) PackCommandOption {
	return packCommandOptionFunc(
		func(body *Body) {
			body.IdempotencyKey = key
		},
	)
}

// WithInstanceID sets the aggregate or process instance ID that is the source
// of messages packed via [Packer.PackEffects].
func WithInstanceID(id string) PackEffectsOption {
	return packEffectsOption(
		func(header *Header) {
			header.Source.InstanceId = id
		},
	)
}

// WithScheduledFor sets the scheduled time of a timeout packed via
// [EffectPacker.PackTimeout].
func WithScheduledFor(t time.Time) PackEffectTimeoutOption {
	return packEffectTimeoutOptionFunc(
		func(body *Body) {
			body.ScheduledFor = timestamppb.New(t)
		},
	)
}

// WithExtension adds x to the envelope's extensions.
//
// Extensions apply only to the envelope being packed; they are not inherited
// by downstream messages in the causal chain.
func WithExtension(x proto.Message) interface {
	PackCommandOption
	PackEffectsOption
	PackEffectCommandOption
	PackEffectEventOption
	PackEffectTimeoutOption
} {
	v := marshalAsAny(x)

	return universalOption{
		applyToBodyFunc: func(body *Body) {
			appendOrReplace(&body.Extensions, v)
		},
		applyToHeaderFunc: func(header *Header) {
			appendOrReplace(&header.Extensions, v)
		},
	}
}

// WithBaggage adds x to the envelope's baggage.
//
// Baggage applies to the envelope being packed and is inherited by downstream
// messages in the causal chain.
func WithBaggage(x proto.Message) interface {
	PackCommandOption
	PackEffectsOption
	PackEffectCommandOption
	PackEffectEventOption
	PackEffectTimeoutOption
} {
	v := marshalAsAny(x)

	return universalOption{
		applyToBodyFunc: func(body *Body) {
			appendOrReplace(&body.Baggage, v)
		},
		applyToHeaderFunc: func(header *Header) {
			appendOrReplace(&header.Baggage, v)
		},
	}
}

// appendOrReplace appends value to values if there is no existing value with
// the same type URL, otherwise it replaces the existing value in place.
func appendOrReplace(values *[]*anypb.Any, value *anypb.Any) {
	for index, existing := range *values {
		if existing.GetTypeUrl() == value.GetTypeUrl() {
			(*values)[index] = value
			return
		}
	}

	*values = append(*values, value)
}

// marshalAsAny returns v as an [*anypb.Any], converting it if necessary. It
// panics if x is nil, if x is an empty [*anypb.Any], or if it cannot be
// marshaled.
func marshalAsAny(x proto.Message) *anypb.Any {
	if x == nil {
		panic("value must not be nil")
	}

	if x, ok := x.(*anypb.Any); ok {
		if err := validateAnyValue(x); err != nil {
			panic("value must not be an empty google.protobuf.Any")
		}

		return x
	}

	v, err := anypb.New(x)
	if err != nil {
		panic(fmt.Sprintf(
			"unable to marshal %T as google.protobuf.Any: %s",
			x,
			err,
		))
	}

	return v
}
