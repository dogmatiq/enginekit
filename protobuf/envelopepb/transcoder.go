package envelopepb

import (
	"mime"
	reflect "reflect"
	"strings"

	"github.com/dogmatiq/enginekit/marshaler"
	"google.golang.org/protobuf/proto"
)

// Transcoder re-encodes messages to different media-types on the fly.
type Transcoder struct {
	// MediaTypes is a map of the message's type to a list of supported
	// media-types, in order of preference.
	MediaTypes map[reflect.Type][]string

	// Marshaler is the marshaler to use to unmarshal and marshal messages.
	Marshaler marshaler.Marshaler
}

// Transcode re-encodes the message in env to one of the supported media-types.
func (t *Transcoder) Transcode(env *Envelope) (*Envelope, bool, error) {
	rt, err := t.Marshaler.UnmarshalTypeFromMediaType(env.MediaType)
	if err != nil {
		return nil, false, err
	}

	supported := t.MediaTypes[rt]

	if len(supported) == 0 {
		return nil, false, nil
	}

	// If the existing encoding is supported by the consumer use the envelope
	// without any re-encoding, even if the existing media-type is not the first
	// preference.
	for _, mediaType := range supported {
		if mediaTypeEqual(env.MediaType, mediaType) {
			return env, true, nil
		}
	}

	m, err := t.Marshaler.Unmarshal(
		marshaler.Packet{
			MediaType: env.MediaType,
			Data:      env.Data,
		},
	)
	if err != nil {
		return nil, false, err
	}

	// Otherwise, attempt to marshal the message using the recipient's requested
	// media-types in order of preference.
	packet, ok, err := t.Marshaler.MarshalAs(m, supported)
	if !ok || err != nil {
		return nil, ok, err
	}

	env = proto.Clone(env).(*Envelope)
	env.MediaType = packet.MediaType
	env.Data = packet.Data

	return env, true, nil
}

func mediaTypeEqual(a, b string) bool {
	baseA, paramsA, err := mime.ParseMediaType(a)
	if err != nil {
		panic(err)
	}

	baseB, paramsB, err := mime.ParseMediaType(b)
	if err != nil {
		panic(err)
	}

	if len(paramsA) != len(paramsB) {
		return false
	}

	if !strings.EqualFold(baseA, baseB) {
		return false
	}

	for k, va := range paramsA {
		k = strings.ToLower(k)
		vb, ok := paramsB[k]
		if !ok {
			return false
		}

		if va != vb {
			return false
		}
	}

	return true
}
