package config

import (
	"fmt"
	"io"
	"strings"

	"github.com/dogmatiq/enginekit/config/internal/renderer"
	"github.com/dogmatiq/enginekit/internal/typename"
	"github.com/dogmatiq/enginekit/optional"
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

func renderEntityDescriptor[T any](
	ren *renderer.Renderer,
	label string,
	src Value[T],
) {
	ren.Print(label)

	if typeName, ok := src.TypeName.TryGet(); ok {
		typeName = typename.Unqualified(typeName)
		typeName = strings.TrimPrefix(typeName, "*")
		ren.Print(":", typeName)
	}
}

func renderHandlerDetails[T any](
	ren *renderer.Renderer,
	h Handler,
	src Value[T],
	isDisabled optional.Optional[bool],
) {
	f, err := validate(h)

	if d, ok := isDisabled.TryGet(); ok && d {
		ren.Print("disabled ")
	}

	renderFidelity(ren, f, err)
	ren.Print(h.HandlerType().String())

	if typeName, ok := src.TypeName.TryGet(); ok {
		ren.Print(" ", typeName)

		if !src.Value.IsPresent() {
			ren.Print(" (runtime type unavailable)")
		}
	}

	ren.Print("\n")
	renderErrors(ren, err)

	for _, i := range h.identities() {
		ren.IndentBullet()
		i.renderDetails(ren)
		ren.Dedent()
	}

	for _, r := range h.routes() {
		ren.IndentBullet()
		r.renderDetails(ren)
		ren.Dedent()
	}
}
