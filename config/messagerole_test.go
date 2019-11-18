package config_test

import (
	. "github.com/dogmatiq/enginekit/config"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("type MessageRole", func() {
	Describe("func MustValidate", func() {
		It("does not panic when the role is valid", func() {
			CommandMessageRole.MustValidate()
			EventMessageRole.MustValidate()
			TimeoutMessageRole.MustValidate()
		})

		It("panics when the role is not valid", func() {
			Expect(func() {
				MessageRole(0).MustValidate()
			}).To(Panic())
		})
	})

	Describe("func Is", func() {
		It("returns true when the role is in the given set", func() {
			Expect(CommandMessageRole.Is(CommandMessageRole, EventMessageRole)).To(BeTrue())
		})

		It("returns false when the role is not in the given set", func() {
			Expect(TimeoutMessageRole.Is(CommandMessageRole, EventMessageRole)).To(BeFalse())
		})
	})

	Describe("func MustBe", func() {
		It("does not panic when the role is in the given set", func() {
			CommandMessageRole.MustBe(CommandMessageRole, EventMessageRole)
		})

		It("panics when the role is not in the given set", func() {
			Expect(func() {
				TimeoutMessageRole.MustBe(CommandMessageRole, EventMessageRole)
			}).To(Panic())
		})
	})

	Describe("func MustNotBe", func() {
		It("does not panic when the role is not in the given set", func() {
			TimeoutMessageRole.MustNotBe(CommandMessageRole, EventMessageRole)
		})

		It("panics when the role is in the given set", func() {
			Expect(func() {
				CommandMessageRole.MustNotBe(CommandMessageRole, EventMessageRole)
			}).To(Panic())
		})
	})

	Describe("func Marker", func() {
		It("returns the correct marker character", func() {
			Expect(CommandMessageRole.Marker()).To(Equal("?"))
			Expect(EventMessageRole.Marker()).To(Equal("!"))
			Expect(TimeoutMessageRole.Marker()).To(Equal("@"))
		})
	})

	Describe("func ShortString()", func() {
		It("returns the role value as a short string", func() {
			Expect(CommandMessageRole.ShortString()).To(Equal("CMD"))
			Expect(EventMessageRole.ShortString()).To(Equal("EVT"))
			Expect(TimeoutMessageRole.ShortString()).To(Equal("TMO"))
		})
	})

	Describe("func String", func() {
		It("returns the role value as a string", func() {
			Expect(CommandMessageRole.String()).To(Equal("command"))
			Expect(EventMessageRole.String()).To(Equal("event"))
			Expect(TimeoutMessageRole.String()).To(Equal("timeout"))
			Expect(MessageRole(0).String()).To(Equal("<invalid message role 0x0>"))
		})
	})

	Describe("func MarshalText()", func() {
		It("marshals the role to text", func() {
			Expect(CommandMessageRole.MarshalText()).To(Equal([]byte("command")))
			Expect(EventMessageRole.MarshalText()).To(Equal([]byte("event")))
			Expect(TimeoutMessageRole.MarshalText()).To(Equal([]byte("timeout")))
		})

		It("returns an error if the role is invalid", func() {
			_, err := MessageRole(0).MarshalText()
			Expect(err).Should(HaveOccurred())
		})
	})

	Describe("func UnmarshalText()", func() {
		It("unmarshals the role from text", func() {
			var r MessageRole

			err := r.UnmarshalText([]byte("command"))
			Expect(err).ShouldNot(HaveOccurred())
			Expect(r).To(Equal(CommandMessageRole))

			err = r.UnmarshalText([]byte("event"))
			Expect(err).ShouldNot(HaveOccurred())
			Expect(r).To(Equal(EventMessageRole))

			err = r.UnmarshalText([]byte("timeout"))
			Expect(err).ShouldNot(HaveOccurred())
			Expect(r).To(Equal(TimeoutMessageRole))
		})

		It("returns an error if the data is invalid", func() {
			var r MessageRole

			err := r.UnmarshalText([]byte("<invalid>"))
			Expect(err).Should(HaveOccurred())
		})
	})

	Describe("func MarshalBinary()", func() {
		It("marshals the role to binary", func() {
			Expect(CommandMessageRole.MarshalBinary()).To(Equal([]byte("C")))
			Expect(EventMessageRole.MarshalBinary()).To(Equal([]byte("E")))
			Expect(TimeoutMessageRole.MarshalBinary()).To(Equal([]byte("T")))
		})

		It("returns an error if the role is invalid", func() {
			_, err := MessageRole(0).MarshalBinary()
			Expect(err).Should(HaveOccurred())
		})
	})

	Describe("func UnmarshalBinary()", func() {
		It("unmarshals the role from binary", func() {
			var r MessageRole

			err := r.UnmarshalBinary([]byte("C"))
			Expect(err).ShouldNot(HaveOccurred())
			Expect(r).To(Equal(CommandMessageRole))

			err = r.UnmarshalBinary([]byte("E"))
			Expect(err).ShouldNot(HaveOccurred())
			Expect(r).To(Equal(EventMessageRole))

			err = r.UnmarshalBinary([]byte("T"))
			Expect(err).ShouldNot(HaveOccurred())
			Expect(r).To(Equal(TimeoutMessageRole))
		})

		It("returns an error if the data is the wrong length", func() {
			var r MessageRole

			err := r.UnmarshalBinary([]byte("<invalid>"))
			Expect(err).Should(HaveOccurred())
		})

		It("returns an error if the data does not contain a valid role", func() {
			var r MessageRole

			err := r.UnmarshalBinary([]byte("X"))
			Expect(err).Should(HaveOccurred())
		})
	})
})
