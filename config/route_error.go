package config

import (
	"fmt"

	"github.com/dogmatiq/enginekit/internal/typename"
	"github.com/dogmatiq/enginekit/message"
)

// MissingRequiredRouteError indicates that a [Handler] is missing one of its
// required route types.
type MissingRequiredRouteError struct {
	RouteType RouteType
}

func (e MissingRequiredRouteError) Error() string {
	return fmt.Sprintf("expected at least one %q route", e.RouteType)
}

// UnexpectedRouteError indicates that a [Handler] is configured with a [Route]
// with a [RouteType] that is not allowed for that handler type.
type UnexpectedRouteError struct {
	UnexpectedRoute Route
}

func (e UnexpectedRouteError) Error() string {
	return fmt.Sprintf("unexpected route: %s", e.UnexpectedRoute)
}

// DuplicateRouteError indicates that a [Handler] is configured with multiple
// routes for the same [MessageType].
type DuplicateRouteError struct {
	RouteType       RouteType
	MessageTypeName string
	DuplicateRoutes []Route
}

func (e DuplicateRouteError) Error() string {
	return fmt.Sprintf(
		"multiple %q routes are configured for %s",
		e.RouteType,
		e.MessageTypeName,
	)
}

// MissingRouteTypeError indicates that a [Route] is missing its [RouteType].
type MissingRouteTypeError struct{}

func (e MissingRouteTypeError) Error() string {
	return "missing route type"
}

// MessageKindMismatchError indicates that a [Route] refers to a [message.Type]
// that has a different [message.Kind] than the route's [RouteType].
type MessageKindMismatchError struct {
	RouteType   RouteType
	MessageType message.Type
}

func (e MessageKindMismatchError) Error() string {
	return fmt.Sprintf(
		"message kind mismatch: %s expects %q, but %s is %q",
		e.RouteType,
		e.RouteType.MessageKind(),
		typename.Get(e.MessageType.ReflectType()),
		e.MessageType.Kind(),
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
