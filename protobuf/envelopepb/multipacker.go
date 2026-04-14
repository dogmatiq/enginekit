package envelopepb

import (
	"fmt"

	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/enginekit/protobuf/uuidpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// A MultiPacker puts messages into a [MultiEnvelope].
type MultiPacker struct {
	generateID func() *uuidpb.UUID
	now        func() *timestamppb.Timestamp
	header     *Header
	bodies     []*Body
	sealed     bool
}

// CausedBy returns a [MultiPacker] that packs messages that share env as their
// cause into a shared [MultiEnvelope].
func (p *Packer) CausedBy(env *Envelope, options ...SourcePackOption) *MultiPacker {
	if env == nil {
		panic("cause must not be nil")
	}

	if err := env.Validate(); err != nil {
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
		CausationId:   env.Body.MessageId,
		CorrelationId: env.Header.CorrelationId,
		Source: &Source{
			Site:        p.Site,
			Application: p.Application,
		},
	}

	for _, opt := range options {
		opt.applyToSource(header.Source)
	}

	if err := header.validate(); err != nil {
		panic(fmt.Errorf("invalid header: %w", err))
	}

	return &MultiPacker{
		generateID: generateID,
		now:        now,
		header:     header,
	}
}

// Pack appends m to the multi-envelope under construction.
func (p *MultiPacker) Pack(m dogma.Message, options ...BodyPackOption) {
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
		opt.applyToBody(body)
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
func (p *MultiPacker) Seal() (*MultiEnvelope, bool) {
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

func (p *MultiPacker) mustNotBeSealed() {
	if p.sealed {
		panic("already sealed")
	}
}
