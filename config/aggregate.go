package config

import (
	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/enginekit/config/internal/renderer"
	"github.com/dogmatiq/enginekit/protobuf/identitypb"
)

// Aggregate represents the (potentially invalid) configuration of a
// [dogma.AggregateMessageHandler] implementation.
type Aggregate struct {
	HandlerTrait[dogma.AggregateMessageHandler]
}

// Identity returns the entity's identity.
//
// It panics if no single valid identity is configured.
func (h *Aggregate) Identity() *identitypb.Identity {
	return h.identity(strictContext(h))
}

// HandlerType returns [HandlerType] of the handler.
func (h *Aggregate) HandlerType() HandlerType {
	return AggregateHandlerType
}

func (h *Aggregate) String() string {
	return RenderDescriptor(h)
}

func (h *Aggregate) renderDescriptor(ren *renderer.Renderer) {
	renderEntityDescriptor(ren, "aggregate", h.Properties.Source)
}

func (h *Aggregate) renderDetails(ren *renderer.Renderer) {
	renderHandlerDetails(ren, h, h.Properties.Source)
}

func (h *Aggregate) disabledFlags() FlagSet[Disabled] {
	return h.Properties.DisabledFlags
}

func (h *Aggregate) identities() []*Identity {
	return h.Properties.Identities
}

func (h *Aggregate) routes() []*Route {
	return h.Properties.Routes
}

func (h *Aggregate) clone() Component {
	clone := &Aggregate{h.Properties}
	cloneInPlace(&clone.Properties.Identities)
	cloneInPlace(&clone.Properties.Routes)
	return clone
}

func (h *Aggregate) normalize(ctx *normalizationContext) {
	normalizeValue(ctx, &h.Properties.Source, &h.Properties.Fidelity)

	normalizeChildren(ctx, h.Properties.Identities)
	normalizeChildren(ctx, h.Properties.Routes)

	reportIdentityErrors(ctx, h.Properties.Identities)
	reportRouteErrors(ctx, h, h.Properties.Routes)
}
