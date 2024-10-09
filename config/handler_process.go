package config

import (
	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/enginekit/optional"
)

// Process represents the (potentially invalid) configuration of a
// [dogma.ProcessMessageHandler] implementation.
type Process struct {
	// Impl contains information about the type that produced the configuration,
	// if available.
	Impl optional.Optional[Implementation[dogma.ProcessMessageHandler]]

	// ConfiguredIdentities is the list of (potentially invalid or duplicated)
	// identities configured for the handler.
	ConfiguredIdentities []Identity

	// ConfiguredRoutes is the list of (potentially invalid, incomplete or
	// duplicated) message routes configured for the handler.
	ConfiguredRoutes []Route

	// IsDisabled is true if the handler was disabled via the configurer.
	IsDisabled bool

	// IsExhaustive is true if the complete configuration was loaded. It may be
	// false, for example, when attempting to load configuration using static
	// analysis, but the code depends on runtime type information.
	IsExhaustive bool
}

func (h Process) String() string {
	return stringify("process", h.Impl, h.ConfiguredIdentities)
}

// Identity returns the entity's identity.
//
// It panics if no single valid identity is configured.
func (h Process) Identity() Identity {
	return normalizedIdentity(h)
}

// Routes returns the routes configured for the handler.
//
// It panics if the routes are incomplete or invalid.
func (h Process) Routes() []Route {
	return normalizedRoutes(h)
}

func (h Process) configuredIdentities() []Identity { return h.ConfiguredIdentities }
func (h Process) configuredRoutes() []Route        { return h.ConfiguredRoutes }

func (h Process) normalize(opts validationOptions) (_ Entity, errs error) {
	normalizeIdentitiesInPlace(opts, h, &errs, &h.ConfiguredIdentities)

	normalizeRoutesInPlace(
		h,
		&errs,
		&h.ConfiguredRoutes,
		map[RouteType]bool{
			HandlesEventRoute:     true,
			ExecutesCommandRoute:  true,
			SchedulesTimeoutRoute: false,
		},
	)

	return h, errs
}
