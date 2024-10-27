package config

import (
	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/enginekit/config/internal/renderer"
	"github.com/dogmatiq/enginekit/optional"
	"github.com/dogmatiq/enginekit/protobuf/identitypb"
)

// Process represents the (potentially invalid) configuration of a
// [dogma.ProcessMessageHandler] implementation.
type Process struct {
	HandlerProperties

	// Source is the instance of the entity from which the configuration was
	// sourced, if available.
	Source optional.Optional[dogma.ProcessMessageHandler]
}

// Identity returns the entity's identity.
//
// It panics if no single valid identity is configured.
func (h *Process) Identity() *identitypb.Identity {
	return resolveIdentity(h)
}

// HandlerType returns [HandlerType] of the handler.
func (h *Process) HandlerType() HandlerType {
	return ProcessHandlerType
}

// RouteSet returns the routes configured for the entity.
//
// It panics if the route configuration is invalid or cannot be determined
// completely.
func (h *Process) RouteSet() RouteSet {
	return resolveRouteSet(strictContext(h), h)
}

// IsDisabled returns true if the handler is disabled.
//
// It panics if the disabled state cannot be determined.
func (h *Process) IsDisabled() bool {
	return resolveIsDisabled(h)
}

func (h *Process) String() string {
	return RenderDescriptor(h)
}

func (h *Process) renderDescriptor(ren *renderer.Renderer) {
	renderEntityDescriptor(ren, h)
}

func (h *Process) renderDetails(ren *renderer.Renderer) {
	renderHandlerDetails(ren, h, h.Source)
}

func (h *Process) cloneHandler() Handler {
	return &Process{
		clone(h.HandlerProperties),
		h.Source,
	}
}

func (h *Process) normalize(ctx *normalizationContext) {
	normalizeHandler(ctx, h, h.Source)
}
