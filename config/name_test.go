package config_test

import (
	. "github.com/dogmatiq/enginekit/config"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ HandlerConfig = &IntegrationConfig{}

var _ = Describe("func IsValidName", func() {
	DescribeTable(
		"it returns true",
		func(n string) {
			Expect(IsValidName(n)).To(BeTrue())
		},
		Entry("ascii", "foo-bar"),
		Entry("unicode", "😀"),
	)

	DescribeTable(
		"it returns false",
		func(n string) {
			Expect(IsValidName(n)).To(BeFalse())
		},
		Entry("empty", ""),
		Entry("non-printable ascii character", "\n"),
		Entry("non-printable unicode character", "\u200B"),
	)
})
