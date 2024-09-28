package json

import (
	"reflect"
	"strings"
)

// portableName returns the portable name for the given type.
func portableName(rt reflect.Type) (string, bool) {
	n := rt.Name()

	if n == "" {
		for rt.Kind() == reflect.Ptr {
			rt = rt.Elem()
		}

		n = rt.Name()
	}

	if n == "" {
		return "", false
	}

	if n[len(n)-1] != ']' {
		return n, true
	}

	bracket := strings.IndexByte(n, '[')

	var w strings.Builder
	w.Grow(len(n))
	w.WriteString(n[:bracket+1])

	params := n[bracket+1 : len(n)-1]

	for i, p := range strings.Split(params, ",") {
		if i > 0 {
			w.WriteByte(',')
		}

		if dot := strings.LastIndexByte(p, '.'); dot == -1 {
			w.WriteString(p)
		} else {
			w.WriteString(p[dot+1:])
		}
	}

	w.WriteByte(']')

	return w.String(), true
}
