# Interface as an entity configurer.

This test ensures that the static analyzer can recognize the type of a handler
when it is used in instantiating a generic handler.

```go au:input
package app

import (
	"context"
	. "github.com/dogmatiq/dogma"
	. "github.com/dogmatiq/enginekit/enginetest/stubs"
)

type GenericIntegration[T any, H IntegrationMessageHandler] struct {
	Handler H
}

func (i *GenericIntegration[T, H]) Configure(c IntegrationConfigurer) {
	i.Handler.Configure(c)
}

func (i *GenericIntegration[T, H]) HandleCommand(
	ctx context.Context,
	s IntegrationCommandScope,
	cmd Command,
) error {
	return i.Handler.HandleCommand(ctx, s, cmd)
}

type integrationHandler struct {}

func (integrationHandler) Configure(c IntegrationConfigurer) {
	c.Identity("<integration>", "abc7c329-c9da-4161-a8e2-6ab45be2dd83")

	routes := []IntegrationRoute{
		HandlesCommand[CommandStub[TypeA]](),
	}

	c.Routes(routes...)
}

func (integrationHandler) HandleCommand(
	_ context.Context,
 	_ IntegrationCommandScope,
  	_ Command,
) error {
	return nil
}

type InstantiatedIntegration = GenericIntegration[struct{}, integrationHandler]

type App struct {
	Integration InstantiatedIntegration
}

func (a App) Configure(c ApplicationConfigurer) {
	c.Identity("<app>", "e522c782-48d2-4c47-a4c9-81e0d7cdeba0")
	c.RegisterIntegration(&a.Integration)
}

```

```au:output
application <app> (e522c782-48d2-4c47-a4c9-81e0d7cdeba0) App

    - integration <integration> (abc7c329-c9da-4161-a8e2-6ab45be2dd83) *InstantiatedIntegration
        handles CommandStub[TypeA]?
```
