# No applications

This test ensures that static analyzer does not fail when the code does not
contain any Dogma applications.

```go au:input
package app

// App looks a lot like a [dogma.Application], but does not actually
// implement the Dogma interface.
type App struct{}

func (a App) Configure(c ApplicationConfigurer) {
    c.Identity("name", "ee6ca834-34a3-4e59-8c36-7aeb796401d7")
}

type ApplicationConfigurer interface {
    Identity(name, key string)
}
```

```au:output
(no applications found)
```
