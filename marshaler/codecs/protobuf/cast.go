package protobuf

import (
	"fmt"
	"reflect"

	"google.golang.org/protobuf/proto"
)

// cast performs a type assertion on v to verify it is a protocol buffers
// message, otherwise it returns a meaningful error.
func cast(v any) (proto.Message, error) {
	if m, ok := v.(proto.Message); ok {
		return m, nil
	}

	return nil, fmt.Errorf(
		"'%s' is not a protocol buffers message",
		reflect.TypeOf(v),
	)
}
