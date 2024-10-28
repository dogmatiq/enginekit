package config

import "github.com/dogmatiq/enginekit/optional"

// A Symbol is a type that uniquely identifies a specific [Flag].
type Symbol interface {
	~struct{ symbol }
}

// symbol is a "marker" struct that is embedded in named types to declare a
// new [Symbol].
type symbol struct{}

// A Flag represents some boolean state of a [Component].
//
// Each type of flag is uniquely identified by a [Symbol].
type Flag[S Symbol] struct {
	Modifications []*FlagModification
}

// A FlagModification is a [Component] that represents a specific point at which
// a flag is set or unset within the configuration.
type FlagModification struct {
	ComponentCommon

	Value optional.Optional[bool]
}

func (m *FlagModification) String() string {
	panic("not implemented")
}
