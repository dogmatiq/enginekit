package staticconfig_test

import (
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/dogmatiq/aureus"
	"github.com/dogmatiq/enginekit/config"
	. "github.com/dogmatiq/enginekit/config/staticconfig"
)

func TestAnalyzer(t *testing.T) {
	aureus.Run(
		t,
		func(w io.Writer, in aureus.Content, out aureus.ContentMetaData) error {
			// Make a temporary directory to write the Go source code.
			//
			// The name is static so that the the test output is deterministic.
			//
			// Additionally, creating the directory within the repository allows
			// the test code to use this repo's go.mod file, ensuring the
			// statically analyzed code uses the same versions of Dogma, etc.
			dir := filepath.Join(
				filepath.Dir(in.File),
				"pkg",
			)
			if err := os.Mkdir(dir, 0700); err != nil {
				return err
			}

			defer os.RemoveAll(dir)

			if err := os.WriteFile(
				filepath.Join(dir, "main.go"),
				[]byte(in.Data),
				0600,
			); err != nil {
				return err
			}

			result := LoadAndAnalyze(dir)

			if len(result.Applications) == 0 {
				if _, err := io.WriteString(w, "(no applications found)\n"); err != nil {
					return err
				}
			}

			for err := range result.Errors() {
				if _, err := io.WriteString(w, err.Error()+"\n"); err != nil {
					return err
				}
			}

			for i, app := range result.Applications {
				if i > 0 {
					if _, err := io.WriteString(w, "\n"); err != nil {
						return err
					}
				}

				if _, err := config.WriteDetails(w, app); err != nil {
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
