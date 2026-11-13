package config_test

import (
	"testing"
	"uuid"

	. "github.com/dogmatiq/enginekit/config"
	"github.com/dogmatiq/protobuf/uuidpb"
)

func TestUnparsedUUID(t *testing.T) {
	t.Run("func String()", func(t *testing.T) {
		t.Run("it returns the string exactly as supplied", func(t *testing.T) {
			for _, want := range []string{
				"F81D4FAE-7DEC-11D0-A765-00A0C91E6BF6",
				"{f81d4fae-7dec-11d0-a765-00a0c91e6bf6}",
			} {
				if got := UnparsedUUID(want).String(); got != want {
					t.Fatalf("unexpected string: got %q, want %q", got, want)
				}
			}
		})
	})

	t.Run("func AsStd()", func(t *testing.T) {
		t.Run("it returns the UUID", func(t *testing.T) {
			s := UnparsedUUID("F81D4FAE-7DEC-11D0-A765-00A0C91E6BF6")

			got := s.AsStd()
			want := uuid.MustParse("f81d4fae-7dec-11d0-a765-00a0c91e6bf6")

			if got != want {
				t.Fatalf("unexpected UUID: got %s, want %s", got, want)
			}
		})

		t.Run("it panics if the string is not a valid UUID", func(t *testing.T) {
			defer expectPanic(t)
			UnparsedUUID("{f81d4fae-7dec-11d0-a765-00a0c91e6bf6}").AsStd()
		})
	})

	t.Run("func TryAsStd()", func(t *testing.T) {
		t.Run("it returns the UUID", func(t *testing.T) {
			s := UnparsedUUID("F81D4FAE-7DEC-11D0-A765-00A0C91E6BF6")

			got, ok := s.TryAsStd()
			want := uuid.MustParse("f81d4fae-7dec-11d0-a765-00a0c91e6bf6")

			if !ok {
				t.Fatal("expected ok to be true")
			}

			if got != want {
				t.Fatalf("unexpected UUID: got %s, want %s", got, want)
			}
		})

		t.Run("it returns false if the string is not a valid UUID", func(t *testing.T) {
			if _, ok := UnparsedUUID("{f81d4fae-7dec-11d0-a765-00a0c91e6bf6}").TryAsStd(); ok {
				t.Fatal("expected ok to be false")
			}
		})
	})

	t.Run("func AsProto()", func(t *testing.T) {
		t.Run("it returns the UUID", func(t *testing.T) {
			s := UnparsedUUID("F81D4FAE-7DEC-11D0-A765-00A0C91E6BF6")

			got := s.AsProto()
			want := uuidpb.MustParse("f81d4fae-7dec-11d0-a765-00a0c91e6bf6")

			if !got.Equal(want) {
				t.Fatalf("unexpected UUID: got %s, want %s", got, want)
			}
		})

		t.Run("it panics if the string is not a valid UUID", func(t *testing.T) {
			defer expectPanic(t)
			UnparsedUUID("{f81d4fae-7dec-11d0-a765-00a0c91e6bf6}").AsProto()
		})
	})

	t.Run("func TryAsProto()", func(t *testing.T) {
		t.Run("it returns the UUID", func(t *testing.T) {
			s := UnparsedUUID("F81D4FAE-7DEC-11D0-A765-00A0C91E6BF6")

			got, ok := s.TryAsProto()
			want := uuidpb.MustParse("f81d4fae-7dec-11d0-a765-00a0c91e6bf6")

			if !ok {
				t.Fatal("expected ok to be true")
			}

			if !got.Equal(want) {
				t.Fatalf("unexpected UUID: got %s, want %s", got, want)
			}
		})

		t.Run("it returns false if the string is not a valid UUID", func(t *testing.T) {
			if _, ok := UnparsedUUID("{f81d4fae-7dec-11d0-a765-00a0c91e6bf6}").TryAsProto(); ok {
				t.Fatal("expected ok to be false")
			}
		})
	})
}

func expectPanic(t *testing.T) {
	t.Helper()

	if recover() == nil {
		t.Fatal("expected a panic")
	}
}
