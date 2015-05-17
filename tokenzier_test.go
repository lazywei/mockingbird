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

	})
})
