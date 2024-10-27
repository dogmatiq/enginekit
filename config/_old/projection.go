package config

import (
	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/enginekit/config/internal/renderer"
	"github.com/dogmatiq/enginekit/internal/typename"
	"github.com/dogmatiq/enginekit/optional"
	"github.com/dogmatiq/enginekit/protobuf/identitypb"
)

// Projection represents the (potentially invalid) configuration of a
// [dogma.ProjectionMessageHandler] implementation.
type Projection struct {
	HandlerProperties

	// Source is the instance of the entity from which the configuration was
	// sourced, if available.
	Source optional.Optional[dogma.ProjectionMessageHandler]

	// DeliveryPolicy is the delivery policy for the handler, if configured.
	DeliveryPolicy optional.Optional[Value[dogma.ProjectionDeliveryPolicy]]
}

// Identity returns the entity's identity.
//
// It panics if no single valid identity is configured.
func (h *Projection) Identity() *identitypb.Identity {
	return resolveIdentity(h)
}

// HandlerType returns [HandlerType] of the handler.
func (h *Projection) HandlerType() HandlerType {
	return ProjectionHandlerType
}

// RouteSet returns the routes configured for the entity.
//
// It panics if the route configuration is invalid or cannot be determined
// completely.
func (h *Projection) RouteSet() RouteSet {
	return resolveRouteSet(strictContext(h), h)
}

// IsDisabled returns true if the handler is disabled.
//
// It panics if the disabled state cannot be determined.
func (h *Projection) IsDisabled() bool {
	return resolveIsDisabled(h)
}

func (h *Projection) String() string {
	return RenderDescriptor(h)
}

func (h *Projection) renderDescriptor(ren *renderer.Renderer) {
	renderEntityDescriptor(ren, h)
}

func (h *Projection) renderDetails(ren *renderer.Renderer) {
	renderHandlerDetails(ren, h, h.Source)

	if p, ok := h.DeliveryPolicy.TryGet(); ok {
		// TODO: https://github.com/dogmatiq/enginekit/issues/55
		if typeName, ok := p.TypeName.TryGet(); ok {
			ren.IndentBullet()

			switch typeName {
			case typename.For[dogma.UnicastProjectionDeliveryPolicy]():
				ren.Printf("unicast delivery policy")
			case typename.For[dogma.BroadcastProjectionDeliveryPolicy]():
				ren.Printf("broadcast delivery policy")
			default:
				ren.Printf("unrecognized delivery policy")
			}

			if !p.Value.IsPresent() {
				ren.Print(" (runtime type unavailable)")
			}

			ren.Indent()
			ren.Print("\n")
		}
	}
}

func (h *Projection) cloneHandler() Handler {
	return &Projection{
		clone(h.HandlerProperties),
		h.Source,
		h.DeliveryPolicy,
	}
}

func (h *Projection) normalize(ctx *normalizationContext) {
	normalizeHandler(ctx, h, h.Source)

	if p, ok := h.DeliveryPolicy.TryGet(); ok {
		normalizeValue(ctx, &p, &h.Fidelity)
		h.DeliveryPolicy = optional.Some(p)
	}
}
