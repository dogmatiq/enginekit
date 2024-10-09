package config

import (
	"errors"
	"fmt"
	"iter"
	"maps"
	"slices"

	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/enginekit/optional"
)

// Application represents the (potentially invalid) configuration of a
// [dogma.Application] implementation.
type Application struct {
	// Impl contains information about the type that produced the
	// configuration, if available.
	Impl optional.Optional[Implementation[dogma.Application]]

	// ConfiguredIdentities is the list of (potentially invalid or duplicated)
	// identities configured for the application.
	ConfiguredIdentities []Identity

	// ConfiguredHandlers is a list of (potentially invalid, incomplete or
	// conflicting) handlers configured for the application.
	ConfiguredHandlers []Handler

	// IsExhaustive is true if the complete configuration was loaded. It may be
	// false, for example, when attempting to load configuration using static
	// analysis, but the code depends on runtime type information.
	IsExhaustive bool
}

func (a Application) String() string {
	return stringify("application", a.Impl, a.ConfiguredIdentities)
}

// Aggregates returns a list of [dogma.AggregateMessageHandler] implementations
// that are registered with the application.
func (a Application) Aggregates() []Aggregate {
	return normalizedHandlers[Aggregate](a.ConfiguredHandlers)
}

// Processes returns a list of [dogma.ProcessMessageHandler] implementations
// that are registered with the application.
func (a Application) Processes() []Process {
	return normalizedHandlers[Process](a.ConfiguredHandlers)
}

// Integrations returns a list of [dogma.IntegrationMessageHandler]
// implementations that are registered with the application.
func (a Application) Integrations() []Integration {
	return normalizedHandlers[Integration](a.ConfiguredHandlers)
}

// Projections returns a list of [dogma.ProjectionMessageHandler]
// implementations that are registered with the application.
func (a Application) Projections() []Projection {
	return normalizedHandlers[Projection](a.ConfiguredHandlers)
}

// Identity returns the entity's identity.
//
// It panics if no single valid identity is configured.
func (a Application) Identity() Identity {
	return normalizedIdentity(a)
}

func (a Application) configuredIdentities() []Identity { return a.ConfiguredIdentities }

func (a Application) normalize(opts validationOptions) (_ Entity, errs error) {
	normalizeIdentitiesInPlace(opts, a, &errs, &a.ConfiguredIdentities)
	normalizeHandlersInPlace(opts, a, &errs, &a.ConfiguredHandlers)
	return a, errs
}

// InvalidHandlerError is an error that occurs when an application contains an
// invalid handler.
type InvalidHandlerError struct {
	Application Application
	Handler     any
	Cause       error
}

func (e InvalidHandlerError) Error() string {
	return fmt.Sprintf("%s contains an invalid handler: %s", e.Application, e.Cause)
}

func (e InvalidHandlerError) Unwrap() error {
	return e.Cause
}

func normalizeHandlersInPlace(
	opts validationOptions,
	app Application,
	errs *error,
	handlers *[]Handler,
) {
	*handlers = slices.Clone(*handlers)

	for i, h := range *handlers {
		norm, err := normalize(opts, h)
		(*handlers)[i] = norm

		if err != nil {
			*errs = errors.Join(
				*errs,
				InvalidHandlerError{app, h, err},
			)
		}
	}

	reportIdentityConflicts(app, errs, *handlers)
	reportRouteConflicts(errs, *handlers)
}

// reportIdentityConflicts appends errors related to handlers that have
// identities that conflict with other handlers or the application itself.
func reportIdentityConflicts(
	app Application,
	errs *error,
	handlers []Handler,
) {
	var (
		conflictingIDs   conflicts[Identity, Entity]
		conflictingNames conflicts[string, Entity]
		conflictingKeys  conflicts[string, Entity]
	)

	entities := []Entity{app}
	for _, h := range handlers {
		entities = append(entities, h)
	}

	for i, ent1 := range entities {
		for _, id1 := range ent1.configuredIdentities() {
			for j, ent2 := range entities[i+1:] {
				for _, id2 := range ent2.configuredIdentities() {
					if conflictingIDs.Add(
						id1, i, ent1,
						id2, j, ent2,
					) {
						continue
					}

					if conflictingKeys.Add(
						id1.Key, i, ent1,
						id2.Key, j, ent2,
					) {
						continue
					}

					// Index 0 is the application, which is allowed to have
					// the same name as one of its handlers.
					if i > 0 {
						conflictingNames.Add(
							id1.Name, i, ent1,
							id2.Name, j, ent2,
						)
					}
				}
			}
		}
	}

	for id, entities := range conflictingIDs.All() {
		*errs = errors.Join(*errs, IdentityConflictError{entities, id})
	}

	for name, entities := range conflictingNames.All() {
		*errs = errors.Join(*errs, IdentityNameConflictError{entities, name})
	}

	for key, entities := range conflictingKeys.All() {
		*errs = errors.Join(*errs, IdentityKeyConflictError{entities, key})
	}
}

func reportRouteConflicts(
	errs *error,
	handlers []Handler,
) {
	var conflictingRoutes conflicts[Route, Handler]

	for i, h1 := range handlers {
		for _, r1 := range h1.configuredRoutes() {
			t, ok := r1.RouteType.TryGet()
			if !ok {
				continue
			}

			if t != HandlesCommandRoute && t != RecordsEventRoute {
				continue
			}

			if !r1.MessageType.IsPresent() {
				continue
			}

			for j, h2 := range handlers[i+1:] {
				for _, r2 := range h2.configuredRoutes() {
					conflictingRoutes.Add(
						r1, i, h1,
						r2, j, h2,
					)
				}
			}
		}
	}

	for route, handlers := range conflictingRoutes.All() {
		*errs = errors.Join(*errs, RouteConflictError{handlers, route})
	}
}

func normalizedHandlers[T Handler]([]Handler) []T {
	panic("not implemented")
}

// conflicts keeps track of conflicting identity components across multiple
// entities.
type conflicts[K comparable, V any] struct {
	m map[K]map[int]V
}

func (t *conflicts[K, V]) Add(
	k1 K, i int, e1 V,
	k2 K, j int, e2 V,
) bool {
	if k1 != k2 {
		return false
	}

	if t.m == nil {
		t.m = map[K]map[int]V{}
	}

	if t.m[k1] == nil {
		t.m[k2] = map[int]V{}
	}

	t.m[k1][i] = e1
	t.m[k1][i+1+j] = e2

	return true
}

func (t *conflicts[K, V]) All() iter.Seq2[K, []V] {
	return func(yield func(K, []V) bool) {
		for k, indices := range t.m {
			sortedIndices := slices.Sorted(maps.Keys(indices))
			sortedValues := make([]V, len(sortedIndices))

			for i, j := range sortedIndices {
				sortedValues[i] = indices[j]
			}

			if !yield(k, sortedValues) {
				return
			}
		}
	}
}
