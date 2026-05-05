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

	env := NewEnvelopeBuilder().
		WithHeader(
			NewHeaderBuilder().
				WithCausationId(id).
				WithCorrelationId(id).
				WithSource(NewSourceBuilder().
					WithSite(p.Site).
					WithApplication(p.Application).
					Build()).
				Build(),
		).
		WithBody(
			NewBodyBuilder().
				WithMessageId(id).
				WithMessage(NewMessageBuilder().
					WithTypeId(uuidpb.MustParse(mt.ID())).
					WithDescription(m.MessageDescription()).
					WithData(data).
					Build()).
				Build(),
		).
		Build()

	for _, opt := range options {
		opt.applyPackCommandOption(env)
	}

	if !env.GetBody().HasCreatedAt() {
		env.GetBody().SetCreatedAt(p.now())
	}

	// For a "single envelope" we normalize extensions and baggage into the body.
	env.GetBody().SetExtensions(flattenAnyValues(env.GetHeader().GetExtensions(), env.GetBody().GetExtensions()))
	env.GetBody().SetBaggage(flattenAnyValues(env.GetHeader().GetBaggage(), env.GetBody().GetBaggage()))
	env.GetHeader().SetExtensions(nil)
	env.GetHeader().SetBaggage(nil)

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

	mt, ok := dogma.RegisteredMessageTypeByID(message.GetTypeId().AsString())
	if !ok {
		return nil, fmt.Errorf(
			"%s is not a registered message type ID",
			message.GetTypeId(),
		)
	}

	m := mt.New()
	if err := m.UnmarshalBinary(message.GetData()); err != nil {
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

// PackEffectDeadlineOption is an option that modifies the behavior of
// [EffectPacker.PackDeadline].
type PackEffectDeadlineOption interface {
	applyPackEffectDeadlineOption(*Body)
}

type (
	packEffectsOption            func(*Header)
	packCommandOptionFunc        func(*Body)
	packEffectDeadlineOptionFunc func(*Body)
	universalOption              struct {
		applyToBodyFunc   func(*Body)
		applyToHeaderFunc func(*Header)
	}
)

func (o packEffectsOption) applyPackEffectsOption(header *Header)               { o(header) }
func (o packCommandOptionFunc) applyPackCommandOption(env *Envelope)            { o(env.GetBody()) }
func (o packEffectDeadlineOptionFunc) applyPackEffectDeadlineOption(body *Body) { o(body) }
func (o universalOption) applyPackCommandOption(env *Envelope)                  { o.applyToBodyFunc(env.GetBody()) }
func (o universalOption) applyPackEffectCommandOption(body *Body)               { o.applyToBodyFunc(body) }
func (o universalOption) applyPackEffectEventOption(body *Body)                 { o.applyToBodyFunc(body) }
func (o universalOption) applyPackEffectDeadlineOption(body *Body)              { o.applyToBodyFunc(body) }
func (o universalOption) applyPackEffectsOption(header *Header)                 { o.applyToHeaderFunc(header) }

// WithIdempotencyKey sets the idempotency key of a command packed via
// [Packer.PackCommand].
func WithIdempotencyKey(key string) PackCommandOption {
	return packCommandOptionFunc(
		func(body *Body) {
			body.SetIdempotencyKey(key)
		},
	)
}

// WithInstanceID sets the aggregate or process instance ID that is the source
// of messages packed via [Packer.PackEffects].
func WithInstanceID(id string) PackEffectsOption {
	return packEffectsOption(
		func(header *Header) {
			header.GetSource().SetInstanceId(id)
		},
	)
}

// WithScheduledFor sets the time at which a deadline is reached when packed
// via [EffectPacker.PackDeadline].
func WithScheduledFor(t time.Time) PackEffectDeadlineOption {
	return packEffectDeadlineOptionFunc(
		func(body *Body) {
			body.SetScheduledFor(timestamppb.New(t))
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
	PackEffectDeadlineOption
} {
	v := marshalAsAny(x)

	return universalOption{
		applyToBodyFunc: func(body *Body) {
			body.SetExtensions(appendOrReplace(body.GetExtensions(), v))
		},
		applyToHeaderFunc: func(header *Header) {
			header.SetExtensions(appendOrReplace(header.GetExtensions(), v))
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
	PackEffectDeadlineOption
} {
	v := marshalAsAny(x)

	return universalOption{
		applyToBodyFunc: func(body *Body) {
			body.SetBaggage(appendOrReplace(body.GetBaggage(), v))
		},
		applyToHeaderFunc: func(header *Header) {
			header.SetBaggage(appendOrReplace(header.GetBaggage(), v))
		},
	}
}

// appendOrReplace appends value to values if there is no existing value with
// the same type URL, otherwise it replaces the existing value in place.
func appendOrReplace(values []*anypb.Any, value *anypb.Any) []*anypb.Any {
	for index, existing := range values {
		if existing.GetTypeUrl() == value.GetTypeUrl() {
			values[index] = value
			return values
		}
	}

	return append(values, value)
}

// marshalAsAny returns v as an [*anypb.Any], converting it if necessary. It
// panics if x is nil, if x is an empty [*anypb.Any], or if it cannot be
// marshaled.
func marshalAsAny(x proto.Message) *anypb.Any {
	if x == nil {
		panic("value must not be nil")
	}

	if x, ok := x.(*anypb.Any); ok {
		if x.GetTypeUrl() == "" {
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
