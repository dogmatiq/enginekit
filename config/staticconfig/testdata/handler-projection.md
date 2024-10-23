# Projection message handler

This test ensures that the static analyzer supports all aspects of configuring
a projection handler.

```au:output au:group=matrix
valid application github.com/dogmatiq/enginekit/config/staticconfig.App (runtime type unavailable)
  - valid identity app/0726ae0d-67e4-4a50-8a19-9f58eae38e51
  - disabled valid projection github.com/dogmatiq/enginekit/config/staticconfig.Projection (runtime type unavailable)
      - valid identity projection/238d7498-515b-44b5-b6a8-914a08762ecc
      - valid handles-event route for github.com/dogmatiq/enginekit/enginetest/stubs.EventStub[github.com/dogmatiq/enginekit/enginetest/stubs.TypeA] (runtime type unavailable)
```

```go au:input au:group=matrix
package app

import "context"
import "github.com/dogmatiq/dogma"
import "github.com/dogmatiq/enginekit/enginetest/stubs"

type Projection struct {}

func (Projection) Configure(c dogma.ProjectionConfigurer) {
    c.Identity("projection", "238d7498-515b-44b5-b6a8-914a08762ecc")
    c.Routes(
        dogma.HandlesEvent[stubs.EventStub[stubs.TypeA]](),
    )
    // c.DeliveryPolicy(dogma.BroadcastProjectionDeliveryPolicy{
    //     PrimaryFirst: true,
    // })
    c.Disable()
}

type App struct{}

func (App) Configure(c dogma.ApplicationConfigurer) {
	c.Identity("app", "0726ae0d-67e4-4a50-8a19-9f58eae38e51")
	c.RegisterProjection(Projection{})
}

func (Projection) HandleEvent(context.Context, []byte, []byte, []byte, dogma.ProjectionEventScope, dogma.Event) (bool, error) { return false, nil }
func (Projection) Compact(context.Context, dogma.ProjectionCompactScope) error { return nil }
func (Projection) ResourceVersion(context.Context, []byte) ([]byte, error)
func (Projection) CloseResource(context.Context, []byte) error
```
