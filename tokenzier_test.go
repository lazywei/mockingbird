package mockingbird_test

import (
	. "github.com/lazywei/mockingbird"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Tokenzier", func() {
	Describe("extract tokens", func() {

		It("should skip string literals", func() {
			expectedResult := []string{"print"}

			Expect(ExtractTokens(`print ""`)).To(Equal(expectedResult))
			Expect(ExtractTokens(`print "Josh"`)).To(Equal(expectedResult))
			Expect(ExtractTokens(`print 'Josh'`)).To(Equal(expectedResult))
			Expect(ExtractTokens(`print "Hello \"Josh\""`)).To(Equal(expectedResult))
			Expect(ExtractTokens(`print 'Hello \'Josh\''`)).To(Equal(expectedResult))
			Expect(ExtractTokens(`print "Hello", "Josh"`)).To(Equal(expectedResult))
			Expect(ExtractTokens(`print 'Hello', 'Josh'`)).To(Equal(expectedResult))
			Expect(ExtractTokens(`print "Hello", "", "Josh"`)).To(Equal(expectedResult))
			Expect(ExtractTokens(`print 'Hello', '', 'Josh'`)).To(Equal(expectedResult))
		})

		It("should skip number literals", func() {
			Expect(ExtractTokens(`1 + 1`)).To(Equal([]string{`+`}))
			Expect(ExtractTokens(`add(123, 45)`)).To(Equal([]string{`add`, `(`, `)`}))
			Expect(ExtractTokens(`0x01 | 0x10`)).To(Equal([]string{`|`}))
			Expect(ExtractTokens(`500.42 * 1.0`)).To(Equal([]string{`*`}))
		})

		It("should extract common operators", func() {
			Expect(ExtractTokens("1 + 1")).To(Equal([]string{`+`}))
			Expect(ExtractTokens("1 - 1")).To(Equal([]string{`-`}))
			Expect(ExtractTokens("1 * 1")).To(Equal([]string{`*`}))
			Expect(ExtractTokens("1 / 1")).To(Equal([]string{`/`}))
			Expect(ExtractTokens("2 % 5")).To(Equal([]string{`%`}))
			Expect(ExtractTokens("1 & 1")).To(Equal([]string{`&`}))
			Expect(ExtractTokens("1 && 1")).To(Equal([]string{`&&`}))
			Expect(ExtractTokens("1 | 1")).To(Equal([]string{`|`}))
			Expect(ExtractTokens("1 || 1")).To(Equal([]string{`||`}))
			Expect(ExtractTokens("1 < 0x01")).To(Equal([]string{`<`}))
			Expect(ExtractTokens("1 << 0x01")).To(Equal([]string{`<<`}))
		})

		It("should skip comments", func() {

			tokens := [][]string{
				[]string{`foo`},
				[]string{`foo`, `bar`},
				[]string{`foo`},
				[]string{`foo`},
				[]string{`foo`},
				[]string{`foo`},
				[]string{`foo`},
				[]string{`foo`},
				[]string{`foo`},
				[]string{`foo`},
				[]string{`%`},
				[]string{`foo`, `bar`},
				[]string{`foo`, `bar`},
			}

			strs := []string{
				"foo\n# Comment",
				"foo\n# Comment\nbar",
				"foo\n// Comment",
				"foo\n-- Comment",
				"foo\n\" Comment",
				"foo /* Comment */",
				"foo /* \nComment\n */",
				"foo <!-- Comment -->",
				"foo {- Comment -}",
				"foo (* Comment *)",
				"2 % 10\n% Comment",
				"foo\n\"\"\"\nComment\n\"\"\"\nbar",
				"foo\n'''\nComment\n'''\nbar",
			}

			for i, str := range strs {
				Expect(ExtractTokens(str)).To(Equal(tokens[i]))
			}

		})

	})
})
