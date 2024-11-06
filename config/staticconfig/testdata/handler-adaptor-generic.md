# Generic handler adaptor

This test ensures that the static analyzer can analyze configuration logic that
is implemented within a type that is passed to a handler as a type parameter.

```au:output au:group=matrix
valid application github.com/dogmatiq/enginekit/config/staticconfig.App (value unavailable)
  - valid identity app/24c4a011-3d9e-493a-95c6-ef9ab059f65f
  - valid integration github.com/dogmatiq/enginekit/config/staticconfig.Adaptor[github.com/dogmatiq/enginekit/config/staticconfig.impl] (value unavailable)
      - valid identity integration/a57834ad-251a-4672-9b82-f2a538a64655
      - valid handles-command route for github.com/dogmatiq/enginekit/enginetest/stubs.CommandStub[github.com/dogmatiq/enginekit/enginetest/stubs.TypeA] (type unavailable)
```

```go au:input au:group=matrix
package app

import "github.com/dogmatiq/dogma"
import "github.com/dogmatiq/enginekit/enginetest/stubs"

type Adaptor[T interface { Configure(dogma.IntegrationConfigurer) }] struct {
	dogma.IntegrationMessageHandler
	Impl T
}

func (a Adaptor[T]) Configure(c dogma.IntegrationConfigurer) {
	a.Impl.Configure(c)
}

type impl struct {}

func (impl) Configure(c dogma.IntegrationConfigurer) {
	c.Identity("integration", "a57834ad-251a-4672-9b82-f2a538a64655")
	c.Routes(dogma.HandlesCommand[stubs.CommandStub[stubs.TypeA]]())
}

type App struct {}

func (App) Configure(c dogma.ApplicationConfigurer) {
	c.Identity("app", "24c4a011-3d9e-493a-95c6-ef9ab059f65f")
	c.RegisterIntegration(Adaptor[impl]{})
}
```
