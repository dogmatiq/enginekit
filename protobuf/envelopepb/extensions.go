package envelopepb

import "google.golang.org/protobuf/proto"

// SetExtension sets x as an extension on b. If an extension with the same
// type URL is already present, it is replaced.
//
// It panics if x is nil.
func SetExtension[
	T interface {
		*E
		proto.Message
	},
	E any,
](body *Body, x proto.Message) {
	body.SetExtensions(
		appendOrReplace(
			body.GetExtensions(),
			marshalAsAny(x),
		),
	)
}

// GetExtension reads an extension matching out's type into out, returning true
// if such an extension was found and successfully unmarshalled.
func GetExtension[
	T interface {
		*E
		proto.Message
	},
	E any,
](body *Body) (T, bool, error) {
	var out T

	for _, ext := range body.GetExtensions() {
		if ext.MessageIs(out) {
			out = new(E)
			return out, true, ext.UnmarshalTo(out)
		}
	}

	return nil, false, nil
}
