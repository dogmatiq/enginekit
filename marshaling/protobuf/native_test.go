package protobuf_test

import (
	"github.com/dogmatiq/enginekit/fixtures"
	"github.com/dogmatiq/enginekit/marshaling/internal/pbfixtures"
	. "github.com/dogmatiq/enginekit/marshaling/protobuf"
	"github.com/golang/protobuf/proto"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("type NativeCodec", func() {
	var codec *NativeCodec

	BeforeEach(func() {
		codec = &NativeCodec{}
	})

	Describe("func MediaType()", func() {
		It("returns the expected media-type", func() {
			Expect(codec.MediaType()).To(Equal("application/vnd.google.protobuf"))
		})
	})

	Describe("func Marshal()", func() {
		It("marshals the value", func() {
			m := &pbfixtures.Message{
				Value: "<value>",
			}

			data, err := codec.Marshal(m)
			Expect(err).ShouldNot(HaveOccurred())

			expected, err := proto.Marshal(m)
			Expect(err).ShouldNot(HaveOccurred())

			Expect(err).ShouldNot(HaveOccurred())
			Expect(data).To(Equal(expected))
		})

		It("returns an error if the type is not a protocol buffers message", func() {
			_, err := codec.Marshal(
				fixtures.MessageA{},
			)
			Expect(err).To(MatchError(
				"'fixtures.MessageA' is not a protocol buffers message",
			))
		})
	})

	Describe("func Unmarshal()", func() {
		It("unmarshals the data", func() {
			m := &pbfixtures.Message{
				Value: "<value>",
			}

			data, err := codec.Marshal(m)
			Expect(err).ShouldNot(HaveOccurred())

			m = &pbfixtures.Message{}
			err = codec.Unmarshal(data, m)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(m).To(Equal(
				&pbfixtures.Message{
					Value: "<value>",
				},
			))
		})

		It("returns an error if the type is not a protocol buffers message", func() {
			m := fixtures.MessageA{}
			err := codec.Unmarshal(nil, m)
			Expect(err).To(MatchError(
				"'fixtures.MessageA' is not a protocol buffers message",
			))
		})
	})
})
