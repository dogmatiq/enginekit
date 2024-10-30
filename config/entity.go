package config

import (
	"fmt"
	"reflect"
	"slices"
	"strings"

	"github.com/dogmatiq/enginekit/internal/typename"
	"github.com/dogmatiq/enginekit/optional"
	"github.com/dogmatiq/enginekit/protobuf/identitypb"
	"github.com/dogmatiq/enginekit/protobuf/uuidpb"
)

// An Entity is a [Component] that that represents the configuration of a Dogma
// entity, which is a type that implements one of the following interfaces:
//
//   - [dogma.Application]
//   - [dogma.AggregateMessageHandler]
//   - [dogma.ProcessMessageHandler]
//   - [dogma.IntegrationMessageHandler]
//   - [dogma.ProjectionMessageHandler]
//
// See [Handler] for a more specific interface that represents the
// configuration of a Dogma handler.
type Entity interface {
	Component

	// Identity returns the entity's identity.
	//
	// It panics the configuration does not specify a singular valid identity.
	Identity() *identitypb.Identity

	// RouteSet returns the routes configured for the entity.
	//
	// It panics if the configuration does not specify a complete set of valid
	// routes for the entity and its constituents.
	RouteSet() RouteSet

	// EntityProperties returns the properties common to all [Entity] types.
	EntityProperties() *EntityCommon
}

// EntityCommon contains the properties common to all [Entity] types.
type EntityCommon struct {
	ComponentCommon

	TypeName           optional.Optional[string]
	IdentityComponents []*Identity
}

// EntityProperties returns the properties common to all [Entity] types.
func (p *EntityCommon) EntityProperties() *EntityCommon {
	return p
}

// UnidentifiedEntityError indicates that an [Entity] has been configured
// without an [Identity].
type UnidentifiedEntityError struct{}

func (e UnidentifiedEntityError) Error() string {
	return "entity has no identity"
}

// AmbiguouslyIdentifiedEntityError indicates that an [Entity] has been
// configured with more than one [Identity].
type AmbiguouslyIdentifiedEntityError struct {
	Identities []*Identity
}

func (e AmbiguouslyIdentifiedEntityError) Error() string {
	return fmt.Sprintf(
		"entity has %d identities",
		len(e.Identities),
	)
}

// EntityUnavailableError indicates that an [Entity] cannot produce the actual
// entity value it represents because there is insufficient runtime type
// information available.
type EntityUnavailableError struct {
	EntityType reflect.Type
}

func (e EntityUnavailableError) Error() string {
	return fmt.Sprintf("%s is unavailable", e.EntityType)
}

func validateEntity[T any](
	ctx *validateContext,
	e Entity,
	source optional.Optional[T],
	funcs ...func(*validateContext),
) {
	validateComponent(
		ctx,
		func(ctx *validateContext) {
			validateEntityIdentities(ctx, e)

			for _, fn := range funcs {
				fn(ctx)
			}

			typeName, hasTypeName := e.EntityProperties().TypeName.TryGet()
			source, hasSource := source.TryGet()

			if hasSource {
				if !hasTypeName {
					ctx.Malformed("Source is present, but TypeName is not")
				} else if typeName != typename.Of(source) {
					ctx.Malformed("TypeName does not match Source: %q != %q", typeName, typename.Of(source))
				}
			}
		},
	)
}

func validateEntityIdentities(ctx *validateContext, e Entity) {
	identities := e.EntityProperties().IdentityComponents

	if len(identities) == 0 {
		ctx.Invalid(UnidentifiedEntityError{})
	} else if len(identities) > 1 {
		ctx.Invalid(AmbiguouslyIdentifiedEntityError{slices.Clone(identities)})
	}

	for _, i := range identities {
		ctx.ValidateChild(i)
	}
}

func resolveIdentity(e Entity) *identitypb.Identity {
	validateEntityIdentities(newResolutionContext(e), e)
	id := e.EntityProperties().IdentityComponents[0]

	return &identitypb.Identity{
		Name: id.Name.Get(),
		Key:  uuidpb.MustParse(id.Key.Get()),
	}
}

func resolveInterface[T any](e Entity, source optional.Optional[T]) T {
	ctx := newResolutionContext(e)

	v, ok := source.TryGet()

	if !ok {
		ctx.Invalid(EntityUnavailableError{reflect.TypeFor[T]()})
	}

	return v
}

func entityLabel(e Entity) string {
	return strings.ToLower(
		reflect.TypeOf(e).Elem().Name(),
	)
}

func stringifyEntity(e Entity) string {
	var w strings.Builder

	w.WriteString(entityLabel(e))

	if n, ok := e.EntityProperties().TypeName.TryGet(); ok {
		n = typename.Unqualified(n)
		n = strings.TrimPrefix(n, "*")
		w.WriteByte(':')
		w.WriteString(n)
	}

	return w.String()
}

func describeEntity[T any](
	ctx *describeContext,
	e Entity,
	source optional.Optional[T],
) {
	ctx.DescribeFidelity()
	ctx.Print(entityLabel(e))

	p := e.EntityProperties()

	if n, ok := p.TypeName.TryGet(); ok {
		ctx.Print(" ")
		ctx.Print(n)

		if !source.IsPresent() {
			ctx.Print(" (value unavailable)")
		}
	}

	ctx.Print("\n")
	ctx.DescribeErrors()

	for _, i := range p.IdentityComponents {
		ctx.DescribeChild(i)
	}
}
