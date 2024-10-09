package config

import "fmt"

// NoIdentityError indicates that an [Entity] has been configured without an
// [Identity].
type NoIdentityError struct {
	Entity Entity
}

func (e NoIdentityError) Error() string {
	return fmt.Sprintf(
		"%s is configured without an identity, Identity() must be called exactly once within Configure()",
		e.Entity,
	)
}

// MultipleIdentitiesError indicates that an [Entity] has been configured with
// more than one [Identity].
type MultipleIdentitiesError struct {
	Entity     Entity
	Identities []Identity
}

func (e MultipleIdentitiesError) Error() string {
	return fmt.Sprintf(
		"%s is configured with multiple identities (%s), Identity() must be called exactly once within Configure()",
		e.Entity,
		renderList(e.Identities),
	)
}

// InvalidIdentityError indicates that an [Entity] has been configured with an
// invalid [Identity].
type InvalidIdentityError struct {
	Entity          Entity
	InvalidIdentity Identity
	Cause           error
}

func (e InvalidIdentityError) Error() string {
	return fmt.Sprintf(
		"%s is configured with an invalid identity (%s): %s",
		e.Entity,
		e.InvalidIdentity,
		e.Cause,
	)
}

func (e InvalidIdentityError) Unwrap() error {
	return e.Cause
}

// MissingRouteError indicates that a [Handler] is missing one of its mandatory
// route types.
type MissingRouteError struct {
	Handler   Handler
	RouteType RouteType
}

func (e MissingRouteError) Error() string {
	return fmt.Sprintf(
		"%s must have at least one %q route",
		e.Handler,
		e.RouteType,
	)
}

// UnexpectedRouteTypeError indicates that a [Handler] is configured with a
// [Route] with a [RouteType] that is not allowed for that handler type.
type UnexpectedRouteTypeError struct {
	Handler         Handler
	UnexpectedRoute Route
}

func (e UnexpectedRouteTypeError) Error() string {
	t := e.UnexpectedRoute.RouteType.Get()
	article := "a"
	if t == ExecutesCommandRoute {
		article = "an"
	}

	return fmt.Sprintf(
		"%s is configured with %s %q route, which is not allowed for that handler type",
		e.Handler,
		article,
		t,
	)
}
