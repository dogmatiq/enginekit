package config

import (
	"fmt"
	"slices"

	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/enginekit/optional"
	"github.com/dogmatiq/enginekit/protobuf/uuidpb"
)

// Application is an [Entity] that represents the configuration of a
// [dogma.Application].
type Application struct {
	EntityCommon[dogma.Application]

	// HandlerComponents is the set of [Handler] components that have been
	// registered with the application.
	HandlerComponents []Handler
}

// RouteSet returns the routes configured for the entity.
//
// It panics if the configuration does not specify a complete set of valid
// routes for the entity and its constituents.
func (a *Application) RouteSet() RouteSet {
	panic("not implemented")
}

func (a *Application) validate(ctx *validateContext) {
	a.EntityCommon.validate(ctx)
	detectIdentityConflicts(ctx, a)
	detectRouteConflicts(ctx, a)

	for _, h := range a.HandlerComponents {
		ctx.ValidateChild(h)
	}
}

func (a *Application) describe(ctx *describeContext) {
	a.EntityCommon.describe(ctx)

	for _, h := range a.HandlerComponents {
		ctx.DescribeChild(h)
	}
}

// IdentityNameConflictError indicates that more than one [QEntity] within the
// same [Application] is shares the same "name" element of an [Identity].
type IdentityNameConflictError struct {
	ConflictingName string
	Entities        []Entity
}

func (e IdentityNameConflictError) Error() string {
	return fmt.Sprintf(
		"identity name %q is shared by %d entities",
		e.ConflictingName,
		len(e.Entities),
	)
}

// IdentityKeyConflictError indicates that more than one [QEntity] within the
// same [Application] is shares the same "key" element of an [Identity].
type IdentityKeyConflictError struct {
	ConflictingKey string
	Entities       []Entity
}

func (e IdentityKeyConflictError) Error() string {
	return fmt.Sprintf(
		"identity key %q is shared by %d entities",
		e.ConflictingKey,
		len(e.Entities),
	)
}

// RouteConflictError indicates that more than one [Handler] within the same
// [Application] is configured with routes for the same [MessageType] in a
// manner that is not permitted.
//
// For example, no two handlers can handle commands of the same type, though any
// number of handlers may handle events of the same type.
type RouteConflictError struct {
	ConflictingRouteType       RouteType
	ConflictingMessageTypeName string
	Handlers                   []Handler
}

func (e RouteConflictError) Error() string {
	return fmt.Sprintf(
		"%s route for %s is shared by %d handlers",
		e.ConflictingRouteType,
		e.ConflictingMessageTypeName,
		len(e.Handlers),
	)
}

// detectIdentityConflicts reports errors related to handlers that have
// identities that conflict with other handlers or the application itself.
func detectIdentityConflicts(ctx *validateContext, app *Application) {
	var (
		byName = map[string][]Entity{}
		byKey  = map[string][]Entity{}
	)

	push := func(m map[string][]Entity, id optional.Optional[string], e Entity) {
		if k, ok := id.TryGet(); ok {
			entities := m[k]
			if !slices.Contains(entities, e) {
				m[k] = append(entities, e)
			}
		}
	}

	normalizeKey := func(k string) string {
		if id, err := uuidpb.Parse(k); err == nil {
			return id.AsString()
		}
		return k
	}

	for _, id := range app.identities() {
		// We don't need to check for conflicts with the application's name
		// because it's allowed to be the same as one of the handler's names.
		push(byKey, optional.Transform(id.Key, normalizeKey), app)
	}

	for _, h := range app.HandlerComponents {
		for _, id := range h.identities() {
			push(byKey, optional.Transform(id.Key, normalizeKey), h)
			push(byName, id.Name, h)
		}
	}

	for name, entities := range byName {
		if len(entities) > 1 {
			ctx.Fail(IdentityNameConflictError{name, entities})
		}
	}

	for key, entities := range byKey {
		if len(entities) > 1 {
			ctx.Fail(IdentityKeyConflictError{key, entities})
		}
	}
}

func detectRouteConflicts(ctx *validateContext, app *Application) {
	byKey := map[routeKey][]Handler{}

	for _, h := range app.HandlerComponents {
		for _, r := range h.routes() {
			if k, ok := r.key(); ok {
				switch k.RouteType {
				case HandlesCommandRouteType, RecordsEventRouteType:
					handlers := byKey[k]
					if !slices.Contains(handlers, h) {
						byKey[k] = append(handlers, h)
					}
				}
			}
		}
	}

	for key, handlers := range byKey {
		if len(handlers) > 1 {
			ctx.Fail(RouteConflictError{
				key.RouteType,
				key.MessageTypeName,
				handlers,
			})
		}
	}
}
