package stubs

import (
	"reflect"

	"github.com/dogmatiq/enginekit/marshaler"
	"github.com/dogmatiq/enginekit/marshaler/codecs/json"
)

// Marshaler is a [marshaler.Marshaler] that can marshal and unmarshal all of
// the data types in this package.
var Marshaler marshaler.Marshaler

func init() {
	var err error
	Marshaler, err = marshaler.New(
		[]reflect.Type{
			reflect.TypeFor[*AggregateRootStub](),
			reflect.TypeFor[*ProcessRootStub](),

			reflect.TypeFor[CommandStub[TypeA]](),
			reflect.TypeFor[CommandStub[TypeB]](),
			reflect.TypeFor[CommandStub[TypeC]](),
			reflect.TypeFor[CommandStub[TypeD]](),
			reflect.TypeFor[CommandStub[TypeE]](),
			reflect.TypeFor[CommandStub[TypeF]](),
			reflect.TypeFor[CommandStub[TypeG]](),
			reflect.TypeFor[CommandStub[TypeH]](),
			reflect.TypeFor[CommandStub[TypeI]](),
			reflect.TypeFor[CommandStub[TypeJ]](),
			reflect.TypeFor[CommandStub[TypeK]](),
			reflect.TypeFor[CommandStub[TypeL]](),
			reflect.TypeFor[CommandStub[TypeM]](),
			reflect.TypeFor[CommandStub[TypeN]](),
			reflect.TypeFor[CommandStub[TypeO]](),
			reflect.TypeFor[CommandStub[TypeP]](),
			reflect.TypeFor[CommandStub[TypeQ]](),
			reflect.TypeFor[CommandStub[TypeR]](),
			reflect.TypeFor[CommandStub[TypeS]](),
			reflect.TypeFor[CommandStub[TypeT]](),
			reflect.TypeFor[CommandStub[TypeU]](),
			reflect.TypeFor[CommandStub[TypeV]](),
			reflect.TypeFor[CommandStub[TypeW]](),
			reflect.TypeFor[CommandStub[TypeX]](),
			reflect.TypeFor[CommandStub[TypeY]](),
			reflect.TypeFor[CommandStub[TypeZ]](),

			reflect.TypeFor[EventStub[TypeA]](),
			reflect.TypeFor[EventStub[TypeB]](),
			reflect.TypeFor[EventStub[TypeC]](),
			reflect.TypeFor[EventStub[TypeD]](),
			reflect.TypeFor[EventStub[TypeE]](),
			reflect.TypeFor[EventStub[TypeF]](),
			reflect.TypeFor[EventStub[TypeG]](),
			reflect.TypeFor[EventStub[TypeH]](),
			reflect.TypeFor[EventStub[TypeI]](),
			reflect.TypeFor[EventStub[TypeJ]](),
			reflect.TypeFor[EventStub[TypeK]](),
			reflect.TypeFor[EventStub[TypeL]](),
			reflect.TypeFor[EventStub[TypeM]](),
			reflect.TypeFor[EventStub[TypeN]](),
			reflect.TypeFor[EventStub[TypeO]](),
			reflect.TypeFor[EventStub[TypeP]](),
			reflect.TypeFor[EventStub[TypeQ]](),
			reflect.TypeFor[EventStub[TypeR]](),
			reflect.TypeFor[EventStub[TypeS]](),
			reflect.TypeFor[EventStub[TypeT]](),
			reflect.TypeFor[EventStub[TypeU]](),
			reflect.TypeFor[EventStub[TypeV]](),
			reflect.TypeFor[EventStub[TypeW]](),
			reflect.TypeFor[EventStub[TypeX]](),
			reflect.TypeFor[EventStub[TypeY]](),
			reflect.TypeFor[EventStub[TypeZ]](),

			reflect.TypeFor[TimeoutStub[TypeA]](),
			reflect.TypeFor[TimeoutStub[TypeB]](),
			reflect.TypeFor[TimeoutStub[TypeC]](),
			reflect.TypeFor[TimeoutStub[TypeD]](),
			reflect.TypeFor[TimeoutStub[TypeE]](),
			reflect.TypeFor[TimeoutStub[TypeF]](),
			reflect.TypeFor[TimeoutStub[TypeG]](),
			reflect.TypeFor[TimeoutStub[TypeH]](),
			reflect.TypeFor[TimeoutStub[TypeI]](),
			reflect.TypeFor[TimeoutStub[TypeJ]](),
			reflect.TypeFor[TimeoutStub[TypeK]](),
			reflect.TypeFor[TimeoutStub[TypeL]](),
			reflect.TypeFor[TimeoutStub[TypeM]](),
			reflect.TypeFor[TimeoutStub[TypeN]](),
			reflect.TypeFor[TimeoutStub[TypeO]](),
			reflect.TypeFor[TimeoutStub[TypeP]](),
			reflect.TypeFor[TimeoutStub[TypeQ]](),
			reflect.TypeFor[TimeoutStub[TypeR]](),
			reflect.TypeFor[TimeoutStub[TypeS]](),
			reflect.TypeFor[TimeoutStub[TypeT]](),
			reflect.TypeFor[TimeoutStub[TypeU]](),
			reflect.TypeFor[TimeoutStub[TypeV]](),
			reflect.TypeFor[TimeoutStub[TypeW]](),
			reflect.TypeFor[TimeoutStub[TypeX]](),
			reflect.TypeFor[TimeoutStub[TypeY]](),
			reflect.TypeFor[TimeoutStub[TypeZ]](),
		},
		[]marshaler.Codec{
			&json.Codec{},
		},
	)
	if err != nil {
		panic(err)
	}
}
