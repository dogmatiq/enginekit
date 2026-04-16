package envelopepb

import (
	"fmt"

	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/enginekit/protobuf/identitypb"
	"github.com/dogmatiq/enginekit/protobuf/uuidpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// An EffectPacker packs messages produced while handling a specific causal
// message into a [MultiEnvelope].
type EffectPacker struct {
	generateID func() *uuidpb.UUID
	now        func() *timestamppb.Timestamp
	header     *Header
	bodies     []*Body
	sealed     bool
}

// PackEffects returns an [EffectPacker] that packs messages produced by h while
// handling cause.
func (p *Packer) PackEffects(
	cause *Envelope,
	h *identitypb.Identity,
	options ...PackEffectsOption,
) *EffectPacker {
	if h == nil {
		panic("handler must not be nil")
	}

	if cause == nil {
		panic("cause must not be nil")
	}

	if err := cause.Validate(); err != nil {
		panic(fmt.Errorf("invalid cause envelope: %w", err))
	}

	generateID := p.GenerateID
	if generateID == nil {
		generateID = uuidpb.Generate
	}

	now := timestamppb.Now
	if p.Now != nil {
		nowTime := p.Now
		now = func() *timestamppb.Timestamp {
			return timestamppb.New(nowTime())
		}
	}

	header := &Header{
		CausationId:   cause.Body.MessageId,
		CorrelationId: cause.Header.CorrelationId,
		Source: &Source{
			Site:        p.Site,
			Application: p.Application,
			Handler:     h,
		},
		Baggage: flattenAnyValues(
			cause.GetHeader().GetBaggage(),
			cause.GetBody().GetBaggage(),
		),
	}

	for _, opt := range options {
		opt.applyPackEffectsOption(header)
	}

	if err := header.validate(); err != nil {
		panic(fmt.Errorf("invalid header: %w", err))
	}

	return &EffectPacker{
		generateID: generateID,
		now:        now,
		header:     header,
	}
}

// PackCommand appends m to the multi-envelope under construction.
func (p *EffectPacker) PackCommand(m dogma.Command, options ...PackEffectCommandOption) {
	packEffectBody(p, m, PackEffectCommandOption.applyPackEffectCommandOption, options...)
}

// PackEvent appends m to the multi-envelope under construction.
func (p *EffectPacker) PackEvent(m dogma.Event, options ...PackEffectEventOption) {
	packEffectBody(p, m, PackEffectEventOption.applyPackEffectEventOption, options...)
}

// PackTimeout appends m to the multi-envelope under construction.
func (p *EffectPacker) PackTimeout(m dogma.Timeout, options ...PackEffectTimeoutOption) {
	packEffectBody(p, m, PackEffectTimeoutOption.applyPackEffectTimeoutOption, options...)
}

func packEffectBody[T any](
	p *EffectPacker,
	m dogma.Message,
	apply func(T, *Body),
	options ...T,
) {
	p.mustNotBeSealed()

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

	body := &Body{
		MessageId: p.generateID(),
		Message: &Message{
			TypeId:      uuidpb.MustParse(mt.ID()),
			Description: m.MessageDescription(),
			Data:        data,
		},
	}

	for _, opt := range options {
		apply(opt, body)
	}

	if body.CreatedAt == nil {
		body.CreatedAt = p.now()
	}

	if err := body.validate(p.header); err != nil {
		panic(fmt.Errorf("invalid body: %w", err))
	}

	p.bodies = append(p.bodies, body)
}

// Seal returns a [MultiEnvelope] containing all packed messages, or false if no
// messages were packed.
func (p *EffectPacker) Seal() (*MultiEnvelope, bool) {
	p.mustNotBeSealed()
	p.sealed = true

	if len(p.bodies) == 0 {
		return nil, false
	}

	return &MultiEnvelope{
		Header: p.header,
		Bodies: p.bodies,
	}, true
}

func (p *EffectPacker) mustNotBeSealed() {
	if p.sealed {
		panic("already sealed")
	}
}
