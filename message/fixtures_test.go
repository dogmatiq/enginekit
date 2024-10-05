package message_test

import "github.com/dogmatiq/dogma"

type (
	// partialMessage implements [dogma.Command], but does not implement any of the
	// more specific interfaces.
	partialMessage struct{}

	// ptrCommand implements [dogma.Command] with a pointer receiver.
	ptrCommand struct{}

	// nonPtrCommand implements [dogma.Command] with a non-pointer
	// receiver.
	nonPtrCommand struct{}
)

func (partialMessage) MessageDescription() string                 { panic("not implemented") }
func (*ptrCommand) MessageDescription() string                    { panic("not implemented") }
func (*ptrCommand) Validate(dogma.CommandValidationScope) error   { panic("not implemented") }
func (nonPtrCommand) MessageDescription() string                  { panic("not implemented") }
func (nonPtrCommand) Validate(dogma.CommandValidationScope) error { panic("not implemented") }
