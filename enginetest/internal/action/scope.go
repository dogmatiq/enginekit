package action

import (
	"time"

	"github.com/dogmatiq/dogma"
)

// Scope is the interface common to all Dogma scope interfaces.
type Scope interface {
	Log(string, ...any)
}

type executor interface {
	Scope
	ExecuteCommand(dogma.Command)
}

type recorder interface {
	Scope
	RecordEvent(dogma.Event)
}

type scheduler interface {
	Scope
	ScheduleTimeout(dogma.Timeout, time.Time)
}

type destroyer interface {
	Scope
	Destroy()
}

type ender interface {
	Scope
	End()
}
