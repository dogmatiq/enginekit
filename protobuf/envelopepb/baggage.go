package envelopepb

import "google.golang.org/protobuf/proto"

// SetBaggage sets x as a baggage value on b. If a baggage value with the same
// type URL is already present, it is replaced.
//
// It panics if x is nil.
func SetBaggage[
	T interface {
		*E
		proto.Message
	},
	E any,
](body *Body, x proto.Message) {
	body.SetBaggage(
		appendOrReplace(
			body.GetBaggage(),
			marshalAsAny(x),
		),
	)
}

// GetBaggage reads a baggage value matching out's type into out, returning true
// if such a value was found and successfully unmarshalled.
func GetBaggage[
	T interface {
		*E
		proto.Message
	},
	E any,
](body *Body) (T, bool, error) {
	var out T

	for _, bag := range body.GetBaggage() {
		if bag.MessageIs(out) {
			out = new(E)
			return out, true, bag.UnmarshalTo(out)
		}
	}

	return nil, false, nil
}
