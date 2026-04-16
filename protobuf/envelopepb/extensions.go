package envelopepb

import (
	anypb "google.golang.org/protobuf/types/known/anypb"
)

// flattenAnyValues merges two slices of [*anypb.Any] values from a [Header] and
// [Body] into a single slice.
//
// Later values override earlier values with the same type URL, so body values
// override header values and the last repeated value in either location wins.
func flattenAnyValues(header, body []*anypb.Any) []*anypb.Any {
	total := len(header) + len(body)
	if total == 0 {
		return nil
	}

	result := make([]*anypb.Any, total)
	start := total

	isSeen := func(v *anypb.Any) bool {
		for _, x := range result[start:] {
			if x.GetTypeUrl() == v.GetTypeUrl() {
				return true
			}
		}

		return false
	}

	appendLastWins := func(values []*anypb.Any) {
		for idx := len(values) - 1; idx >= 0; idx-- {
			v := values[idx]
			if !isSeen(v) {
				start--
				result[start] = v
			}
		}
	}

	appendLastWins(body)
	appendLastWins(header)

	return result[start:]
}
