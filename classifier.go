package mockingbird

import "github.com/gonum/matrix/mat64"

type Classifier interface {
	Fit()
	Predict()
}

type NaiveBayes struct {
}

func NewNaiveBayes() *NaiveBayes {
	return &NaiveBayes{}
}

func (nb *NaiveBayes) Fit(X, y *mat64.Dense) {
}
