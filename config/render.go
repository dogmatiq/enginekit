package config

import (
	"fmt"
	"io"
	"strings"

	"github.com/dogmatiq/enginekit/config/internal/renderer"
	"github.com/dogmatiq/enginekit/internal/typename"
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
	r := &renderer.Renderer{Target: w}
	c.renderDescriptor(r)
	return r.Done()
}

// WriteDetails writes a detailed human-readable description of c to w.
func WriteDetails(w io.Writer, c Component) (int, error) {
	r := &renderer.Renderer{Target: w}
	c.renderDetails(r)
	return r.Done()
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
	src Value[T],
) {
	typeName := typename.Unqualified(
		src.TypeName.Get(),
	)

	ren.Print(
		label,
		":",
		strings.TrimPrefix(typeName, "*"),
	)
}

func renderFidelity(r *renderer.Renderer, f Fidelity, errs []error) {
	if f&Incomplete != 0 {
		r.Print("incomplete ")
	} else if len(errs) == 0 {
		r.Print("valid ")
	} else {
		r.Print("invalid ")
	}

	if f&Speculative != 0 {
		r.Print("speculative ")
	}
}

func renderErrors(r *renderer.Renderer, errs []error) {
	for _, err := range errs {
		r.IndentBullet()
		r.Print(err.Error(), "\n")
		r.Dedent()
	}
}
