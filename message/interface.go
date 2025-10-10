package message

import (
	"fmt"
	"reflect"

	"github.com/dogmatiq/dogma"
)

var (
	commandInterface = reflect.TypeFor[dogma.Command]()
	eventInterface   = reflect.TypeFor[dogma.Event]()
	timeoutInterface = reflect.TypeFor[dogma.Timeout]()
)

// tryInterfaceOf returns the message interface that r implements, or nil if it
// does not implement any of the message interfaces.
func tryInterfaceOf(r reflect.Type) reflect.Type {
	if r.Implements(commandInterface) {
		return commandInterface
	}

	if r.Implements(eventInterface) {
		return eventInterface
	}

	if r.Implements(timeoutInterface) {
		return timeoutInterface
	}

	return nil
}

// interfaceOf returns the message interface that r implements, or nil if it
// does not implement any of the message interfaces.
func interfaceOf(r reflect.Type) reflect.Type {
	if r == nil {
		panic("message type must not be nil")
	}

	i := tryInterfaceOf(r)
	if i == nil {
		panic(fmt.Sprintf("%s does not implement dogma.Command, dogma.Event or dogma.Timeout", r))
	}

	if r.Kind() != reflect.Pointer {
		p := reflect.PointerTo(r)
		panic(fmt.Sprintf("%s implements %s, but message implementations must use pointer receivers, use %s instead", r, i, p))
	}

	return i
}

func guardAgainstNonMessage(r reflect.Type) {
	interfaceOf(r)
}
