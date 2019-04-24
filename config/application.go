package config

import (
	"context"

	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/enginekit/message"
)

// ApplicationConfig represents the configuration of an entire Dogma application.
type ApplicationConfig struct {
	// Application is the application that the configuration applies to.
	Application dogma.Application

	// ApplicationName is the application's name, as specified in the dogma.App struct.
	ApplicationName string

	// Handlers is a map of handler name to their respective configuration.
	Handlers map[string]HandlerConfig

	// Roles is a map of message type to the role it performs within the
	// application.
	Roles message.RoleMap

	// Consumers is a map of message type to the names of the handlers that
	// consume messages of that type.
	Consumers map[message.Type][]string

	// Producers is a map of message type to the name of the handlers that
	// produce messages of that type.
	Producers map[message.Type][]string
}

// NewApplicationConfig returns a new application config for the given application.
func NewApplicationConfig(app dogma.Application) (*ApplicationConfig, error) {
	cfg := &ApplicationConfig{
		Application: app,
		Handlers:    map[string]HandlerConfig{},
		Roles:       message.RoleMap{},
		Consumers:   map[message.Type][]string{},
		Producers:   map[message.Type][]string{},
	}

	c := &applicationConfigurer{
		cfg: cfg,
	}

	if err := catch(func() {
		app.Configure(c)
	}); err != nil {
		return nil, err
	}

	if c.cfg.ApplicationName == "" {
		return nil, errorf(
			"%T.Configure() did not call ApplicationConfigurer.Name()",
			app,
		)
	}
	return cfg, nil
}

// Name returns the application name.
func (c *ApplicationConfig) Name() string {
	return c.ApplicationName
}

// Accept calls v.VisitApplicationConfig(ctx, c).
func (c *ApplicationConfig) Accept(ctx context.Context, v Visitor) error {
	return v.VisitApplicationConfig(ctx, c)
}

// register adds the given handler configuration to the app configuration.
func (c *ApplicationConfig) register(cfg HandlerConfig) {
	n := cfg.Name()

	if x, ok := c.Handlers[n]; ok {
		panicf(
			"%s can not use the handler name %#v, because it is already used by %s",
			cfg.HandlerReflectType(),
			n,
			x.HandlerReflectType(),
		)
	}

	for t, r := range cfg.ConsumedMessageTypes() {
		c.checkMessageType(cfg, t, r)

		if r == message.CommandRole {
			if x, ok := c.Consumers[t]; ok {
				panicf(
					"the %#v handler can not consume %s commands because they are already consumed by %#v",
					cfg.Name(),
					t,
					x[0],
				)
			}
		}
	}

	for t, r := range cfg.ProducedMessageTypes() {
		c.checkMessageType(cfg, t, r)

		if r == message.EventRole {
			if x, ok := c.Producers[t]; ok {
				panicf(
					"the %#v handler can not produce %s events because they are already produced by %#v",
					cfg.Name(),
					t,
					x[0],
				)
			}
		}
	}

	c.Handlers[n] = cfg

	for t, r := range cfg.ConsumedMessageTypes() {
		c.Roles.Add(t, r)
		c.Consumers[t] = append(c.Consumers[t], cfg.Name())
	}

	for t, r := range cfg.ProducedMessageTypes() {
		c.Roles.Add(t, r)
		c.Producers[t] = append(c.Producers[t], cfg.Name())
	}
}

// checkMessageType panics if the message type has already been registered with
// a different role.
func (c *ApplicationConfig) checkMessageType(
	cfg HandlerConfig,
	t message.Type,
	r message.Role,
) {
	x, ok := c.Roles[t]
	if !ok || x == r {
		return
	}

	var h string

	if n, ok := c.Consumers[t]; ok {
		h = n[0]
	} else if n, ok := c.Producers[t]; ok {
		h = n[0]
	}

	panicf(
		"the %#v handler configures %s as a %s but %#v configures it as a %s",
		cfg.Name(),
		t,
		r,
		h,
		x,
	)
}

// applicationConfigurer is an implementation of dogma.ApplicationConfigurer
// that builds an ApplicationConfig value.
type applicationConfigurer struct {
	cfg *ApplicationConfig
}

func (c *applicationConfigurer) Name(n string) {
	if c.cfg.ApplicationName != "" {
		panicf(
			`%T.Configure() has already called ApplicationConfigurer.Name(%#v)`,
			c.cfg.Application,
			c.cfg.ApplicationName,
		)
	}

	if !IsValidName(n) {
		panicf(
			`%T.Configure() called ApplicationConfigurer.Name(%#v) with an invalid name`,
			c.cfg.Application,
			n,
		)
	}

	c.cfg.ApplicationName = n
}

// RegisterAggregate configures the engine to route messages to h.
func (c *applicationConfigurer) RegisterAggregate(h dogma.AggregateMessageHandler) {
	cfg, err := NewAggregateConfig(h)
	if err != nil {
		panic(err)
	}

	c.cfg.register(cfg)
}

// RegisterProcess configures the engine to route messages to h.
func (c *applicationConfigurer) RegisterProcess(h dogma.ProcessMessageHandler) {
	cfg, err := NewProcessConfig(h)
	if err != nil {
		panic(err)
	}

	c.cfg.register(cfg)
}

// RegisterIntegration configures the engine to route messages to h.
func (c *applicationConfigurer) RegisterIntegration(h dogma.IntegrationMessageHandler) {
	cfg, err := NewIntegrationConfig(h)
	if err != nil {
		panic(err)
	}

	c.cfg.register(cfg)
}

// RegisterProjection configures the engine to route messages to h.
func (c *applicationConfigurer) RegisterProjection(h dogma.ProjectionMessageHandler) {
	cfg, err := NewProjectionConfig(h)
	if err != nil {
		panic(err)
	}

	c.cfg.register(cfg)
}
