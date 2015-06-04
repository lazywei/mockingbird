package main

import "gopkg.in/alecthomas/kingpin.v2"

var (
	debug    = kingpin.Flag("debug", "enable debug mode").Default("false").Bool()
	serverIP = kingpin.Flag("server", "server address").Default("127.0.0.1").IP()

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
)

func main() {
	switch kingpin.Parse() {
	// Register user
	case "collectRosetta":
		CollectRosetta(*rosettaRootPath, *rosettaDestPath)

	// Post message
	case "convertLibsvm":
		ConvertLibsvm(*libsvmSamplePath)
	}
}
