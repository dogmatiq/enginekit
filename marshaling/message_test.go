package marshaling_test

import (
	"reflect"

	"github.com/dogmatiq/enginekit/fixtures"
	. "github.com/dogmatiq/enginekit/marshaling"
	"github.com/dogmatiq/enginekit/marshaling/json"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Context("messages", func() {
	var marshaler *Marshaler

	BeforeEach(func() {
		var err error
		marshaler, err = NewMarshaler(
			[]reflect.Type{
				reflect.TypeOf(fixtures.MessageA{}),
			},
			[]Codec{
				&json.Codec{},
			},
		)
		Expect(err).ShouldNot(HaveOccurred())
	})

	Describe("func MarshalMessageType()", func() {
		It("marshals the type name using the marshaler", func() {
			n, err := MarshalMessageType(
				marshaler,
				fixtures.MessageAType,
			)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(n).To(Equal("MessageA"))
		})

		It("returns an error if the type is not registered", func() {
			_, err := MarshalMessageType(
				marshaler,
				fixtures.MessageCType,
			)
			Expect(err).To(MatchError(
				"no codecs support the 'fixtures.MessageC' type",
			))
		})
	})

	Describe("func MustMarshalMessageType()", func() {
		It("marshals the type name using the marshaler", func() {
			n := MustMarshalMessageType(
				marshaler,
				fixtures.MessageAType,
			)
			Expect(n).To(Equal("MessageA"))
		})

		It("panics if marshaling fails", func() {
			Expect(func() {
				MustMarshalMessageType(
					marshaler,
					fixtures.MessageCType,
				)

			}).To(Panic())
		})
	})

	Describe("func UnmarshalMessageType()", func() {
		It("returns the message type", func() {
			t, err := UnmarshalMessageType(
				marshaler,
				"MessageA",
			)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(t).To(Equal(fixtures.MessageAType))
		})

		It("returns an error if the type is not registered", func() {
			_, err := UnmarshalMessageType(
				marshaler,
				"MessageC",
			)
			Expect(err).To(MatchError(
				"the portable type name 'MessageC' is not recognized",
			))
		})
	})

	Describe("func MustUnmarshalMessageType()", func() {
		It("returns the message type", func() {
			t := MustUnmarshalMessageType(
				marshaler,
				"MessageA",
			)
			Expect(t).To(Equal(fixtures.MessageAType))
		})

		It("panics if the type is not registered", func() {
			Expect(func() {
				MustUnmarshalMessageType(
					marshaler,
					"MessageC",
				)
			}).To(Panic())
		})
	})

	Describe("func UnmarshalMessageTypeFromMediaType()", func() {
		It("returns the message type", func() {
			t, err := UnmarshalMessageTypeFromMediaType(
				marshaler,
				"application/json; type=MessageA",
			)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(t).To(Equal(fixtures.MessageAType))
		})

		It("returns an error if the type is not registered", func() {
			_, err := UnmarshalMessageTypeFromMediaType(
				marshaler,
				"application/json; type=MessageC",
			)
			Expect(err).To(MatchError(
				"the portable type name 'MessageC' is not recognized",
			))
		})
	})

	Describe("func MustUnmarshalMessageTypeFromMediaType()", func() {
		It("returns the message type", func() {
			t := MustUnmarshalMessageTypeFromMediaType(
				marshaler,
				"application/json; type=MessageA",
			)
			Expect(t).To(Equal(fixtures.MessageAType))
		})

		It("returns an error if the type is not registered", func() {
			Expect(func() {
				MustUnmarshalMessageTypeFromMediaType(
					marshaler,
					"application/json; type=MessageC",
				)
			}).To(Panic())
		})
	})

	Describe("func MarshalMessage()", func() {
		It("marshals the message using the marshaler", func() {
			p, err := MarshalMessage(
				marshaler,
				fixtures.MessageA{
					Value: "<value>",
				},
			)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(p.MediaType).To(Equal("application/json; type=MessageA"))
			Expect(p.Data).To(Equal([]byte(`{"Value":"\u003cvalue\u003e"}`)))
		})
	})

	Describe("func MustMarshalMessage()", func() {
		It("marshals the message using the marshaler", func() {
			p := MustMarshalMessage(
				marshaler,
				fixtures.MessageA{
					Value: "<value>",
				},
			)
			Expect(p.MediaType).To(Equal("application/json; type=MessageA"))
			Expect(p.Data).To(Equal([]byte(`{"Value":"\u003cvalue\u003e"}`)))
		})

		It("panics if marshaling fails", func() {
			Expect(func() {
				MustMarshalMessage(
					marshaler,
					fixtures.MessageC{},
				)

			}).To(Panic())
		})
	})

	Describe("func UnmarshalMessage()", func() {
		It("unmarshals the message using the marshaler", func() {
			m, err := UnmarshalMessage(
				marshaler,
				Packet{
					"application/json; type=MessageA",
					[]byte(`{"Value":"\u003cvalue\u003e"}`),
				},
			)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(m).To(Equal(
				fixtures.MessageA{
					Value: "<value>",
				},
			))
		})
	})
})
