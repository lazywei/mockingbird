package mockingbird

import (
	"log"
	"math"

	"github.com/gonum/matrix/mat64"
	"gopkg.in/yaml.v2"
)

type Prediction struct {
	Label    int
	Language string
	Score    float64
}

type Classifier interface {
	Fit()
	Predict()
}

type NaiveBayes struct {
	params nbParams
}

type nbParams struct {
	TokensTotal        int
	LangsTotal         int
	LangsCount         map[int]int
	TokensTotalPerLang map[int]int
	TokenCountPerLang  map[int](map[int]int)
}

func NewNaiveBayes() *NaiveBayes {
	return &NaiveBayes{}
}

func (nb *NaiveBayes) Fit(X, y *mat64.Dense) {
	nSamples, nFeatures := X.Dims()

	tokensTotal := 0
	langsTotal, _ := y.Dims()

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
	}

	params := nbParams{
		TokensTotal:        tokensTotal,
		LangsTotal:         langsTotal,
		LangsCount:         langsCount,
		TokensTotalPerLang: tokensTotalPerLang,
		TokenCountPerLang:  tokenCountPerLang,
	}

	nb.params = params
}

func (nb *NaiveBayes) GetParams() (
	tokensTotal int,
	langsTotal int,
	langsCount map[int]int,
	tokensTotalPerLang map[int]int,
	tokenCountPerLang map[int](map[int]int)) {

	tokensTotal = nb.params.TokensTotal
	langsTotal = nb.params.LangsTotal
	langsCount = nb.params.LangsCount
	tokensTotalPerLang = nb.params.TokensTotalPerLang
	tokenCountPerLang = nb.params.TokenCountPerLang

	return
}

func (nb *NaiveBayes) ToYaml() string {
	params := nb.params
	d, err := yaml.Marshal(&params)
	if err != nil {
		log.Fatal(err)
	}
	return string(d)
}

func FromYaml(yamlStr string) *NaiveBayes {
	params := nbParams{}
	err := yaml.Unmarshal([]byte(yamlStr), &params)
	if err != nil {
		log.Fatal(err)
	}
	nb := NewNaiveBayes()
	nb.params = params
	return nb
}

func (nb *NaiveBayes) Predict(X *mat64.Dense) []Prediction {
	nSamples, _ := X.Dims()

	prediction := []Prediction{}

	for i := 0; i < nSamples; i++ {
		scores := map[int]float64{}
		for langIdx, _ := range nb.params.LangsCount {
			scores[langIdx] = nb.tokensProba(X.Row(nil, i), langIdx) + nb.langProba(langIdx)
		}

		bestScore := scores[0]
		bestLangIdx := 0

		for langIdx, score := range scores {
			if score > bestScore {
				bestScore = score
				bestLangIdx = langIdx
			}
		}

		prediction = append(prediction, Prediction{
			Label:    bestLangIdx,
			Language: "TODO: PENDING",
			Score:    bestScore,
		})
	}

	return prediction
}

func (nb *NaiveBayes) tokensProba(dataArr []float64, langIdx int) float64 {
	result := 0.0

	for tokenIdx, nTokens := range dataArr {
		// Equivalent to:
		// for i = 0 to nTokens
		//     result += log(tokenProba(tokenIdx, langIdx))
		result = result + math.Log(
			math.Pow(nb.tokenProba(tokenIdx, langIdx), nTokens))
	}

	return result
}

func (nb *NaiveBayes) tokenProba(tokenIdx int, langIdx int) float64 {
	tokenCount, ok := nb.params.TokenCountPerLang[langIdx][tokenIdx]
	proba := 0.0
	if tokenCount > 0 && ok {
		proba = float64(tokenCount) / float64(nb.params.TokensTotalPerLang[langIdx])
	} else {
		proba = 1.0 / float64(nb.params.TokensTotal)
	}
	return proba
}

func (nb *NaiveBayes) langProba(langIdx int) float64 {
	return math.Log(float64(nb.params.LangsCount[langIdx]) / float64(nb.params.LangsTotal))
}

func histogram(dataArr []float64) map[int]int {
	results := map[int]int{}

	for _, val := range dataArr {
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
