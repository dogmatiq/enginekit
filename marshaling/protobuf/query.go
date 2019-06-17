package protobuf

import (
	"reflect"

	"github.com/dogmatiq/enginekit/marshaling"
	"github.com/golang/protobuf/proto"
)

type codec struct{}

// Query returns the capabilities of the codec for the given types.
func (codec) Query(types []reflect.Type) marshaling.CodecCapabilities {
	caps := marshaling.CodecCapabilities{
		Types: map[reflect.Type]string{},
	}

	for _, rt := range types {
		z := reflect.Zero(rt).Interface()

		if m, ok := z.(proto.Message); ok {
			if n := proto.MessageName(m); n != "" {
				caps.Types[rt] = n
			}
		}
	}

	return caps
}
