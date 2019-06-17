package fixtures

import (
	"reflect"

	"github.com/dogmatiq/enginekit/marshaling"
	"github.com/dogmatiq/enginekit/marshaling/json"
)

// Marshaler is a marshaler that is aware of the message and aggregate/process
// root fixture types. It uses the JSON codec.
var Marshaler *marshaling.Marshaler

var (
	// MessageA1Packet is a packet containing MessageA1 in JSON format.
	MessageA1Packet marshaling.Packet

	// MessageA2Packet is a packet containing MessageA2 in JSON format.
	MessageA2Packet marshaling.Packet

	// MessageA3Packet is a packet containing MessageA3 in JSON format.
	MessageA3Packet marshaling.Packet
)

var (
	// MessageB1Packet is a packet containing MessageB1 in JSON format.
	MessageB1Packet marshaling.Packet

	// MessageB2Packet is a packet containing MessageB2 in JSON format.
	MessageB2Packet marshaling.Packet

	// MessageB3Packet is a packet containing MessageB3 in JSON format.
	MessageB3Packet marshaling.Packet
)

var (
	// MessageC1Packet is a packet containing MessageC1 in JSON format.
	MessageC1Packet marshaling.Packet

	// MessageC2Packet is a packet containing MessageC2 in JSON format.
	MessageC2Packet marshaling.Packet

	// MessageC3Packet is a packet containing MessageC3 in JSON format.
	MessageC3Packet marshaling.Packet
)

var (
	// MessageD1Packet is a packet containing MessageD1 in JSON format.
	MessageD1Packet marshaling.Packet

	// MessageD2Packet is a packet containing MessageD2 in JSON format.
	MessageD2Packet marshaling.Packet

	// MessageD3Packet is a packet containing MessageD3 in JSON format.
	MessageD3Packet marshaling.Packet
)

var (
	// MessageE1Packet is a packet containing MessageE1 in JSON format.
	MessageE1Packet marshaling.Packet

	// MessageE2Packet is a packet containing MessageE2 in JSON format.
	MessageE2Packet marshaling.Packet

	// MessageE3Packet is a packet containing MessageE3 in JSON format.
	MessageE3Packet marshaling.Packet
)

var (
	// MessageF1Packet is a packet containing MessageF1 in JSON format.
	MessageF1Packet marshaling.Packet

	// MessageF2Packet is a packet containing MessageF2 in JSON format.
	MessageF2Packet marshaling.Packet

	// MessageF3Packet is a packet containing MessageF3 in JSON format.
	MessageF3Packet marshaling.Packet
)

var (
	// MessageG1Packet is a packet containing MessageG1 in JSON format.
	MessageG1Packet marshaling.Packet

	// MessageG2Packet is a packet containing MessageG2 in JSON format.
	MessageG2Packet marshaling.Packet

	// MessageG3Packet is a packet containing MessageG3 in JSON format.
	MessageG3Packet marshaling.Packet
)

var (
	// MessageH1Packet is a packet containing MessageH1 in JSON format.
	MessageH1Packet marshaling.Packet

	// MessageH2Packet is a packet containing MessageH2 in JSON format.
	MessageH2Packet marshaling.Packet

	// MessageH3Packet is a packet containing MessageH3 in JSON format.
	MessageH3Packet marshaling.Packet
)

var (
	// MessageI1Packet is a packet containing MessageI1 in JSON format.
	MessageI1Packet marshaling.Packet

	// MessageI2Packet is a packet containing MessageI2 in JSON format.
	MessageI2Packet marshaling.Packet

	// MessageI3Packet is a packet containing MessageI3 in JSON format.
	MessageI3Packet marshaling.Packet
)

var (
	// MessageJ1Packet is a packet containing MessageJ1 in JSON format.
	MessageJ1Packet marshaling.Packet

	// MessageJ2Packet is a packet containing MessageJ2 in JSON format.
	MessageJ2Packet marshaling.Packet

	// MessageJ3Packet is a packet containing MessageJ3 in JSON format.
	MessageJ3Packet marshaling.Packet
)

var (
	// MessageK1Packet is a packet containing MessageK1 in JSON format.
	MessageK1Packet marshaling.Packet

	// MessageK2Packet is a packet containing MessageK2 in JSON format.
	MessageK2Packet marshaling.Packet

	// MessageK3Packet is a packet containing MessageK3 in JSON format.
	MessageK3Packet marshaling.Packet
)

var (
	// MessageL1Packet is a packet containing MessageL1 in JSON format.
	MessageL1Packet marshaling.Packet

	// MessageL2Packet is a packet containing MessageL2 in JSON format.
	MessageL2Packet marshaling.Packet

	// MessageL3Packet is a packet containing MessageL3 in JSON format.
	MessageL3Packet marshaling.Packet
)

var (
	// MessageM1Packet is a packet containing MessageM1 in JSON format.
	MessageM1Packet marshaling.Packet

	// MessageM2Packet is a packet containing MessageM2 in JSON format.
	MessageM2Packet marshaling.Packet

	// MessageM3Packet is a packet containing MessageM3 in JSON format.
	MessageM3Packet marshaling.Packet
)

var (
	// MessageN1Packet is a packet containing MessageN1 in JSON format.
	MessageN1Packet marshaling.Packet

	// MessageN2Packet is a packet containing MessageN2 in JSON format.
	MessageN2Packet marshaling.Packet

	// MessageN3Packet is a packet containing MessageN3 in JSON format.
	MessageN3Packet marshaling.Packet
)

var (
	// MessageO1Packet is a packet containing MessageO1 in JSON format.
	MessageO1Packet marshaling.Packet

	// MessageO2Packet is a packet containing MessageO2 in JSON format.
	MessageO2Packet marshaling.Packet

	// MessageO3Packet is a packet containing MessageO3 in JSON format.
	MessageO3Packet marshaling.Packet
)

var (
	// MessageP1Packet is a packet containing MessageP1 in JSON format.
	MessageP1Packet marshaling.Packet

	// MessageP2Packet is a packet containing MessageP2 in JSON format.
	MessageP2Packet marshaling.Packet

	// MessageP3Packet is a packet containing MessageP3 in JSON format.
	MessageP3Packet marshaling.Packet
)

var (
	// MessageQ1Packet is a packet containing MessageQ1 in JSON format.
	MessageQ1Packet marshaling.Packet

	// MessageQ2Packet is a packet containing MessageQ2 in JSON format.
	MessageQ2Packet marshaling.Packet

	// MessageQ3Packet is a packet containing MessageQ3 in JSON format.
	MessageQ3Packet marshaling.Packet
)

var (
	// MessageR1Packet is a packet containing MessageR1 in JSON format.
	MessageR1Packet marshaling.Packet

	// MessageR2Packet is a packet containing MessageR2 in JSON format.
	MessageR2Packet marshaling.Packet

	// MessageR3Packet is a packet containing MessageR3 in JSON format.
	MessageR3Packet marshaling.Packet
)

var (
	// MessageS1Packet is a packet containing MessageS1 in JSON format.
	MessageS1Packet marshaling.Packet

	// MessageS2Packet is a packet containing MessageS2 in JSON format.
	MessageS2Packet marshaling.Packet

	// MessageS3Packet is a packet containing MessageS3 in JSON format.
	MessageS3Packet marshaling.Packet
)

var (
	// MessageT1Packet is a packet containing MessageT1 in JSON format.
	MessageT1Packet marshaling.Packet

	// MessageT2Packet is a packet containing MessageT2 in JSON format.
	MessageT2Packet marshaling.Packet

	// MessageT3Packet is a packet containing MessageT3 in JSON format.
	MessageT3Packet marshaling.Packet
)

var (
	// MessageU1Packet is a packet containing MessageU1 in JSON format.
	MessageU1Packet marshaling.Packet

	// MessageU2Packet is a packet containing MessageU2 in JSON format.
	MessageU2Packet marshaling.Packet

	// MessageU3Packet is a packet containing MessageU3 in JSON format.
	MessageU3Packet marshaling.Packet
)

var (
	// MessageV1Packet is a packet containing MessageV1 in JSON format.
	MessageV1Packet marshaling.Packet

	// MessageV2Packet is a packet containing MessageV2 in JSON format.
	MessageV2Packet marshaling.Packet

	// MessageV3Packet is a packet containing MessageV3 in JSON format.
	MessageV3Packet marshaling.Packet
)

var (
	// MessageW1Packet is a packet containing MessageW1 in JSON format.
	MessageW1Packet marshaling.Packet

	// MessageW2Packet is a packet containing MessageW2 in JSON format.
	MessageW2Packet marshaling.Packet

	// MessageW3Packet is a packet containing MessageW3 in JSON format.
	MessageW3Packet marshaling.Packet
)

var (
	// MessageX1Packet is a packet containing MessageX1 in JSON format.
	MessageX1Packet marshaling.Packet

	// MessageX2Packet is a packet containing MessageX2 in JSON format.
	MessageX2Packet marshaling.Packet

	// MessageX3Packet is a packet containing MessageX3 in JSON format.
	MessageX3Packet marshaling.Packet
)

var (
	// MessageY1Packet is a packet containing MessageY1 in JSON format.
	MessageY1Packet marshaling.Packet

	// MessageY2Packet is a packet containing MessageY2 in JSON format.
	MessageY2Packet marshaling.Packet

	// MessageY3Packet is a packet containing MessageY3 in JSON format.
	MessageY3Packet marshaling.Packet
)

var (
	// MessageZ1Packet is a packet containing MessageZ1 in JSON format.
	MessageZ1Packet marshaling.Packet

	// MessageZ2Packet is a packet containing MessageZ2 in JSON format.
	MessageZ2Packet marshaling.Packet

	// MessageZ3Packet is a packet containing MessageZ3 in JSON format.
	MessageZ3Packet marshaling.Packet
)

func init() {
	m, err := marshaling.NewMarshaler(
		[]reflect.Type{
			reflect.TypeOf(&AggregateRoot{}),
			reflect.TypeOf(&ProcessRoot{}),
			MessageAType.ReflectType(),
			MessageBType.ReflectType(),
			MessageCType.ReflectType(),
			MessageDType.ReflectType(),
			MessageEType.ReflectType(),
			MessageFType.ReflectType(),
			MessageGType.ReflectType(),
			MessageHType.ReflectType(),
			MessageIType.ReflectType(),
			MessageJType.ReflectType(),
			MessageKType.ReflectType(),
			MessageLType.ReflectType(),
			MessageMType.ReflectType(),
			MessageNType.ReflectType(),
			MessageOType.ReflectType(),
			MessagePType.ReflectType(),
			MessageQType.ReflectType(),
			MessageRType.ReflectType(),
			MessageSType.ReflectType(),
			MessageTType.ReflectType(),
			MessageUType.ReflectType(),
			MessageVType.ReflectType(),
			MessageWType.ReflectType(),
			MessageXType.ReflectType(),
			MessageYType.ReflectType(),
			MessageZType.ReflectType(),
		},
		[]marshaling.Codec{
			&json.Codec{},
		},
	)
	if err != nil {
		panic(err)
	}

	Marshaler = m

	MessageA1Packet = marshaling.MustMarshalMessage(m, MessageA1)
	MessageA2Packet = marshaling.MustMarshalMessage(m, MessageA2)
	MessageA3Packet = marshaling.MustMarshalMessage(m, MessageA3)

	MessageB1Packet = marshaling.MustMarshalMessage(m, MessageB1)
	MessageB2Packet = marshaling.MustMarshalMessage(m, MessageB2)
	MessageB3Packet = marshaling.MustMarshalMessage(m, MessageB3)

	MessageC1Packet = marshaling.MustMarshalMessage(m, MessageC1)
	MessageC2Packet = marshaling.MustMarshalMessage(m, MessageC2)
	MessageC3Packet = marshaling.MustMarshalMessage(m, MessageC3)

	MessageD1Packet = marshaling.MustMarshalMessage(m, MessageD1)
	MessageD2Packet = marshaling.MustMarshalMessage(m, MessageD2)
	MessageD3Packet = marshaling.MustMarshalMessage(m, MessageD3)

	MessageE1Packet = marshaling.MustMarshalMessage(m, MessageE1)
	MessageE2Packet = marshaling.MustMarshalMessage(m, MessageE2)
	MessageE3Packet = marshaling.MustMarshalMessage(m, MessageE3)

	MessageF1Packet = marshaling.MustMarshalMessage(m, MessageF1)
	MessageF2Packet = marshaling.MustMarshalMessage(m, MessageF2)
	MessageF3Packet = marshaling.MustMarshalMessage(m, MessageF3)

	MessageG1Packet = marshaling.MustMarshalMessage(m, MessageG1)
	MessageG2Packet = marshaling.MustMarshalMessage(m, MessageG2)
	MessageG3Packet = marshaling.MustMarshalMessage(m, MessageG3)

	MessageH1Packet = marshaling.MustMarshalMessage(m, MessageH1)
	MessageH2Packet = marshaling.MustMarshalMessage(m, MessageH2)
	MessageH3Packet = marshaling.MustMarshalMessage(m, MessageH3)

	MessageI1Packet = marshaling.MustMarshalMessage(m, MessageI1)
	MessageI2Packet = marshaling.MustMarshalMessage(m, MessageI2)
	MessageI3Packet = marshaling.MustMarshalMessage(m, MessageI3)

	MessageJ1Packet = marshaling.MustMarshalMessage(m, MessageJ1)
	MessageJ2Packet = marshaling.MustMarshalMessage(m, MessageJ2)
	MessageJ3Packet = marshaling.MustMarshalMessage(m, MessageJ3)

	MessageK1Packet = marshaling.MustMarshalMessage(m, MessageK1)
	MessageK2Packet = marshaling.MustMarshalMessage(m, MessageK2)
	MessageK3Packet = marshaling.MustMarshalMessage(m, MessageK3)

	MessageL1Packet = marshaling.MustMarshalMessage(m, MessageL1)
	MessageL2Packet = marshaling.MustMarshalMessage(m, MessageL2)
	MessageL3Packet = marshaling.MustMarshalMessage(m, MessageL3)

	MessageM1Packet = marshaling.MustMarshalMessage(m, MessageM1)
	MessageM2Packet = marshaling.MustMarshalMessage(m, MessageM2)
	MessageM3Packet = marshaling.MustMarshalMessage(m, MessageM3)

	MessageN1Packet = marshaling.MustMarshalMessage(m, MessageN1)
	MessageN2Packet = marshaling.MustMarshalMessage(m, MessageN2)
	MessageN3Packet = marshaling.MustMarshalMessage(m, MessageN3)

	MessageO1Packet = marshaling.MustMarshalMessage(m, MessageO1)
	MessageO2Packet = marshaling.MustMarshalMessage(m, MessageO2)
	MessageO3Packet = marshaling.MustMarshalMessage(m, MessageO3)

	MessageP1Packet = marshaling.MustMarshalMessage(m, MessageP1)
	MessageP2Packet = marshaling.MustMarshalMessage(m, MessageP2)
	MessageP3Packet = marshaling.MustMarshalMessage(m, MessageP3)

	MessageQ1Packet = marshaling.MustMarshalMessage(m, MessageQ1)
	MessageQ2Packet = marshaling.MustMarshalMessage(m, MessageQ2)
	MessageQ3Packet = marshaling.MustMarshalMessage(m, MessageQ3)

	MessageR1Packet = marshaling.MustMarshalMessage(m, MessageR1)
	MessageR2Packet = marshaling.MustMarshalMessage(m, MessageR2)
	MessageR3Packet = marshaling.MustMarshalMessage(m, MessageR3)

	MessageS1Packet = marshaling.MustMarshalMessage(m, MessageS1)
	MessageS2Packet = marshaling.MustMarshalMessage(m, MessageS2)
	MessageS3Packet = marshaling.MustMarshalMessage(m, MessageS3)

	MessageT1Packet = marshaling.MustMarshalMessage(m, MessageT1)
	MessageT2Packet = marshaling.MustMarshalMessage(m, MessageT2)
	MessageT3Packet = marshaling.MustMarshalMessage(m, MessageT3)

	MessageU1Packet = marshaling.MustMarshalMessage(m, MessageU1)
	MessageU2Packet = marshaling.MustMarshalMessage(m, MessageU2)
	MessageU3Packet = marshaling.MustMarshalMessage(m, MessageU3)

	MessageV1Packet = marshaling.MustMarshalMessage(m, MessageV1)
	MessageV2Packet = marshaling.MustMarshalMessage(m, MessageV2)
	MessageV3Packet = marshaling.MustMarshalMessage(m, MessageV3)

	MessageW1Packet = marshaling.MustMarshalMessage(m, MessageW1)
	MessageW2Packet = marshaling.MustMarshalMessage(m, MessageW2)
	MessageW3Packet = marshaling.MustMarshalMessage(m, MessageW3)

	MessageX1Packet = marshaling.MustMarshalMessage(m, MessageX1)
	MessageX2Packet = marshaling.MustMarshalMessage(m, MessageX2)
	MessageX3Packet = marshaling.MustMarshalMessage(m, MessageX3)

	MessageY1Packet = marshaling.MustMarshalMessage(m, MessageY1)
	MessageY2Packet = marshaling.MustMarshalMessage(m, MessageY2)
	MessageY3Packet = marshaling.MustMarshalMessage(m, MessageY3)

	MessageZ1Packet = marshaling.MustMarshalMessage(m, MessageZ1)
	MessageZ2Packet = marshaling.MustMarshalMessage(m, MessageZ2)
	MessageZ3Packet = marshaling.MustMarshalMessage(m, MessageZ3)
}
