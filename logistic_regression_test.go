package mockingbird_test

import (
	"github.com/lazywei/liblinear"
	. "github.com/lazywei/mockingbird"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Naive Bayes", func() {
	X, y := liblinear.ReadLibsvm("test_fixture/test_samples.libsvm", false)

	lr := NewLogisticRegression()

	Describe("Prediction", func() {
		lr.Fit(X, y)

		It("should predict", func() {

			expectedPreds := []struct {
				label int
				score float64
			}{

				{0, 0.528771},
				{0, 0.41582},
				{1, 0.53485},
				{1, 0.577725},
				{2, 0.973495},
				{3, 0.351709},
				{4, 0.815347},
				{4, 0.252707},
				{5, 0.427446},
				{5, 0.27975},
				{5, 0.701517},
				{6, 0.275777},
				{7, 0.164584},
				{8, 0.164584},
				{9, 0.328818},
				{9, 0.217835},
				{9, 0.238034},
				{9, 0.238034},
				{10, 0.248419},
				{10, 0.205646},
				{10, 0.248419},
				{10, 0.248419},
			}
			for i, pred := range lr.Predict(X) {
				Expect(pred.Label).To(Equal(expectedPreds[i].label))
				Expect(pred.Score).To(BeNumerically("~", expectedPreds[i].score, 1e5))
			}

		})

	})
})
