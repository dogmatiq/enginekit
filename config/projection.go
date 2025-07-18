package config

import (
	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/enginekit/optional"
	"github.com/dogmatiq/enginekit/protobuf/identitypb"
)

// Projection is a [Handler] that represents the configuration of a
// [dogma.ProjectionMessageHandler].
type Projection struct {
	HandlerCommon
	DeliveryPolicyComponents []*ProjectionDeliveryPolicy
	Source                   optional.Optional[dogma.ProjectionMessageHandler]
}

// Identity returns the entity's identity.
//
// It panics the configuration does not specify a singular valid identity.
func (h *Projection) Identity() *identitypb.Identity {
	return resolveIdentity(h)
}

// HandlerType returns the [HandlerType] of the handler.
func (h *Projection) HandlerType() HandlerType {
	return ProjectionHandlerType
}

// RouteSet returns the routes configured for the entity.
//
// It panics if the configuration does not specify a complete set of valid
// routes for the entity and its constituents.
func (h *Projection) RouteSet() RouteSet {
	return resolveRouteSet(h)
}

// IsDisabled returns true if the handler is disabled.
//
// It panics if the configuration does not specify unambiguously whether the
// handler is enabled or disabled.
func (h *Projection) IsDisabled() bool {
	return resolveIsDisabled(h)
}

// DeliveryPolicy returns the delivery policy for the handler.
//
// It returns the last delivery policy component that was configured, as each
// call to [dogma.ProjectionConfigurer.DeliveryPolicy] replaces the previous
// value.
//
// If no delivery policy has been configured, it returns
// [dogma.UnicastProjectionDeliveryPolicy], which is the default.
//
// It panics if the configuration does not specify valid delivery policies.
func (h *Projection) DeliveryPolicy() dogma.ProjectionDeliveryPolicy {
	ctx := newResolutionContext(h, false)

	for _, p := range h.DeliveryPolicyComponents {
		ctx.ValidateChild(p)
	}

	n := len(h.DeliveryPolicyComponents)
	if n == 0 {
		return dogma.UnicastProjectionDeliveryPolicy{}
	}

	return h.DeliveryPolicyComponents[n-1].Interface()
}

// Interface returns the [dogma.Application] that the entity represents.
func (h *Projection) Interface() dogma.ProjectionMessageHandler {
	return resolveInterface(h, h.Source)
}

func (h *Projection) String() string {
	return stringifyEntity(h)
}

func (h *Projection) validate(ctx *validateContext) {
	validateHandler(ctx, h, h.Source)

	for _, p := range h.DeliveryPolicyComponents {
		ctx.ValidateChild(p)
	}
}

func (h *Projection) describe(ctx *describeContext) {
	describeHandler(ctx, h, h.Source)

	for _, p := range h.DeliveryPolicyComponents {
		ctx.DescribeChild(p)
	}
}
