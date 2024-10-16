package config

import (
	"fmt"
	"strings"

	"github.com/dogmatiq/enginekit/internal/ioutil"
)

// RenderDescriptor returns a one-line, human-readable string that attempts to
// uniquely identity the component based on the the information available. It
// may or may not consist of information from an [Identity].
func RenderDescriptor(c Component) string {
	w := &strings.Builder{}
	p := &ioutil.Renderer{Target: w}

	c.renderDescriptor(p)

	if _, err := p.Done(); err != nil {
		panic(err)
	}

	return w.String()
}

// RenderDetails returns a multi-line, human-readable string that describes the
// component in detail.
func RenderDetails(Component) string {
	panic("not implemented")
}

// renderList returns a human-readable list of items.
func renderList[T any](items []T) string {
	var s string

	for i, item := range items {
		if i == len(items)-1 {
			s += " and "
		} else if i > 0 {
			s += ", "
		}
		s += fmt.Sprint(item)
	}

	return s
}

func renderEntityDescriptor[T any](
	ren *ioutil.Renderer,
	label string,
	e Entity,
	src Value[T],
) {
	ren.Print(label)

	desc := ""

	if n, ok := src.TypeName.TryGet(); ok {
		desc = strings.TrimPrefix(n, "*")
	}

	if desc == "" {
		for _, id := range e.identities() {
			if norm, err := Normalize(id); err == nil {
				desc = norm.String()
				break
			}
			desc = id.String()
		}
	}

	if desc != "" {
		ren.Print(":", desc)
	}
}
