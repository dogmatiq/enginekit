package config

import (
	"iter"
	"maps"
	"slices"

	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/enginekit/config/internal/renderer"
	"github.com/dogmatiq/enginekit/optional"
	"github.com/dogmatiq/enginekit/protobuf/identitypb"
)

// Application represents the (potentially invalid) configuration of a
// [dogma.Application] implementation.
type Application struct {
	EntityProperties

	// Handlers is a list of handlers registered with the application.
	HandlerComponents []Handler

	// Source is the instance of the entity from which the configuration
	// was sourced, if available.
	Source optional.Optional[dogma.Application]
}

// Identity returns the entity's identity.
//
// It panics if no single valid identity is configured.
func (a *Application) Identity() *identitypb.Identity {
	return resolveIdentity(a)
}

// Handlers returns the list of handlers configured for the application.
//
// It panics if the handlers are incomplete or invalid.
func (a *Application) Handlers() []Handler {
	handlers := slices.Clone(a.HandlerComponents)
	for i, h := range handlers {
		handlers[i] = clone(h)
	}

	normalizeChildren(strictContext(a), handlers...)

	return handlers
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
	ctx := strictContext(a)
	var set RouteSet

	for _, h := range a.HandlerComponents {
		set.merge(resolveRouteSet(ctx.NewChild(h), h))
	}

	return set
}

func (a *Application) String() string {
	return RenderDescriptor(a)
}

func (a *Application) renderDescriptor(ren *renderer.Renderer) {
	renderEntityDescriptor(ren, a)
}

func (a *Application) renderDetails(ren *renderer.Renderer) {
	renderEntityDetails(ren, a, a.Source)

	for _, h := range a.HandlerComponents {
		ren.IndentBullet()
		h.renderDetails(ren)
		ren.Dedent()
	}
}

func (a *Application) normalize(ctx *normalizationContext) {
	normalizeEntity(ctx, a, a.Source)
	normalizeChildren(ctx, a.HandlerComponents...)

	reportIdentityConflicts(ctx, a)
	reportRouteConflicts(ctx, a)
}

func (a *Application) clone() any {
	return &Application{
		clone(a.EntityProperties),
		cloneSlice(a.HandlerComponents),
		a.Source,
	}
}

// reportIdentityConflicts appends errors related to handlers that have
// identities that conflict with other handlers or the application itself.
func reportIdentityConflicts(ctx *normalizationContext, app *Application) {
	var (
		conflictingNames conflictDetector[string, QEntity]
		conflictingKeys  conflictDetector[string, QEntity]
	)

	entities := []QEntity{app}
	for _, h := range app.HandlerComponents {
		entities = append(entities, h)
	}

	for i, ent1 := range entities {
		for _, id1 := range ent1.CommonEntityProperties().IdentityComponents {
			k1, hasK1 := id1.Key.TryGet()
			n1, hasN1 := id1.Name.TryGet()

			for j, ent2 := range entities[i+1:] {
				for _, id2 := range ent2.CommonEntityProperties().IdentityComponents {
					k2, hasK2 := id2.Key.TryGet()
					n2, hasN2 := id2.Name.TryGet()

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

func reportRouteConflicts(ctx *normalizationContext, app *Application) {
	var conflictingRoutes conflictDetector[routeKey, Handler]

	for i, h1 := range app.HandlerComponents {
		for _, r1 := range h1.CommonHandlerProperties().RouteComponents {
			k1, ok := routeKeyOf(r1)
			if !ok {
				continue
			}

			if k1.RouteType != HandlesCommandRouteType && k1.RouteType != RecordsEventRouteType {
				continue
			}

			for j, h2 := range app.HandlerComponents[i+1:] {
				for _, r2 := range h2.CommonHandlerProperties().RouteComponents {
					k2, ok := routeKeyOf(r2)
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
