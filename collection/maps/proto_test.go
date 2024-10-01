package maps_test

import (
	"testing"

	. "github.com/dogmatiq/enginekit/collection/maps"
	. "github.com/dogmatiq/enginekit/internal/stubs"
	"google.golang.org/protobuf/proto"
	"pgregory.net/rapid"
)

func TestProtoMap(t *testing.T) {
	testMap(
		t,
		NewProto[*ProtoStubA, int],
		func(a, b *ProtoStubA) bool { return proto.Equal(a, b) },
		rapid.Custom(
			func(t *rapid.T) *ProtoStubA {
				return &ProtoStubA{
					Value: rapid.String().Draw(t, "value"),
				}
			},
		),
	)
}
