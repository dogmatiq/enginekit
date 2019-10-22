package config_test

import (
	"fmt"

	. "github.com/dogmatiq/enginekit/config"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("type Identity", func() {
	Describe("func Validate()", func() {
		DescribeTable(
			"it returns nil if the name and key are valid",
			func(v string) {
				i := Identity{v, v}
				Expect(i.Validate()).ShouldNot(HaveOccurred())
			},
			Entry("ascii", "foo-bar"),
			Entry("unicode", "😀"),
		)

		invalidEntries := []TableEntry{
			Entry("empty", ""),
			Entry("non-printable ascii character (newline)", "\n"),
			Entry("non-printable ascii character (space)", " "),
			Entry("non-printable unicode character", "\u200B"),
		}

		DescribeTable(
			"it returns an error if the name is invalid",
			func(v string) {
				i := Identity{v, "<key>"}
				Expect(i.Validate()).Should(MatchError(
					fmt.Sprintf(
						"invalid name %#v, names must be non-empty, printable UTF-8 strings with no whitespace",
						v,
					),
				))
			},
			invalidEntries...,
		)

		DescribeTable(
			"it returns an error if the key is invalid",
			func(v string) {
				i := Identity{"<name>", v}
				Expect(i.Validate()).To(MatchError(
					fmt.Sprintf(
						"invalid key %#v, keys must be non-empty, printable UTF-8 strings with no whitespace",
						v,
					),
				))
			},
			invalidEntries...,
		)
	})
})
