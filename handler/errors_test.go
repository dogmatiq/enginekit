package handler_test

import (
	"github.com/dogmatiq/enginekit/fixtures"
	. "github.com/dogmatiq/enginekit/handler"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("type EmptyInstanceIDError", func() {
	Describe("func Error", func() {
		It("returns a meaningful error description", func() {
			err := EmptyInstanceIDError{
				HandlerName: "<name>",
				HandlerKey:  "<key>",
				HandlerType: AggregateType,
				Message:     fixtures.MessageA1,
			}

			Expect(err.Error()).To(Equal(
				"the '<name>' aggregate message handler attempted to route a fixtures.MessageA message to an empty instance ID",
			))
		})
	})
})

var _ = Describe("type NilRootError", func() {
	Describe("func Error", func() {
		It("returns a meaningful error description", func() {
			err := NilRootError{
				HandlerName: "<name>",
				HandlerKey:  "<key>",
				HandlerType: AggregateType,
			}

			Expect(err.Error()).To(Equal(
				"the '<name>' aggregate message handler returned a nil root from New()",
			))
		})
	})
})

var _ = Describe("type EventNotRecordedError", func() {
	Describe("func Error", func() {
		When("the instance was created", func() {
			It("returns a meaningful error description", func() {
				err := EventNotRecordedError{
					HandlerName:  "<name>",
					HandlerKey:   "<key>",
					WasDestroyed: false,
					Message:      fixtures.MessageA1,
					InstanceID:   "<instance>",
				}

				Expect(err.Error()).To(Equal(
					"the '<name>' aggregate message handler created the '<instance>' instance without recording an event while handling a fixtures.MessageA command",
				))
			})
		})

		When("the instance was destroyed", func() {
			It("returns a meaningful error description", func() {
				err := EventNotRecordedError{
					HandlerName:  "<name>",
					HandlerKey:   "<key>",
					WasDestroyed: true,
					Message:      fixtures.MessageA1,
					InstanceID:   "<instance>",
				}

				Expect(err.Error()).To(Equal(
					"the '<name>' aggregate message handler destroyed the '<instance>' instance without recording an event while handling a fixtures.MessageA command",
				))
			})
		})
	})
})

var _ = Describe("type UnexpectedMessageError", func() {
	Describe("func Error", func() {
		It("returns a meaningful error description", func() {
			err := UnexpectedMessageError{
				HandlerName: "<name>",
				HandlerKey:  "<key>",
				HandlerType: AggregateType,
				Message:     fixtures.MessageA1,
				InstanceID:  "<instance>",
			}

			Expect(err.Error()).To(Equal(
				"the '<name>' aggregate message handler does not expect fixtures.MessageA messages",
			))
		})
	})
})
