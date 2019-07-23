package config_test

import (
	"github.com/dogmatiq/dogma"
	. "github.com/dogmatiq/enginekit/config"
	"github.com/dogmatiq/enginekit/fixtures"
	"github.com/dogmatiq/enginekit/message"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ Config = &ApplicationConfig{}

var _ = Describe("type ApplicationConfig", func() {
	Describe("func NewApplicationConfig", func() {
		var (
			aggregate   *fixtures.AggregateMessageHandler
			process     *fixtures.ProcessMessageHandler
			integration *fixtures.IntegrationMessageHandler
			projection  *fixtures.ProjectionMessageHandler
			app         *fixtures.Application
		)

		BeforeEach(func() {
			aggregate = &fixtures.AggregateMessageHandler{
				ConfigureFunc: func(c dogma.AggregateConfigurer) {
					c.Identity("<aggregate-name>", "<aggregate-key>")
					c.ConsumesCommandType(fixtures.MessageA{})
					c.ProducesEventType(fixtures.MessageE{})
				},
			}

			process = &fixtures.ProcessMessageHandler{
				ConfigureFunc: func(c dogma.ProcessConfigurer) {
					c.Identity("<process-name>", "<process-key>")
					c.ConsumesEventType(fixtures.MessageB{})
					c.ConsumesEventType(fixtures.MessageE{}) // shared with <projection-name>
					c.ProducesCommandType(fixtures.MessageC{})
					c.SchedulesTimeoutType(fixtures.MessageT{})
				},
			}

			integration = &fixtures.IntegrationMessageHandler{
				ConfigureFunc: func(c dogma.IntegrationConfigurer) {
					c.Identity("<integration-name>", "<integration-key>")
					c.ConsumesCommandType(fixtures.MessageC{})
					c.ProducesEventType(fixtures.MessageF{})
				},
			}

			projection = &fixtures.ProjectionMessageHandler{
				ConfigureFunc: func(c dogma.ProjectionConfigurer) {
					c.Identity("<projection-name>", "<projection-key>")
					c.ConsumesEventType(fixtures.MessageD{})
					c.ConsumesEventType(fixtures.MessageE{}) // shared with <process-name>
				},
			}

			app = &fixtures.Application{
				ConfigureFunc: func(c dogma.ApplicationConfigurer) {
					c.Identity("<application-name>", "<application-key>")
					c.RegisterAggregate(aggregate)
					c.RegisterProcess(process)
					c.RegisterIntegration(integration)
					c.RegisterProjection(projection)
				},
			}
		})

		When("the configuration is valid", func() {
			var cfg *ApplicationConfig

			BeforeEach(func() {
				var err error
				cfg, err = NewApplicationConfig(app)
				Expect(err).ShouldNot(HaveOccurred())
			})

			It("sets the app name", func() {
				Expect(cfg.ApplicationName).To(Equal("<application-name>"))
			})

			It("sets the app key", func() {
				Expect(cfg.ApplicationKey).To(Equal("<application-key>"))
			})

			It("sets the roles map", func() {
				Expect(cfg.Roles).To(Equal(
					message.RoleMap{
						fixtures.MessageAType: message.CommandRole,
						fixtures.MessageBType: message.EventRole,
						fixtures.MessageCType: message.CommandRole,
						fixtures.MessageDType: message.EventRole,
						fixtures.MessageEType: message.EventRole,
						fixtures.MessageFType: message.EventRole,
						fixtures.MessageTType: message.TimeoutRole,
					},
				))
			})

			It("sets the consumers map", func() {
				Expect(cfg.Consumers).To(Equal(
					map[message.Type][]HandlerConfig{
						fixtures.MessageAType: {cfg.HandlersByName["<aggregate-name>"]},
						fixtures.MessageBType: {cfg.HandlersByName["<process-name>"]},
						fixtures.MessageCType: {cfg.HandlersByName["<integration-name>"]},
						fixtures.MessageDType: {cfg.HandlersByName["<projection-name>"]},
						fixtures.MessageEType: {cfg.HandlersByName["<process-name>"], cfg.HandlersByName["<projection-name>"]},
						fixtures.MessageTType: {cfg.HandlersByName["<process-name>"]},
					},
				))
			})

			It("sets the producers map", func() {
				Expect(cfg.Producers).To(Equal(
					map[message.Type][]HandlerConfig{
						fixtures.MessageCType: {cfg.HandlersByName["<process-name>"]},
						fixtures.MessageEType: {cfg.HandlersByName["<aggregate-name>"]},
						fixtures.MessageFType: {cfg.HandlersByName["<integration-name>"]},
						fixtures.MessageTType: {cfg.HandlersByName["<process-name>"]},
					},
				))
			})

			Describe("func Name()", func() {
				It("returns the app name", func() {
					Expect(cfg.Name()).To(Equal("<application-name>"))
				})
			})

			Describe("func Key()", func() {
				It("returns the app key", func() {
					Expect(cfg.Key()).To(Equal("<application-key>"))
				})
			})
		})

		When("the app does not configure an identity", func() {
			BeforeEach(func() {
				app.ConfigureFunc = nil
			})

			It("returns a descriptive error", func() {
				_, err := NewApplicationConfig(app)

				Expect(err).To(Equal(
					Error(
						"*fixtures.Application.Configure() did not call ApplicationConfigurer.Identity()",
					),
				))
			})
		})

		When("the app configures multiple identities", func() {
			BeforeEach(func() {
				app.ConfigureFunc = func(c dogma.ApplicationConfigurer) {
					c.Identity("<application-name-1>", "<application-key-1>")
					c.Identity("<application-name-2>", "<application-key-2>")
				}
			})

			It("returns a descriptive error", func() {
				_, err := NewApplicationConfig(app)

				Expect(err).To(Equal(
					Error(
						`*fixtures.Application.Configure() has already called ApplicationConfigurer.Identity("<application-name-1>", "<application-key-1>")`,
					),
				))
			})
		})

		When("the app configures an invalid application name", func() {
			BeforeEach(func() {
				app.ConfigureFunc = func(c dogma.ApplicationConfigurer) {
					c.Identity("\t \n", "<application-key>")
				}
			})

			It("returns a descriptive error", func() {
				_, err := NewApplicationConfig(app)

				Expect(err).To(Equal(
					Error(
						`*fixtures.Application.Configure() called ApplicationConfigurer.Identity() with an invalid name "\t \n"`,
					),
				))
			})
		})

		When("the app configures an invalid application key", func() {
			BeforeEach(func() {
				app.ConfigureFunc = func(c dogma.ApplicationConfigurer) {
					c.Identity("<application-name>", "\t \n")
				}
			})

			It("returns a descriptive error", func() {
				_, err := NewApplicationConfig(app)

				Expect(err).To(Equal(
					Error(
						`*fixtures.Application.Configure() called ApplicationConfigurer.Identity() with an invalid key "\t \n"`,
					),
				))
			})
		})

		When("the app contains an invalid handler configurations", func() {
			It("returns an error when an aggregate is misconfigured", func() {
				aggregate.ConfigureFunc = nil

				_, err := NewApplicationConfig(app)

				Expect(err).Should(HaveOccurred())
			})

			It("returns an error when a process is misconfigured", func() {
				process.ConfigureFunc = nil

				_, err := NewApplicationConfig(app)

				Expect(err).Should(HaveOccurred())
			})

			It("returns an error when an integration is misconfigured", func() {
				integration.ConfigureFunc = nil

				_, err := NewApplicationConfig(app)

				Expect(err).Should(HaveOccurred())
			})

			It("returns an error when a projection is misconfigured", func() {
				projection.ConfigureFunc = nil

				_, err := NewApplicationConfig(app)

				Expect(err).Should(HaveOccurred())
			})
		})

		When("the app contains conflicting handler identities", func() {
			It("returns an error when an aggregate handler name is in conflict", func() {
				aggregate.ConfigureFunc = func(c dogma.AggregateConfigurer) {
					c.Identity("<process-name>", "<aggregate-key>") // conflict!
					c.ConsumesCommandType(fixtures.MessageA{})
					c.ProducesEventType(fixtures.MessageE{})
				}

				app.ConfigureFunc = func(c dogma.ApplicationConfigurer) {
					c.Identity("<application-name>", "<application-key>")
					c.RegisterProcess(process)
					c.RegisterAggregate(aggregate) // register the conflicting aggregate last
				}

				_, err := NewApplicationConfig(app)

				Expect(err).To(Equal(
					Error(
						`*fixtures.AggregateMessageHandler can not use the handler name "<process-name>", because it is already used by *fixtures.ProcessMessageHandler`,
					),
				))
			})

			It("returns an error when an aggregate handler key is in conflict", func() {
				aggregate.ConfigureFunc = func(c dogma.AggregateConfigurer) {
					c.Identity("<aggregate-name>", "<process-key>") // conflict!
					c.ConsumesCommandType(fixtures.MessageA{})
					c.ProducesEventType(fixtures.MessageE{})
				}

				app.ConfigureFunc = func(c dogma.ApplicationConfigurer) {
					c.Identity("<application-name>", "<application-key>")
					c.RegisterProcess(process)
					c.RegisterAggregate(aggregate) // register the conflicting aggregate last
				}

				_, err := NewApplicationConfig(app)

				Expect(err).To(Equal(
					Error(
						`*fixtures.AggregateMessageHandler can not use the handler key "<process-key>", because it is already used by *fixtures.ProcessMessageHandler`,
					),
				))
			})

			It("returns an error when a process handler name is in conflict", func() {
				process.ConfigureFunc = func(c dogma.ProcessConfigurer) {
					c.Identity("<aggregate-name>", "<process-key>") // conflict!
					c.ConsumesEventType(fixtures.MessageB{})
					c.ProducesCommandType(fixtures.MessageC{})
				}

				_, err := NewApplicationConfig(app)

				Expect(err).To(Equal(
					Error(
						`*fixtures.ProcessMessageHandler can not use the handler name "<aggregate-name>", because it is already used by *fixtures.AggregateMessageHandler`,
					),
				))
			})

			It("returns an error when a process handler key is in conflict", func() {
				process.ConfigureFunc = func(c dogma.ProcessConfigurer) {
					c.Identity("<process-name>", "<aggregate-key>") // conflict!
					c.ConsumesEventType(fixtures.MessageB{})
					c.ProducesCommandType(fixtures.MessageC{})
				}

				_, err := NewApplicationConfig(app)

				Expect(err).To(Equal(
					Error(
						`*fixtures.ProcessMessageHandler can not use the handler key "<aggregate-key>", because it is already used by *fixtures.AggregateMessageHandler`,
					),
				))
			})

			It("returns an error when an integration handler name is in conflict", func() {
				integration.ConfigureFunc = func(c dogma.IntegrationConfigurer) {
					c.Identity("<process-name>", "<integration-key>") // conflict!
					c.ConsumesCommandType(fixtures.MessageC{})
					c.ProducesEventType(fixtures.MessageF{})
				}

				_, err := NewApplicationConfig(app)

				Expect(err).To(Equal(
					Error(
						`*fixtures.IntegrationMessageHandler can not use the handler name "<process-name>", because it is already used by *fixtures.ProcessMessageHandler`,
					),
				))
			})

			It("returns an error when an integration handler key is in conflict", func() {
				integration.ConfigureFunc = func(c dogma.IntegrationConfigurer) {
					c.Identity("<integration-name>", "<process-key>") // conflict!
					c.ConsumesCommandType(fixtures.MessageC{})
					c.ProducesEventType(fixtures.MessageF{})
				}

				_, err := NewApplicationConfig(app)

				Expect(err).To(Equal(
					Error(
						`*fixtures.IntegrationMessageHandler can not use the handler key "<process-key>", because it is already used by *fixtures.ProcessMessageHandler`,
					),
				))
			})

			It("returns an error when a projection handler name is in conflict", func() {
				projection.ConfigureFunc = func(c dogma.ProjectionConfigurer) {
					c.Identity("<integration-name>", "<projection-key>") // conflict!
					c.ConsumesEventType(fixtures.MessageD{})
				}

				_, err := NewApplicationConfig(app)

				Expect(err).To(Equal(
					Error(
						`*fixtures.ProjectionMessageHandler can not use the handler name "<integration-name>", because it is already used by *fixtures.IntegrationMessageHandler`,
					),
				))
			})

			It("returns an error when a projection handler key is in conflict", func() {
				projection.ConfigureFunc = func(c dogma.ProjectionConfigurer) {
					c.Identity("<projection-name>", "<integration-key>") // conflict!
					c.ConsumesEventType(fixtures.MessageD{})
				}

				_, err := NewApplicationConfig(app)

				Expect(err).To(Equal(
					Error(
						`*fixtures.ProjectionMessageHandler can not use the handler key "<integration-key>", because it is already used by *fixtures.IntegrationMessageHandler`,
					),
				))
			})
		})

		It("returns an error when the app contains multiple consumers of the same command", func() {
			integration.ConfigureFunc = func(c dogma.IntegrationConfigurer) {
				c.Identity("<integration-name>", "<integration-key>")
				c.ConsumesCommandType(fixtures.MessageA{}) // conflict with <aggregate-name>
				c.ProducesEventType(fixtures.MessageF{})
			}

			_, err := NewApplicationConfig(app)

			Expect(err).To(Equal(
				Error(
					`the "<integration-name>" handler can not consume fixtures.MessageA commands because they are already consumed by "<aggregate-name>"`,
				),
			))
		})

		It("returns an error when the app contains multiple producers of the same event", func() {
			integration.ConfigureFunc = func(c dogma.IntegrationConfigurer) {
				c.Identity("<integration-name>", "<integration-key>")
				c.ConsumesCommandType(fixtures.MessageC{})
				c.ProducesEventType(fixtures.MessageE{}) // conflict with <aggregate-name>
			}

			_, err := NewApplicationConfig(app)

			Expect(err).To(Equal(
				Error(
					`the "<integration-name>" handler can not produce fixtures.MessageE events because they are already produced by "<aggregate-name>"`,
				),
			))
		})

		It("does not return an error when the app contains multiple processes that schedule the same timeout", func() {
			process1 := &fixtures.ProcessMessageHandler{
				ConfigureFunc: func(c dogma.ProcessConfigurer) {
					c.Identity("<process-name-1>", "<process-key-1>")
					c.ConsumesEventType(fixtures.MessageB{})
					c.ProducesCommandType(fixtures.MessageC{})
					c.SchedulesTimeoutType(fixtures.MessageT{})
				},
			}

			process2 := &fixtures.ProcessMessageHandler{
				ConfigureFunc: func(c dogma.ProcessConfigurer) {
					c.Identity("<process-name-2>", "<process-key-2>")
					c.ConsumesEventType(fixtures.MessageB{})
					c.ProducesCommandType(fixtures.MessageC{})
					c.SchedulesTimeoutType(fixtures.MessageT{})
				},
			}

			app := &fixtures.Application{
				ConfigureFunc: func(c dogma.ApplicationConfigurer) {
					c.Identity("<application-name>", "<application-key>")
					c.RegisterProcess(process1)
					c.RegisterProcess(process2)
				},
			}

			_, err := NewApplicationConfig(app)

			Expect(err).ShouldNot(HaveOccurred())
		})

		When("multiple handlers use a single message type in differing roles", func() {
			It("returns an error when a conflict occurs with a consumed message", func() {
				process.ConfigureFunc = func(c dogma.ProcessConfigurer) {
					c.Identity("<process-name>", "<process-key>")
					c.ConsumesEventType(fixtures.MessageA{}) // conflict with <aggregate-name>
					c.ProducesCommandType(fixtures.MessageC{})
				}

				_, err := NewApplicationConfig(app)

				Expect(err).To(Equal(
					Error(
						`the "<process-name>" handler configures fixtures.MessageA as an event but "<aggregate-name>" configures it as a command`,
					),
				))
			})

			It("returns an error when a conflict occurs with a produced message", func() {
				process.ConfigureFunc = func(c dogma.ProcessConfigurer) {
					c.Identity("<process-name>", "<process-key>")
					c.ConsumesEventType(fixtures.MessageB{})
					c.ProducesCommandType(fixtures.MessageE{}) // conflict with <aggregate-name>
				}

				_, err := NewApplicationConfig(app)

				Expect(err).To(Equal(
					Error(
						`the "<process-name>" handler configures fixtures.MessageE as a command but "<aggregate-name>" configures it as an event`,
					),
				))
			})
		})
	})
})
