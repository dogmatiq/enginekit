package config_test

import (
	"github.com/dogmatiq/dogma"
	. "github.com/dogmatiq/enginekit/config"
	"github.com/dogmatiq/enginekit/fixtures"
	handlerkit "github.com/dogmatiq/enginekit/handler"
	"github.com/dogmatiq/enginekit/message"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ HandlerConfig = &ProjectionConfig{}

var _ = Describe("type ProjectionConfig", func() {
	Describe("func NewProjectionConfig", func() {
		var handler *fixtures.ProjectionMessageHandler

		BeforeEach(func() {
			handler = &fixtures.ProjectionMessageHandler{
				ConfigureFunc: func(c dogma.ProjectionConfigurer) {
					c.Name("<name>")
					c.AcceptsEventType(fixtures.MessageA{})
					c.AcceptsEventType(fixtures.MessageB{})
				},
			}
		})

		When("the configuration is valid", func() {
			var cfg *ProjectionConfig

			BeforeEach(func() {
				var err error
				cfg, err = NewProjectionConfig(handler)
				Expect(err).ShouldNot(HaveOccurred())
			})

			It("the handler name is set", func() {
				Expect(cfg.HandlerName).To(Equal("<name>"))
			})

			Describe("func Name()", func() {
				It("returns the handler name", func() {
					Expect(cfg.Name()).To(Equal("<name>"))
				})
			})

			Describe("func HandlerType()", func() {
				It("returns handler.ProjectionType", func() {
					Expect(cfg.HandlerType()).To(Equal(handlerkit.ProjectionType))
				})
			})

			Describe("func ConsumedMessageTypes()", func() {
				It("returns the expected message types", func() {
					Expect(cfg.ConsumedMessageTypes()).To(Equal(
						map[message.Type]message.Role{
							fixtures.MessageAType: message.EventRole,
							fixtures.MessageBType: message.EventRole,
						},
					))
				})
			})

			Describe("func ProducedMessageTypes()", func() {
				It("returns an empty map", func() {
					Expect(cfg.ProducedMessageTypes()).To(BeEmpty())
				})
			})
		})

		DescribeTable(
			"when the configuration is invalid",
			func(
				msg string,
				fn func(dogma.ProjectionConfigurer),
			) {
				handler.ConfigureFunc = fn

				_, err := NewProjectionConfig(handler)
				Expect(err).Should(HaveOccurred())

				if msg != "" {
					Expect(err).To(MatchError(msg))
				}
			},
			Entry(
				"when the handler does not configure anything",
				"", // any error
				nil,
			),
			Entry(
				"when the handler does not configure a name",
				`*fixtures.ProjectionMessageHandler.Configure() did not call ProjectionConfigurer.Name()`,
				func(c dogma.ProjectionConfigurer) {
					c.AcceptsEventType(fixtures.MessageA{})
				},
			),
			Entry(
				"when the handler configures multiple names",
				`*fixtures.ProjectionMessageHandler.Configure() has already called ProjectionConfigurer.Name("<name>")`,
				func(c dogma.ProjectionConfigurer) {
					c.Name("<name>")
					c.Name("<other>")
					c.AcceptsEventType(fixtures.MessageA{})
				},
			),
			Entry(
				"when the handler configures an invalid name",
				`*fixtures.ProjectionMessageHandler.Configure() called ProjectionConfigurer.Name("\t \n") with an invalid name`,
				func(c dogma.ProjectionConfigurer) {
					c.Name("\t \n")
					c.AcceptsEventType(fixtures.MessageA{})
				},
			),
			Entry(
				"when the handler does not configure any accepted event types",
				`*fixtures.ProjectionMessageHandler.Configure() did not call ProjectionConfigurer.AcceptsEventType()`,
				func(c dogma.ProjectionConfigurer) {
					c.Name("<name>")
				},
			),
			Entry(
				"when the handler configures the same accepted event type multiple times",
				`*fixtures.ProjectionMessageHandler.Configure() has already called ProjectionConfigurer.AcceptsEventType(fixtures.MessageA)`,
				func(c dogma.ProjectionConfigurer) {
					c.Name("<name>")
					c.AcceptsEventType(fixtures.MessageA{})
					c.AcceptsEventType(fixtures.MessageA{})
				},
			),
		)
	})
})
