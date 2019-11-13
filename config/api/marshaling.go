package api

import (
	"github.com/dogmatiq/enginekit/config"
	"github.com/dogmatiq/enginekit/config/api/internal/pb"
	"github.com/dogmatiq/enginekit/handler"
	"github.com/dogmatiq/enginekit/identity"
	"github.com/dogmatiq/enginekit/marshaling"
	"github.com/dogmatiq/enginekit/message"
)

type unmarshalError string

// marshalIdentity marshals an identity.Identity to its protocol buffers
// representation.
func marshalIdentity(in identity.Identity) *pb.Identity {
	return &pb.Identity{
		Name: in.Name,
		Key:  in.Key,
	}
}

// unmarshalIdentity unmarshals an identity.Identity from its protocol buffers
// representation.
func unmarshalIdentity(in *pb.Identity) identity.Identity {
	return identity.MustNew(in.Name, in.Key)
}

// marshalApplication marshals a config.ApplicationConfig to its protocol
// buffers representation.
func marshalApplication(m *marshaling.Marshaler, in *config.ApplicationConfig) *pb.ApplicationConfig {
	out := &pb.ApplicationConfig{
		Identity: marshalIdentity(in.Identity()),
	}

	for _, cfg := range in.HandlersByKey {
		out.Handlers = append(
			out.Handlers,
			marshalHandler(m, cfg),
		)
	}

	return out
}

// unmarshalApplication marshals a config.ApplicationConfig from its protocol
// buffers representation.
func unmarshalApplication(m *marshaling.Marshaler, in *pb.ApplicationConfig) *config.ApplicationConfig {
	out := &config.ApplicationConfig{
		ApplicationIdentity: unmarshalIdentity(in.GetIdentity()),
		HandlersByName:      map[string]config.HandlerConfig{},
		HandlersByKey:       map[string]config.HandlerConfig{},
		Roles:               message.RoleMap{},
		Consumers:           map[message.Type][]config.HandlerConfig{},
		Producers:           map[message.Type][]config.HandlerConfig{},
	}

	for _, h := range in.GetHandlers() {
		cfg := unmarshalHandler(m, h)

		out.HandlersByName[cfg.Identity().Name] = cfg
		out.HandlersByKey[cfg.Identity().Key] = cfg

		for mt, r := range cfg.ConsumedMessageTypes() {
			out.Roles[mt] = r
			out.Consumers[mt] = append(out.Consumers[mt], cfg)
		}

		for mt, r := range cfg.ProducedMessageTypes() {
			out.Roles[mt] = r
			out.Producers[mt] = append(out.Producers[mt], cfg)
		}
	}

	return out
}

// marshalHandler marshals a config.HandlerConfig to its protocol buffers
// representation.
func marshalHandler(m *marshaling.Marshaler, in config.HandlerConfig) *pb.HandlerConfig {
	return &pb.HandlerConfig{
		Identity: marshalIdentity(in.Identity()),
		Type:     string(in.HandlerType()),
		Consumed: marshalRoleMap(m, in.ConsumedMessageTypes()),
		Produced: marshalRoleMap(m, in.ProducedMessageTypes()),
	}
}

// unmarshalHandler marshals a config.HandlerConfig to its protocol buffers
// representation.
func unmarshalHandler(m *marshaling.Marshaler, in *pb.HandlerConfig) config.HandlerConfig {
	t := handler.Type(in.Type)
	t.MustValidate()

	i := unmarshalIdentity(in.GetIdentity())
	c := unmarshalRoleMap(m, in.Consumed)
	p := unmarshalRoleMap(m, in.Produced)

	switch t {
	case handler.AggregateType:
		return &config.AggregateConfig{
			HandlerIdentity: i,
			Consumed:        c,
			Produced:        p,
		}
	case handler.ProcessType:
		return &config.ProcessConfig{
			HandlerIdentity: i,
			Consumed:        c,
			Produced:        p,
		}
	case handler.IntegrationType:
		return &config.IntegrationConfig{
			HandlerIdentity: i,
			Consumed:        c,
			Produced:        p,
		}
	default: // case handler.ProjectionType:
		return &config.ProjectionConfig{
			HandlerIdentity: i,
			Consumed:        c,
		}
	}
}

// marshalRoleMap marshals a message.RoleMap to its protocol buffers
// representation.
func marshalRoleMap(m *marshaling.Marshaler, in message.RoleMap) map[string]string {
	var out map[string]string

	for mt, r := range in {
		if out == nil {
			out = map[string]string{}
		}

		k, err := marshaling.MarshalMessageType(m, mt)
		if err != nil {
			panic(err)
		}

		out[k] = string(r)
	}

	return out
}

// unmarshalRoleMap marshals a message.RoleMap from its protocol buffers
// representation.
func unmarshalRoleMap(m *marshaling.Marshaler, in map[string]string) message.RoleMap {
	var out message.RoleMap

	for mt, r := range in {
		if out == nil {
			out = message.RoleMap{}
		}

		k, err := marshaling.UnmarshalMessageType(m, mt)
		if err != nil {
			panic(err)
		}

		v := message.Role(r)
		v.MustValidate()

		out[k] = v
	}

	return out
}
