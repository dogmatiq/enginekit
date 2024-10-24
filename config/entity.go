package config

import (
	"github.com/dogmatiq/enginekit/protobuf/identitypb"
	"github.com/dogmatiq/enginekit/protobuf/uuidpb"
)

// An Entity is a [Component] that represents the configuration of some
// configurable Dogma entity; that is, any type with a Configure() method that
// accepts one of the Dogma "configurer" interfaces.
type Entity interface {
	Component

	// Identity returns the entity's identity.
	//
	// It panics if no single valid identity is configured.
	Identity() *identitypb.Identity

	// RouteSet returns the routes configured for the entity.
	//
	// It panics if the route configuration is incomplete or invalid.
	RouteSet() RouteSet

	identities() []*Identity
}

// EntityTrait is a partial implementation of [Entity].
type EntityTrait[T any] struct {
	ComponentTrait

	// Source describes the type and value that produced the configuration.
	Source Value[T]

	// Identities is the list of identities configured for the handler.
	Identities []*Identity
}

// Identity returns the entity's identity.
//
// It panics if no single valid identity is configured.
func (e EntityTrait[T]) identity(ctx *normalizationContext) *identitypb.Identity {
	identities := clone(e.Identities)
	normalizeChildren(ctx, identities)
	reportIdentityErrors(ctx, identities)

	id := identities[0].AsConfigured

	return &identitypb.Identity{
		Name: id.Name.Get(),
		Key:  uuidpb.MustParse(id.Key.Get()),
	}
}

// Interface returns the instance that the configuration represents, or panics
// if it is not available.
func (e EntityTrait[T]) Interface() T {
	return e.Source.Value.Get()
}
