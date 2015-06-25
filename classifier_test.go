package mockingbird_test

import (
	"github.com/davecgh/go-spew/spew"
	. "github.com/lazywei/mockingbird"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Naive Bayes", func() {
	X, y := ReadLibsvm("test_fixture/samples.libsvm")

	/* _, nFeatures := X.Dims() */
	/* X, _ = X.View(0, 0, 3, nFeatures).(*mat64.Dense) */
	/* y, _ = y.View(0, 0, 3, 1).(*mat64.Dense) */

	nb := NewNaiveBayes()

	Describe("Fit", func() {

		nb.Fit(X, y)
		tokensTotal, langsTotal, langsCount, tokensTotalPerLang, tokenCountPerLang := nb.GetParams()

		It("should count tokens and languages", func() {
			Expect(tokensTotal).To(Equal(2933))
			Expect(langsTotal).To(Equal(4))
		})

		It("should count samples for each languages", func() {
			Expect(langsCount[0]).To(Equal(10))
			Expect(langsCount[1]).To(Equal(10))
			Expect(langsCount[2]).To(Equal(6))
			Expect(langsCount[3]).To(Equal(4))
		})

		It("should count total number of tokens for each languages", func() {
			Expect(tokensTotalPerLang[0]).To(Equal(878))
			Expect(tokensTotalPerLang[1]).To(Equal(1713))
			Expect(tokensTotalPerLang[2]).To(Equal(212))
			Expect(tokensTotalPerLang[3]).To(Equal(130))
		})

		It("should count number of each token for each languages", func() {
			// We only test first 200 "token counts" here
			expected := []int{
				226, 122, 5, 23, 4, 14, 7, 7, 4, 6, 4, 2, 3,
				8, 2, 2, 3, 2, 22, 2, 32, 21, 1, 2, 20, 1,
				1, 1, 13, 4, 4, 1, 2, 2, 5, 1, 1, 1, 1,
				2, 1, 1, 1, 7, 2, 1, 2, 3, 2, 2, 2, 4,
				1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 9, 7,
				8, 2, 20, 3, 5, 3, 3, 4, 1, 2, 4, 1, 1,
				2, 2, 4, 2, 1, 1, 3, 2, 12, 2, 4, 2, 6,
				1, 2, 3, 3, 2, 1, 3, 2, 2, 2, 1, 5, 2,
				2, 4, 2, 2, 2, 3, 2, 2, 3, 3, 3, 1, 1,
				2, 3, 1, 2, 2, 2, 1, 2, 2, 4, 1, 3, 3,
				1, 1, 1, 2, 1, 1, 1, 3, 1, 5, 1, 1, 1,
				1, 2, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
				1, 1, 1, 1, 1, 2, 2, 1, 2, 1, 1, 1, 2,
				1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0}

			for i, expectedVal := range expected {
				Expect(tokenCountPerLang[0][i]).To(Equal(expectedVal))
			}

			// We only test first 100 "token counts" here
			expected = []int{
				0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 19, 0, 0, 7, 0, 0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 0,
				0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}

			for i, expectedVal := range expected {
				Expect(tokenCountPerLang[3][i]).To(Equal(expectedVal))
			}

		})

	})

	Describe("Prediction", func() {
		nb.Fit(X, y)

		It("should predict", func() {
			spew.Dump(nb.Predict(X))
		})

	})
})
