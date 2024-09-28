package test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"google.golang.org/protobuf/testing/protocmp"
)

// Expect compares two values and fails the test if they are different.
func Expect[T any](
	t *testing.T,
	failMessage string,
	got, want T,
	options ...cmp.Option,
) {
	t.Helper()

	options = append(
		[]cmp.Option{
			protocmp.Transform(),
			cmpopts.EquateEmpty(),
			cmpopts.EquateErrors(),
		},
		options...,
	)

	if diff := cmp.Diff(
		want,
		got,
		options...,
	); diff != "" {
		t.Log(failMessage)
		t.Fatal(diff)
	}
}
