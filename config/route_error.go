package config

import (
	"fmt"
	"strings"

	"github.com/dogmatiq/enginekit/config/internal/renderer"
	"github.com/dogmatiq/enginekit/internal/typename"
	"github.com/dogmatiq/enginekit/message"
)

// MissingRequiredRouteError indicates that a [Handler] is missing one of its
// required route types.
type MissingRequiredRouteError struct {
	RouteType RouteType
}

func (e MissingRequiredRouteError) Error() string {
	return fmt.Sprintf("no %q routes are configured", e.RouteType)
}

// UnexpectedRouteTypeError indicates that a [Handler] is configured with a [Route]
// with a [RouteType] that is not allowed for that handler type.
type UnexpectedRouteTypeError struct {
	UnexpectedRoute *Route
}

func (e UnexpectedRouteTypeError) Error() string {
	w := &strings.Builder{}

	fmt.Fprintf(w, "unexpected %s route", e.UnexpectedRoute.RouteType())

	if name, ok := e.UnexpectedRoute.AsConfigured.MessageTypeName.TryGet(); ok {
		fmt.Fprintf(w, " for %s", name)
	}

	return w.String()
}

// DuplicateRouteError indicates that a [Handler] is configured with multiple
// routes for the same [MessageType].
type DuplicateRouteError struct {
	RouteType       RouteType
	MessageTypeName string
	DuplicateRoutes []*Route
}

func (e DuplicateRouteError) Error() string {
	return fmt.Sprintf(
		"multiple %q routes are configured for %s",
		e.RouteType,
		e.MessageTypeName,
	)
}

// MessageKindMismatchError indicates that a [Route] refers to a [message.Type]
// that has a different [message.Kind] than the route's [RouteType].
type MessageKindMismatchError struct {
	RouteType   RouteType
	MessageType message.Type
}

func (e MessageKindMismatchError) Error() string {
	return renderer.Inflect(
		"unexpected message kind: %s is a %s, expected a %s",
		typename.Get(e.MessageType.ReflectType()),
		e.MessageType.Kind(),
		e.RouteType.MessageKind(),
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
