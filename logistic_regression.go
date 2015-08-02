package mockingbird

import (
	"github.com/gonum/matrix/mat64"
	"github.com/lazywei/liblinear"
)

// type Prediction struct {
// 	Label    int
// 	Language string
// 	Score    float64
// }

type LogisticRegression struct {
	model *liblinear.Model
}

func NewLogisticRegression() *LogisticRegression {
	return &LogisticRegression{}
}

func (lr *LogisticRegression) Fit(X, y *mat64.Dense) {
	model := liblinear.Train(X, y, 1, liblinear.L2R_LR, 0.01, 1, 0.1, nil)
	lr.model = model
}

func (lr *LogisticRegression) Predict(X *mat64.Dense) []Prediction {
	nSamples, _ := X.Dims()

	prediction := []Prediction{}

	for i := 0; i < nSamples; i++ {
		scores := liblinear.PredictProba(lr.model, X)
		_, nClasses := scores.Dims()

		bestScore := scores.At(i, 0)
		bestLangIdx := 0

		for langIdx := 0; langIdx < nClasses; langIdx++ {
			score := scores.At(i, langIdx)
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

// func (nb *NaiveBayes) ToGob() string {
// 	var output bytes.Buffer

// 	params := nb.params

// 	enc := gob.NewEncoder(&output)

// 	err := enc.Encode(params)
// 	if err != nil {
// 		log.Fatal("encode error:", err)
// 	}

// 	return output.String()
// }

// func NewNaiveBayesFromGob(gobStr string) *NaiveBayes {
// 	var params nbParams

// 	input := bytes.NewBufferString(gobStr)

// 	dec := gob.NewDecoder(input)

// 	err := dec.Decode(&params)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	nb := NewNaiveBayes()
// 	nb.params = params
// 	return nb
// }
