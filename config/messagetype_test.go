package config_test

import (
	"reflect"

	. "github.com/dogmatiq/enginekit/config"
	"github.com/dogmatiq/enginekit/fixtures"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("type MessageType", func() {
	Describe("func NewMessageType", func() {
		It("returns the message type if the reflect type implements the dogma.Message interface", func() {
			rt := reflect.TypeOf(fixtures.MessageA1)
			mt, ok := NewMessageType(rt)
			Expect(mt).To(Equal(MessageTypeOf(fixtures.MessageA1)))
			Expect(ok).To(BeTrue())
		})
	})

	Describe("func MessageTypeOf", func() {
		It("returns values that compare as equal for messages of the same type", func() {
			tb := MessageTypeOf(fixtures.MessageA1)
			ta := MessageTypeOf(fixtures.MessageA1)

			Expect(ta).To(Equal(tb))
			Expect(ta == tb).To(BeTrue()) // explicitly check the pointers for standard equality comparability
		})

		It("returns values that do not compare as equal for messages of different types", func() {
			ta := MessageTypeOf(fixtures.MessageA1)
			tb := MessageTypeOf(fixtures.MessageB1)

			Expect(ta).NotTo(Equal(tb))
			Expect(ta != tb).To(BeTrue()) // explicitly check the pointers for standard equality comparability
		})
	})

	Describe("func ReflectType", func() {
		It("returns the reflect.Type for the message", func() {
			mt := MessageTypeOf(fixtures.MessageA1)
			rt := reflect.TypeOf(fixtures.MessageA1)

			Expect(mt.ReflectType()).To(BeIdenticalTo(rt))
		})
	})

	Describe("func String", func() {
		It("returns the package-qualified type name", func() {
			t := MessageTypeOf(fixtures.MessageA1)

			Expect(t.String()).To(Equal(
				"fixtures.MessageA",
			))
		})

		It("returns the package-qualified type name for pointer types", func() {
			t := MessageTypeOf(&fixtures.MessageA1)

			Expect(t.String()).To(Equal(
				"*fixtures.MessageA",
			))
		})

		It("supports anonymous types", func() {
			t := MessageTypeOf(struct{ fixtures.MessageA }{})

			Expect(t.String()).To(Equal("<anonymous>"))
		})
	})
})
