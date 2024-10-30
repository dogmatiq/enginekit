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
	panic("not implemented")
}

// Interface returns the [dogma.Application] that the entity represents.
func (h *Projection) Interface() dogma.ProjectionMessageHandler {
	return resolveInterface(h, h.Source)
}

func (h *Projection) String() string {
	return stringifyEntity(h)
}

func (h *Projection) identities() []*Identity {
	return h.IdentityComponents
}

func (h *Projection) routes() []*Route {
	return h.RouteComponents
}

func (h *Projection) validate(ctx *validateContext) {
	validateHandler(
		ctx,
		h,
		h.Source,
		func(ctx *validateContext) {
			// TODO: delivery policies
			panic("not implemented")
		},
	)
}

func (h *Projection) describe(ctx *describeContext) {
	describeHandler(ctx, h, h.Source)
}

// ProjectionDeliveryPolicyType is an enumeration of the different types of
// projection delivery policies.
type ProjectionDeliveryPolicyType int

const (
	// UnicastProjectionDeliveryPolicyType is the [ProjectionDeliveryPolicyType]
	// for [dogma.UnicastProjectionDeliveryPolicy].
	UnicastProjectionDeliveryPolicyType ProjectionDeliveryPolicyType = iota

	// BroadcastProjectionDeliveryPolicyType is the
	// [ProjectionDeliveryPolicyType] for
	// [dogma.BroadcastProjectionDeliveryPolicy].
	BroadcastProjectionDeliveryPolicyType
)

// ProjectionDeliveryPolicy is a [Component] that represents the configuration
// of a [dogma.ProjectionDeliveryPolicy].
type ProjectionDeliveryPolicy struct {
	ComponentCommon

	DeliveryPolicyType      optional.Optional[ProjectionDeliveryPolicyType]
	BroadcastToPrimaryFirst optional.Optional[bool]
}

func (p *ProjectionDeliveryPolicy) String() string {
	panic("not implemented")
}

func (p *ProjectionDeliveryPolicy) validate(*validateContext) {
	panic("not implemented")
}

func (p *ProjectionDeliveryPolicy) describe(*describeContext) {
	panic("not implemented")
}
