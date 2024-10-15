package config

import (
	"iter"
	"maps"
	"slices"

	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/enginekit/optional"
	"github.com/dogmatiq/enginekit/protobuf/identitypb"
)

// ApplicationAsConfigured contains the raw unvalidated properties of an
// [Application].
type ApplicationAsConfigured struct {
	// Source describes the type and value that produced the configuration, if
	// available.
	Source optional.Optional[Value[dogma.Application]]

	// Identities is the list of identities configured for the application.
	Identities []*Identity

	// Handlers is a list of handlers registered with the application.
	Handlers []Handler

	// Fidelity describes the configuration's accuracy in comparison to the
	// actual configuration that would be used at runtime.
	Fidelity Fidelity
}

// Application represents the (potentially invalid) configuration of a
// [dogma.Application] implementation.
type Application struct {
	AsConfigured ApplicationAsConfigured
}

func (a *Application) String() string {
	return renderEntity("application", a, a.AsConfigured.Source)
}

// Identity returns the entity's identity.
//
// It panics if no single valid identity is configured.
func (a *Application) Identity() *identitypb.Identity {
	return finalizeIdentity(newFinalizeContext(a), a)
}

// Fidelity returns information about how well the configuration represents
// the actual configuration that would be used at runtime.
func (a *Application) Fidelity() Fidelity {
	return a.AsConfigured.Fidelity
}

// Interface returns the [dogma.Application] instance that the configuration
// represents, or panics if it is not available.
func (a *Application) Interface() dogma.Application {
	return a.AsConfigured.Source.Get().Value.Get()
}

// Handlers returns the list of handlers configured for the application.
//
// It panics if the handlers are incomplete or invalid.
func (a *Application) Handlers() []Handler {
	return normalizeHandlers(newFinalizeContext(a), a)
}

// HandlerByName returns the [Handler] with the given name, or false if no such
// handler has been configured.
//
// It panics if the handlers are incomplete or invalid.
func (a *Application) HandlerByName(name string) (Handler, bool) {
	for _, h := range a.Handlers() {
		if h.Identity().Name == name {
			return h, true
		}
	}

	return nil, false
}

// RouteSet returns the routes configured for the entity.
//
// It panics if the route configuration is incomplete or invalid.
func (a *Application) RouteSet() RouteSet {
	ctx := newFinalizeContext(a)
	var set RouteSet

	for _, h := range a.AsConfigured.Handlers {
		set.merge(finalizeRouteSet(ctx.NewChild(h), h))
	}

	return set
}

func (a *Application) identitiesAsConfigured() []*Identity {
	return a.AsConfigured.Identities
}

func (a *Application) normalize(ctx *normalizeContext) Component {
	a.AsConfigured.Fidelity, a.AsConfigured.Source = normalizeValue(ctx, a.AsConfigured.Fidelity, a.AsConfigured.Source)
	a.AsConfigured.Identities = normalizeIdentities(ctx, a)
	a.AsConfigured.Handlers = normalizeHandlers(ctx, a)

	detectIdentityConflicts(ctx, a)
	detectRouteConflicts(ctx, a)

	return a
}

func normalizeHandlers(ctx *normalizeContext, app *Application) []Handler {
	handlers := slices.Clone(app.AsConfigured.Handlers)

	for i, h := range handlers {
		handlers[i] = normalize(ctx, h)
	}

	return handlers
}

// detectIdentityConflicts appends errors related to handlers that have
// identities that conflict with other handlers or the application itself.
func detectIdentityConflicts(ctx *normalizeContext, app *Application) {
	var (
		conflictingNames conflictDetector[string, Entity]
		conflictingKeys  conflictDetector[string, Entity]
	)

	entities := []Entity{app}
	for _, h := range app.AsConfigured.Handlers {
		entities = append(entities, h)
	}

	for i, ent1 := range entities {
		for _, id1 := range ent1.identitiesAsConfigured() {
			k1, hasK1 := id1.AsConfigured.Key.TryGet()
			n1, hasN1 := id1.AsConfigured.Name.TryGet()

			for j, ent2 := range entities[i+1:] {
				for _, id2 := range ent2.identitiesAsConfigured() {
					k2, hasK2 := id2.AsConfigured.Key.TryGet()
					n2, hasN2 := id2.AsConfigured.Name.TryGet()

					if hasK1 && hasK2 {
						conflictingKeys.Add(
							k1, i, ent1,
							k2, j, ent2,
						)
					}

					// Index 0 is the application, which is allowed to have the
					// same name as one of its handlers.
					if i > 0 && hasN1 && hasN2 {
						conflictingNames.Add(
							n1, i, ent1,
							n2, j, ent2,
						)
					}
				}
			}
		}
	}

	for name, entities := range conflictingNames.All() {
		ctx.Fail(IdentityNameConflictError{entities, name})
	}

	for key, entities := range conflictingKeys.All() {
		ctx.Fail(IdentityKeyConflictError{entities, key})
	}
}

func detectRouteConflicts(ctx *normalizeContext, app *Application) {
	var conflictingRoutes conflictDetector[routeKey, Handler]

	for i, h1 := range app.AsConfigured.Handlers {
		for _, r1 := range h1.routesAsConfigured() {
			k1, ok := r1.key()
			if !ok {
				continue
			}

			if k1.RouteType != HandlesCommandRouteType && k1.RouteType != RecordsEventRouteType {
				continue
			}

			for j, h2 := range app.AsConfigured.Handlers[i+1:] {
				for _, r2 := range h2.routesAsConfigured() {
					k2, ok := r2.key()
					if !ok {
						continue
					}

					conflictingRoutes.Add(
						k1, i, h1,
						k2, j, h2,
					)
				}
			}
		}
	}

	for key, handlers := range conflictingRoutes.All() {
		ctx.Fail(ConflictingRouteError{handlers, key.RouteType, key.MessageTypeName})
	}
}

type conflictDetector[T comparable, S any] struct {
	m map[T]map[int]S
}

func (t *conflictDetector[T, S]) Add(
	v1 T, i int, src1 S,
	v2 T, j int, src2 S,
) bool {
	if v1 != v2 {
		return false
	}

	if t.m == nil {
		t.m = map[T]map[int]S{}
	}

	if t.m[v1] == nil {
		t.m[v2] = map[int]S{}
	}

	t.m[v1][i] = src1
	t.m[v1][i+1+j] = src2

	return true
}

func (t *conflictDetector[T, S]) All() iter.Seq2[T, []S] {
	return func(yield func(T, []S) bool) {
		for v, indices := range t.m {
			sortedIndices := slices.Sorted(maps.Keys(indices))
			sortedComponents := make([]S, len(sortedIndices))

			for i, j := range sortedIndices {
				sortedComponents[i] = indices[j]
			}

			if !yield(v, sortedComponents) {
				return
			}
		}
	}
}
