package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strconv"

	"github.com/lazywei/mockingbird"
)

type SparseFeatures map[int]int

func ConvertLibsvm(samplePath string, outputFilePath string) {

	// var start time.Time
	// var elapsed time.Duration
	langsIdx := map[string]int{}
	tokensIdx := map[string]int{}

	bagOfWords := []SparseFeatures{}
	labels := []int{}

	langDirs := getSubEntries(samplePath, true, nil)

	for _, langDir := range langDirs {
		fmt.Println(langDir)

		codeFiles := getSubEntries(langDir, false, nil)
		langName := getLastSegInPath(langDir)

		if _, ok := langsIdx[langName]; !ok {
			langsIdx[langName] = len(langsIdx)
		}

		for _, codeFile := range codeFiles {
			sparseFeatures := SparseFeatures{}

			fileContent, err := ioutil.ReadFile(codeFile)
			if err != nil {
				panic(err)
			}

			tokens := mockingbird.ExtractTokens(string(fileContent))

			for _, token := range tokens {

				if _, ok := tokensIdx[token]; !ok {
					tokensIdx[token] = len(tokensIdx)
				}

				sparseFeatures[tokensIdx[token]] += 1
			}

			bagOfWords = append(bagOfWords, sparseFeatures)
			labels = append(labels, langsIdx[langName])

			// early stop, remove this when ExtractTokens is efficient enough
			break
		}
	}

	outputFile, err := os.Create(outputFilePath)
	if err != nil {
		fmt.Println("Got error when trying to create libsvm format output file")
		panic(err)
	}
	defer outputFile.Close()

	for i, label := range labels {
		outputFile.WriteString(libsvmFmt(label, bagOfWords[i]) + "\n")
	}
}

func sortedKeys(sparseFeatures SparseFeatures) []int {
	sk := make([]int, 0, len(sparseFeatures))
	for i := range sparseFeatures {
		sk = append(sk, i)
	}
	sort.Ints(sk)
	return sk
}

func libsvmFmt(label int, sparseFeatures SparseFeatures) string {
	rtn := strconv.Itoa(label)
	for _, i := range sortedKeys(sparseFeatures) {
		rtn = rtn + " " + strconv.Itoa(i) + ":" + strconv.Itoa(sparseFeatures[i])
	}
	return rtn
}
