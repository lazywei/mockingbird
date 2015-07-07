package mockingbird_test

import (
	. "github.com/lazywei/mockingbird"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Naive Bayes", func() {
	X, y := ReadLibsvm("test_fixture/test_samples.libsvm")

	/* _, nFeatures := X.Dims() */
	/* X, _ = X.View(0, 0, 3, nFeatures).(*mat64.Dense) */
	/* y, _ = y.View(0, 0, 3, 1).(*mat64.Dense) */

	nb := NewNaiveBayes()

	Describe("Fit", func() {

		nb.Fit(X, y)
		tokensTotal, langsTotal, langsCount, tokensTotalPerLang, tokenCountPerLang := nb.GetParams()

		It("should count tokens and languages", func() {
			Expect(tokensTotal).To(Equal(238))
			Expect(langsTotal).To(Equal(22))
		})

		It("should count samples for each languages", func() {
			Expect(langsCount).To(Equal(map[int]int{
				0: 2, 1: 2, 2: 1,
				3: 1, 4: 2, 5: 3,
				6: 1, 7: 1, 8: 1,
				9: 4, 10: 4}))
		})

		It("should count total number of tokens for each languages", func() {
			Expect(tokensTotalPerLang).To(Equal(map[int]int{
				0: 24, 1: 31, 2: 96,
				3: 6, 4: 20, 5: 33,
				6: 5, 7: 2, 8: 2,
				9: 12, 10: 7,
			}))
		})

		It("should count number of each token for each languages", func() {
			expected := []int{
				1, 1, 1, 1, 3, 3, 1, 1, 3, 1, 1, 1, 2, 1, 1, 1, 1, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0}

			for i, expectedVal := range expected {
				Expect(tokenCountPerLang[0][i]).To(Equal(expectedVal))
			}

			expected = []int{
				0, 0, 2, 1, 2, 2, 2, 0, 2, 1, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0,
				0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0,
				1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 3, 1, 1,
				2, 1, 2, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0}

			for i, expectedVal := range expected {
				Expect(tokenCountPerLang[5][i]).To(Equal(expectedVal))
			}

		})

	})

	Describe("Prediction", func() {
		nb.Fit(X, y)

		It("should predict", func() {

			expectedPreds := []struct {
				label int
				score float64
			}{
				{0, -40.298975165660956},
				{0, -29.496302349153602},
				{1, -44.92853869111085},
				{1, -53.84420594344641},
				{2, -291.5183688683531},
				{3, -11.069010546486865},
				{4, -35.37466680850959},
				{4, -10.468801361586188},
				{5, -26.093289645514158},
				{5, -13.493553760768126},
				{5, -67.38900486121871},
				{6, -11.138232015528818},
				{7, -4.477336814478206},
				{8, -4.477336814478206},
				{9, -12.337521871950372},
				{9, -8.466320861042481},
				{9, -4.882801922586371},
				{9, -4.882801922586371},
				{10, -4.210274029229161},
				{10, -2.264363880173848},
				{10, -4.210274029229161},
				{10, -4.210274029229161},
			}
			for i, pred := range nb.Predict(X) {
				Expect(pred.Label).To(Equal(expectedPreds[i].label))
				Expect(pred.Score).To(BeNumerically("~", expectedPreds[i].score))
			}

		})

	})
})
