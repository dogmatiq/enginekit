package sets_test

import (
	"testing"

	. "github.com/dogmatiq/enginekit/collection/sets"
	. "github.com/dogmatiq/enginekit/internal/stubs"
	"google.golang.org/protobuf/proto"
	"pgregory.net/rapid"
)

func TestProtoSet(t *testing.T) {
	testSet(
		t,
		NewProto[*ProtoStubA],
		func(a, b *ProtoStubA) bool { return proto.Equal(a, b) },
		func(m *ProtoStubA) bool { return len(m.Value)%2 == 0 },
		rapid.Custom(
			func(t *rapid.T) *ProtoStubA {
				return &ProtoStubA{
					Value: rapid.String().Draw(t, "value"),
				}
			},
		),
	)
}
