package config

import (
	"iter"
	"reflect"
	"strings"

	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/enginekit/internal/enum"
	"github.com/dogmatiq/enginekit/optional"
)

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

// ProjectionDeliveryPolicyTypes returns a sequence that yields all valid
// [ProjectionDeliveryPolicyType] values.
func ProjectionDeliveryPolicyTypes() iter.Seq[ProjectionDeliveryPolicyType] {
	return enum.Range(UnicastProjectionDeliveryPolicyType, BroadcastProjectionDeliveryPolicyType)
}

func (t ProjectionDeliveryPolicyType) String() string {
	return enum.String(
		t,
		"unicast",
		"broadcast",
	)
}

// SwitchByProjectionDeliveryPolicyType invokes one of the provided functions
// based on t.
//
// It provides a compile-time guarantee that all possible values are handled,
// even if new [ProjectionDeliveryPolicyType] values are added in the future.
//
// It panics if the function associated with t is nil, or if t is not a valid
// [ProjectionDeliveryPolicyType].
func SwitchByProjectionDeliveryPolicyType(
	t ProjectionDeliveryPolicyType,
	unicast func(),
	broadcast func(),
) {
	enum.Switch(
		t,
		unicast,
		broadcast,
	)
}

// MapByProjectionDeliveryPolicyType maps t to a value of type T.
//
// It provides a compile-time guarantee that all possible values are handled,
// even if new [ProjectionDeliveryPolicyType] values are added in the future.
//
// It panics if t is not a valid [ProjectionDeliveryPolicyType].
func MapByProjectionDeliveryPolicyType[T any](
	t ProjectionDeliveryPolicyType,
	unicast T,
	broadcast T,
) T {
	return enum.Map(
		t,
		unicast,
		broadcast,
	)
}

// ProjectionDeliveryPolicy is a [Component] that represents the configuration
// of a [dogma.ProjectionDeliveryPolicy].
type ProjectionDeliveryPolicy struct {
	// Fidelity reports how faithfully the [Component] describes a complete
	// configuration that can be used to execute an application.
	Fidelity Fidelity

	DeliveryPolicyType optional.Optional[ProjectionDeliveryPolicyType]

	Broadcast struct {
		PrimaryFirst optional.Optional[bool]
	}
}

// Interface returns the [dogma.ProjectionDeliveryPolicy] that the [Component]
// describes.
func (p *ProjectionDeliveryPolicy) Interface() dogma.ProjectionDeliveryPolicy {
	ctx := newResolutionContext(p)
	p.validateForExecution(ctx)

	return MapByProjectionDeliveryPolicyType(
		p.DeliveryPolicyType.Get(),
		func() dogma.ProjectionDeliveryPolicy {
			return dogma.UnicastProjectionDeliveryPolicy{}
		},
		func() dogma.ProjectionDeliveryPolicy {
			return dogma.BroadcastProjectionDeliveryPolicy{
				PrimaryFirst: p.Broadcast.PrimaryFirst.Get(),
			}
		},
	)()
}

func (p *ProjectionDeliveryPolicy) String() string {
	var w strings.Builder

	w.WriteString("delivery-policy")

	if t, ok := p.DeliveryPolicyType.TryGet(); ok {
		w.WriteByte(':')
		w.WriteString(t.String())
	}

	return w.String()
}

func (p *ProjectionDeliveryPolicy) validate(ctx *validateContext) {
	validateFidelity(ctx, p.Fidelity)
	p.validateForExecution(ctx)
}

func (p *ProjectionDeliveryPolicy) validateForExecution(ctx *validateContext) {
	if !ctx.Options.ForExecution {
		return
	}

	if t, ok := p.DeliveryPolicyType.TryGet(); ok {
		if t != BroadcastProjectionDeliveryPolicyType {
			return
		}

		if p.Broadcast.PrimaryFirst.IsPresent() {
			return
		}
	}

	ctx.Invalid(ValueUnavailableError{reflect.TypeFor[dogma.ProjectionDeliveryPolicy]()})
}

func (p *ProjectionDeliveryPolicy) describe(ctx *describeContext) {
	describeFidelity(ctx, p.Fidelity)

	t, ok := p.DeliveryPolicyType.TryGet()

	if ok {
		ctx.Print(t.String(), " ")
	}

	ctx.Print("delivery policy")

	if ok && t == BroadcastProjectionDeliveryPolicyType {
		if v, _ := p.Broadcast.PrimaryFirst.TryGet(); v {
			ctx.Print(" (primary first)")
		}
	}

	ctx.Print("\n")
}
