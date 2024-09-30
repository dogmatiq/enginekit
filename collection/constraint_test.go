package collection_test

import "cmp"

type elem int

func (e elem) Compare(other elem) int {
	return cmp.Compare(e, other)
}
