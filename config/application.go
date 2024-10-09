package config

import (
	"fmt"
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

// Identity returns the entity's identity.
//
// It panics if no single valid identity is configured.
func (a Application) Identity() Identity {
	return normalizedIdentity(a)
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
	Handlers         []Handler
	ConflictingRoute Route
}

func (e ConflictingRouteError) Error() string {
	rt := e.ConflictingRoute.RouteType.Get()

	verb := "handled"
	switch rt {
	case ExecutesCommandRoute:
		verb = "executed"
	case RecordsEventRoute:
		verb = "recorded"
	case SchedulesTimeoutRoute:
		verb = "scheduled"
	}

	return fmt.Sprintf(
		"handlers have conflicting %q routes: %s is %s by %s",
		rt,
		e.ConflictingRoute.MessageType.Get(),
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
	var conflictingRoutes conflictDetector[Route, Handler]

	for i, h1 := range app.ConfiguredHandlers {
		for _, r1 := range h1.routes() {
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

			for j, h2 := range app.ConfiguredHandlers[i+1:] {
				for _, r2 := range h2.routes() {
					conflictingRoutes.Add(
						r1, i, h1,
						r2, j, h2,
					)
				}
			}
		}
	}

	for route, handlers := range conflictingRoutes.All() {
		ctx.Fail(ConflictingRouteError{handlers, route})
	}
}
