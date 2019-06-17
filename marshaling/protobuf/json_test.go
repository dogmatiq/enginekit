package protobuf_test

import (
	"github.com/dogmatiq/enginekit/fixtures"
	"github.com/dogmatiq/enginekit/marshaling/internal/pbfixtures"
	. "github.com/dogmatiq/enginekit/marshaling/protobuf"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("type JSONCodec", func() {
	var codec *JSONCodec

	BeforeEach(func() {
		codec = &JSONCodec{}
	})

	Describe("func MediaType()", func() {
		It("returns the expected media-type", func() {
			Expect(codec.MediaType()).To(Equal("application/vnd.google.protobuf+json"))
		})
	})

	Describe("func Marshal()", func() {
		It("marshals the value", func() {
			data, err := codec.Marshal(
				&pbfixtures.Message{
					Value: "<value>",
				},
			)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(string(data)).To(Equal(`{"value":"\u003cvalue\u003e"}`))
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
			data := []byte(`{"value":"\u003cvalue\u003e"}`)

			m := &pbfixtures.Message{}
			err := codec.Unmarshal(data, m)
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
