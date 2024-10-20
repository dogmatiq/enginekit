# Applications with no handlers

This test ensures that the static analyzer includes Dogma applications that have
no handlers.

## With non-pointer receiver

```au:output au:group=matrix au:group="non-pointer"
valid application github.com/dogmatiq/enginekit/config/staticconfig/testdata.App (runtime type unavailable)
  - valid identity app/8a6baab1-ee64-402e-a081-e43f4bebc243
```

```go au:input au:group=matrix au:group="non-pointer"
package app

import "github.com/dogmatiq/dogma"

type App struct{}

func (App) Configure(c dogma.ApplicationConfigurer) {
	c.Identity("app", "8a6baab1-ee64-402e-a081-e43f4bebc243")
}
```

## With pointer receiver

```au:output au:group=matrix au:group="pointer"
valid application *github.com/dogmatiq/enginekit/config/staticconfig/testdata/empty-app.App (runtime type unavailable)
  - valid identity app/d196eb7a-bad4-4826-8763-db1111882fbd
```

```go au:input au:group=matrix au:group="pointer"
package app

import "github.com/dogmatiq/dogma"

type App struct{}

func (*App) Configure(c dogma.ApplicationConfigurer) {
	c.Identity("app", "d196eb7a-bad4-4826-8763-db1111882fbd")
}
```
