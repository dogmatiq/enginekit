package config

import (
	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/enginekit/optional"
)

// Projection is a [Handler] that represents the configuration of a
// [dogma.ProjectionMessageHandler].
type Projection struct {
	HandlerCommon[dogma.ProjectionMessageHandler]
	DeliveryPolicyComponents []*ProjectionDeliveryPolicy
}

// HandlerType returns [HandlerType] of the handler.
func (h *Projection) HandlerType() HandlerType {
	return ProjectionHandlerType
}

func (h *Projection) validate(ctx *validateContext) {
	h.HandlerCommon.validate(ctx, h.HandlerType())
	panic("not implemented")
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

func (p *ProjectionDeliveryPolicy) describe(*describeContext) {
	panic("not implemented")
}
