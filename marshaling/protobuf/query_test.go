package protobuf

import (
	"reflect"

	"github.com/dogmatiq/enginekit/fixtures"
	"github.com/dogmatiq/enginekit/marshaling/internal/pbfixtures"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("type TextCodec", func() {
	var codec codec

	Describe("func Query()", func() {
		It("uses the protocol name as the portable type", func() {
			rt := reflect.TypeOf(&pbfixtures.Message{})

			caps := codec.Query(
				[]reflect.Type{rt},
			)

			Expect(caps.Types[rt]).To(Equal("dogmatiq.enginekit.marshaling.pbfixtures.Message"))
		})

		It("excludes non-protocol-buffers types", func() {
			rt := reflect.TypeOf(fixtures.MessageA{})

			caps := codec.Query(
				[]reflect.Type{rt},
			)

			Expect(caps.Types).To(BeEmpty())
		})
	})
})
