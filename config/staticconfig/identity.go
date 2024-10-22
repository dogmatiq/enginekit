package staticconfig

import (
	"github.com/dogmatiq/enginekit/config"
	"github.com/dogmatiq/enginekit/config/internal/configbuilder"
	"github.com/dogmatiq/enginekit/config/staticconfig/internal/ssax"
)

func analyzeIdentityCall(
	b configbuilder.EntityBuilder,
	call configurerCall,
) {
	b.Identity(func(b *configbuilder.IdentityBuilder) {
		b.UpdateFidelity(call.Fidelity)

		if name, ok := ssax.AsString(call.Args[0]).TryGet(); ok {
			b.SetName(name)
		} else {
			b.UpdateFidelity(config.Incomplete)
		}

		if key, ok := ssax.AsString(call.Args[1]).TryGet(); ok {
			b.SetKey(key)
		} else {
			b.UpdateFidelity(config.Incomplete)
		}
	})
}
