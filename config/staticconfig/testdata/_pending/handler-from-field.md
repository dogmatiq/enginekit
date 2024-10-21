# Handler from struct field

This test ensures that the static analyzer can recognized the type of a handler
when it is registered using the value of a struct field, rather than constructed
inline.

```go au:input au:group=matrix
package app

import (
    "context"
    . "github.com/dogmatiq/dogma"
    . "github.com/dogmatiq/enginekit/enginetest/stubs"
)

type App struct {
	Field Handler
}

func (a App) Configure(c ApplicationConfigurer) {
	c.Identity("<app>", "7468a57f-20f0-4d11-9aad-48fcd553a908")
	c.RegisterIntegration(a.Field)
}

type Handler struct{}

func (Handler) Configure(c IntegrationConfigurer) {
	c.Identity("<integration>", "195ede4a-3f26-4d19-a8fe-41b2a5f92d06")
	c.Routes(HandlesCommand[CommandStub[TypeA]]())
}

func (Handler) HandleCommand(
    context.Context,
    IntegrationCommandScope,
    Command,
) error{
    return nil
}
```

```au:output au:group=matrix
application <app> (7468a57f-20f0-4d11-9aad-48fcd553a908) App

    - integration <integration> (195ede4a-3f26-4d19-a8fe-41b2a5f92d06) Handler
        handles CommandStub[TypeA]?
```
