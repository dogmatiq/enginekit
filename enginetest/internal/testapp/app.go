package testapp

import (
	"github.com/dogmatiq/dogma"
)

// App is a Dogma application that is used to test the engine.
type App struct {
	Events EventProjection
}

// Configure configures the Dogma application.
func (a *App) Configure(c dogma.ApplicationConfigurer) {
	c.Identity("enginetest", "861916bb-e09b-4027-90d2-139722be331a")

	c.RegisterProjection(&a.Events)
	c.RegisterIntegration(&actionExecutor{})

	c.RegisterIntegration(&integrationA{})
	c.RegisterIntegration(&integrationB{})

	c.RegisterProcess(&processA{})
}
