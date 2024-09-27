package marshaler_test

import (
	"reflect"
	"slices"
	"testing"

	"github.com/dogmatiq/enginekit/marshaler"
	. "github.com/dogmatiq/enginekit/marshaler"
	"github.com/dogmatiq/enginekit/marshaler/codecs/json"
	"github.com/dogmatiq/enginekit/marshaler/codecs/protobuf"
	"github.com/dogmatiq/enginekit/marshaler/internal/stubs1"
	"github.com/dogmatiq/enginekit/marshaler/internal/stubs2"
	"google.golang.org/protobuf/proto"
)

func TestNew(t *testing.T) {
	t.Run("it returns an error if multiple codecs used the same media-type", func(t *testing.T) {
		_, err := New(
			[]reflect.Type{
				reflect.TypeFor[*stubs1.ProtoMessage](),
			},
			[]Codec{
				&json.Codec{},
				&json.Codec{},
			},
		)

		if err == nil {
			t.Fatalf("expected an error")
		}

		got := err.Error()
		want := "multiple codecs use the 'application/json' media-type"

		if got != want {
			t.Fatalf("unexpected error: got %q, want %q", got, want)
		}
	})

	t.Run("it excludes types with conflicting portable type names", func(t *testing.T) {
		marshaler, err := New(
			[]reflect.Type{
				// These two types have conflicting portable type names when
				// using the JSON codec (because it uses the unqualified type
				// name), but not when using the protobuf codec (because it uses
				// the fully-qualified protocol buffers package and message
				// name).
				reflect.TypeFor[*stubs1.ProtoMessage](),
				reflect.TypeFor[*stubs2.ProtoMessage](),
			},
			[]Codec{
				// We prefer the JSON codec, to verify that it is skipped due to
				// the ambiguity in portable type names.
				&json.Codec{},
				&protobuf.DefaultNativeCodec,
			},
		)
		if err != nil {
			t.Fatal(err)
		}

		// Now we encode one of the types with a conflicting name, and verify
		// that the marshaler chooses a codec that is unambiguous, rather than
		// returning an error.
		p, err := marshaler.Marshal(
			&stubs1.ProtoMessage{
				Value: "<value>",
			},
		)
		if err != nil {
			t.Fatal(err)
		}

		got := p.MediaType
		want := "application/vnd.google.protobuf; type=dogmatiq.enginekit.marshaler.stubs1.ProtoMessage"

		if got != want {
			t.Fatalf("unexpected media type: got %q, want %q", got, want)
		}
	})

	t.Run("it returns an error if types with conflicting portable type names are excluded by all codecs", func(t *testing.T) {
		_, err := New(
			[]reflect.Type{
				reflect.TypeFor[*stubs1.ProtoMessage](),
				reflect.TypeFor[*stubs2.ProtoMessage](),
			},
			[]Codec{
				&json.Codec{},
			},
		)
		if err == nil {
			t.Fatalf("expected an error")
		}

		got := err.Error()
		want := []string{
			"naming conflicts occurred within all of the codecs that support the '*stubs1.ProtoMessage' type",
			"naming conflicts occurred within all of the codecs that support the '*stubs2.ProtoMessage' type",
		}

		if !slices.Contains(want, got) {
			t.Fatalf("unexpected error: got %q, want one of %v", got, want)
		}
	})

	t.Run("it returns an error if there are unsupported types", func(t *testing.T) {
		_, err := New(
			[]reflect.Type{
				reflect.TypeFor[*stubs1.ProtoMessage](),
				reflect.TypeFor[int](),
			},
			[]Codec{
				&protobuf.DefaultJSONCodec,
			},
		)
		if err == nil {
			t.Fatalf("expected an error")
		}

		got := err.Error()
		want := "no codecs support the 'int' type"

		if got != want {
			t.Fatalf("unexpected error: got %q, want %q", got, want)
		}
	})
}

func TestMarshaler(t *testing.T) {
	type Value struct {
		Value string `json:"value"`
	}

	marshaler, err := New(
		[]reflect.Type{
			reflect.TypeFor[*stubs1.ProtoMessage](),
			reflect.TypeFor[Value](),
		},
		[]marshaler.Codec{
			&protobuf.DefaultNativeCodec,
			&protobuf.DefaultJSONCodec,
			&json.Codec{},
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("func MarshalType()", func(t *testing.T) {
		t.Run("it returns the portable type name", func(t *testing.T) {
			got, err := marshaler.MarshalType(
				reflect.TypeFor[*stubs1.ProtoMessage](),
			)
			if err != nil {
				t.Fatal(err)
			}

			want := "dogmatiq.enginekit.marshaler.stubs1.ProtoMessage"

			if got != want {
				t.Fatalf("unexpected type name: got %q, want %q", got, want)
			}
		})

		t.Run("it returns an error if the type is not supported", func(t *testing.T) {
			_, err := marshaler.MarshalType(
				reflect.TypeFor[float64](),
			)
			if err == nil {
				t.Fatalf("expected an error")
			}

			got := err.Error()
			want := "no codecs support the 'float64' type"

			if got != want {
				t.Fatalf("unexpected error: got %q, want %q", got, want)
			}
		})
	})

	t.Run("func UnmarshalType()", func(t *testing.T) {
		t.Run("it returns the reflection type", func(t *testing.T) {
			got, err := marshaler.UnmarshalType("dogmatiq.enginekit.marshaler.stubs1.ProtoMessage")
			if err != nil {
				t.Fatal(err)
			}

			want := reflect.TypeFor[*stubs1.ProtoMessage]()

			if got != want {
				t.Fatalf("unexpected type: got %v, want %v", got, want)
			}
		})

		t.Run("it returns an error if the type name is not recognized", func(t *testing.T) {
			_, err := marshaler.UnmarshalType("float64")
			if err == nil {
				t.Fatalf("expected an error")
			}

			got := err.Error()
			want := "the portable type name 'float64' is not recognized"

			if got != want {
				t.Fatalf("unexpected error: got %q, want %q", got, want)
			}
		})
	})

	t.Run("func Marshal()", func(t *testing.T) {
		t.Run("it marshals using the first suitable codec", func(t *testing.T) {
			got, err := marshaler.Marshal(Value{"<value>"})
			if err != nil {
				t.Fatal(err)
			}

			want := Packet{
				MediaType: "application/json; type=Value",
				Data:      []byte(`{"value":"\u003cvalue\u003e"}`),
			}

			if !reflect.DeepEqual(got, want) {
				t.Fatalf("unexpected packet: got %v, want %v", got, want)
			}

			got, err = marshaler.Marshal(&stubs1.ProtoMessage{Value: "<value>"})
			if err != nil {
				t.Fatal(err)
			}

			want = Packet{
				MediaType: "application/vnd.google.protobuf; type=dogmatiq.enginekit.marshaler.stubs1.ProtoMessage",
				Data:      []byte{10, 7, 60, 118, 97, 108, 117, 101, 62},
			}

			if !reflect.DeepEqual(got, want) {
				t.Fatalf("unexpected packet: got %v, want %v", got, want)
			}
		})

		t.Run("returns an error if the codec fails", func(t *testing.T) {
			_, err := marshaler.Marshal(
				&stubs1.ProtoMessage{
					Value: string([]byte{0xfe}),
				},
			)
			if err == nil {
				t.Fatalf("expected an error")
			}
		})

		t.Run("returns an error if the type is not supported", func(t *testing.T) {
			_, err := marshaler.Marshal(123.45)
			if err == nil {
				t.Fatalf("expected an error")
			}

			got := err.Error()
			want := "no codecs support the 'float64' type"

			if got != want {
				t.Fatalf("unexpected error: got %q, want %q", got, want)
			}
		})

	})

	t.Run("func MarshalAs()", func(t *testing.T) {
		t.Run("it marshals using the codec associated with the given media type", func(t *testing.T) {
			cases := []struct {
				Value  any
				Packet Packet
			}{
				{
					Value{"<value>"},
					Packet{
						"application/json; type=Value",
						[]byte(`{"value":"\u003cvalue\u003e"}`),
					},
				},
				{
					&stubs1.ProtoMessage{Value: "<value>"},
					Packet{
						"application/vnd.google.protobuf; type=dogmatiq.enginekit.marshaler.stubs1.ProtoMessage",
						[]byte{10, 7, 60, 118, 97, 108, 117, 101, 62},
					},
				},
				{
					&stubs1.ProtoMessage{Value: "<value>"},
					Packet{
						"application/json; type=ProtoMessage",
						[]byte(`{"value":"\u003cvalue\u003e"}`),
					},
				},
			}

			for _, c := range cases {
				t.Run(c.Packet.MediaType, func(t *testing.T) {
					got, ok, err := marshaler.MarshalAs(
						c.Value,
						[]string{c.Packet.MediaType},
					)
					if err != nil {
						t.Fatal(err)
					}
					if !ok {
						t.Fatalf("expected ok to be true")
					}

					if !reflect.DeepEqual(got, c.Packet) {
						t.Fatalf("unexpected packet: got %v, want %v", got, c.Packet)
					}
				})
			}
		})

		t.Run("it marshals using the codec associated with the highest priority media-type", func(t *testing.T) {
			got, ok, err := marshaler.MarshalAs(
				&stubs1.ProtoMessage{
					Value: "<value>",
				},
				[]string{
					"application/json; type=ProtoMessage",
					"application/vnd.google.protobuf; type=dogmatiq.enginekit.marshaler.stubs1.ProtoMessage",
				},
			)
			if err != nil {
				t.Fatal(err)
			}
			if !ok {
				t.Fatalf("expected ok to be true")
			}

			want := Packet{
				"application/json; type=ProtoMessage",
				[]byte(`{"value":"\u003cvalue\u003e"}`),
			}

			if !reflect.DeepEqual(got, want) {
				t.Fatalf("unexpected packet: got %v, want %v", got, want)
			}
		})

		t.Run("it ignores unsupported media-types", func(t *testing.T) {
			p, ok, err := marshaler.MarshalAs(
				Value{"<value>"},
				[]string{
					"application/binary; type=Value",
					"application/json; type=Value",
				},
			)
			if err != nil {
				t.Fatal(err)
			}
			if !ok {
				t.Fatalf("expected ok to be true")
			}

			got := p.MediaType
			want := "application/json; type=Value"

			if got != want {
				t.Fatalf("unexpected media type: got %q, want %q", got, want)
			}
		})

		t.Run("it returns an error if the media-type is malformed", func(t *testing.T) {
			_, _, err := marshaler.MarshalAs(
				Value{"<value>"},
				[]string{"<malformed>"},
			)
			if err == nil {
				t.Fatalf("expected an error")
			}
		})

		t.Run("it returns an error if the codec fails", func(t *testing.T) {
			_, _, err := marshaler.MarshalAs(
				&stubs1.ProtoMessage{
					Value: string([]byte{0xfe}),
				},
				[]string{
					"application/vnd.google.protobuf; type=dogmatiq.enginekit.marshaler.stubs1.ProtoMessage",
				},
			)
			if err == nil {
				t.Fatalf("expected an error")
			}
		})

		t.Run("it returns false if the media-type is not supported", func(t *testing.T) {
			_, ok, err := marshaler.MarshalAs(
				Value{"<value>"},
				[]string{
					"application/binary; type=Value",
				},
			)
			if err != nil {
				t.Fatal(err)
			}
			if ok {
				t.Fatalf("expected ok to be false")
			}
		})

		t.Run("it returns false if the portable name in the media-type does not match the value's type", func(t *testing.T) {
			_, ok, err := marshaler.MarshalAs(
				Value{"<value>"},
				[]string{
					"application/json; type=MismatchedType",
				},
			)
			if err != nil {
				t.Fatal(err)
			}
			if ok {
				t.Fatalf("expected ok to be false")
			}
		})

		t.Run("it panics if no media-types are provided", func(t *testing.T) {
			defer func() {
				got := recover()
				want := "at least one media-type must be provided"

				if got != want {
					t.Fatalf("unexpected panic: got %q, want %q", got, want)
				}
			}()

			marshaler.MarshalAs(
				&stubs1.ProtoMessage{},
				nil,
			)
		})
	})

	t.Run("func Unmarshal()", func(t *testing.T) {
		t.Run("it unmarshals using the first suitable codec", func(t *testing.T) {
			cases := []struct {
				Packet  Packet
				IsEqual func(any) bool
			}{
				{
					Packet{
						"application/json; type=Value",
						[]byte(`{"value":"\u003cvalue\u003e"}`),
					},
					func(v any) bool {
						return v == Value{"<value>"}
					},
				},
				{
					Packet{
						"application/vnd.google.protobuf; type=dogmatiq.enginekit.marshaler.stubs1.ProtoMessage",
						[]byte{10, 7, 60, 118, 97, 108, 117, 101, 62},
					},
					func(v any) bool {
						m, ok := v.(proto.Message)
						return ok && proto.Equal(m, &stubs1.ProtoMessage{Value: "<value>"})
					},
				},
			}

			for _, c := range cases {
				t.Run(c.Packet.MediaType, func(t *testing.T) {
					got, err := marshaler.Unmarshal(c.Packet)
					if err != nil {
						t.Fatal(err)
					}

					if !c.IsEqual(got) {
						t.Fatalf("unexpected value: got %v", got)
					}
				})
			}
		})

		t.Run("it returns an error if the media-type is not supported", func(t *testing.T) {
			_, err := marshaler.Unmarshal(Packet{
				MediaType: "text/plain; type=MessageA",
			})
			if err == nil {
				t.Fatalf("expected an error")
			}

			got := err.Error()
			want := "no codecs support the 'text/plain' media-type"

			if got != want {
				t.Fatalf("unexpected error: got %q, want %q", got, want)
			}
		})

		t.Run("it returns an error if the media-type is malformed", func(t *testing.T) {
			_, err := marshaler.Unmarshal(Packet{
				MediaType: "<malformed>",
			})
			if err == nil {
				t.Fatalf("expected an error")
			}
		})

		t.Run("it returns an error if the media-type does not specify a type parameter", func(t *testing.T) {
			_, err := marshaler.Unmarshal(Packet{
				MediaType: "application/json",
			})
			if err == nil {
				t.Fatalf("expected an error")
			}

			got := err.Error()
			want := "the media-type 'application/json' does not specify a 'type' parameter"

			if got != want {
				t.Fatalf("unexpected error: got %q, want %q", got, want)
			}
		})

		t.Run("it returns an error if the type is not supported", func(t *testing.T) {
			_, err := marshaler.Unmarshal(Packet{
				MediaType: "application/json; type=MessageC",
			})
			if err == nil {
				t.Fatalf("expected an error")
			}

			got := err.Error()
			want := "the portable type name 'MessageC' is not recognized"

			if got != want {
				t.Fatalf("unexpected error: got %q, want %q", got, want)
			}
		})

		t.Run("it returns an error if the codec fails", func(t *testing.T) {
			_, err := marshaler.Unmarshal(Packet{
				MediaType: "application/json; type=MessageA",
				Data:      []byte("{"),
			})
			if err == nil {
				t.Fatalf("expected an error")
			}
		})
	})

	t.Run("func MediaTypesFor", func(t *testing.T) {
		t.Run("it returns media types in order of codec priority", func(t *testing.T) {
			got := marshaler.MediaTypesFor(reflect.TypeFor[*stubs1.ProtoMessage]())
			want := []string{
				"application/vnd.google.protobuf; type=dogmatiq.enginekit.marshaler.stubs1.ProtoMessage",
				"application/vnd.google.protobuf+json; type=dogmatiq.enginekit.marshaler.stubs1.ProtoMessage",
				"application/json; type=ProtoMessage",
			}

			if !reflect.DeepEqual(got, want) {
				t.Fatalf("unexpected media types: got %v, want %v", got, want)
			}
		})

		t.Run("it returns an empty slice when given an unsupported message type", func(t *testing.T) {
			got := marshaler.MediaTypesFor(reflect.TypeFor[float64]())
			want := []string(nil)

			if !reflect.DeepEqual(got, want) {
				t.Fatalf("unexpected media types: got %v, want %v", got, want)
			}
		})
	})
}
