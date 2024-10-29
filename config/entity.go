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

	identities() []*Identity
}

// EntityCommon is a partial implementation of [Entity].
type EntityCommon[T any] struct {
	ComponentCommon

	SourceTypeName     optional.Optional[string]
	Source             optional.Optional[T]
	IdentityComponents []*Identity
}

// Identity returns the entity's identity.
//
// It panics the configuration does not specify a singular valid identity.
func (e *EntityCommon[T]) Identity() *identitypb.Identity {
	e.validateIdentities(nil)

	id := e.IdentityComponents[0]

	return &identitypb.Identity{
		Name: id.Name.Get(),
		Key:  uuidpb.MustParse(id.Key.Get()),
	}
}

// Interface returns the entity represented by the configuration, if available.
func (e *EntityCommon[T]) Interface() T {
	if v, ok := e.Source.TryGet(); ok {
		return v
	}
	panic(EntityUnavailableError{reflect.TypeFor[T]()})
}

func (e *EntityCommon[T]) String() string {
	var w strings.Builder

	w.WriteString(entityClass[T]())

	if typeName, ok := e.SourceTypeName.TryGet(); ok {
		typeName = typename.Unqualified(typeName)
		typeName = strings.TrimPrefix(typeName, "*")
		w.WriteByte(':')
		w.WriteString(typeName)
	}

	return w.String()
}

func (e *EntityCommon[T]) identities() []*Identity {
	return e.IdentityComponents
}

func (e *EntityCommon[T]) validate(ctx *validateContext) {
	e.ComponentCommon.validate(ctx)
	e.validateIdentities(ctx)
}

func (e *EntityCommon[T]) validateIdentities(ctx *validateContext) {
	if len(e.IdentityComponents) == 0 {
		ctx.Fail(UnidentifiedEntityError{})
	} else if len(e.IdentityComponents) > 1 {
		ctx.Fail(AmbiguouslyIdentifiedEntityError{slices.Clone(e.IdentityComponents)})
	}

	for _, i := range e.IdentityComponents {
		ctx.ValidateChild(i)
	}
}

func (e *EntityCommon[T]) describe(ctx *describeContext) {
	ctx.DescribeFidelity()
	ctx.Print(entityClass[T]())

	if typeName, ok := e.SourceTypeName.TryGet(); ok {
		ctx.Print(" ")
		ctx.Print(typeName)

		if !e.Source.IsPresent() {
			ctx.Print(" (value unavailable)")
		}
	}

	ctx.Print("\n")
	ctx.DescribeErrors()

	for _, i := range e.IdentityComponents {
		ctx.DescribeChild(i)
	}
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

// entityClass returns the "class" of an entity, e.g. "application",
// "aggregate", etc. based on the Dogma interface type.
func entityClass[T any]() string {
	var w strings.Builder

	for i, r := range reflect.TypeFor[T]().Name() {
		if r >= 'A' && r <= 'Z' {
			if i != 0 {
				break
			}
			w.WriteRune(r - 'A' + 'a')
		} else {
			w.WriteRune(r)
		}
	}

	return w.String()
}