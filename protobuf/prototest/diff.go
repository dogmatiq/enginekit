package prototest

import (
	"fmt"
	"strings"

	"github.com/google/go-cmp/cmp"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/testing/protocmp"
)

// T is the interface that testing.T implements.
type T interface {
	Helper()
	Errorf(format string, args ...any)
}

// Equal fails t if got and want are not equal according to [proto.Equal].
//
// It uses [protocmp.Transform] to produce field-by-field comparison, providing
// meaningful error messages that show exactly which fields differ between the
// two messages.
func Equal[M proto.Message](
	t T,
	got, want M,
) {
	t.Helper()

	if proto.Equal(got, want) {
		return
	}

	diff := cmp.Diff(want, got, protocmp.Transform())

	w := &strings.Builder{}

	fmt.Fprintln(w, "protocol buffer messages are not equal")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "--- want ---")
	fmt.Fprintf(w, "%v\n", want)
	fmt.Fprintln(w)
	fmt.Fprintln(w, "--- got ---")
	fmt.Fprintf(w, "%v\n", got)
	fmt.Fprintln(w)
	fmt.Fprintln(w, "--- diff (-want +got) ---")
	fmt.Fprint(w, diff)

	t.Errorf("%s", w.String())
}
