package envelopepb_test

import (
	"testing"
	"time"

	"github.com/dogmatiq/dogma"
	. "github.com/dogmatiq/enginekit/enginetest/stubs"
	. "github.com/dogmatiq/enginekit/internal/test"
	. "github.com/dogmatiq/enginekit/protobuf/envelopepb"
	"github.com/dogmatiq/enginekit/protobuf/identitypb"
	"github.com/dogmatiq/enginekit/protobuf/uuidpb"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func TestPacker_PackAndUnpack(t *testing.T) {
	id := uuidpb.Generate()
	now := time.Now()

	packer := &Packer{
		Site: &identitypb.Identity{
			Name: "site",
			Key:  uuidpb.Generate(),
		},
		Application: &identitypb.Identity{
			Name: "app",
			Key:  uuidpb.Generate(),
		},
		GenerateID: func() *uuidpb.UUID {
			return id
		},
		Now: func() time.Time {
			return now
		},
	}

	got := packer.PackCommand(CommandA1)

	if err := got.Validate(); err != nil {
		t.Fatalf("packer produced an invalid envelope: %v", err)
	}

	want := &Envelope{
		Header: &Header{
			CausationId:   id,
			CorrelationId: id,
			Source: &Source{
				Site:        packer.Site,
				Application: packer.Application,
			},
		},
		Body: &Body{
			MessageId: id,
			CreatedAt: timestamppb.New(now),
			Message: &Message{
				Description: `command(stubs.TypeA:A1, valid)`,
				TypeId:      uuidpb.MustParse(MessageTypeID[*CommandStub[TypeA]]()),
				Data:        []byte(`{"content":"A1"}`),
			},
		},
	}

	Expect(
		t,
		"unexpected envelope",
		got,
		want,
	)

	gotMessage, err := Unpack(want)
	if err != nil {
		t.Fatal(err)
	}

	Expect(
		t,
		"unexpected message",
		gotMessage,
		dogma.Message(CommandA1),
	)

	t.Run("func PackCommand()", func(t *testing.T) {
		t.Run("it panics if passed an unregistered message", func(t *testing.T) {
			ExpectPanic(
				t,
				"*envelopepb_test.T is not a registered message type",
				func() {
					type T struct{ dogma.Command }
					packer.PackCommand(&T{})
				},
			)
		})

		t.Run("it panics if the message cannot be marshaled", func(t *testing.T) {
			ExpectPanic(
				t,
				"unable to marshal *envelopepb_test.T: json: unsupported type: func()",
				func() {
					type T struct{ CommandStub[func()] }
					dogma.RegisterCommand[*T]("622003a4-01a5-4c77-8a4c-cb36b51994e7")
					packer.PackCommand(&T{})
				},
			)
		})

		t.Run("it panics if the envelope is invalid", func(t *testing.T) {
			before := packer.Site
			packer.Site = &identitypb.Identity{} // invalid
			defer func() { packer.Site = before }()

			ExpectPanic(
				t,
				"invalid header: invalid source: invalid site (/00000000-0000-0000-0000-000000000000): invalid name: must be between 1 and 255 bytes",
				func() {
					packer.PackCommand(CommandA1)
				},
			)
		})
	})

	t.Run("func Unpack()", func(t *testing.T) {
		t.Run("it returns an error if the message type is not registered", func(t *testing.T) {
			env := &Envelope{
				Body: &Body{
					Message: &Message{
						Description: "<description>",
						TypeId:      uuidpb.MustParse("f1816a71-3593-4771-8d8b-327650571288"),
						Data:        []byte(`{"content":"A1"}`),
					},
				},
			}

			want := "f1816a71-3593-4771-8d8b-327650571288 is not a registered message type ID"

			if _, err := Unpack(env); err == nil {
				t.Fatal("expected an error")
			} else if err.Error() != want {
				t.Fatalf("unexpected error: got %q, want %q", err, want)
			}
		})

		t.Run("it returns an error if the message cannot be unmarshaled", func(t *testing.T) {
			env := &Envelope{
				Body: &Body{
					Message: &Message{
						Description: "<description>",
						TypeId:      uuidpb.MustParse(MessageTypeID[*CommandStub[TypeA]]()),
						Data:        []byte(`}`),
					},
				},
			}

			want := "unable to unmarshal *stubs.CommandStub[github.com/dogmatiq/enginekit/enginetest/stubs.TypeA]: invalid character '}' looking for beginning of value"

			if _, err := Unpack(env); err == nil {
				t.Fatal("expected an error")
			} else if err.Error() != want {
				t.Fatalf("unexpected error: got %q, want %q", err, want)
			}
		})
	})
}
func TestWithIdempotencyKey(t *testing.T) {
	packer := &Packer{
		Application: &identitypb.Identity{
			Name: "app",
			Key:  uuidpb.Generate(),
		},
	}

	got := packer.PackCommand(CommandA1, WithIdempotencyKey("test-key"))

	if got.Body.IdempotencyKey != "test-key" {
		t.Fatalf("unexpected idempotency key: got %q, want %q", got.Body.IdempotencyKey, "test-key")
	}
}

func TestWithExtension(t *testing.T) {
	t.Run("it panics if x is nil", func(t *testing.T) {
		ExpectPanic(
			t,
			"value must not be nil",
			func() {
				WithExtension(nil)
			},
		)
	})

	t.Run("it panics if x is an empty any", func(t *testing.T) {
		ExpectPanic(
			t,
			"value must not be an empty google.protobuf.Any",
			func() {
				var x *anypb.Any
				WithExtension(x)
			},
		)

		ExpectPanic(
			t,
			"value must not be an empty google.protobuf.Any",
			func() {
				WithExtension(&anypb.Any{})
			},
		)
	})

	t.Run("it accepts empty serialized values", func(t *testing.T) {
		packer := &Packer{
			Application: &identitypb.Identity{
				Name: "app",
				Key:  uuidpb.Generate(),
			},
		}

		got := packer.PackCommand(CommandA1, WithExtension(wrapperspb.String("")))

		if len(got.Body.Extensions) != 1 {
			t.Fatalf("unexpected extension count: got %d, want 1", len(got.Body.Extensions))
		}

		if got.Body.Extensions[0].GetTypeUrl() != "type.googleapis.com/google.protobuf.StringValue" {
			t.Fatalf(
				"unexpected extension type URL: got %q, want %q",
				got.Body.Extensions[0].GetTypeUrl(),
				"type.googleapis.com/google.protobuf.StringValue",
			)
		}

		if len(got.Body.Extensions[0].GetValue()) != 0 {
			t.Fatalf(
				"unexpected extension payload length: got %d, want 0",
				len(got.Body.Extensions[0].GetValue()),
			)
		}
	})

	t.Run("it accepts typed nil messages", func(t *testing.T) {
		packer := &Packer{
			Application: &identitypb.Identity{
				Name: "app",
				Key:  uuidpb.Generate(),
			},
		}

		var value *wrapperspb.StringValue
		got := packer.PackCommand(CommandA1, WithExtension(value))

		if len(got.Body.Extensions) != 1 {
			t.Fatalf("unexpected extension count: got %d, want 1", len(got.Body.Extensions))
		}

		if got.Body.Extensions[0].GetTypeUrl() != "type.googleapis.com/google.protobuf.StringValue" {
			t.Fatalf(
				"unexpected extension type URL: got %q, want %q",
				got.Body.Extensions[0].GetTypeUrl(),
				"type.googleapis.com/google.protobuf.StringValue",
			)
		}

		if len(got.Body.Extensions[0].GetValue()) != 0 {
			t.Fatalf(
				"unexpected extension payload length: got %d, want 0",
				len(got.Body.Extensions[0].GetValue()),
			)
		}
	})

	t.Run("it adds x to the extensions", func(t *testing.T) {
		packer := &Packer{
			Application: &identitypb.Identity{
				Name: "app",
				Key:  uuidpb.Generate(),
			},
		}

		extension := wrapperspb.String("extension")
		want, err := anypb.New(extension)
		if err != nil {
			t.Fatal(err)
		}

		got := packer.PackCommand(CommandA1, WithExtension(extension))

		if got.Header.Extensions != nil {
			t.Fatalf("unexpected header extensions: got %#v, want nil", got.Header.Extensions)
		}

		Expect(
			t,
			"unexpected extensions",
			got.Body.Extensions,
			[]*anypb.Any{want},
		)
	})

	t.Run("it keeps only the last value for a type URL", func(t *testing.T) {
		packer := &Packer{
			Application: &identitypb.Identity{
				Name: "app",
				Key:  uuidpb.Generate(),
			},
		}

		want, err := anypb.New(wrapperspb.String("second"))
		if err != nil {
			t.Fatal(err)
		}

		got := packer.PackCommand(
			CommandA1,
			WithExtension(wrapperspb.String("first")),
			WithExtension(wrapperspb.String("second")),
		)

		Expect(
			t,
			"unexpected extensions",
			got.Body.Extensions,
			[]*anypb.Any{want},
		)
	})

	t.Run("it replaces an existing value in place", func(t *testing.T) {
		packer := &Packer{
			Application: &identitypb.Identity{
				Name: "app",
				Key:  uuidpb.Generate(),
			},
		}

		wantString, err := anypb.New(wrapperspb.String("second"))
		if err != nil {
			t.Fatal(err)
		}

		wantInt, err := anypb.New(wrapperspb.Int64(42))
		if err != nil {
			t.Fatal(err)
		}

		got := packer.PackCommand(
			CommandA1,
			WithExtension(wrapperspb.String("first")),
			WithExtension(wrapperspb.Int64(42)),
			WithExtension(wrapperspb.String("second")),
		)

		Expect(
			t,
			"unexpected extensions",
			got.Body.Extensions,
			[]*anypb.Any{wantString, wantInt},
		)
	})
}

func TestWithBaggage(t *testing.T) {
	t.Run("it panics if x is nil", func(t *testing.T) {
		ExpectPanic(
			t,
			"value must not be nil",
			func() {
				WithBaggage(nil)
			},
		)
	})

	t.Run("it panics if x is an empty any", func(t *testing.T) {
		ExpectPanic(
			t,
			"value must not be an empty google.protobuf.Any",
			func() {
				WithBaggage(&anypb.Any{})
			},
		)
	})

	t.Run("it adds x to the baggage", func(t *testing.T) {
		packer := &Packer{
			Application: &identitypb.Identity{
				Name: "app",
				Key:  uuidpb.Generate(),
			},
		}

		baggage := wrapperspb.String("baggage")
		want, err := anypb.New(baggage)
		if err != nil {
			t.Fatal(err)
		}

		got := packer.PackCommand(CommandA1, WithBaggage(baggage))

		if got.Header.Baggage != nil {
			t.Fatalf("unexpected header baggage: got %#v, want nil", got.Header.Baggage)
		}

		Expect(
			t,
			"unexpected baggage",
			got.Body.Baggage,
			[]*anypb.Any{want},
		)
	})

	t.Run("it keeps only the last value for a type URL", func(t *testing.T) {
		packer := &Packer{
			Application: &identitypb.Identity{
				Name: "app",
				Key:  uuidpb.Generate(),
			},
		}

		want, err := anypb.New(wrapperspb.String("second"))
		if err != nil {
			t.Fatal(err)
		}

		got := packer.PackCommand(
			CommandA1,
			WithBaggage(wrapperspb.String("first")),
			WithBaggage(wrapperspb.String("second")),
		)

		Expect(
			t,
			"unexpected baggage",
			got.Body.Baggage,
			[]*anypb.Any{want},
		)
	})

	t.Run("it replaces an existing value in place", func(t *testing.T) {
		packer := &Packer{
			Application: &identitypb.Identity{
				Name: "app",
				Key:  uuidpb.Generate(),
			},
		}

		wantString, err := anypb.New(wrapperspb.String("second"))
		if err != nil {
			t.Fatal(err)
		}

		wantInt, err := anypb.New(wrapperspb.Int64(42))
		if err != nil {
			t.Fatal(err)
		}

		got := packer.PackCommand(
			CommandA1,
			WithBaggage(wrapperspb.String("first")),
			WithBaggage(wrapperspb.Int64(42)),
			WithBaggage(wrapperspb.String("second")),
		)

		Expect(
			t,
			"unexpected baggage",
			got.Body.Baggage,
			[]*anypb.Any{wantString, wantInt},
		)
	})
}
