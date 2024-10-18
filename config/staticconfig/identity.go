package staticconfig

import (
	"go/constant"

	"github.com/dogmatiq/enginekit/config"
	"github.com/dogmatiq/enginekit/config/internal/configbuilder"
	"golang.org/x/tools/go/ssa"
)

func analyzeIdentityCall(
	b configbuilder.EntityBuilder,
	call configurerCall,
) {
	b.Identity(func(b *configbuilder.IdentityBuilder) {
		b.UpdateFidelity(call.Fidelity)

		if name, ok := call.Args[0].(*ssa.Const); ok {
			b.SetName(constant.StringVal(name.Value))
		} else {
			b.UpdateFidelity(config.Incomplete)
		}

		if key, ok := call.Args[1].(*ssa.Const); ok {
			b.SetKey(constant.StringVal(key.Value))
		} else {
			b.UpdateFidelity(config.Incomplete)
		}
	})
}
