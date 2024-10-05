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

	if i := tryInterfaceOf(r); i != nil {
		if r.Kind() == reflect.Pointer {
			e := r.Elem()
			if i := tryInterfaceOf(e); i != nil {
				panic(fmt.Sprintf("%s does not use a pointer receiver to implement %s, use %s instead", r, i, e))
			}
		}

		return i
	}

	p := reflect.PointerTo(r)
	if i := tryInterfaceOf(p); i != nil {
		panic(fmt.Sprintf("%s uses a pointer receiver to implement %s, use %s instead", r, i, p))
	}

	panic(fmt.Sprintf("%s does not implement dogma.Command, dogma.Event or dogma.Timeout", r))
}

func guardAgainstNonMessage(r reflect.Type) {
	interfaceOf(r)
}
