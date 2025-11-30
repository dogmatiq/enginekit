package prototest

import (
	"google.golang.org/protobuf/proto"
)

type T interface {
}

// Equal fails t if got and want are not equal according to [proto.Equal].
func Equal[M proto.Message](
	t T,
	got, want M,
) {
	if proto.Equal(got, want) {
		return
	}

	// TODO: use protoreflect to test field-by-field to produce meaningful error emssages.

}
