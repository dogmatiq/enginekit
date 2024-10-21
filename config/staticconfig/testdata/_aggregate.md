# Aggregate

This test ensures that the static analyzer supports all aspects of configuring
an aggregate.

```au:output au:group=matrix
valid application github.com/dogmatiq/enginekit/config/staticconfig.App (runtime type unavailable)
  - valid identity app/0726ae0d-67e4-4a50-8a19-9f58eae38e51
  - disabled valid aggregate github.com/dogmatiq/enginekit/config/staticconfig.Aggregate (runtime type unavailable)
      - valid identity aggregate/916e5e95-70c4-4823-9de2-0f7389d18b4f
      - incomplete route
      - incomplete route
```

```go au:input au:group=matrix
package app

import "github.com/dogmatiq/dogma"
import "github.com/dogmatiq/enginekit/enginetest/stubs"

type Aggregate struct {}

func (Aggregate) Configure(c dogma.AggregateConfigurer) {
    c.Identity("aggregate", "916e5e95-70c4-4823-9de2-0f7389d18b4f")
    c.Routes(
        dogma.HandlesCommand[stubs.CommandStub[stubs.TypeA]](),
        dogma.RecordsEvent[stubs.EventStub[stubs.TypeA]](),
    )
    c.Disable()
}

type App struct{}

func (App) Configure(c dogma.ApplicationConfigurer) {
	c.Identity("app", "0726ae0d-67e4-4a50-8a19-9f58eae38e51")
	c.RegisterAggregate(Aggregate{})
}

func (Aggregate) New() dogma.AggregateRoot { return nil }
func (Aggregate) RouteCommandToInstance(dogma.Command) string { return "" }
func (Aggregate) HandleCommand(dogma.AggregateRoot, dogma.AggregateCommandScope, dogma.Command) {}
```
