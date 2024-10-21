package staticconfig

import (
	"go/constant"

	"github.com/dogmatiq/enginekit/config"
	"github.com/dogmatiq/enginekit/config/internal/configbuilder"
)

func analyzeIdentityCall(
	b configbuilder.EntityBuilder,
	call configurerCall,
) {
	b.Identity(func(b *configbuilder.IdentityBuilder) {
		b.UpdateFidelity(call.Fidelity)

		if name := staticValue(call.Args[0]); name != nil {
			b.SetName(constant.StringVal(name))
		} else {
			b.UpdateFidelity(config.Incomplete)
		}

		if key := staticValue(call.Args[1]); key != nil {
			b.SetKey(constant.StringVal(key))
		} else {
			b.UpdateFidelity(config.Incomplete)
		}
	})
}
