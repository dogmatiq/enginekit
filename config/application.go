package config

import (
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
	// ConfigurationSource contains information about the type and value that
	// produced the configuration, if available.
	ConfigurationSource optional.Optional[Source[dogma.Application]]

	// ConfiguredIdentities is the list of (potentially invalid or duplicated)
	// identities configured for the application.
	ConfiguredIdentities []Identity

	// ConfiguredHandlers is a list of (potentially invalid, incomplete or
	// conflicting) handlers configured for the application.
	ConfiguredHandlers []Handler

	// ConfigurationIsExhaustive is true if the entire configuration was loaded.
	ConfigurationIsExhaustive bool
}

func (a Application) String() string {
	return stringify("application", a, a.ConfigurationSource)
}

// Identity returns the entity's identity.
//
// It panics if no single valid identity is configured.
func (a Application) Identity() Identity {
	return normalizedIdentity(a)
}

// IsExhaustive returns true if the entire configuration was loaded.
func (a Application) IsExhaustive() bool {
	return a.ConfigurationIsExhaustive
}

// Interface returns the [dogma.Application] instance that the configuration
// represents, or panics if it is not available.
func (a Application) Interface() dogma.Application {
	return a.ConfigurationSource.Get().Value.Get()
}

// HandlerByName returns the [Handler] with the given name, or false if no such
// handler has been configured.
func (a Application) HandlerByName(name string) (Handler, bool) {
	ctx := &normalizationContext{
		Component: a,
	}

	handlers := normalizeHandlers(ctx, a)

	if err := ctx.Err(); err != nil {
		panic(err)
	}

	for _, h := range handlers {
		if h.Identity().Name == name {
			return h, true
		}
	}

	return nil, false
}

func (a Application) identities() []Identity {
	return a.ConfiguredIdentities
}

func (a Application) normalize(ctx *normalizationContext) Component {
	a.ConfiguredIdentities = normalizeIdentities(ctx, a)
	a.ConfiguredHandlers = normalizeHandlers(ctx, a)

	detectIdentityConflicts(ctx, a)
	detectRouteConflicts(ctx, a)

	return a
}

// IdentityConflictError indicates that more than one [Entity] within the same
// [Application] shares the same [Identity].
type IdentityConflictError struct {
	Entities            []Entity
	ConflictingIdentity Identity
}

func (e IdentityConflictError) Error() string {
	return fmt.Sprintf(
		"entities have conflicting identities: %s is shared by %s",
		e.ConflictingIdentity,
		renderList(e.Entities),
	)
}

// IdentityNameConflictError indicates that more than one [Entity] within the
// same [Application] is shares the same "name" component of an [Identity].
type IdentityNameConflictError struct {
	Entities        []Entity
	ConflictingName string
}

func (e IdentityNameConflictError) Error() string {
	return fmt.Sprintf(
		"entities have conflicting identities: the %q name is shared by %s",
		e.ConflictingName,
		renderList(e.Entities),
	)
}

// IdentityKeyConflictError indicates that more than one [Entity] within the
// same [Application] is shares the same "key" component of an [Identity].
type IdentityKeyConflictError struct {
	Entities       []Entity
	ConflictingKey string
}

func (e IdentityKeyConflictError) Error() string {
	return fmt.Sprintf(
		"entities have conflicting identities: the %q key is shared by %s",
		e.ConflictingKey,
		renderList(e.Entities),
	)
}

// ConflictingRouteError indicates that more than one [Handler] within the same
// [Application] is configured with routes for the same [MessageType] in a
// manner that is not permitted.
//
// For example, no two handlers can handle commands of the same type, though any
// number of handlers may handle events of the same type.
type ConflictingRouteError struct {
	Handlers                   []Handler
	ConflictingRouteType       RouteType
	ConflictingMessageTypeName string
}

func (e ConflictingRouteError) Error() string {
	verb := "handled"
	switch e.ConflictingRouteType {
	case ExecutesCommandRouteType:
		verb = "executed"
	case RecordsEventRouteType:
		verb = "recorded"
	case SchedulesTimeoutRouteType:
		verb = "scheduled"
	}

	return fmt.Sprintf(
		"handlers have conflicting %q routes: %s is %s by %s",
		e.ConflictingRouteType,
		e.ConflictingMessageTypeName,
		verb,
		renderList(e.Handlers),
	)
}

func normalizeHandlers(ctx *normalizationContext, a Application) []Handler {
	handlers := slices.Clone(a.ConfiguredHandlers)

	for i, h := range handlers {
		handlers[i] = normalize(ctx, h)
	}

	return handlers
}

// detectIdentityConflicts appends errors related to handlers that have
// identities that conflict with other handlers or the application itself.
func detectIdentityConflicts(ctx *normalizationContext, app Application) {
	var (
		conflictingIDs   conflictDetector[Identity, Entity]
		conflictingNames conflictDetector[string, Entity]
		conflictingKeys  conflictDetector[string, Entity]
	)

	entities := []Entity{app}
	for _, h := range app.ConfiguredHandlers {
		entities = append(entities, h)
	}

	for i, ent1 := range entities {
		for _, id1 := range ent1.identities() {
			for j, ent2 := range entities[i+1:] {
				for _, id2 := range ent2.identities() {
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
		ctx.Fail(IdentityConflictError{entities, id})
	}

	for name, entities := range conflictingNames.All() {
		ctx.Fail(IdentityNameConflictError{entities, name})
	}

	for key, entities := range conflictingKeys.All() {
		ctx.Fail(IdentityKeyConflictError{entities, key})
	}
}

func detectRouteConflicts(ctx *normalizationContext, app Application) {
	var conflictingRoutes conflictDetector[routeKey, Handler]

	for i, h1 := range app.ConfiguredHandlers {
		for _, r1 := range h1.routes() {
			k1, ok := r1.key()
			if !ok {
				continue
			}

			if k1.RouteType != HandlesCommandRouteType && k1.RouteType != RecordsEventRouteType {
				continue
			}

			for j, h2 := range app.ConfiguredHandlers[i+1:] {
				for _, r2 := range h2.routes() {
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
