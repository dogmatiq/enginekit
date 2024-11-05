package configbuilder

import (
	"fmt"

	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/enginekit/config"
	"github.com/dogmatiq/enginekit/optional"
)

// ProjectionDeliveryPolicy returns a new [dogma.ProjectionDeliveryPolicy] as
// configured by fn.
func ProjectionDeliveryPolicy(fn func(*ProjectionDeliveryPolicyBuilder)) *config.ProjectionDeliveryPolicy {
	x := &ProjectionDeliveryPolicyBuilder{}
	fn(x)
	return x.Done()
}

// ProjectionDeliveryPolicyBuilder constructs a
// [config.ProjectionDeliveryPolicy].
type ProjectionDeliveryPolicyBuilder struct {
	target config.ProjectionDeliveryPolicy
}

// AsPerDeliveryPolicy configures the builder to use the same properties as p.
func (b *ProjectionDeliveryPolicyBuilder) AsPerDeliveryPolicy(p dogma.ProjectionDeliveryPolicy) {
	switch p := p.(type) {
	case dogma.UnicastProjectionDeliveryPolicy:
		b.target.DeliveryPolicyType = optional.Some(config.UnicastProjectionDeliveryPolicyType)
	case dogma.BroadcastProjectionDeliveryPolicy:
		b.target.DeliveryPolicyType = optional.Some(config.BroadcastProjectionDeliveryPolicyType)
		b.target.Broadcast.PrimaryFirst = optional.Some(p.PrimaryFirst)
	default:
		panic("unsupported projection delivery policy type")
	}
}

// Type sets the type of the delivery policy.
func (b *ProjectionDeliveryPolicyBuilder) Type(t config.ProjectionDeliveryPolicyType) {
	b.target.DeliveryPolicyType = optional.Some(t)
}

// BroadcastToPrimaryFirst sets the value of the "broadcast to primary first"
// property of a [config.BroadcastProjectionDeliveryPolicyType].
func (b *ProjectionDeliveryPolicyBuilder) BroadcastToPrimaryFirst(v bool) {
	b.target.DeliveryPolicyType = optional.Some(config.BroadcastProjectionDeliveryPolicyType)
	b.target.Broadcast.PrimaryFirst = optional.Some(v)
}

// Partial marks the compomnent as partially configured.
func (b *ProjectionDeliveryPolicyBuilder) Partial(format string, args ...any) {
	b.target.IsPartialReasons = append(b.target.IsPartialReasons, fmt.Sprintf(format, args...))
}

// Speculative marks the component as speculative.
func (b *ProjectionDeliveryPolicyBuilder) Speculative() {
	b.target.IsSpeculative = true
}

// Done completes the configuration of the policy.
func (b *ProjectionDeliveryPolicyBuilder) Done() *config.ProjectionDeliveryPolicy {
	return &b.target
}
