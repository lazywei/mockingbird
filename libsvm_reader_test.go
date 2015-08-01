package mockingbird_test

import (
	. "github.com/lazywei/mockingbird"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("LibsvmReader", func() {
	Describe("ReadLibsvm", func() {

		It("should read libsvm format file", func() {
			X, y := ReadLibsvm("test_fixture/test_samples.libsvm", false)

			nSamples, nFeatures := X.Dims()
			Expect(nSamples).To(Equal(22))
			Expect(nFeatures).To(Equal(91))

			nSamplesY, nColsY := y.Dims()
			Expect(nSamplesY).To(Equal(nSamples))
			Expect(nColsY).To(Equal(1))
		})

	})
})
