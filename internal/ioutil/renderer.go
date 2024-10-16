package ioutil

import (
	"fmt"
	"io"
	"strings"
)

type Renderer struct {
	Target io.Writer

	count int
	err   error

	depth        int
	bullet       bool
	continuation bool
}

func (r *Renderer) Indent() {
	r.depth++
}

func (r *Renderer) IndentBullet() {
	r.depth++
	r.bullet = true
}

func (r *Renderer) Dedent() {
	r.depth--
	r.bullet = false
}

func (r *Renderer) Print(str ...string) {
	for _, s := range str {
		r.print(s)
	}
}

func (r *Renderer) Printf(format string, args ...any) {
	r.print(fmt.Sprintf(format, args...))
}

func (r *Renderer) Done() (int, error) {
	return r.count, r.err
}

func (r *Renderer) print(s string) {
	for s != "" && r.err == nil {
		if !r.continuation {
			if r.depth > 0 {
				for range r.depth - 1 {
					r.write("  ")
				}

				if r.bullet {
					r.write("• ")
					r.bullet = false
				} else {
					r.write("  ")
				}
			}

			r.continuation = true
		}

		i := strings.IndexByte(s, '\n')

		if i == -1 {
			// There are no more line break characters, simply write the
			// remainder of the buffer and we're done.
			r.write(s)
			return
		}

		// Write the remainder of this line including the line break character,
		// and trim the beginning of the buffer.
		r.write(s[:i+1])
		s = s[i+1:]
		r.continuation = false
	}
}

func (r *Renderer) write(s string) {
	if r.err == nil {
		n, err := io.WriteString(r.Target, s)
		r.count += n
		r.err = err
	}
}
