package configbuilder

import (
	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/enginekit/config"
)

// Route builds [config.Route] values.
type Route struct {
}

func (b *Route) Set(r dogma.Route) {
}

func (b *Route) Get() config.Route {
	return config.Route{}
}
