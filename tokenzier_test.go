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

			cmtTests := []struct {
				str    string
				tokens []string
			}{
				{"foo\n# Comment", []string{`foo`}},
				{"foo\n# Comment\nbar", []string{`foo`, `bar`}},
				{"foo\n// Comment", []string{`foo`}},
				{"foo\n-- Comment", []string{`foo`}},
				{"foo\n\" Comment", []string{`foo`}},
				{"foo /* Comment */", []string{`foo`}},
				{"foo /* \nComment\n */", []string{`foo`}},
				{"foo <!-- Comment -->", []string{`foo`}},
				{"foo {- Comment -}", []string{`foo`}},
				{"foo (* Comment *)", []string{`foo`}},
				{"2 % 10\n% Comment", []string{`%`}},
				{"foo\n\"\"\"\nComment\n\"\"\"\nbar", []string{`foo`, `bar`}},
				{"foo\n'''\nComment\n'''\nbar", []string{`foo`, `bar`}},
			}

			for _, cmtTest := range cmtTests {
				Expect(ExtractTokens(cmtTest.str)).To(Equal(cmtTest.tokens))
			}

		})

		It("should extract SGML tokens", func() {

			sgmlTests := []struct {
				str    string
				tokens []string
			}{
				{"<html></html>", []string{"<html>", "</html>"}},
				{"<div id></div>", []string{"<div>", "id", "</div>"}},
				{"<div id=foo></div>", []string{"<div>", "id=", "</div>"}},
				{"<div id class></div>", []string{"<div>", "id", "class", "</div>"}},
				{"<div id=\"foo bar\"></div>", []string{"<div>", "id=", "</div>"}},
				{"<div id='foo bar'></div>", []string{"<div>", "id=", "</div>"}},
				{"<?xml version=\"1.0\"?>", []string{"<?xml>", "version="}},
			}

			for _, sgmlTest := range sgmlTests {
				Expect(ExtractTokens(sgmlTest.str)).To(Equal(sgmlTest.tokens))
			}
		})

	})
})
