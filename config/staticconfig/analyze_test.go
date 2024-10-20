package staticconfig_test

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/dogmatiq/aureus"
	"github.com/dogmatiq/enginekit/config"
	. "github.com/dogmatiq/enginekit/config/staticconfig"
)

func TestAnalyzer(t *testing.T) {
	aureus.Run(
		t,
		func(t *testing.T, in aureus.Input, out aureus.Output) error {
			t.Parallel()

			cwd, _ := os.Getwd()

			// Create a temporary directory to write the Go source code, but
			// create it within this Go module so that it uses the same version
			// of Dogma, etc.
			dir, err := os.MkdirTemp(
				filepath.Join(cwd, "testdata"),
				"aureus-",
			)
			if err != nil {
				return err
			}
			defer os.RemoveAll(dir)

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

			for err := range result.Errors() {
				message := err.Error()
				message = strings.ReplaceAll(message, dir+"/", "")
				message = strings.ReplaceAll(message, dir, "")

				if _, err := io.WriteString(out, "ERROR: "+message+"\n"); err != nil {
					return err
				}
			}

			if len(result.Applications) == 0 {
				if _, err := io.WriteString(out, "(no applications found)\n"); err != nil {
					return err
				}
			}

			for i, app := range result.Applications {
				if i > 0 {
					if _, err := io.WriteString(out, "\n"); err != nil {
						return err
					}
				}

				// Render the details of the application.
				details := config.RenderDetails(app)

				// Remove the random portion of the temporary directory name
				// so that the test output is deterministic.
				details = strings.ReplaceAll(
					details,
					"/"+filepath.Base(dir)+".",
					".",
				)

				if _, err := io.WriteString(out, details); err != nil {
					return err
				}
			}

			return nil
		},
	)

	// t.Run("should parse multiple packages contain applications", func(t *testing.T) {
	// 	apps := FromDir("testdata/multiple-apps-in-pkgs")

	// 	if len(apps) != 2 {
	// 		t.Fatalf("expected 2 applications, got %d", len(apps))
	// 	}

	// 	if expected, actual := "<app-first>",
	// 		apps[0].Identity().Name; expected != actual {
	// 		t.Fatalf(
	// 			"unexpected application name: want %s, got %s",
	// 			expected,
	// 			actual,
	// 		)
	// 	}

	// 	if expected, actual := "b754902b-47c8-48fc-84d2-d920c9cbdaec",
	// 		apps[0].Identity().Key; expected != actual {
	// 		t.Fatalf(
	// 			"unexpected application key: want %s, got %s",
	// 			expected,
	// 			actual,
	// 		)
	// 	}

	// 	if expected, actual := "<app-second>",
	// 		apps[1].Identity().Name; expected != actual {
	// 		t.Fatalf(
	// 			"unexpected application name: want %s, got %s",
	// 			expected,
	// 			actual,
	// 		)
	// 	}

	// 	if expected, actual := "bfaf2a16-23a0-495d-8098-051d77635822",
	// 		apps[1].Identity().Key; expected != actual {
	// 		t.Fatalf(
	// 			"unexpected application key: want %s, got %s",
	// 			expected,
	// 			actual,
	// 		)
	// 	}
	// })

	// t.Run("should parse all application-level messages", func(t *testing.T) {
	// 	apps := FromDir("testdata/app-level-messages")

	// 	if len(apps) != 1 {
	// 		t.Fatalf("expected 1 application, got %d", len(apps))
	// 	}

	// 	contains := func(
	// 		mn message.Name,
	// 		mk message.Kind,
	// 		iterator iter.Seq2[message.Name, message.Kind],
	// 	) bool {
	// 		for k, v := range iterator {
	// 			if k == mn && v == mk {
	// 				return true
	// 			}
	// 		}
	// 		return false
	// 	}

	// 	if !contains(
	// 		message.NameFor[CommandStub[TypeA]](),
	// 		message.CommandKind,
	// 		apps[0].MessageNames().Consumed(),
	// 	) {
	// 		t.Fatal("expected consumed TypeA command message")
	// 	}

	// 	if !contains(
	// 		message.NameFor[EventStub[TypeA]](),
	// 		message.EventKind,
	// 		apps[0].MessageNames().Consumed(),
	// 	) {
	// 		t.Fatal("expected consumed TypeA event message")
	// 	}

	// 	if !contains(
	// 		message.NameFor[EventStub[TypeC]](),
	// 		message.EventKind,
	// 		apps[0].MessageNames().Consumed(),
	// 	) {
	// 		t.Fatal("expected consumed TypeC event message")
	// 	}

	// 	if !contains(
	// 		message.NameFor[TimeoutStub[TypeA]](),
	// 		message.TimeoutKind,
	// 		apps[0].MessageNames().Consumed(),
	// 	) {
	// 		t.Fatal("expected consumed TypeA timeout message")
	// 	}

	// 	if !contains(
	// 		message.NameFor[EventStub[TypeA]](),
	// 		message.EventKind,
	// 		apps[0].MessageNames().Produced(),
	// 	) {
	// 		t.Fatal("expected produced TypeA event message")
	// 	}

	// 	if !contains(
	// 		message.NameFor[CommandStub[TypeB]](),
	// 		message.CommandKind,
	// 		apps[0].MessageNames().Produced(),
	// 	) {
	// 		t.Fatal("expected produced TypeB command message")
	// 	}

	// 	if !contains(
	// 		message.NameFor[TimeoutStub[TypeA]](),
	// 		message.TimeoutKind,
	// 		apps[0].MessageNames().Produced(),
	// 	) {
	// 		t.Fatal("expected produced TypeA timeout message")
	// 	}

	// 	if !contains(
	// 		message.NameFor[EventStub[TypeB]](),
	// 		message.EventKind,
	// 		apps[0].MessageNames().Produced(),
	// 	) {
	// 		t.Fatal("expected produced TypeB event message")
	// 	}
	// })
}
