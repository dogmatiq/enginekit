package api

import (
	"github.com/dogmatiq/enginekit/config"
	"github.com/dogmatiq/enginekit/config/api/internal/pb"
	"github.com/dogmatiq/marshalkit"
)

type unmarshalError string

// marshalIdentity marshals an config.Identity to its protocol buffers
// representation.
func marshalIdentity(in config.Identity) *pb.Identity {
	return &pb.Identity{
		Name: in.Name,
		Key:  in.Key,
	}
}

// unmarshalIdentity unmarshals an config.Identity from its protocol buffers
// representation.
func unmarshalIdentity(in *pb.Identity) config.Identity {
	return config.MustNewIdentity(in.Name, in.Key)
}

// marshalApplication marshals a config.ApplicationConfig to its protocol
// buffers representation.
func marshalApplication(m *marshalkit.Marshaler, in *config.ApplicationConfig) *pb.ApplicationConfig {
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

// unmarshalApplication unmarshals a config.ApplicationConfig from its protocol
// buffers representation.
func unmarshalApplication(m *marshalkit.Marshaler, in *pb.ApplicationConfig) *config.ApplicationConfig {
	out := &config.ApplicationConfig{
		ApplicationIdentity: unmarshalIdentity(in.GetIdentity()),
		HandlersByName:      map[string]config.HandlerConfig{},
		HandlersByKey:       map[string]config.HandlerConfig{},
		Roles:               config.MessageRoleMap{},
		Consumers:           map[config.MessageType][]config.HandlerConfig{},
		Producers:           map[config.MessageType][]config.HandlerConfig{},
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
func marshalHandler(m *marshalkit.Marshaler, in config.HandlerConfig) *pb.HandlerConfig {
	t, err := in.HandlerType().MarshalBinary()
	assertOk(err)

	return &pb.HandlerConfig{
		Identity: marshalIdentity(in.Identity()),
		Type:     t,
		Consumed: marshalRoleMap(m, in.ConsumedMessageTypes()),
		Produced: marshalRoleMap(m, in.ProducedMessageTypes()),
	}
}

// unmarshalHandler unmarshals a config.HandlerConfig from its protocol buffers
// representation.
func unmarshalHandler(m *marshalkit.Marshaler, in *pb.HandlerConfig) config.HandlerConfig {
	var t config.HandlerType
	err := t.UnmarshalBinary(in.Type)
	assertOk(err)

	i := unmarshalIdentity(in.GetIdentity())
	c := unmarshalRoleMap(m, in.Consumed)
	p := unmarshalRoleMap(m, in.Produced)

	switch t {
	case config.AggregateHandlerType:
		return &config.AggregateConfig{
			HandlerIdentity: i,
			Consumed:        c,
			Produced:        p,
		}
	case config.ProcessHandlerType:
		return &config.ProcessConfig{
			HandlerIdentity: i,
			Consumed:        c,
			Produced:        p,
		}
	case config.IntegrationHandlerType:
		return &config.IntegrationConfig{
			HandlerIdentity: i,
			Consumed:        c,
			Produced:        p,
		}
	default: // config.ProjectionHandlerType:
		return &config.ProjectionConfig{
			HandlerIdentity: i,
			Consumed:        c,
		}
	}
}

// marshalRoleMap marshals a config.MessageRoleMap to its protocol buffers
// representation.
func marshalRoleMap(m *marshalkit.Marshaler, in config.MessageRoleMap) map[string][]byte {
	var out map[string][]byte

	for mt, r := range in {
		if out == nil {
			out = map[string][]byte{}
		}

		k, err := config.MarshalMessageType(m, mt)
		assertOk(err)
		v, err := r.MarshalBinary()
		assertOk(err)

		out[k] = v
	}

	return out
}

// unmarshalRoleMap unmarshals a config.MessageRoleMap from its protocol buffers
// representation.
func unmarshalRoleMap(m *marshalkit.Marshaler, in map[string][]byte) config.MessageRoleMap {
	var out config.MessageRoleMap

	var v config.MessageRole

	for mt, r := range in {
		if out == nil {
			out = config.MessageRoleMap{}
		}

		k, err := config.UnmarshalMessageType(m, mt)
		assertOk(err)
		err = v.UnmarshalBinary(r)
		assertOk(err)

		out[k] = v
	}

	return out
}

type sentinel struct {
	cause error
}

func assertOk(err error) {
	if err != nil {
		panic(sentinel{err})
	}
}

func catch(err *error) {
	switch r := recover().(type) {
	case error:
		*err = r
	case nil:
		return
	default:
		panic(r)
	}
}
