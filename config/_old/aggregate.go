package config

import (
	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/enginekit/config/internal/renderer"
	"github.com/dogmatiq/enginekit/optional"
	"github.com/dogmatiq/enginekit/protobuf/identitypb"
)

// Aggregate represents the (potentially invalid) configuration of a
// [dogma.AggregateMessageHandler] implementation.
type Aggregate struct {
	HandlerProperties

	// Source is the instance of the entity from which the configuration was
	// sourced, if available.
	Source optional.Optional[dogma.AggregateMessageHandler]
}

// Identity returns the entity's identity.
//
// It panics if no single valid identity is configured.
func (h *Aggregate) Identity() *identitypb.Identity {
	return resolveIdentity(h)
}

// HandlerType returns [HandlerType] of the handler.
func (h *Aggregate) HandlerType() HandlerType {
	return AggregateHandlerType
}

// RouteSet returns the routes configured for the entity.
//
// It panics if the route configuration is invalid or cannot be determined
// completely.
func (h *Aggregate) RouteSet() RouteSet {
	return resolveRouteSet(strictContext(h), h)
}

// IsDisabled returns true if the handler is disabled.
//
// It panics if the disabled state cannot be determined.
func (h *Aggregate) IsDisabled() bool {
	return resolveIsDisabled(h)
}

func (h *Aggregate) String() string {
	return RenderDescriptor(h)
}

func (h *Aggregate) renderDescriptor(ren *renderer.Renderer) {
	renderEntityDescriptor(ren, h)
}

func (h *Aggregate) renderDetails(ren *renderer.Renderer) {
	renderHandlerDetails(ren, h, h.Source)
}

func (h *Aggregate) normalize(ctx *normalizationContext) {
	normalizeHandler(ctx, h, h.Source)
}

func (h *Aggregate) clone() any {
	return &Aggregate{
		clone(h.HandlerProperties),
		h.Source,
	}
}
