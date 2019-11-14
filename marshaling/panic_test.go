package marshaling_test

import (
	"errors"

	. "github.com/dogmatiq/enginekit/marshaling"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("func Recover", func() {
	It("assigns marshaling errors to the given pointer", func() {
		err := func() (err error) {
			defer Recover(&err)

			panic(
				PanicSentinel{
					Cause: errors.New("<error>"),
				},
			)
		}()

		Expect(err).To(MatchError("<error>"))
	})

	It("does not recover from unrelated panics", func() {
		var value interface{}

		func() {
			defer func() {
				value = recover()
			}()

			func() (err error) {
				defer Recover(&err)
				panic("<value>") // not a MustPanicSentinel
			}()
		}()

		Expect(value).To(Equal("<value>"))
	})

	It("does not panic when no panic occurs", func() {
		err := func() (err error) {
			defer Recover(&err)
			return nil
		}()

		Expect(err).ShouldNot(HaveOccurred())
	})

	It("panics when passed a nil pointer", func() {
		Expect(func() {
			Recover(nil)
		}).To(Panic())
	})
})
