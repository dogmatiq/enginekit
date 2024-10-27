package config

import (
	"fmt"
	"io"
	"reflect"
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

func renderEntityDescriptor(r *renderer.Renderer, e QEntity) {
	r.Print(strings.ToLower(reflect.TypeOf(e).Name()))

	if typeName, ok := e.CommonEntityProperties().TypeName.TryGet(); ok {
		typeName = typename.Unqualified(typeName)
		typeName = strings.TrimPrefix(typeName, "*")
		r.Print(":", typeName)
	}
}

func renderEntityDetails[T any](
	r *renderer.Renderer,
	e QEntity,
	src optional.Optional[T],
) {
	p := e.CommonEntityProperties()

	f, err := validate(e)

	renderFidelity(r, f, err)
	r.Print(strings.ToLower(reflect.TypeOf(e).Name()))

	if typeName, ok := p.TypeName.TryGet(); ok {
		r.Print(" ", typeName)

		if !src.IsPresent() {
			r.Print(" (runtime type unavailable)")
		}
	}

	r.Print("\n")
	renderErrors(r, err)

	for _, id := range p.IdentityComponents {
		r.IndentBullet()
		id.renderDetails(r)
		r.Dedent()
	}
}

func renderHandlerDetails[T any](
	r *renderer.Renderer,
	h Handler,
	src optional.Optional[T],
) {
	renderEntityDetails(r, h, src)

	p := h.CommonHandlerProperties()

	for _, route := range p.RouteComponents {
		r.IndentBullet()
		route.renderDetails(r)
		r.Dedent()
	}

	r.IndentBullet()
	p.DisabledFlag.renderDetails(r)
	r.Dedent()
}
