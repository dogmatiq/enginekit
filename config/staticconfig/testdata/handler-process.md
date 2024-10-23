# Process message handler

This test ensures that the static analyzer supports all aspects of configuring
a process handler.

```au:output au:group=matrix
valid application github.com/dogmatiq/enginekit/config/staticconfig.App (runtime type unavailable)
  - valid identity app/0726ae0d-67e4-4a50-8a19-9f58eae38e51
  - disabled valid process github.com/dogmatiq/enginekit/config/staticconfig.Process (runtime type unavailable)
      - valid identity process/4ff1b1c1-5c64-49bc-a547-c13f5bafad7d
      - valid handles-event route for github.com/dogmatiq/enginekit/enginetest/stubs.EventStub[github.com/dogmatiq/enginekit/enginetest/stubs.TypeA] (runtime type unavailable)
      - valid executes-command route for github.com/dogmatiq/enginekit/enginetest/stubs.CommandStub[github.com/dogmatiq/enginekit/enginetest/stubs.TypeA] (runtime type unavailable)
      - valid schedules-timeout route for github.com/dogmatiq/enginekit/enginetest/stubs.TimeoutStub[github.com/dogmatiq/enginekit/enginetest/stubs.TypeA] (runtime type unavailable)
```

```go au:input au:group=matrix
package app

import "context"
import "github.com/dogmatiq/dogma"
import "github.com/dogmatiq/enginekit/enginetest/stubs"

type Process struct {}

func (Process) Configure(c dogma.ProcessConfigurer) {
    c.Identity("process", "4ff1b1c1-5c64-49bc-a547-c13f5bafad7d")
    c.Routes(
        dogma.HandlesEvent[stubs.EventStub[stubs.TypeA]](),
        dogma.ExecutesCommand[stubs.CommandStub[stubs.TypeA]](),
        dogma.SchedulesTimeout[stubs.TimeoutStub[stubs.TypeA]](),
    )
    c.Disable()
}

type App struct{}

func (App) Configure(c dogma.ApplicationConfigurer) {
	c.Identity("app", "0726ae0d-67e4-4a50-8a19-9f58eae38e51")
	c.RegisterProcess(Process{})
}

func (Process) New() dogma.ProcessRoot { return nil }
func (Process) RouteEventToInstance(context.Context, dogma.Event) (string, bool, error) { return "", false, nil }
func (Process) HandleEvent(context.Context, dogma.ProcessRoot, dogma.ProcessEventScope, dogma.Event) error { return nil }
func (Process) HandleTimeout(context.Context, dogma.ProcessRoot, dogma.ProcessTimeoutScope, dogma.Timeout) error { return nil }
```
