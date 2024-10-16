package config

import (
	"fmt"
	"slices"

	"github.com/dogmatiq/enginekit/internal/ioutil"
	"github.com/dogmatiq/enginekit/protobuf/identitypb"
)

// A Component is some element of the configuration of a Dogma application.
type Component interface {
	fmt.Stringer

	// Fidelity returns information about how well the configuration represents
	// the actual configuration that would be used at runtime.
	Fidelity() Fidelity

	renderDescriptor(*ioutil.Renderer)
	renderDetails(*ioutil.Renderer)

	clone() Component
	normalize(*normalizationContext)
}

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

// A Handler is a specialization of [Entity] that represents configuration of a
// Dogma message handler.
type Handler interface {
	Entity

	// HandlerType returns [HandlerType] of the handler.
	HandlerType() HandlerType

	// IsDisabled returns true if the handler was disabled via the configurer.
	IsDisabled() bool

	routes() []*Route
}

// Fidelity is a bit-field that describes how well a [Component] configuration
// represents the actual configuration that would be used at runtime.
//
// Importantly, it does not describe the validity of the configuration itself.
type Fidelity int

const (
	// Immaculate is the [Fidelity] value that indicates the configuration is an
	// exact match for the actual configuration that would be used at runtime.
	Immaculate Fidelity = 0

	// Speculative is a [Fidelity] flag that indicates that the [Component] is
	// only present in the configuration under certain conditions, and that
	// those conditions could not be evaluated at configuration time.
	Speculative Fidelity = 1 << iota

	// Incomplete is a [Fidelity] flag that indicates that the [Component] has
	// some configuration that could not be resolved accurately at configuration
	// time.
	//
	// Most commonly this is occurs during static analysis of code that uses
	// interfaces that cannot be followed statically.
	//
	// Its absence means that all of the _available_ configuration logic was
	// applied; it does not imply that all _mandatory_ configuration is present.
	Incomplete
)

var (
	_ Component = (*Identity)(nil)
	_ Component = (*Route)(nil)

	_ Entity = (*Application)(nil)

	_ Handler = (*Aggregate)(nil)
	_ Handler = (*Process)(nil)
	_ Handler = (*Integration)(nil)
	_ Handler = (*Projection)(nil)
)

func clone[T Component](components []T) []T {
	clones := slices.Clone(components)

	for i, c := range components {
		clones[i] = c.clone().(T)
	}

	return clones
}

func cloneInPlace[T Component](components *[]T) {
	*components = clone(*components)
}
