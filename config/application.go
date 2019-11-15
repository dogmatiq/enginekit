package config

import (
	"context"

	"github.com/dogmatiq/dogma"
)

// ApplicationConfig represents the configuration of an entire Dogma application.
type ApplicationConfig struct {
	// Application is the application that the configuration applies to.
	// It is nil if the config was obtained via the config API.
	Application dogma.Application

	// ApplicationIdentity is the application's identity, as specified by its
	// Configure() method.
	ApplicationIdentity Identity

	// handlersByName is a map of handler name to their respective configuration.
	HandlersByName map[string]HandlerConfig

	// HandlersByKey is a map of handler key to their respective configuration.
	HandlersByKey map[string]HandlerConfig

	// Roles is a map of message type to the role it performs within the
	// application.
	Roles MessageRoleMap

	// Consumers is a map of message type to the handler configs of the handlers
	// that consume messages of that type.
	Consumers map[MessageType][]HandlerConfig

	// Producers is a map of message type to the handler configs of the handlers
	// that produce messages of that type.
	Producers map[MessageType][]HandlerConfig
}

// NewApplicationConfig returns a new application config for the given application.
func NewApplicationConfig(app dogma.Application) (*ApplicationConfig, error) {
	cfg := &ApplicationConfig{
		Application:    app,
		HandlersByName: map[string]HandlerConfig{},
		HandlersByKey:  map[string]HandlerConfig{},
		Roles:          MessageRoleMap{},
		Consumers:      map[MessageType][]HandlerConfig{},
		Producers:      map[MessageType][]HandlerConfig{},
	}

	c := &applicationConfigurer{
		cfg: cfg,
	}

	if err := catch(func() {
		app.Configure(c)
	}); err != nil {
		return nil, err
	}

	if c.cfg.ApplicationIdentity.IsZero() {
		return nil, errorf(
			"%T.Configure() did not call ApplicationConfigurer.Identity()",
			app,
		)
	}
	return cfg, nil
}

// Identity returns the application identity.
func (c *ApplicationConfig) Identity() Identity {
	return c.ApplicationIdentity
}

// Accept calls v.VisitApplicationConfig(ctx, c).
func (c *ApplicationConfig) Accept(ctx context.Context, v Visitor) error {
	return v.VisitApplicationConfig(ctx, c)
}

// register adds the given handler configuration to the app configuration.
func (c *ApplicationConfig) register(cfg HandlerConfig) {
	i := cfg.Identity()

	if x, ok := c.HandlersByName[i.Name]; ok {
		panicf(
			"%s can not use the handler name %#v, because it is already used by %s",
			cfg.HandlerReflectType(),
			i.Name,
			x.HandlerReflectType(),
		)
	}

	if x, ok := c.HandlersByKey[i.Key]; ok {
		panicf(
			"%s can not use the handler key %#v, because it is already used by %s",
			cfg.HandlerReflectType(),
			i.Key,
			x.HandlerReflectType(),
		)
	}

	for t, r := range cfg.ConsumedMessageTypes() {
		c.checkMessageType(cfg, t, r)

		if r == CommandMessageRole {
			if x, ok := c.Consumers[t]; ok {
				panicf(
					"the %#v handler can not consume %s commands because they are already consumed by %#v",
					i.Name,
					t,
					x[0].Identity().Name,
				)
			}
		}
	}

	for t, r := range cfg.ProducedMessageTypes() {
		c.checkMessageType(cfg, t, r)

		if r == EventMessageRole {
			if x, ok := c.Producers[t]; ok {
				panicf(
					"the %#v handler can not produce %s events because they are already produced by %#v",
					i.Name,
					t,
					x[0].Identity().Name,
				)
			}
		}
	}

	c.HandlersByName[i.Name] = cfg
	c.HandlersByKey[i.Key] = cfg

	for t, r := range cfg.ConsumedMessageTypes() {
		c.Roles.Add(t, r)
		c.Consumers[t] = append(c.Consumers[t], cfg)
	}

	for t, r := range cfg.ProducedMessageTypes() {
		c.Roles.Add(t, r)
		c.Producers[t] = append(c.Producers[t], cfg)
	}
}

// checkMessageType panics if the message type has already been registered with
// a different role.
func (c *ApplicationConfig) checkMessageType(
	cfg HandlerConfig,
	t MessageType,
	r MessageRole,
) {
	x, ok := c.Roles[t]
	if !ok || x == r {
		return
	}

	var h string

	if hc, ok := c.Consumers[t]; ok {
		h = hc[0].Identity().Name
	} else if hc, ok := c.Producers[t]; ok {
		h = hc[0].Identity().Name
	}

	panicf(
		"the %#v handler configures %s as a %s but %#v configures it as a %s",
		cfg.Identity().Name,
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

func (c *applicationConfigurer) Identity(n, k string) {
	if !c.cfg.ApplicationIdentity.IsZero() {
		panicf(
			`%T.Configure() has already called ApplicationConfigurer.Identity(%#v, %#v)`,
			c.cfg.Application,
			c.cfg.ApplicationIdentity.Name,
			c.cfg.ApplicationIdentity.Key,
		)
	}

	i, err := NewIdentity(n, k)
	if err != nil {
		panicf(
			`%T.Configure() called ApplicationConfigurer.Identity() with an %s`,
			c.cfg.Application,
			err,
		)
	}

	c.cfg.ApplicationIdentity = i
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
