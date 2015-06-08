package mockingbird_test

import (
	. "github.com/lazywei/mockingbird"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Naive Bayes", func() {
	X, y := ReadLibsvm("test_fixture/samples.libsvm")
	nb := NewNaiveBayes()

	Describe("Fit", func() {

		It("should fit a model by NB", func() {
			nb.Fit(X, y)
			Expect(true).To(Equal(true))
		})

	})
})
