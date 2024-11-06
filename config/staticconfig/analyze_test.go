package staticconfig_test

import (
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/dogmatiq/aureus"
	"github.com/dogmatiq/enginekit/config"
	. "github.com/dogmatiq/enginekit/config/staticconfig"
)

func TestAnalyzer(t *testing.T) {
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	// Create a single directory for the Go source code used as Aureus test
	// inputs.
	//
	// Since it's under the testdata directory it is ignored by Go's tooling,
	// but it is still subject to the same go.mod file, and hence the same
	// version of Dogma, etc.
	outputDir := filepath.Join(
		cwd,
		"testdata",
		".aureus",
		strconv.Itoa(os.Getpid()),
	)
	if err := os.MkdirAll(outputDir, 0700); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		os.RemoveAll(outputDir)
	})

	aureus.Run(
		t,
		func(t *testing.T, in aureus.Input, out aureus.Output) error {
			t.Parallel()

			dir, err := os.MkdirTemp(outputDir, "aureus-")
			if err != nil {
				return err
			}
			t.Cleanup(func() {
				os.RemoveAll(dir)
			})

			f, err := os.Create(filepath.Join(dir, "main.go"))
			if err != nil {
				return err
			}
			defer f.Close()

			if _, err := io.Copy(f, in); err != nil {
				return err
			}

			if err := f.Close(); err != nil {
				return err
			}

			result := LoadAndAnalyze(dir)

			hasErrors := false
			for err := range result.Errors() {
				hasErrors = true

				message := err.Error()
				message = strings.ReplaceAll(message, dir+"/", "")
				message = strings.ReplaceAll(message, dir, "")

				if _, err := io.WriteString(out, message+"\n"); err != nil {
					return err
				}
			}

			if !hasErrors && len(result.Applications) == 0 {
				_, err := io.WriteString(out, "(no applications found)\n")
				return err
			}

			for i, app := range result.Applications {
				if i > 0 {
					if _, err := io.WriteString(out, "\n"); err != nil {
						return err
					}
				}

				// Render the details of the application.
				err := config.Validate(app)
				desc := config.Description(
					app,
					config.WithValidationResult(err),
				)

				// Remove the random portion of the temporary directory name
				// so that the test output is deterministic.
				rel, _ := filepath.Rel(cwd, dir)
				desc = strings.ReplaceAll(
					desc,
					"/"+rel+".",
					".",
				)

				if _, err := io.WriteString(out, desc); err != nil {
					return err
				}
			}

			return nil
		},
	)
}
