package main

import (
	mb "github.com/lazywei/mockingbird"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	// debug    = kingpin.Flag("debug", "enable debug mode").Default("false").Bool()
	// serverIP = kingpin.Flag("server", "server address").Default("127.0.0.1").IP()

	train       = kingpin.Command("train", "Train Classifier")
	trainSample = train.
			Flag("sample", "Path for samples (in libsvm format)").
			Default("samples.libsvm").String()
	trainOutput = train.
			Flag("output", "Path for saving trained model").
			Default("model").String()

	predict = kingpin.Command("predict", "Predict via trained Classifier")

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

func main() {
	switch kingpin.Parse() {
	case "train":
		X, y := mb.ReadLibsvm(*trainSample)
		nb := mb.NewNaiveBayes()
		nb.Fit(X, y)

	case "predict":

	case "collectRosetta":
		CollectRosetta(*rosettaRootPath, *rosettaDestPath)

	case "convertLibsvm":
		ConvertLibsvm(*libsvmSamplePath, *libsvmOutputFilePath)
	}
}
