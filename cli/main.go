package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/davecgh/go-spew/spew"
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

	predict      = kingpin.Command("predict", "Predict via trained Classifier")
	predictModel = predict.
			Flag("model", "Path for loading saved model").
			Default("./model/naive_bayes.yaml").String()
	predictTestData = predict.
			Flag("data", "Path for testing data (in libsvm format)").
			Required().String()

	collectRosetta  = kingpin.Command("collectRosetta", "Collect RosettaCodeData into proper structure.")
	rosettaRootPath = collectRosetta.
			Arg("rootPath", "Path to RosettaCodeData root").
			Required().String()
	rosettaDestPath = collectRosetta.
			Arg("destPath", "Path for storing converted RosettaCodeData").
			Required().String()

	convertLibsvm    = kingpin.Command("convertLibsvm", "Convert collected Rosetta data to BoW in libsvm format")
	libsvmSamplePath = convertLibsvm.
				Arg("samplePath", "Path to collected samples").
				Required().String()
	libsvmOutputFilePath = convertLibsvm.
				Arg("outputFilePath", "Path for converted output samples").
				Required().String()
)

type P struct {
	X    int
	Name map[int]int
}

func main() {
	switch kingpin.Parse() {
	case "train":
		X, y := mb.ReadLibsvm(*trainSample)
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

	case "predict":
		fmt.Println("Model Loading ...")
		gobStr, err := ioutil.ReadFile(*predictModel)
		fmt.Println("Model Loaded")

		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Model Initiating ...")
		nb := mb.NewNaiveBayesFromGob(string(gobStr))
		fmt.Println("Model Initiated ...")

		fmt.Println("Data Loading ...")
		X, _ := mb.ReadLibsvm(*predictTestData)
		fmt.Println("Data Loaded")

		labels := []int{}
		for _, y := range nb.Predict(X) {
			labels = append(labels, y.Label)
		}
		spew.Dump(labels)

	case "collectRosetta":
		CollectRosetta(*rosettaRootPath, *rosettaDestPath)

	case "convertLibsvm":
		ConvertLibsvm(*libsvmSamplePath, *libsvmOutputFilePath)
	}
}
