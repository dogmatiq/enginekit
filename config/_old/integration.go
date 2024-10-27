package config

import (
	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/enginekit/config/internal/renderer"
	"github.com/dogmatiq/enginekit/optional"
	"github.com/dogmatiq/enginekit/protobuf/identitypb"
)

// Integration represents the (potentially invalid) configuration of a
// [dogma.IntegrationMessageHandler] implementation.
type Integration struct {
	HandlerProperties

	// Source is the instance of the entity from which the configuration was
	// sourced, if available.
	Source optional.Optional[dogma.IntegrationMessageHandler]
}

// Identity returns the entity's identity.
//
// It panics if no single valid identity is configured.
func (h *Integration) Identity() *identitypb.Identity {
	return resolveIdentity(h)
}

// HandlerType returns [HandlerType] of the handler.
func (h *Integration) HandlerType() HandlerType {
	return IntegrationHandlerType
}

// RouteSet returns the routes configured for the entity.
//
// It panics if the route configuration is invalid or cannot be determined
// completely.
func (h *Integration) RouteSet() RouteSet {
	return resolveRouteSet(strictContext(h), h)
}

// IsDisabled returns true if the handler is disabled.
//
// It panics if the disabled state cannot be determined.
func (h *Integration) IsDisabled() bool {
	return resolveIsDisabled(h)
}

func (h *Integration) String() string {
	return RenderDescriptor(h)
}

func (h *Integration) renderDescriptor(ren *renderer.Renderer) {
	renderEntityDescriptor(ren, h)
}

func (h *Integration) renderDetails(ren *renderer.Renderer) {
	renderHandlerDetails(ren, h, h.Source)
}

func (h *Integration) cloneHandler() Handler {
	return &Integration{
		clone(h.HandlerProperties),
		h.Source,
	}
}

func (h *Integration) normalize(ctx *normalizationContext) {
	normalizeHandler(ctx, h, h.Source)
}
