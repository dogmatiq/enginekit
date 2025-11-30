package prototest_test

import (
	"fmt"
	"strings"
	"testing"

	. "github.com/dogmatiq/enginekit/protobuf/prototest"
	"github.com/dogmatiq/enginekit/protobuf/uuidpb"
)

func TestEqual_whenMessagesAreEqual(t *testing.T) {
	t.Parallel()

	mt := &mockT{}

	got := &uuidpb.UUID{Upper: 0x123, Lower: 0x456}
	want := &uuidpb.UUID{Upper: 0x123, Lower: 0x456}

	Equal(mt, got, want)

	if mt.failed {
		t.Errorf("expected test to pass, but Errorf was called: %s", mt.message)
	}
}

func TestEqual_whenMessagesAreDifferent(t *testing.T) {
	t.Parallel()

	mt := &mockT{}

	got := &uuidpb.UUID{Upper: 0x123, Lower: 0x456}
	want := &uuidpb.UUID{Upper: 0x789, Lower: 0xabc}

	Equal(mt, got, want)

	if !mt.failed {
		t.Error("expected test to fail, but Errorf was not called")
	}

	if mt.message == "" {
		t.Error("expected a non-empty error message")
	}

	// Check that the message contains useful field-level information
	if !strings.Contains(mt.message, "upper") || !strings.Contains(mt.message, "lower") {
		t.Errorf("expected error message to contain field names, got: %s", mt.message)
	}
}

func TestEqual_whenMessagesHavePartialDifferences(t *testing.T) {
	t.Parallel()

	mt := &mockT{}

	got := &uuidpb.UUID{Upper: 0x123, Lower: 0x456}
	want := &uuidpb.UUID{Upper: 0x123, Lower: 0xabc}

	Equal(mt, got, want)

	if !mt.failed {
		t.Error("expected test to fail, but Errorf was not called")
	}

	// Check that the message contains the diff
	if !strings.Contains(mt.message, "lower") {
		t.Errorf("expected error message to contain 'lower' field difference, got: %s", mt.message)
	}
}

// mockT implements the T interface for testing purposes.
type mockT struct {
	failed  bool
	message string
}

func (m *mockT) Helper() {}

func (m *mockT) Errorf(format string, args ...any) {
	m.failed = true
	m.message = fmt.Sprintf(format, args...)
}
