package config

import (
	"fmt"
	"io"
	"strings"

	"github.com/dogmatiq/enginekit/config/internal/renderer"
)

// RenderDescriptor returns a one-line human-readable description of c.
//
// The descriptor attempts to uniquely identify the component based on the
// information available. It may or may not consist of information from an
// [Identity] component.
func RenderDescriptor(c Component) string {
	w := &strings.Builder{}
	if _, err := WriteDescriptor(w, c); err != nil {
		panic(err)
	}
	return w.String()
}

// RenderDetails returns a detailed human-readable description of c.
func RenderDetails(c Component) string {
	w := &strings.Builder{}
	if _, err := WriteDetails(w, c); err != nil {
		panic(err)
	}
	return w.String()
}

// WriteDescriptor writes a one-line human-readable description of c to w.
//
// The descriptor attempts to uniquely identify the component based on the
// information available. It may or may not consist of information from an
// [Identity] component.
func WriteDescriptor(w io.Writer, c Component) (int, error) {
	p := &renderer.Renderer{Target: w}
	c.renderDescriptor(p)
	return p.Done()
}

// WriteDetails writes a detailed human-readable description of c to w.
func WriteDetails(w io.Writer, c Component) (int, error) {
	p := &renderer.Renderer{Target: w}
	c.renderDetails(p)
	return p.Done()
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
	ren *renderer.Renderer,
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
