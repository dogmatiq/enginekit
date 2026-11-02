package envelope

import (
	"github.com/dogmatiq/protobuf/envelopepb"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

// packAsAny returns x as an [*anypb.Any], panicking if it is nil or cannot be
// packed.
//
// Options pack their value when they are created, rather than when they are
// applied, so that an invalid value is reported where the option is created and
// is packed only once.
func packAsAny(x proto.Message) *anypb.Any {
	body := &envelopepb.Body{}
	envelopepb.SetExtension(body, x)
	return body.GetExtensions()[0]
}
