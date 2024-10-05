package message_test

import (
	"errors"
	"testing"

	"github.com/dogmatiq/dogma"
	. "github.com/dogmatiq/enginekit/enginetest/stubs"
	. "github.com/dogmatiq/enginekit/message"
)

func TestSwitchByKind(t *testing.T) {
	cases := []struct {
		Kind Kind
		Want string
	}{
		{CommandKind, "command"},
		{EventKind, "event"},
		{TimeoutKind, "timeout"},
	}

	for _, c := range cases {
		var result string

		SwitchByKind(
			c.Kind,
			func() { result = "command" },
			func() { result = "event" },
			func() { result = "timeout" },
		)

		if result != c.Want {
			t.Fatalf("unexpected result: got %q, want %q", result, c.Want)
		}
	}

	t.Run("it panics when the associated function is nil", func(t *testing.T) {
		cases := []struct {
			Kind Kind
			Want string
		}{
			{CommandKind, `no case function was provided for the "command" kind`},
			{EventKind, `no case function was provided for the "event" kind`},
			{TimeoutKind, `no case function was provided for the "timeout" kind`},
		}

		for _, c := range cases {
			func() {
				defer func() {
					if got := recover(); got != c.Want {
						t.Fatalf("unexpected panic: got %q, want %q", got, c.Want)
					}
				}()

				SwitchByKind(c.Kind, nil, nil, nil)
			}()
		}
	})

	t.Run("it panics when the kind is invalid", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Fatal("expected a panic")
			}
		}()

		SwitchByKind(Kind(-1), nil, nil, nil)
	})
}

func TestMapByKind(t *testing.T) {
	cases := []struct {
		Kind Kind
		Want string
	}{
		{CommandKind, "command"},
		{EventKind, "event"},
		{TimeoutKind, "timeout"},
	}

	for _, c := range cases {
		got := MapByKind(
			c.Kind,
			"command",
			"event",
			"timeout",
		)

		if got != c.Want {
			t.Fatalf("unexpected result: got %q, want %q", got, c.Want)
		}
	}

	t.Run("it panics when the kind is invalid", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Fatal("expected a panic")
			}
		}()

		MapByKind(Kind(-1), "command", "event", "timeout")
	})
}

func TestSwitchByKindOf(t *testing.T) {
	t.Run("when the message is a command", func(t *testing.T) {
		var message dogma.Command

		SwitchByKindOf(
			CommandA1,
			func(m dogma.Command) { message = m },
			func(m dogma.Event) { t.Fatal("unexpected call to event case") },
			func(m dogma.Timeout) { t.Fatal("unexpected call to timeout case") },
		)

		if message != CommandA1 {
			t.Fatal("command case was not called with the expected message")
		}
	})

	t.Run("when the message is an event", func(t *testing.T) {
		var message dogma.Event

		SwitchByKindOf(
			EventA1,
			func(m dogma.Command) { t.Fatal("unexpected call to command case") },
			func(m dogma.Event) { message = m },
			func(m dogma.Timeout) { t.Fatal("unexpected call to timeout case") },
		)

		if message != EventA1 {
			t.Fatal("event case was not called with the expected message")
		}
	})

	t.Run("when the message is a timeout", func(t *testing.T) {
		var message dogma.Timeout

		SwitchByKindOf(
			TimeoutA1,
			func(m dogma.Command) { t.Fatal("unexpected call to command case") },
			func(m dogma.Event) { t.Fatal("unexpected call to event case") },
			func(m dogma.Timeout) { message = m },
		)

		if message != TimeoutA1 {
			t.Fatal("timeout case was not called with the expected message")
		}
	})

	t.Run("it panics when the associated function is nil", func(t *testing.T) {
		cases := []struct {
			Message dogma.Message
			Want    string
		}{
			{CommandA1, `no case function was provided for dogma.Command messages`},
			{EventA1, `no case function was provided for dogma.Event messages`},
			{TimeoutA1, `no case function was provided for dogma.Timeout messages`},
		}

		for _, c := range cases {
			func() {
				defer func() {
					if got := recover(); got != c.Want {
						t.Fatalf("unexpected panic: got %q, want %q", got, c.Want)
					}
				}()

				SwitchByKindOf(c.Message, nil, nil, nil)
			}()
		}
	})

	t.Run("it panics if the message does not implement any of the more specific interfaces", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Fatal("expected a panic")
			}
		}()

		SwitchByKindOf(
			nil,
			func(m dogma.Command) { t.Fatal("unexpected call to command case") },
			func(m dogma.Event) { t.Fatal("unexpected call to event case") },
			func(m dogma.Timeout) { t.Fatal("unexpected call to timeout case") },
		)
	})
}

func TestMapByKindOf(t *testing.T) {
	cases := []struct {
		Message dogma.Message
		Want    string
	}{
		{CommandA1, "command"},
		{EventA1, "event"},
		{TimeoutA1, "timeout"},
	}

	for _, c := range cases {
		got := MapByKindOf(
			c.Message,
			func(m dogma.Command) string { return "command" },
			func(m dogma.Event) string { return "event" },
			func(m dogma.Timeout) string { return "timeout" },
		)

		if got != c.Want {
			t.Fatalf("unexpected result: got %q, want %q", got, c.Want)
		}
	}

	t.Run("it panics when the associated function is nil", func(t *testing.T) {
		cases := []struct {
			Message dogma.Message
			Want    string
		}{
			{CommandA1, `no case function was provided for dogma.Command messages`},
			{EventA1, `no case function was provided for dogma.Event messages`},
			{TimeoutA1, `no case function was provided for dogma.Timeout messages`},
		}

		for _, c := range cases {
			func() {
				defer func() {
					if got := recover(); got != c.Want {
						t.Fatalf("unexpected panic: got %q, want %q", got, c.Want)
					}
				}()

				MapByKindOf[int](c.Message, nil, nil, nil)
			}()
		}
	})

	t.Run("it panics if the message does not implement any of the more specific interfaces", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Fatal("expected a panic")
			}
		}()

		MapByKindOf(
			nil,
			func(m dogma.Command) string { t.Fatal("unexpected call to command case"); return "" },
			func(m dogma.Event) string { t.Fatal("unexpected call to event case"); return "" },
			func(m dogma.Timeout) string { t.Fatal("unexpected call to timeout case"); return "" },
		)
	})
}

func TestMapByKindOfWithErr(t *testing.T) {
	cases := []struct {
		Message dogma.Message
		Want    string
	}{
		{CommandA1, "command"},
		{EventA1, "event"},
		{TimeoutA1, "timeout"},
	}

	for _, c := range cases {
		got, gotErr := MapByKindOfWithErr(
			c.Message,
			func(m dogma.Command) (string, error) { return "command", errors.New("command") },
			func(m dogma.Event) (string, error) { return "event", errors.New("event") },
			func(m dogma.Timeout) (string, error) { return "timeout", errors.New("timeout") },
		)

		if got != c.Want {
			t.Fatalf("unexpected result: got %q, want %q", got, c.Want)
		}

		if gotErr == nil {
			t.Fatal("expected an error")
		}

		if gotErr.Error() != c.Want {
			t.Fatalf("unexpected error: got %q, want %q", gotErr, c.Want)
		}
	}

	t.Run("it panics when the associated function is nil", func(t *testing.T) {
		cases := []struct {
			Message dogma.Message
			Want    string
		}{
			{CommandA1, `no case function was provided for dogma.Command messages`},
			{EventA1, `no case function was provided for dogma.Event messages`},
			{TimeoutA1, `no case function was provided for dogma.Timeout messages`},
		}

		for _, c := range cases {
			func() {
				defer func() {
					if got := recover(); got != c.Want {
						t.Fatalf("unexpected panic: got %q, want %q", got, c.Want)
					}
				}()

				MapByKindOfWithErr[int](c.Message, nil, nil, nil)
			}()
		}
	})

	t.Run("it panics if the message does not implement any of the more specific interfaces", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Fatal("expected a panic")
			}
		}()

		MapByKindOfWithErr(
			nil,
			func(m dogma.Command) (string, error) { t.Fatal("unexpected call to command case"); return "", nil },
			func(m dogma.Event) (string, error) { t.Fatal("unexpected call to event case"); return "", nil },
			func(m dogma.Timeout) (string, error) { t.Fatal("unexpected call to timeout case"); return "", nil },
		)
	})
}
