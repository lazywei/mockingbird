package mockingbird

import "github.com/gonum/matrix/mat64"

type Classifier interface {
	Fit()
	Predict()
}

type NaiveBayes struct {
	tokensTotal        int
	langsTotal         int
	langsCount         map[int]int
	tokensTotalPerLang map[int]int
	tokenCountPerLang  map[int](map[int]int)
}

func NewNaiveBayes() *NaiveBayes {
	return &NaiveBayes{}
}

func (nb *NaiveBayes) Fit(X, y *mat64.Dense) {
	nSamples, nFeatures := X.Dims()

	tokensTotal := 0
	langsTotal := len(uniqVals(y.Col(nil, 0)))

	langsCount := histogram(y.Col(nil, 0))

	tokensTotalPerLang := map[int]int{}
	tokenCountPerLang := map[int](map[int]int){}

	for i := 0; i < nSamples; i++ {
		langIdx := int(y.At(i, 0))

		for j := 0; j < nFeatures; j++ {
			tokensTotal += int(X.At(i, j))
			tokensTotalPerLang[langIdx] += int(X.At(i, j))

			if _, ok := tokenCountPerLang[langIdx]; !ok {
				tokenCountPerLang[langIdx] = map[int]int{}
			}
			tokenCountPerLang[langIdx][j] += int(X.At(i, j))
		}

		// if _, ok := tokenCountPerLang[langIdx]; !ok {
		// 	tokenCountPerLang[langIdx] = mat64.NewVector(nFeatures, nil)
		// }

		// tokenCountPerLang[langIdx].AddVec(
		// 	tokenCountPerLang[langIdx],
		// 	X.RowView(0),
		// )
	}

	nb.tokensTotal = tokensTotal
	nb.langsTotal = langsTotal
	nb.langsCount = langsCount
	nb.tokensTotalPerLang = tokensTotalPerLang
	nb.tokenCountPerLang = tokenCountPerLang
}

func (nb *NaiveBayes) GetParams() (
	int, int, map[int]int, map[int]int, map[int](map[int]int)) {

	return nb.tokensTotal, nb.langsTotal,
		nb.langsCount, nb.tokensTotalPerLang, nb.tokenCountPerLang
}

func histogram(data_arr []float64) map[int]int {
	results := map[int]int{}

	for _, val := range data_arr {
		results[int(val)] += 1
	}

	return results
}

func uniqVals(data_arr []float64) []float64 {
	results := map[float64]bool{}

	for _, val := range data_arr {
		if _, ok := results[val]; !ok {
			results[val] = true
		}
	}

	uniqResults := []float64{}

	for k := range results {
		uniqResults = append(uniqResults, k)
	}

	return uniqResults
}

func tokenProba() float64 {
	return 0
}

func langProba() float64 {
	return 0
}
