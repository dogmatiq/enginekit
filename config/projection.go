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
func (p *Projection) HandlerType() HandlerType {
	return ProjectionHandlerType
}

// ProjectionDeliveryPolicy is a [Component] that represents the configuration
// of a [dogma.ProjectionDeliveryPolicy].
type ProjectionDeliveryPolicy struct {
	ComponentCommon

	DeliveryPolicyType      optional.Optional[ProjectionDeliveryPolicyType]
	BroadcastToPrimaryFirst optional.Optional[bool]
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
