package marshaling_test

import (
	"reflect"

	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/enginekit/config"
	"github.com/dogmatiq/enginekit/fixtures"
	. "github.com/dogmatiq/enginekit/marshaling"
	"github.com/dogmatiq/enginekit/marshaling/internal/pbfixtures"
	"github.com/dogmatiq/enginekit/marshaling/json"
	"github.com/dogmatiq/enginekit/marshaling/protobuf"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("type Marshaler", func() {
	var marshaler *Marshaler

	BeforeEach(func() {
		var err error
		marshaler, err = NewMarshaler(
			[]reflect.Type{
				reflect.TypeOf(&pbfixtures.Message{}),
				reflect.TypeOf(fixtures.MessageA{}),
				reflect.TypeOf(fixtures.MessageB{}),
			},
			[]Codec{
				&protobuf.NativeCodec{},
				&protobuf.JSONCodec{},
				&json.Codec{},
			},
		)
		Expect(err).ShouldNot(HaveOccurred())
	})

	Describe("func NewMarshaler()", func() {
		It("returns an error if multiple codecs used the same media-type", func() {
			_, err := NewMarshaler(
				[]reflect.Type{
					reflect.TypeOf(fixtures.MessageA{}),
				},
				[]Codec{
					&json.Codec{},
					&json.Codec{},
				},
			)
			Expect(err).To(MatchError(
				"multiple codecs use the 'application/json' media-type",
			))
		})

		It("returns an error if there conflicting portable type names", func() {
			_, err := NewMarshaler(
				[]reflect.Type{
					reflect.TypeOf(fixtures.MessageA{}),
					reflect.TypeOf(&fixtures.MessageA{}),
				},
				[]Codec{
					&json.Codec{},
				},
			)
			Expect(err).To(Or(
				MatchError(
					"the type name 'MessageA' is used by both 'fixtures.MessageA' and '*fixtures.MessageA'",
				),
				MatchError(
					"the type name 'MessageA' is used by both '*fixtures.MessageA' and 'fixtures.MessageA'",
				),
			))
		})

		It("returns an error if there are unsupported types", func() {
			_, err := NewMarshaler(
				[]reflect.Type{
					reflect.TypeOf(&pbfixtures.Message{}),
					reflect.TypeOf(fixtures.MessageA{}),
				},
				[]Codec{
					&protobuf.JSONCodec{},
				},
			)
			Expect(err).To(MatchError(
				"no codecs support the 'fixtures.MessageA' type",
			))
		})
	})

	Describe("func NewMarshaler()", func() {
		var cfg *config.ApplicationConfig

		BeforeEach(func() {
			app := &fixtures.Application{
				ConfigureFunc: func(c dogma.ApplicationConfigurer) {
					c.Identity("<application-name>", "<application-key>")

					c.RegisterIntegration(&fixtures.IntegrationMessageHandler{
						ConfigureFunc: func(c dogma.IntegrationConfigurer) {
							c.Identity("<integration-name>", "<integration-key>")
							c.ConsumesCommandType(fixtures.MessageC{})
							c.ProducesEventType(fixtures.MessageE{})
						},
					})
				},
			}

			var err error
			cfg, err = config.NewApplicationConfig(app)
			Expect(err).ShouldNot(HaveOccurred())
		})

		It("returns the expected marshaler", func() {
			m, err := NewMarshalerForApplication(
				cfg,
				[]Codec{
					&protobuf.NativeCodec{},
					&protobuf.JSONCodec{},
					&protobuf.TextCodec{},
					&json.Codec{},
				},
			)
			Expect(err).ShouldNot(HaveOccurred())

			expected, err := NewMarshaler(
				[]reflect.Type{
					fixtures.MessageCType.ReflectType(),
					fixtures.MessageEType.ReflectType(),
				},
				[]Codec{
					&protobuf.NativeCodec{},
					&protobuf.JSONCodec{},
					&protobuf.TextCodec{},
					&json.Codec{},
				},
			)

			Expect(m).To(Equal(expected))
		})

		It("returns the an error if there is an error constructing the marshaler", func() {
			_, err := NewMarshalerForApplication(
				cfg,
				[]Codec{
					&protobuf.NativeCodec{}, // fixture messages are not protobufs
				},
			)
			Expect(err).Should(HaveOccurred())
		})
	})

	Describe("func MarshalType()", func() {
		It("returns the portable type name", func() {
			n, err := marshaler.MarshalType(
				reflect.TypeOf(fixtures.MessageA{}),
			)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(n).To(Equal("MessageA"))
		})

		It("returns an error if the type is not supported", func() {
			_, err := marshaler.MarshalType(
				reflect.TypeOf(fixtures.MessageC{}),
			)
			Expect(err).To(MatchError(
				"no codecs support the 'fixtures.MessageC' type",
			))
		})
	})

	Describe("func UnmarshalType()", func() {
		It("returns the reflection type", func() {
			t, err := marshaler.UnmarshalType("MessageA")
			Expect(err).ShouldNot(HaveOccurred())
			Expect(t).To(Equal(reflect.TypeOf(fixtures.MessageA{})))
		})

		It("returns an error if the type name is not recognized", func() {
			_, err := marshaler.UnmarshalType("MessageC")
			Expect(err).To(MatchError(
				"the portable type name 'MessageC' is not recognized",
			))
		})
	})

	Describe("func UnmarshalTypeFromMediaType()", func() {
		It("returns the reflection type", func() {
			t, err := marshaler.UnmarshalTypeFromMediaType(
				"application/vnd.google.protobuf; type=MessageA",
			)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(t).To(Equal(reflect.TypeOf(fixtures.MessageA{})))
		})

		It("returns an error if the type name is not recognized", func() {
			_, err := marshaler.UnmarshalTypeFromMediaType(
				"application/vnd.google.protobuf; type=MessageC",
			)
			Expect(err).To(MatchError(
				"the portable type name 'MessageC' is not recognized",
			))
		})
	})

	Describe("func Marshal()", func() {
		It("marshals using the first suitable codec", func() {
			p, err := marshaler.Marshal(fixtures.MessageA{})
			Expect(err).ShouldNot(HaveOccurred())
			Expect(p.MediaType).To(Equal("application/json; type=MessageA"))

			p, err = marshaler.Marshal(&pbfixtures.Message{})
			Expect(err).ShouldNot(HaveOccurred())
			Expect(p.MediaType).To(Equal("application/vnd.google.protobuf; type=dogmatiq.enginekit.marshaling.pbfixtures.Message"))
		})

		It("returns an error if the codec fails", func() {
			_, err := marshaler.Marshal(
				&pbfixtures.Message{
					Value: string([]byte{0xfe}),
				},
			)
			Expect(err).Should(HaveOccurred())
		})

		It("returns an error if the type is not supported", func() {
			_, err := marshaler.Marshal(fixtures.MessageC{})
			Expect(err).To(MatchError(
				"no codecs support the 'fixtures.MessageC' type",
			))
		})
	})

	Describe("func Unmarshal()", func() {
		It("marshals using the first suitable codec", func() {
			v, err := marshaler.Unmarshal(
				Packet{
					"application/json; type=MessageA",
					[]byte("{}"),
				},
			)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(v).To(Equal(fixtures.MessageA{}))

			v, err = marshaler.Unmarshal(
				Packet{
					"application/vnd.google.protobuf+json; type=dogmatiq.enginekit.marshaling.pbfixtures.Message",
					[]byte("{}"),
				},
			)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(v).To(Equal(&pbfixtures.Message{}))
		})

		It("returns an error if the media-type is not supported", func() {
			_, err := marshaler.Unmarshal(Packet{"text/plain", nil})
			Expect(err).To(MatchError(
				"no codecs support the 'text/plain' media-type",
			))
		})

		It("returns an error if the media-type is malformed", func() {
			_, err := marshaler.Unmarshal(Packet{"<malformed>", nil})
			Expect(err).Should(HaveOccurred())
		})

		It("returns an error if the media-type does not specify a type parameter", func() {
			_, err := marshaler.Unmarshal(Packet{"application/json", nil})
			Expect(err).Should(MatchError(
				"the media-type 'application/json' does not specify a 'type' parameter",
			))
		})

		It("returns an error if the type is not supported", func() {
			_, err := marshaler.Unmarshal(Packet{"application/json; type=MessageC", nil})
			Expect(err).Should(MatchError(
				"the portable type name 'MessageC' is not recognized",
			))
		})

		It("returns an error if the codec fails", func() {
			_, err := marshaler.Unmarshal(
				Packet{
					"application/json; type=MessageA",
					[]byte("{"),
				},
			)
			Expect(err).Should(HaveOccurred())
		})
	})
})
