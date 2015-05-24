package scanner_test

import (
	. "github.com/lazywei/mockingbird/scanner"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Scanner", func() {
	Describe("Scan", func() {

		Context("when there is matched string", func() {
			It("should return the string, and ok=true", func() {
				s := NewScanner("This is an example.")
				rtn, ok := s.Scan(`This`)
				Expect(rtn).To(Equal("This"))
				Expect(ok).To(Equal(true))
			})

			It("should scan forward for next Scan", func() {
				s := NewScanner("This is an example.")
				s.Scan(`This`)
				rtn, ok := s.Scan(`\s`)

				Expect(rtn).To(Equal(" "))
				Expect(ok).To(Equal(true))

				rtn, ok = s.Scan(`is`)

				Expect(rtn).To(Equal("is"))
				Expect(ok).To(Equal(true))
			})
		})

		Context("when there is no matched string", func() {
			It("should return empty string, and ok=false", func() {
				s := NewScanner("This is an example.")
				rtn, ok := s.Scan(`is`)
				Expect(rtn).To(Equal(""))
				Expect(ok).To(Equal(false))
			})
		})

	})

	Describe("ScanUntil", func() {

		Context("when there is matched string", func() {
			It("should return the whole string until the match, and ok=true", func() {
				s := NewScanner("This is an example.")
				rtn, ok := s.ScanUntil(`an`)
				Expect(rtn).To(Equal("This is an"))
				Expect(ok).To(Equal(true))
			})

			It("should scan forward for next ScanUntil", func() {
				s := NewScanner("This is an example.")
				s.ScanUntil(`an`)
				rtn, ok := s.ScanUntil(`amp`)

				Expect(rtn).To(Equal(" examp"))
				Expect(ok).To(Equal(true))
			})
		})

	})

	Describe("Getch", func() {

		Context("when not yet EOF", func() {
			It("should return the current char, ok=true, and move forward", func() {
				s := NewScanner("foo")
				rtn, ok := s.Getch()
				Expect(rtn).To(Equal("f"))
				Expect(ok).To(Equal(true))

				rtn, ok = s.Getch()
				Expect(rtn).To(Equal("o"))
				Expect(ok).To(Equal(true))

				rtn, ok = s.Getch()
				Expect(rtn).To(Equal("o"))
				Expect(ok).To(Equal(true))
			})
		})

		Context("when at EOF", func() {
			It("should return empty string, ok=false", func() {
				s := NewScanner("f")
				s.Getch()

				rtn, ok := s.Getch()
				Expect(rtn).To(Equal(""))
				Expect(ok).To(Equal(false))
			})
		})

	})

	Describe("Peek", func() {

		Context("when not yet EOF", func() {
			It("should return the next n-length string", func() {
				s := NewScanner("foo")
				rtn := s.Peek(1)
				Expect(rtn).To(Equal("f"))

				rtn = s.Peek(2)
				Expect(rtn).To(Equal("fo"))

				s.Getch()
				rtn = s.Peek(2)
				Expect(rtn).To(Equal("oo"))
			})
		})

		Context("when length exceeds the limit", func() {
			It("should return the string until the end", func() {
				s := NewScanner("foo")
				rtn := s.Peek(100)
				Expect(rtn).To(Equal("foo"))
			})
		})

		Context("when at EOF", func() {
			It("should return empty string", func() {
				s := NewScanner("f")
				s.Getch()

				rtn := s.Peek(1)
				Expect(rtn).To(Equal(""))
			})
		})

	})

	Describe("SkipUntil", func() {

		Context("when there is matched string", func() {
			It("should return number of bytes and move to the match, and ok=true", func() {
				s := NewScanner("This is an example.")
				rtn, ok := s.SkipUntil(`an`)
				Expect(rtn).To(Equal(10))
				Expect(ok).To(Equal(true))
			})

			It("should scan forward for next SkipUntil", func() {
				s := NewScanner("This is an example.")
				s.SkipUntil(`an`)
				rtn, ok := s.SkipUntil(`amp`)

				Expect(rtn).To(Equal(6))
				Expect(ok).To(Equal(true))
			})
		})

	})
})
