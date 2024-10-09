package config

import "fmt"

// InvalidHandlerError is an error that occurs when an [Application] contains an
// [Handler] with an invalid configuration.
type InvalidHandlerError struct {
	Application Application
	Handler     Handler
	Cause       error
}

func (e InvalidHandlerError) Error() string {
	return fmt.Sprintf("%s contains an invalid handler: %s", e.Application, e.Cause)
}

func (e InvalidHandlerError) Unwrap() error {
	return e.Cause
}

// IdentityConflictError indicates that more than one [Entity] within the same
// [Application] shares the same [Identity].
type IdentityConflictError struct {
	Entities            []Entity
	ConflictingIdentity Identity
}

func (e IdentityConflictError) Error() string {
	return fmt.Sprintf(
		"%s have the same identity (%s), which is not allowed",
		renderList(e.Entities),
		e.ConflictingIdentity,
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
		"%s have the same identity name (%s), which is not allowed",
		renderList(e.Entities),
		e.ConflictingName,
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
		"%s have the same identity key (%s), which is not allowed",
		renderList(e.Entities),
		e.ConflictingKey,
	)
}

// RouteConflictError indicates that more than one [Handler] within the same
// [Application] is configured with conflicting routes for the same
// [MessageType].
type RouteConflictError struct {
	Handlers         []Handler
	ConflictingRoute Route
}

func (e RouteConflictError) Error() string {
	return fmt.Sprintf(
		"%s have %q routes for the same %s type (%s), which is not allowed",
		renderList(e.Handlers),
		e.ConflictingRoute.RouteType.Get(),
		e.ConflictingRoute.MessageType.Get().Kind,
		e.ConflictingRoute.MessageType.Get().TypeName,
	)
}
