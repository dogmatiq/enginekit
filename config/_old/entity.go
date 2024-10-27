package config

import (
	"reflect"

	"github.com/dogmatiq/enginekit/internal/typename"
	"github.com/dogmatiq/enginekit/optional"
	"github.com/dogmatiq/enginekit/protobuf/identitypb"
)

// An QEntity is a [Component] that represents the configuration of a Dogma
// entity of type T; that is, any type with a Configure() method that accepts
// one of the Dogma "configurer" interfaces.
type QEntity[T any] interface {
	Component

	// Identity returns the entity's identity.
	//
	// It panics if no single valid identity is configured.
	Identity() *identitypb.Identity

	// RouteSet returns the routes configured for the entity.
	//
	// It panics if the route configuration is invalid or cannot be determined
	// completely.
	RouteSet() RouteSet

	// CommonEntityProperties returns the (possibly invalid or incomplete)
	// properties of the entity.
	CommonEntityProperties() *EntityProperties
}

// EntityProperties contains the properties common to all [QEntity]
// implementations.
type EntityProperties struct {
	ComponentProperties

	// TypeName is the fully-qualified name of the concrete Go type that
	// implements the entity.
	TypeName optional.Optional[string]

	// Identities is the list of identities configured for the handler.
	IdentityComponents []*Identity
}

// CommonEntityProperties returns the (possibly invalid or incomplete)
// properties of the entity.
func (p *EntityProperties) CommonEntityProperties() *EntityProperties {
	return p
}

func (p EntityProperties) clone() any {
	cloneInPlace(&p.ComponentProperties)
	cloneSliceInPlace(&p.IdentityComponents)
	return p
}

func normalizeEntity[T any](
	ctx *normalizationContext,
	e QEntity,
	v optional.Optional[T],
) {
	p := e.CommonEntityProperties()

	normalizeEntityTypeName(ctx, p, v)
	normalizeChildren(ctx, p.IdentityComponents...)

	reportIdentityErrors(ctx, p.IdentityComponents)
}

func normalizeEntityTypeName[T any](
	ctx *normalizationContext,
	p *EntityProperties,
	v optional.Optional[T],
) {
	name, nameOK := p.TypeName.TryGet()
	value, valueOK := v.TryGet()

	if !nameOK {
		p.Fidelity |= Incomplete
	}

	if valueOK {
		nameFromValue := typename.Get(reflect.TypeOf(value))
		if nameOK && name != nameFromValue {
			ctx.Fail(TypeNameMismatchError{name, nameFromValue})
		}

		p.TypeName = optional.Some(nameFromValue)
	} else if ctx.Options.RequireValues {
		ctx.Fail(RuntimeValueUnavailableError{reflect.TypeFor[T]()})
	}
}
