# Multiple handlers of the same type

This test ensures that the static analyzer supports multiple handlers of the
same handler type.

```au:output au:group=matrix
valid application github.com/dogmatiq/enginekit/config/staticconfig.App (value unavailable)
  - valid identity app/c0d4a0fc-2075-4a41-a7bf-7d1870dc0de9
  - valid aggregate github.com/dogmatiq/enginekit/config/staticconfig.One (value unavailable)
      - valid identity one/62e0efa9-c5a0-4b5c-a237-9b51533a6963
      - valid handles-command route for github.com/dogmatiq/enginekit/enginetest/stubs.CommandStub[github.com/dogmatiq/enginekit/enginetest/stubs.TypeA] (type unavailable)
      - valid records-event route for github.com/dogmatiq/enginekit/enginetest/stubs.EventStub[github.com/dogmatiq/enginekit/enginetest/stubs.TypeA] (type unavailable)
  - valid aggregate github.com/dogmatiq/enginekit/config/staticconfig.Two (value unavailable)
      - valid identity two/0c3e2f49-acd0-4d82-800d-5d6d839535de
      - valid handles-command route for github.com/dogmatiq/enginekit/enginetest/stubs.CommandStub[github.com/dogmatiq/enginekit/enginetest/stubs.TypeB] (type unavailable)
      - valid records-event route for github.com/dogmatiq/enginekit/enginetest/stubs.EventStub[github.com/dogmatiq/enginekit/enginetest/stubs.TypeB] (type unavailable)
```

```go au:input au:group=matrix
package app

import "github.com/dogmatiq/dogma"
import "github.com/dogmatiq/enginekit/enginetest/stubs"

type (
    One struct { dogma.AggregateMessageHandler }
    Two struct { dogma.AggregateMessageHandler }
)

func (One) Configure(c dogma.AggregateConfigurer) {
    c.Identity("one", "62e0efa9-c5a0-4b5c-a237-9b51533a6963")
    c.Routes(
        dogma.HandlesCommand[stubs.CommandStub[stubs.TypeA]](),
        dogma.RecordsEvent[stubs.EventStub[stubs.TypeA]](),
    )
}


func (Two) Configure(c dogma.AggregateConfigurer) {
    c.Identity("two", "0c3e2f49-acd0-4d82-800d-5d6d839535de")
    c.Routes(
        dogma.HandlesCommand[stubs.CommandStub[stubs.TypeB]](),
        dogma.RecordsEvent[stubs.EventStub[stubs.TypeB]](),
    )
}

type App struct{}

func (App) Configure(c dogma.ApplicationConfigurer) {
	c.Identity("app", "c0d4a0fc-2075-4a41-a7bf-7d1870dc0de9")
	c.RegisterAggregate(One{})
    c.RegisterAggregate(Two{})
}
```
