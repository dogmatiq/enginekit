package stubs

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"regexp"

	"github.com/dogmatiq/dogma"
)

// CommandStub is a test implementation of [dogma.Command].
type CommandStub[T any] struct {
	Content         T      `json:"content,omitempty"`
	ValidationError string `json:"validation_error,omitempty"`
}

// MessageDescription returns a description of the command.
func (s *CommandStub[T]) MessageDescription() string {
	validity := "valid"
	if s.ValidationError != "" {
		validity = "invalid: " + s.ValidationError
	}
	return fmt.Sprintf(
		"command(%T:%v, %s)",
		s.Content,
		s.Content,
		validity,
	)
}

// Validate returns a non-nil error if c.Invalid is not empty.
func (s *CommandStub[T]) Validate(dogma.CommandValidationScope) error {
	if s.ValidationError != "" {
		return errors.New(s.ValidationError)
	}
	return nil
}

// MarshalBinary implements [encoding.BinaryMarshaler].
func (s *CommandStub[T]) MarshalBinary() ([]byte, error) {
	return json.Marshal(s)
}

// UnmarshalBinary implements [encoding.BinaryUnmarshaler].
func (s *CommandStub[T]) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, s)
}

// EventStub is a test implementation of [dogma.Event].
type EventStub[T any] struct {
	Content         T      `json:"content,omitempty"`
	ValidationError string `json:"validation_error,omitempty"`
}

// MessageDescription returns a description of the command.
func (s *EventStub[T]) MessageDescription() string {
	validity := "valid"
	if s.ValidationError != "" {
		validity = "invalid: " + s.ValidationError
	}
	return fmt.Sprintf(
		"event(%T:%v, %s)",
		s.Content,
		s.Content,
		validity,
	)
}

// Validate returns a non-nil error if c.Invalid is not empty.
func (s *EventStub[T]) Validate(dogma.EventValidationScope) error {
	if s.ValidationError != "" {
		return errors.New(s.ValidationError)
	}
	return nil
}

// MarshalBinary implements [encoding.BinaryMarshaler].
func (s *EventStub[T]) MarshalBinary() ([]byte, error) {
	return json.Marshal(s)
}

// UnmarshalBinary implements [encoding.BinaryUnmarshaler].
func (s *EventStub[T]) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, s)
}

// TimeoutStub is a test implementation of [dogma.Test].
type TimeoutStub[T any] struct {
	Content         T      `json:"content,omitempty"`
	ValidationError string `json:"validation_error,omitempty"`
}

// MessageDescription returns a description of the command.
func (s *TimeoutStub[T]) MessageDescription() string {
	validity := "valid"
	if s.ValidationError != "" {
		validity = "invalid: " + s.ValidationError
	}
	return fmt.Sprintf(
		"timeout(%T:%v, %s)",
		s.Content,
		s.Content,
		validity,
	)
}

// Validate returns a non-nil error if c.Invalid is not empty.
func (s *TimeoutStub[T]) Validate(dogma.TimeoutValidationScope) error {
	if s.ValidationError != "" {
		return errors.New(s.ValidationError)
	}
	return nil
}

// MarshalBinary implements [encoding.BinaryMarshaler].
func (s *TimeoutStub[T]) MarshalBinary() ([]byte, error) {
	return json.Marshal(s)
}

// UnmarshalBinary implements [encoding.BinaryUnmarshaler].
func (s *TimeoutStub[T]) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, s)
}

type (
	// TypeA is a named type used as a type parameter for [CommandStub],
	// [EventStub] and [TimeoutStub] to provide a unique type.
	TypeA string
	// TypeB is a named type used as a type parameter for [CommandStub],
	// [EventStub] and [TimeoutStub] to provide a unique type.
	TypeB string
	// TypeC is a named type used as a type parameter for [CommandStub],
	// [EventStub] and [TimeoutStub] to provide a unique type.
	TypeC string
	// TypeD is a named type used as a type parameter for [CommandStub],
	// [EventStub] and [TimeoutStub] to provide a unique type.
	TypeD string
	// TypeE is a named type used as a type parameter for [CommandStub],
	// [EventStub] and [TimeoutStub] to provide a unique type.
	TypeE string
	// TypeF is a named type used as a type parameter for [CommandStub],
	// [EventStub] and [TimeoutStub] to provide a unique type.
	TypeF string
	// TypeG is a named type used as a type parameter for [CommandStub],
	// [EventStub] and [TimeoutStub] to provide a unique type.
	TypeG string
	// TypeH is a named type used as a type parameter for [CommandStub],
	// [EventStub] and [TimeoutStub] to provide a unique type.
	TypeH string
	// TypeI is a named type used as a type parameter for [CommandStub],
	// [EventStub] and [TimeoutStub] to provide a unique type.
	TypeI string
	// TypeJ is a named type used as a type parameter for [CommandStub],
	// [EventStub] and [TimeoutStub] to provide a unique type.
	TypeJ string
	// TypeK is a named type used as a type parameter for [CommandStub],
	// [EventStub] and [TimeoutStub] to provide a unique type.
	TypeK string
	// TypeL is a named type used as a type parameter for [CommandStub],
	// [EventStub] and [TimeoutStub] to provide a unique type.
	TypeL string
	// TypeM is a named type used as a type parameter for [CommandStub],
	// [EventStub] and [TimeoutStub] to provide a unique type.
	TypeM string
	// TypeN is a named type used as a type parameter for [CommandStub],
	// [EventStub] and [TimeoutStub] to provide a unique type.
	TypeN string
	// TypeO is a named type used as a type parameter for [CommandStub],
	// [EventStub] and [TimeoutStub] to provide a unique type.
	TypeO string
	// TypeP is a named type used as a type parameter for [CommandStub],
	// [EventStub] and [TimeoutStub] to provide a unique type.
	TypeP string
	// TypeQ is a named type used as a type parameter for [CommandStub],
	// [EventStub] and [TimeoutStub] to provide a unique type.
	TypeQ string
	// TypeR is a named type used as a type parameter for [CommandStub],
	// [EventStub] and [TimeoutStub] to provide a unique type.
	TypeR string
	// TypeS is a named type used as a type parameter for [CommandStub],
	// [EventStub] and [TimeoutStub] to provide a unique type.
	TypeS string
	// TypeT is a named type used as a type parameter for [CommandStub],
	// [EventStub] and [TimeoutStub] to provide a unique type.
	TypeT string
	// TypeU is a named type used as a type parameter for [CommandStub],
	// [EventStub] and [TimeoutStub] to provide a unique type.
	TypeU string
	// TypeV is a named type used as a type parameter for [CommandStub],
	// [EventStub] and [TimeoutStub] to provide a unique type.
	TypeV string
	// TypeW is a named type used as a type parameter for [CommandStub],
	// [EventStub] and [TimeoutStub] to provide a unique type.
	TypeW string
	// TypeX is a named type used as a type parameter for [CommandStub],
	// [EventStub] and [TimeoutStub] to provide a unique type.
	TypeX string
	// TypeY is a named type used as a type parameter for [CommandStub],
	// [EventStub] and [TimeoutStub] to provide a unique type.
	TypeY string
	// TypeZ is a named type used as a type parameter for [CommandStub],
	// [EventStub] and [TimeoutStub] to provide a unique type.
	TypeZ string
)

var (
	// CommandA1 is command of type [CommandA] with content "A1".
	CommandA1 = &CommandStub[TypeA]{Content: "A1"}
	// CommandA2 is command of type [CommandA] with content "A2".
	CommandA2 = &CommandStub[TypeA]{Content: "A2"}
	// CommandA3 is command of type [CommandA] with content "A3".
	CommandA3 = &CommandStub[TypeA]{Content: "A3"}

	// CommandB1 is command of type [CommandB] with content "B1".
	CommandB1 = &CommandStub[TypeB]{Content: "B1"}
	// CommandB2 is command of type [CommandB] with content "B2".
	CommandB2 = &CommandStub[TypeB]{Content: "B2"}
	// CommandB3 is command of type [CommandB] with content "B3".
	CommandB3 = &CommandStub[TypeB]{Content: "B3"}

	// CommandC1 is command of type [CommandC] with content "C1".
	CommandC1 = &CommandStub[TypeC]{Content: "C1"}
	// CommandC2 is command of type [CommandC] with content "C2".
	CommandC2 = &CommandStub[TypeC]{Content: "C2"}
	// CommandC3 is command of type [CommandC] with content "C3".
	CommandC3 = &CommandStub[TypeC]{Content: "C3"}

	// CommandD1 is command of type [CommandD] with content "D1".
	CommandD1 = &CommandStub[TypeD]{Content: "D1"}
	// CommandD2 is command of type [CommandD] with content "D2".
	CommandD2 = &CommandStub[TypeD]{Content: "D2"}
	// CommandD3 is command of type [CommandD] with content "D3".
	CommandD3 = &CommandStub[TypeD]{Content: "D3"}

	// CommandE1 is command of type [CommandE] with content "E1".
	CommandE1 = &CommandStub[TypeE]{Content: "E1"}
	// CommandE2 is command of type [CommandE] with content "E2".
	CommandE2 = &CommandStub[TypeE]{Content: "E2"}
	// CommandE3 is command of type [CommandE] with content "E3".
	CommandE3 = &CommandStub[TypeE]{Content: "E3"}

	// CommandF1 is command of type [CommandF] with content "F1".
	CommandF1 = &CommandStub[TypeF]{Content: "F1"}
	// CommandF2 is command of type [CommandF] with content "F2".
	CommandF2 = &CommandStub[TypeF]{Content: "F2"}
	// CommandF3 is command of type [CommandF] with content "F3".
	CommandF3 = &CommandStub[TypeF]{Content: "F3"}

	// CommandG1 is command of type [CommandG] with content "G1".
	CommandG1 = &CommandStub[TypeG]{Content: "G1"}
	// CommandG2 is command of type [CommandG] with content "G2".
	CommandG2 = &CommandStub[TypeG]{Content: "G2"}
	// CommandG3 is command of type [CommandG] with content "G3".
	CommandG3 = &CommandStub[TypeG]{Content: "G3"}

	// CommandH1 is command of type [CommandH] with content "H1".
	CommandH1 = &CommandStub[TypeH]{Content: "H1"}
	// CommandH2 is command of type [CommandH] with content "H2".
	CommandH2 = &CommandStub[TypeH]{Content: "H2"}
	// CommandH3 is command of type [CommandH] with content "H3".
	CommandH3 = &CommandStub[TypeH]{Content: "H3"}

	// CommandI1 is command of type [CommandI] with content "I1".
	CommandI1 = &CommandStub[TypeI]{Content: "I1"}
	// CommandI2 is command of type [CommandI] with content "I2".
	CommandI2 = &CommandStub[TypeI]{Content: "I2"}
	// CommandI3 is command of type [CommandI] with content "I3".
	CommandI3 = &CommandStub[TypeI]{Content: "I3"}

	// CommandJ1 is command of type [CommandJ] with content "J1".
	CommandJ1 = &CommandStub[TypeJ]{Content: "J1"}
	// CommandJ2 is command of type [CommandJ] with content "J2".
	CommandJ2 = &CommandStub[TypeJ]{Content: "J2"}
	// CommandJ3 is command of type [CommandJ] with content "J3".
	CommandJ3 = &CommandStub[TypeJ]{Content: "J3"}

	// CommandK1 is command of type [CommandK] with content "K1".
	CommandK1 = &CommandStub[TypeK]{Content: "K1"}
	// CommandK2 is command of type [CommandK] with content "K2".
	CommandK2 = &CommandStub[TypeK]{Content: "K2"}
	// CommandK3 is command of type [CommandK] with content "K3".
	CommandK3 = &CommandStub[TypeK]{Content: "K3"}

	// CommandL1 is command of type [CommandL] with content "L1".
	CommandL1 = &CommandStub[TypeL]{Content: "L1"}
	// CommandL2 is command of type [CommandL] with content "L2".
	CommandL2 = &CommandStub[TypeL]{Content: "L2"}
	// CommandL3 is command of type [CommandL] with content "L3".
	CommandL3 = &CommandStub[TypeL]{Content: "L3"}

	// CommandM1 is command of type [CommandM] with content "M1".
	CommandM1 = &CommandStub[TypeM]{Content: "M1"}
	// CommandM2 is command of type [CommandM] with content "M2".
	CommandM2 = &CommandStub[TypeM]{Content: "M2"}
	// CommandM3 is command of type [CommandM] with content "M3".
	CommandM3 = &CommandStub[TypeM]{Content: "M3"}

	// CommandN1 is command of type [CommandN] with content "N1".
	CommandN1 = &CommandStub[TypeN]{Content: "N1"}
	// CommandN2 is command of type [CommandN] with content "N2".
	CommandN2 = &CommandStub[TypeN]{Content: "N2"}
	// CommandN3 is command of type [CommandN] with content "N3".
	CommandN3 = &CommandStub[TypeN]{Content: "N3"}

	// CommandO1 is command of type [CommandO] with content "O1".
	CommandO1 = &CommandStub[TypeO]{Content: "O1"}
	// CommandO2 is command of type [CommandO] with content "O2".
	CommandO2 = &CommandStub[TypeO]{Content: "O2"}
	// CommandO3 is command of type [CommandO] with content "O3".
	CommandO3 = &CommandStub[TypeO]{Content: "O3"}

	// CommandP1 is command of type [CommandP] with content "P1".
	CommandP1 = &CommandStub[TypeP]{Content: "P1"}
	// CommandP2 is command of type [CommandP] with content "P2".
	CommandP2 = &CommandStub[TypeP]{Content: "P2"}
	// CommandP3 is command of type [CommandP] with content "P3".
	CommandP3 = &CommandStub[TypeP]{Content: "P3"}

	// CommandQ1 is command of type [CommandQ] with content "Q1".
	CommandQ1 = &CommandStub[TypeQ]{Content: "Q1"}
	// CommandQ2 is command of type [CommandQ] with content "Q2".
	CommandQ2 = &CommandStub[TypeQ]{Content: "Q2"}
	// CommandQ3 is command of type [CommandQ] with content "Q3".
	CommandQ3 = &CommandStub[TypeQ]{Content: "Q3"}

	// CommandR1 is command of type [CommandR] with content "R1".
	CommandR1 = &CommandStub[TypeR]{Content: "R1"}
	// CommandR2 is command of type [CommandR] with content "R2".
	CommandR2 = &CommandStub[TypeR]{Content: "R2"}
	// CommandR3 is command of type [CommandR] with content "R3".
	CommandR3 = &CommandStub[TypeR]{Content: "R3"}

	// CommandS1 is command of type [CommandS] with content "S1".
	CommandS1 = &CommandStub[TypeS]{Content: "S1"}
	// CommandS2 is command of type [CommandS] with content "S2".
	CommandS2 = &CommandStub[TypeS]{Content: "S2"}
	// CommandS3 is command of type [CommandS] with content "S3".
	CommandS3 = &CommandStub[TypeS]{Content: "S3"}

	// CommandT1 is command of type [CommandT] with content "T1".
	CommandT1 = &CommandStub[TypeT]{Content: "T1"}
	// CommandT2 is command of type [CommandT] with content "T2".
	CommandT2 = &CommandStub[TypeT]{Content: "T2"}
	// CommandT3 is command of type [CommandT] with content "T3".
	CommandT3 = &CommandStub[TypeT]{Content: "T3"}

	// CommandU1 is command of type [CommandU] with content "U1".
	CommandU1 = &CommandStub[TypeU]{Content: "U1"}
	// CommandU2 is command of type [CommandU] with content "U2".
	CommandU2 = &CommandStub[TypeU]{Content: "U2"}
	// CommandU3 is command of type [CommandU] with content "U3".
	CommandU3 = &CommandStub[TypeU]{Content: "U3"}

	// CommandV1 is command of type [CommandV] with content "V1".
	CommandV1 = &CommandStub[TypeV]{Content: "V1"}
	// CommandV2 is command of type [CommandV] with content "V2".
	CommandV2 = &CommandStub[TypeV]{Content: "V2"}
	// CommandV3 is command of type [CommandV] with content "V3".
	CommandV3 = &CommandStub[TypeV]{Content: "V3"}

	// CommandW1 is command of type [CommandW] with content "W1".
	CommandW1 = &CommandStub[TypeW]{Content: "W1"}
	// CommandW2 is command of type [CommandW] with content "W2".
	CommandW2 = &CommandStub[TypeW]{Content: "W2"}
	// CommandW3 is command of type [CommandW] with content "W3".
	CommandW3 = &CommandStub[TypeW]{Content: "W3"}

	// CommandX1 is command of type [CommandX] with content "X1".
	CommandX1 = &CommandStub[TypeX]{Content: "X1"}
	// CommandX2 is command of type [CommandX] with content "X2".
	CommandX2 = &CommandStub[TypeX]{Content: "X2"}
	// CommandX3 is command of type [CommandX] with content "X3".
	CommandX3 = &CommandStub[TypeX]{Content: "X3"}

	// CommandY1 is command of type [CommandY] with content "Y1".
	CommandY1 = &CommandStub[TypeY]{Content: "Y1"}
	// CommandY2 is command of type [CommandY] with content "Y2".
	CommandY2 = &CommandStub[TypeY]{Content: "Y2"}
	// CommandY3 is command of type [CommandY] with content "Y3".
	CommandY3 = &CommandStub[TypeY]{Content: "Y3"}

	// CommandZ1 is command of type [CommandZ] with content "Z1".
	CommandZ1 = &CommandStub[TypeZ]{Content: "Z1"}
	// CommandZ2 is command of type [CommandZ] with content "Z2".
	CommandZ2 = &CommandStub[TypeZ]{Content: "Z2"}
	// CommandZ3 is command of type [CommandZ] with content "Z3".
	CommandZ3 = &CommandStub[TypeZ]{Content: "Z3"}
)

var (
	// EventA1 is event of type [EventA] with content "A1".
	EventA1 = &EventStub[TypeA]{Content: "A1"}
	// EventA2 is event of type [EventA] with content "A2".
	EventA2 = &EventStub[TypeA]{Content: "A2"}
	// EventA3 is event of type [EventA] with content "A3".
	EventA3 = &EventStub[TypeA]{Content: "A3"}

	// EventB1 is event of type [EventB] with content "B1".
	EventB1 = &EventStub[TypeB]{Content: "B1"}
	// EventB2 is event of type [EventB] with content "B2".
	EventB2 = &EventStub[TypeB]{Content: "B2"}
	// EventB3 is event of type [EventB] with content "B3".
	EventB3 = &EventStub[TypeB]{Content: "B3"}

	// EventC1 is event of type [EventC] with content "C1".
	EventC1 = &EventStub[TypeC]{Content: "C1"}
	// EventC2 is event of type [EventC] with content "C2".
	EventC2 = &EventStub[TypeC]{Content: "C2"}
	// EventC3 is event of type [EventC] with content "C3".
	EventC3 = &EventStub[TypeC]{Content: "C3"}

	// EventD1 is event of type [EventD] with content "D1".
	EventD1 = &EventStub[TypeD]{Content: "D1"}
	// EventD2 is event of type [EventD] with content "D2".
	EventD2 = &EventStub[TypeD]{Content: "D2"}
	// EventD3 is event of type [EventD] with content "D3".
	EventD3 = &EventStub[TypeD]{Content: "D3"}

	// EventE1 is event of type [EventE] with content "E1".
	EventE1 = &EventStub[TypeE]{Content: "E1"}
	// EventE2 is event of type [EventE] with content "E2".
	EventE2 = &EventStub[TypeE]{Content: "E2"}
	// EventE3 is event of type [EventE] with content "E3".
	EventE3 = &EventStub[TypeE]{Content: "E3"}

	// EventF1 is event of type [EventF] with content "F1".
	EventF1 = &EventStub[TypeF]{Content: "F1"}
	// EventF2 is event of type [EventF] with content "F2".
	EventF2 = &EventStub[TypeF]{Content: "F2"}
	// EventF3 is event of type [EventF] with content "F3".
	EventF3 = &EventStub[TypeF]{Content: "F3"}

	// EventG1 is event of type [EventG] with content "G1".
	EventG1 = &EventStub[TypeG]{Content: "G1"}
	// EventG2 is event of type [EventG] with content "G2".
	EventG2 = &EventStub[TypeG]{Content: "G2"}
	// EventG3 is event of type [EventG] with content "G3".
	EventG3 = &EventStub[TypeG]{Content: "G3"}

	// EventH1 is event of type [EventH] with content "H1".
	EventH1 = &EventStub[TypeH]{Content: "H1"}
	// EventH2 is event of type [EventH] with content "H2".
	EventH2 = &EventStub[TypeH]{Content: "H2"}
	// EventH3 is event of type [EventH] with content "H3".
	EventH3 = &EventStub[TypeH]{Content: "H3"}

	// EventI1 is event of type [EventI] with content "I1".
	EventI1 = &EventStub[TypeI]{Content: "I1"}
	// EventI2 is event of type [EventI] with content "I2".
	EventI2 = &EventStub[TypeI]{Content: "I2"}
	// EventI3 is event of type [EventI] with content "I3".
	EventI3 = &EventStub[TypeI]{Content: "I3"}

	// EventJ1 is event of type [EventJ] with content "J1".
	EventJ1 = &EventStub[TypeJ]{Content: "J1"}
	// EventJ2 is event of type [EventJ] with content "J2".
	EventJ2 = &EventStub[TypeJ]{Content: "J2"}
	// EventJ3 is event of type [EventJ] with content "J3".
	EventJ3 = &EventStub[TypeJ]{Content: "J3"}

	// EventK1 is event of type [EventK] with content "K1".
	EventK1 = &EventStub[TypeK]{Content: "K1"}
	// EventK2 is event of type [EventK] with content "K2".
	EventK2 = &EventStub[TypeK]{Content: "K2"}
	// EventK3 is event of type [EventK] with content "K3".
	EventK3 = &EventStub[TypeK]{Content: "K3"}

	// EventL1 is event of type [EventL] with content "L1".
	EventL1 = &EventStub[TypeL]{Content: "L1"}
	// EventL2 is event of type [EventL] with content "L2".
	EventL2 = &EventStub[TypeL]{Content: "L2"}
	// EventL3 is event of type [EventL] with content "L3".
	EventL3 = &EventStub[TypeL]{Content: "L3"}

	// EventM1 is event of type [EventM] with content "M1".
	EventM1 = &EventStub[TypeM]{Content: "M1"}
	// EventM2 is event of type [EventM] with content "M2".
	EventM2 = &EventStub[TypeM]{Content: "M2"}
	// EventM3 is event of type [EventM] with content "M3".
	EventM3 = &EventStub[TypeM]{Content: "M3"}

	// EventN1 is event of type [EventN] with content "N1".
	EventN1 = &EventStub[TypeN]{Content: "N1"}
	// EventN2 is event of type [EventN] with content "N2".
	EventN2 = &EventStub[TypeN]{Content: "N2"}
	// EventN3 is event of type [EventN] with content "N3".
	EventN3 = &EventStub[TypeN]{Content: "N3"}

	// EventO1 is event of type [EventO] with content "O1".
	EventO1 = &EventStub[TypeO]{Content: "O1"}
	// EventO2 is event of type [EventO] with content "O2".
	EventO2 = &EventStub[TypeO]{Content: "O2"}
	// EventO3 is event of type [EventO] with content "O3".
	EventO3 = &EventStub[TypeO]{Content: "O3"}

	// EventP1 is event of type [EventP] with content "P1".
	EventP1 = &EventStub[TypeP]{Content: "P1"}
	// EventP2 is event of type [EventP] with content "P2".
	EventP2 = &EventStub[TypeP]{Content: "P2"}
	// EventP3 is event of type [EventP] with content "P3".
	EventP3 = &EventStub[TypeP]{Content: "P3"}

	// EventQ1 is event of type [EventQ] with content "Q1".
	EventQ1 = &EventStub[TypeQ]{Content: "Q1"}
	// EventQ2 is event of type [EventQ] with content "Q2".
	EventQ2 = &EventStub[TypeQ]{Content: "Q2"}
	// EventQ3 is event of type [EventQ] with content "Q3".
	EventQ3 = &EventStub[TypeQ]{Content: "Q3"}

	// EventR1 is event of type [EventR] with content "R1".
	EventR1 = &EventStub[TypeR]{Content: "R1"}
	// EventR2 is event of type [EventR] with content "R2".
	EventR2 = &EventStub[TypeR]{Content: "R2"}
	// EventR3 is event of type [EventR] with content "R3".
	EventR3 = &EventStub[TypeR]{Content: "R3"}

	// EventS1 is event of type [EventS] with content "S1".
	EventS1 = &EventStub[TypeS]{Content: "S1"}
	// EventS2 is event of type [EventS] with content "S2".
	EventS2 = &EventStub[TypeS]{Content: "S2"}
	// EventS3 is event of type [EventS] with content "S3".
	EventS3 = &EventStub[TypeS]{Content: "S3"}

	// EventT1 is event of type [EventT] with content "T1".
	EventT1 = &EventStub[TypeT]{Content: "T1"}
	// EventT2 is event of type [EventT] with content "T2".
	EventT2 = &EventStub[TypeT]{Content: "T2"}
	// EventT3 is event of type [EventT] with content "T3".
	EventT3 = &EventStub[TypeT]{Content: "T3"}

	// EventU1 is event of type [EventU] with content "U1".
	EventU1 = &EventStub[TypeU]{Content: "U1"}
	// EventU2 is event of type [EventU] with content "U2".
	EventU2 = &EventStub[TypeU]{Content: "U2"}
	// EventU3 is event of type [EventU] with content "U3".
	EventU3 = &EventStub[TypeU]{Content: "U3"}

	// EventV1 is event of type [EventV] with content "V1".
	EventV1 = &EventStub[TypeV]{Content: "V1"}
	// EventV2 is event of type [EventV] with content "V2".
	EventV2 = &EventStub[TypeV]{Content: "V2"}
	// EventV3 is event of type [EventV] with content "V3".
	EventV3 = &EventStub[TypeV]{Content: "V3"}

	// EventW1 is event of type [EventW] with content "W1".
	EventW1 = &EventStub[TypeW]{Content: "W1"}
	// EventW2 is event of type [EventW] with content "W2".
	EventW2 = &EventStub[TypeW]{Content: "W2"}
	// EventW3 is event of type [EventW] with content "W3".
	EventW3 = &EventStub[TypeW]{Content: "W3"}

	// EventX1 is event of type [EventX] with content "X1".
	EventX1 = &EventStub[TypeX]{Content: "X1"}
	// EventX2 is event of type [EventX] with content "X2".
	EventX2 = &EventStub[TypeX]{Content: "X2"}
	// EventX3 is event of type [EventX] with content "X3".
	EventX3 = &EventStub[TypeX]{Content: "X3"}

	// EventY1 is event of type [EventY] with content "Y1".
	EventY1 = &EventStub[TypeY]{Content: "Y1"}
	// EventY2 is event of type [EventY] with content "Y2".
	EventY2 = &EventStub[TypeY]{Content: "Y2"}
	// EventY3 is event of type [EventY] with content "Y3".
	EventY3 = &EventStub[TypeY]{Content: "Y3"}

	// EventZ1 is event of type [EventZ] with content "Z1".
	EventZ1 = &EventStub[TypeZ]{Content: "Z1"}
	// EventZ2 is event of type [EventZ] with content "Z2".
	EventZ2 = &EventStub[TypeZ]{Content: "Z2"}
	// EventZ3 is event of type [EventZ] with content "Z3".
	EventZ3 = &EventStub[TypeZ]{Content: "Z3"}
)

var (
	// TimeoutA1 is a timeout message of type [TimeoutA] with content "A1".
	TimeoutA1 = &TimeoutStub[TypeA]{Content: "A1"}
	// TimeoutA2 is a timeout message of type [TimeoutA] with content "A2".
	TimeoutA2 = &TimeoutStub[TypeA]{Content: "A2"}
	// TimeoutA3 is a timeout message of type [TimeoutA] with content "A3".
	TimeoutA3 = &TimeoutStub[TypeA]{Content: "A3"}

	// TimeoutB1 is a timeout message of type [TimeoutB] with content "B1".
	TimeoutB1 = &TimeoutStub[TypeB]{Content: "B1"}
	// TimeoutB2 is a timeout message of type [TimeoutB] with content "B2".
	TimeoutB2 = &TimeoutStub[TypeB]{Content: "B2"}
	// TimeoutB3 is a timeout message of type [TimeoutB] with content "B3".
	TimeoutB3 = &TimeoutStub[TypeB]{Content: "B3"}

	// TimeoutC1 is a timeout message of type [TimeoutC] with content "C1".
	TimeoutC1 = &TimeoutStub[TypeC]{Content: "C1"}
	// TimeoutC2 is a timeout message of type [TimeoutC] with content "C2".
	TimeoutC2 = &TimeoutStub[TypeC]{Content: "C2"}
	// TimeoutC3 is a timeout message of type [TimeoutC] with content "C3".
	TimeoutC3 = &TimeoutStub[TypeC]{Content: "C3"}

	// TimeoutD1 is a timeout message of type [TimeoutD] with content "D1".
	TimeoutD1 = &TimeoutStub[TypeD]{Content: "D1"}
	// TimeoutD2 is a timeout message of type [TimeoutD] with content "D2".
	TimeoutD2 = &TimeoutStub[TypeD]{Content: "D2"}
	// TimeoutD3 is a timeout message of type [TimeoutD] with content "D3".
	TimeoutD3 = &TimeoutStub[TypeD]{Content: "D3"}

	// TimeoutE1 is a timeout message of type [TimeoutE] with content "E1".
	TimeoutE1 = &TimeoutStub[TypeE]{Content: "E1"}
	// TimeoutE2 is a timeout message of type [TimeoutE] with content "E2".
	TimeoutE2 = &TimeoutStub[TypeE]{Content: "E2"}
	// TimeoutE3 is a timeout message of type [TimeoutE] with content "E3".
	TimeoutE3 = &TimeoutStub[TypeE]{Content: "E3"}

	// TimeoutF1 is a timeout message of type [TimeoutF] with content "F1".
	TimeoutF1 = &TimeoutStub[TypeF]{Content: "F1"}
	// TimeoutF2 is a timeout message of type [TimeoutF] with content "F2".
	TimeoutF2 = &TimeoutStub[TypeF]{Content: "F2"}
	// TimeoutF3 is a timeout message of type [TimeoutF] with content "F3".
	TimeoutF3 = &TimeoutStub[TypeF]{Content: "F3"}

	// TimeoutG1 is a timeout message of type [TimeoutG] with content "G1".
	TimeoutG1 = &TimeoutStub[TypeG]{Content: "G1"}
	// TimeoutG2 is a timeout message of type [TimeoutG] with content "G2".
	TimeoutG2 = &TimeoutStub[TypeG]{Content: "G2"}
	// TimeoutG3 is a timeout message of type [TimeoutG] with content "G3".
	TimeoutG3 = &TimeoutStub[TypeG]{Content: "G3"}

	// TimeoutH1 is a timeout message of type [TimeoutH] with content "H1".
	TimeoutH1 = &TimeoutStub[TypeH]{Content: "H1"}
	// TimeoutH2 is a timeout message of type [TimeoutH] with content "H2".
	TimeoutH2 = &TimeoutStub[TypeH]{Content: "H2"}
	// TimeoutH3 is a timeout message of type [TimeoutH] with content "H3".
	TimeoutH3 = &TimeoutStub[TypeH]{Content: "H3"}

	// TimeoutI1 is a timeout message of type [TimeoutI] with content "I1".
	TimeoutI1 = &TimeoutStub[TypeI]{Content: "I1"}
	// TimeoutI2 is a timeout message of type [TimeoutI] with content "I2".
	TimeoutI2 = &TimeoutStub[TypeI]{Content: "I2"}
	// TimeoutI3 is a timeout message of type [TimeoutI] with content "I3".
	TimeoutI3 = &TimeoutStub[TypeI]{Content: "I3"}

	// TimeoutJ1 is a timeout message of type [TimeoutJ] with content "J1".
	TimeoutJ1 = &TimeoutStub[TypeJ]{Content: "J1"}
	// TimeoutJ2 is a timeout message of type [TimeoutJ] with content "J2".
	TimeoutJ2 = &TimeoutStub[TypeJ]{Content: "J2"}
	// TimeoutJ3 is a timeout message of type [TimeoutJ] with content "J3".
	TimeoutJ3 = &TimeoutStub[TypeJ]{Content: "J3"}

	// TimeoutK1 is a timeout message of type [TimeoutK] with content "K1".
	TimeoutK1 = &TimeoutStub[TypeK]{Content: "K1"}
	// TimeoutK2 is a timeout message of type [TimeoutK] with content "K2".
	TimeoutK2 = &TimeoutStub[TypeK]{Content: "K2"}
	// TimeoutK3 is a timeout message of type [TimeoutK] with content "K3".
	TimeoutK3 = &TimeoutStub[TypeK]{Content: "K3"}

	// TimeoutL1 is a timeout message of type [TimeoutL] with content "L1".
	TimeoutL1 = &TimeoutStub[TypeL]{Content: "L1"}
	// TimeoutL2 is a timeout message of type [TimeoutL] with content "L2".
	TimeoutL2 = &TimeoutStub[TypeL]{Content: "L2"}
	// TimeoutL3 is a timeout message of type [TimeoutL] with content "L3".
	TimeoutL3 = &TimeoutStub[TypeL]{Content: "L3"}

	// TimeoutM1 is a timeout message of type [TimeoutM] with content "M1".
	TimeoutM1 = &TimeoutStub[TypeM]{Content: "M1"}
	// TimeoutM2 is a timeout message of type [TimeoutM] with content "M2".
	TimeoutM2 = &TimeoutStub[TypeM]{Content: "M2"}
	// TimeoutM3 is a timeout message of type [TimeoutM] with content "M3".
	TimeoutM3 = &TimeoutStub[TypeM]{Content: "M3"}

	// TimeoutN1 is a timeout message of type [TimeoutN] with content "N1".
	TimeoutN1 = &TimeoutStub[TypeN]{Content: "N1"}
	// TimeoutN2 is a timeout message of type [TimeoutN] with content "N2".
	TimeoutN2 = &TimeoutStub[TypeN]{Content: "N2"}
	// TimeoutN3 is a timeout message of type [TimeoutN] with content "N3".
	TimeoutN3 = &TimeoutStub[TypeN]{Content: "N3"}

	// TimeoutO1 is a timeout message of type [TimeoutO] with content "O1".
	TimeoutO1 = &TimeoutStub[TypeO]{Content: "O1"}
	// TimeoutO2 is a timeout message of type [TimeoutO] with content "O2".
	TimeoutO2 = &TimeoutStub[TypeO]{Content: "O2"}
	// TimeoutO3 is a timeout message of type [TimeoutO] with content "O3".
	TimeoutO3 = &TimeoutStub[TypeO]{Content: "O3"}

	// TimeoutP1 is a timeout message of type [TimeoutP] with content "P1".
	TimeoutP1 = &TimeoutStub[TypeP]{Content: "P1"}
	// TimeoutP2 is a timeout message of type [TimeoutP] with content "P2".
	TimeoutP2 = &TimeoutStub[TypeP]{Content: "P2"}
	// TimeoutP3 is a timeout message of type [TimeoutP] with content "P3".
	TimeoutP3 = &TimeoutStub[TypeP]{Content: "P3"}

	// TimeoutQ1 is a timeout message of type [TimeoutQ] with content "Q1".
	TimeoutQ1 = &TimeoutStub[TypeQ]{Content: "Q1"}
	// TimeoutQ2 is a timeout message of type [TimeoutQ] with content "Q2".
	TimeoutQ2 = &TimeoutStub[TypeQ]{Content: "Q2"}
	// TimeoutQ3 is a timeout message of type [TimeoutQ] with content "Q3".
	TimeoutQ3 = &TimeoutStub[TypeQ]{Content: "Q3"}

	// TimeoutR1 is a timeout message of type [TimeoutR] with content "R1".
	TimeoutR1 = &TimeoutStub[TypeR]{Content: "R1"}
	// TimeoutR2 is a timeout message of type [TimeoutR] with content "R2".
	TimeoutR2 = &TimeoutStub[TypeR]{Content: "R2"}
	// TimeoutR3 is a timeout message of type [TimeoutR] with content "R3".
	TimeoutR3 = &TimeoutStub[TypeR]{Content: "R3"}

	// TimeoutS1 is a timeout message of type [TimeoutS] with content "S1".
	TimeoutS1 = &TimeoutStub[TypeS]{Content: "S1"}
	// TimeoutS2 is a timeout message of type [TimeoutS] with content "S2".
	TimeoutS2 = &TimeoutStub[TypeS]{Content: "S2"}
	// TimeoutS3 is a timeout message of type [TimeoutS] with content "S3".
	TimeoutS3 = &TimeoutStub[TypeS]{Content: "S3"}

	// TimeoutT1 is a timeout message of type [TimeoutT] with content "T1".
	TimeoutT1 = &TimeoutStub[TypeT]{Content: "T1"}
	// TimeoutT2 is a timeout message of type [TimeoutT] with content "T2".
	TimeoutT2 = &TimeoutStub[TypeT]{Content: "T2"}
	// TimeoutT3 is a timeout message of type [TimeoutT] with content "T3".
	TimeoutT3 = &TimeoutStub[TypeT]{Content: "T3"}

	// TimeoutU1 is a timeout message of type [TimeoutU] with content "U1".
	TimeoutU1 = &TimeoutStub[TypeU]{Content: "U1"}
	// TimeoutU2 is a timeout message of type [TimeoutU] with content "U2".
	TimeoutU2 = &TimeoutStub[TypeU]{Content: "U2"}
	// TimeoutU3 is a timeout message of type [TimeoutU] with content "U3".
	TimeoutU3 = &TimeoutStub[TypeU]{Content: "U3"}

	// TimeoutV1 is a timeout message of type [TimeoutV] with content "V1".
	TimeoutV1 = &TimeoutStub[TypeV]{Content: "V1"}
	// TimeoutV2 is a timeout message of type [TimeoutV] with content "V2".
	TimeoutV2 = &TimeoutStub[TypeV]{Content: "V2"}
	// TimeoutV3 is a timeout message of type [TimeoutV] with content "V3".
	TimeoutV3 = &TimeoutStub[TypeV]{Content: "V3"}

	// TimeoutW1 is a timeout message of type [TimeoutW] with content "W1".
	TimeoutW1 = &TimeoutStub[TypeW]{Content: "W1"}
	// TimeoutW2 is a timeout message of type [TimeoutW] with content "W2".
	TimeoutW2 = &TimeoutStub[TypeW]{Content: "W2"}
	// TimeoutW3 is a timeout message of type [TimeoutW] with content "W3".
	TimeoutW3 = &TimeoutStub[TypeW]{Content: "W3"}

	// TimeoutX1 is a timeout message of type [TimeoutX] with content "X1".
	TimeoutX1 = &TimeoutStub[TypeX]{Content: "X1"}
	// TimeoutX2 is a timeout message of type [TimeoutX] with content "X2".
	TimeoutX2 = &TimeoutStub[TypeX]{Content: "X2"}
	// TimeoutX3 is a timeout message of type [TimeoutX] with content "X3".
	TimeoutX3 = &TimeoutStub[TypeX]{Content: "X3"}

	// TimeoutY1 is a timeout message of type [TimeoutY] with content "Y1".
	TimeoutY1 = &TimeoutStub[TypeY]{Content: "Y1"}
	// TimeoutY2 is a timeout message of type [TimeoutY] with content "Y2".
	TimeoutY2 = &TimeoutStub[TypeY]{Content: "Y2"}
	// TimeoutY3 is a timeout message of type [TimeoutY] with content "Y3".
	TimeoutY3 = &TimeoutStub[TypeY]{Content: "Y3"}

	// TimeoutZ1 is a timeout message of type [TimeoutZ] with content "Z1".
	TimeoutZ1 = &TimeoutStub[TypeZ]{Content: "Z1"}
	// TimeoutZ2 is a timeout message of type [TimeoutZ] with content "Z2".
	TimeoutZ2 = &TimeoutStub[TypeZ]{Content: "Z2"}
	// TimeoutZ3 is a timeout message of type [TimeoutZ] with content "Z3".
	TimeoutZ3 = &TimeoutStub[TypeZ]{Content: "Z3"}
)

var namePattern = regexp.MustCompile(`(Command|Event|Timeout)Stub\[[^]]+\.Type(.)\]`)

// MessageTypeID returns the RFC 4122 UUID for a message stub of type T.
//
// T must be one of [CommandStub], [EventStub], or [TimeoutStub], with a type
// parameter of [TypeA] to [TypeZ], otherwise the function panics.
func MessageTypeID[
	T interface {
		dogma.Message
		*E
	},
	E any,
]() string {
	typeName := reflect.TypeFor[T]().String()
	matches := namePattern.FindStringSubmatch(typeName)

	if matches == nil {
		panic("cannot generate message type ID for non-stub message: " + typeName)
	}

	var prefix string
	switch matches[1] {
	case "Command":
		prefix = "c"
	case "Event":
		prefix = "e"
	case "Timeout":
		prefix = "7" // 7 looks a little like `t`, for timeout, right?
	}

	letter := 0xa + matches[2][0] - 'A' // convert letter to 0-25 range

	return fmt.Sprintf("%s0000000-0000-4000-8000-%012x", prefix, letter)
}

func init() {
	dogma.RegisterCommand[*CommandStub[TypeA]](MessageTypeID[*CommandStub[TypeA]]())
	dogma.RegisterCommand[*CommandStub[TypeB]](MessageTypeID[*CommandStub[TypeB]]())
	dogma.RegisterCommand[*CommandStub[TypeC]](MessageTypeID[*CommandStub[TypeC]]())
	dogma.RegisterCommand[*CommandStub[TypeD]](MessageTypeID[*CommandStub[TypeD]]())
	dogma.RegisterCommand[*CommandStub[TypeE]](MessageTypeID[*CommandStub[TypeE]]())
	dogma.RegisterCommand[*CommandStub[TypeF]](MessageTypeID[*CommandStub[TypeF]]())
	dogma.RegisterCommand[*CommandStub[TypeG]](MessageTypeID[*CommandStub[TypeG]]())
	dogma.RegisterCommand[*CommandStub[TypeH]](MessageTypeID[*CommandStub[TypeH]]())
	dogma.RegisterCommand[*CommandStub[TypeI]](MessageTypeID[*CommandStub[TypeI]]())
	dogma.RegisterCommand[*CommandStub[TypeJ]](MessageTypeID[*CommandStub[TypeJ]]())
	dogma.RegisterCommand[*CommandStub[TypeK]](MessageTypeID[*CommandStub[TypeK]]())
	dogma.RegisterCommand[*CommandStub[TypeL]](MessageTypeID[*CommandStub[TypeL]]())
	dogma.RegisterCommand[*CommandStub[TypeM]](MessageTypeID[*CommandStub[TypeM]]())
	dogma.RegisterCommand[*CommandStub[TypeN]](MessageTypeID[*CommandStub[TypeN]]())
	dogma.RegisterCommand[*CommandStub[TypeO]](MessageTypeID[*CommandStub[TypeO]]())
	dogma.RegisterCommand[*CommandStub[TypeP]](MessageTypeID[*CommandStub[TypeP]]())
	dogma.RegisterCommand[*CommandStub[TypeQ]](MessageTypeID[*CommandStub[TypeQ]]())
	dogma.RegisterCommand[*CommandStub[TypeR]](MessageTypeID[*CommandStub[TypeR]]())
	dogma.RegisterCommand[*CommandStub[TypeS]](MessageTypeID[*CommandStub[TypeS]]())
	dogma.RegisterCommand[*CommandStub[TypeT]](MessageTypeID[*CommandStub[TypeT]]())
	dogma.RegisterCommand[*CommandStub[TypeU]](MessageTypeID[*CommandStub[TypeU]]())
	dogma.RegisterCommand[*CommandStub[TypeV]](MessageTypeID[*CommandStub[TypeV]]())
	dogma.RegisterCommand[*CommandStub[TypeW]](MessageTypeID[*CommandStub[TypeW]]())
	dogma.RegisterCommand[*CommandStub[TypeX]](MessageTypeID[*CommandStub[TypeX]]())
	dogma.RegisterCommand[*CommandStub[TypeY]](MessageTypeID[*CommandStub[TypeY]]())
	dogma.RegisterCommand[*CommandStub[TypeZ]](MessageTypeID[*CommandStub[TypeZ]]())

	dogma.RegisterEvent[*EventStub[TypeA]](MessageTypeID[*EventStub[TypeA]]())
	dogma.RegisterEvent[*EventStub[TypeB]](MessageTypeID[*EventStub[TypeB]]())
	dogma.RegisterEvent[*EventStub[TypeC]](MessageTypeID[*EventStub[TypeC]]())
	dogma.RegisterEvent[*EventStub[TypeD]](MessageTypeID[*EventStub[TypeD]]())
	dogma.RegisterEvent[*EventStub[TypeE]](MessageTypeID[*EventStub[TypeE]]())
	dogma.RegisterEvent[*EventStub[TypeF]](MessageTypeID[*EventStub[TypeF]]())
	dogma.RegisterEvent[*EventStub[TypeG]](MessageTypeID[*EventStub[TypeG]]())
	dogma.RegisterEvent[*EventStub[TypeH]](MessageTypeID[*EventStub[TypeH]]())
	dogma.RegisterEvent[*EventStub[TypeI]](MessageTypeID[*EventStub[TypeI]]())
	dogma.RegisterEvent[*EventStub[TypeJ]](MessageTypeID[*EventStub[TypeJ]]())
	dogma.RegisterEvent[*EventStub[TypeK]](MessageTypeID[*EventStub[TypeK]]())
	dogma.RegisterEvent[*EventStub[TypeL]](MessageTypeID[*EventStub[TypeL]]())
	dogma.RegisterEvent[*EventStub[TypeM]](MessageTypeID[*EventStub[TypeM]]())
	dogma.RegisterEvent[*EventStub[TypeN]](MessageTypeID[*EventStub[TypeN]]())
	dogma.RegisterEvent[*EventStub[TypeO]](MessageTypeID[*EventStub[TypeO]]())
	dogma.RegisterEvent[*EventStub[TypeP]](MessageTypeID[*EventStub[TypeP]]())
	dogma.RegisterEvent[*EventStub[TypeQ]](MessageTypeID[*EventStub[TypeQ]]())
	dogma.RegisterEvent[*EventStub[TypeR]](MessageTypeID[*EventStub[TypeR]]())
	dogma.RegisterEvent[*EventStub[TypeS]](MessageTypeID[*EventStub[TypeS]]())
	dogma.RegisterEvent[*EventStub[TypeT]](MessageTypeID[*EventStub[TypeT]]())
	dogma.RegisterEvent[*EventStub[TypeU]](MessageTypeID[*EventStub[TypeU]]())
	dogma.RegisterEvent[*EventStub[TypeV]](MessageTypeID[*EventStub[TypeV]]())
	dogma.RegisterEvent[*EventStub[TypeW]](MessageTypeID[*EventStub[TypeW]]())
	dogma.RegisterEvent[*EventStub[TypeX]](MessageTypeID[*EventStub[TypeX]]())
	dogma.RegisterEvent[*EventStub[TypeY]](MessageTypeID[*EventStub[TypeY]]())
	dogma.RegisterEvent[*EventStub[TypeZ]](MessageTypeID[*EventStub[TypeZ]]())

	dogma.RegisterTimeout[*TimeoutStub[TypeA]](MessageTypeID[*TimeoutStub[TypeA]]())
	dogma.RegisterTimeout[*TimeoutStub[TypeB]](MessageTypeID[*TimeoutStub[TypeB]]())
	dogma.RegisterTimeout[*TimeoutStub[TypeC]](MessageTypeID[*TimeoutStub[TypeC]]())
	dogma.RegisterTimeout[*TimeoutStub[TypeD]](MessageTypeID[*TimeoutStub[TypeD]]())
	dogma.RegisterTimeout[*TimeoutStub[TypeE]](MessageTypeID[*TimeoutStub[TypeE]]())
	dogma.RegisterTimeout[*TimeoutStub[TypeF]](MessageTypeID[*TimeoutStub[TypeF]]())
	dogma.RegisterTimeout[*TimeoutStub[TypeG]](MessageTypeID[*TimeoutStub[TypeG]]())
	dogma.RegisterTimeout[*TimeoutStub[TypeH]](MessageTypeID[*TimeoutStub[TypeH]]())
	dogma.RegisterTimeout[*TimeoutStub[TypeI]](MessageTypeID[*TimeoutStub[TypeI]]())
	dogma.RegisterTimeout[*TimeoutStub[TypeJ]](MessageTypeID[*TimeoutStub[TypeJ]]())
	dogma.RegisterTimeout[*TimeoutStub[TypeK]](MessageTypeID[*TimeoutStub[TypeK]]())
	dogma.RegisterTimeout[*TimeoutStub[TypeL]](MessageTypeID[*TimeoutStub[TypeL]]())
	dogma.RegisterTimeout[*TimeoutStub[TypeM]](MessageTypeID[*TimeoutStub[TypeM]]())
	dogma.RegisterTimeout[*TimeoutStub[TypeN]](MessageTypeID[*TimeoutStub[TypeN]]())
	dogma.RegisterTimeout[*TimeoutStub[TypeO]](MessageTypeID[*TimeoutStub[TypeO]]())
	dogma.RegisterTimeout[*TimeoutStub[TypeP]](MessageTypeID[*TimeoutStub[TypeP]]())
	dogma.RegisterTimeout[*TimeoutStub[TypeQ]](MessageTypeID[*TimeoutStub[TypeQ]]())
	dogma.RegisterTimeout[*TimeoutStub[TypeR]](MessageTypeID[*TimeoutStub[TypeR]]())
	dogma.RegisterTimeout[*TimeoutStub[TypeS]](MessageTypeID[*TimeoutStub[TypeS]]())
	dogma.RegisterTimeout[*TimeoutStub[TypeT]](MessageTypeID[*TimeoutStub[TypeT]]())
	dogma.RegisterTimeout[*TimeoutStub[TypeU]](MessageTypeID[*TimeoutStub[TypeU]]())
	dogma.RegisterTimeout[*TimeoutStub[TypeV]](MessageTypeID[*TimeoutStub[TypeV]]())
	dogma.RegisterTimeout[*TimeoutStub[TypeW]](MessageTypeID[*TimeoutStub[TypeW]]())
	dogma.RegisterTimeout[*TimeoutStub[TypeX]](MessageTypeID[*TimeoutStub[TypeX]]())
	dogma.RegisterTimeout[*TimeoutStub[TypeY]](MessageTypeID[*TimeoutStub[TypeY]]())
	dogma.RegisterTimeout[*TimeoutStub[TypeZ]](MessageTypeID[*TimeoutStub[TypeZ]]())
}
