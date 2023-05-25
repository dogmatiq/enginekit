package action

import (
	"github.com/dogmatiq/dogma"
)

// Actionable is a dogma message that provides a set of actions to execute.
type Actionable interface {
	dogma.Message
	GetActions() []*Action
}

// Run runs the actions defined by m against the given scope.
func Run(
	s Scope,
	m Actionable,
) error {
	type behavior interface {
		do(Scope) error
	}

	for _, act := range m.GetActions() {
		if err := act.GetBehavior().(behavior).do(s); err != nil {
			return err
		}
	}

	return nil
}
