package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/davecgh/go-spew/spew"
	"github.com/lazywei/liblinear"
	mb "github.com/lazywei/mockingbird"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	train       = kingpin.Command("train", "Train Classifier")
	trainSample = train.
			Flag("sample", "Path for samples (in libsvm format)").
			Default("samples.libsvm").String()
	trainOutput = train.
			Flag("output", "Path for saving trained model").
			Default("model").String()
	trainSolver = train.
			Flag("solver", "0 = NaiveBayes, 1 = LogisticRegression").
			Default("0").Int()

	predict      = kingpin.Command("predict", "Predict via trained Classifier")
	predictModel = predict.
			Flag("model", "Path for loading saved model").
			Default("./model/naive_bayes.gob").String()
	predictTestData = predict.
			Flag("data", "Path for testing data (in libsvm format)").
			Required().String()
	predictSolver = predict.
			Flag("solver", "0 = NaiveBayes, 1 = LogisticRegression").
			Default("0").Int()

	collectRosetta  = kingpin.Command("collectRosetta", "Collect RosettaCodeData into proper structure.")
	rosettaRootPath = collectRosetta.
			Arg("rootPath", "Path to RosettaCodeData root").
			Required().String()
	rosettaDestPath = collectRosetta.
			Arg("destPath", "Path for storing converted RosettaCodeData").
			Required().String()

	convertLibsvm    = kingpin.Command("convertLibsvm", "Convert collected Rosetta data to BoW in libsvm format")
	libsvmSamplePath = convertLibsvm.
				Arg("samplePath", "Path for collected samples").
				Required().String()
	libsvmOutputDirPath = convertLibsvm.
				Arg("outputDirPath", "Path for saving converted output samples and BOW params").
				Required().String()
	libsvmBowPath = convertLibsvm.
			Flag("bowPath", "Path for predefined bag-of-words params for constructing libsvm format").
			Default("").String()
)

func main() {
	switch kingpin.Parse() {
	case "train":
		X, y := liblinear.ReadLibsvm(*trainSample, false)

		switch *trainSolver {
		case 0:
			nb := mb.NewNaiveBayes()
			nb.Fit(X, y)

			os.MkdirAll(*trainOutput, 0755)
			err := ioutil.WriteFile(
				filepath.Join(*trainOutput, "naive_bayes.gob"),
				[]byte(nb.ToGob()),
				0644,
			)

			if err != nil {
				log.Fatal(err)
			}
		case 1:
			lr := mb.NewLogisticRegression()
			lr.Fit(X, y)
			lr.SaveModel(filepath.Join(*trainOutput, "lr.model"))

		default:
			fmt.Println("Unsupported Solver Type: %v", *trainSolver)
		}

	case "predict":
		fmt.Println("Data Loading ...")
		X, _ := liblinear.ReadLibsvm(*predictTestData, false)
		fmt.Println("Data Loaded")

		switch *predictSolver {
		case 0:
			fmt.Println("Model Loading ...")
			gobStr, err := ioutil.ReadFile(*predictModel)
			fmt.Println("Model Loaded")

			if err != nil {
				log.Fatal(err)
			}

			fmt.Println("Model Initiating ...")
			nb := mb.NewNaiveBayesFromGob(string(gobStr))
			fmt.Println("Model Initiated ...")

			labels := []int{}
			for _, y := range nb.Predict(X) {
				labels = append(labels, y.Label)
			}
			spew.Dump(labels)
		case 1:
			lr := mb.NewLogisticRegressionFromModel(*predictModel)
			labels := []int{}
			for _, y := range lr.Predict(X) {
				labels = append(labels, y.Label)
			}
			spew.Dump(labels)

		default:

			fmt.Println("Unsupported Solver Type: %v", *predictSolver)
		}

	case "collectRosetta":
		CollectRosetta(*rosettaRootPath, *rosettaDestPath)

	case "convertLibsvm":
		if *libsvmBowPath != "" {
			ConvertLibsvmWithBow(*libsvmSamplePath, *libsvmOutputDirPath, *libsvmBowPath)
		} else {
			ConvertLibsvm(*libsvmSamplePath, *libsvmOutputDirPath)
		}
	}
}
