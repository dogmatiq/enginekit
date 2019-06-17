package protobuf

import (
	"fmt"
	"reflect"

	"github.com/golang/protobuf/proto"
)

// cast performs a type assertion on v to verify it is a protocol buffers
// message, otherwise it returns a meaningful error.
func cast(v interface{}) (proto.Message, error) {
	if m, ok := v.(proto.Message); ok {
		return m, nil
	}

	return nil, fmt.Errorf(
		"'%s' is not a protocol buffers message",
		reflect.TypeOf(v),
	)
}
