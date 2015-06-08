package mockingbird

type Classifier interface {
	Fit()
	Predict()
}

type NaiveBayes struct {
}

func NewNaiveBayes() *NaiveBayes {
	return &NaiveBayes{}
}

func (nb *NaiveBayes) Fit() {
}
