package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"

	"github.com/lazywei/mockingbird"
)

type SparseFeatures map[int]int

type BowParams struct {
	LangsMapping  map[string]int
	TokensMapping map[string]int
}

func ConvertLibsvm(samplePath string, outputDirPath string) {

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

		// early stop, remove this when ExtractTokens is efficient enough
		earlyStopCnt := 0

		for _, codeFile := range codeFiles {
			sparseFeatures := SparseFeatures{}

			fileContent, err := ioutil.ReadFile(codeFile)
			if err != nil {
				panic(err)
			}

			tokens := mockingbird.ExtractTokens(string(fileContent))

			if len(tokens) == 0 {
				fmt.Println("Tokens are empty:", codeFile)
				break
			}

			for _, token := range tokens {

				if _, ok := tokensIdx[token]; !ok {
					tokensIdx[token] = len(tokensIdx)
				}

				sparseFeatures[tokensIdx[token]] += 1
			}

			bagOfWords = append(bagOfWords, sparseFeatures)
			labels = append(labels, langsIdx[langName])

			// early stop, remove this when ExtractTokens is efficient enough
			earlyStopCnt += 1
			if earlyStopCnt >= 10 {
				break
			}
		}
	}

	outputFile, err := os.Create(filepath.Join(outputDirPath, "samples.libsvm"))
	bowParamsFile, err := os.Create(filepath.Join(outputDirPath, "bow.gob"))
	if err != nil {
		fmt.Println("Got error when trying to create libsvm format output file")
		panic(err)
	}
	defer outputFile.Close()
	defer bowParamsFile.Close()

	for i, label := range labels {
		outputFile.WriteString(libsvmFmt(label, bagOfWords[i]) + "\n")
	}

	var output bytes.Buffer
	enc := gob.NewEncoder(&output)

	err = enc.Encode(BowParams{langsIdx, tokensIdx})
	if err != nil {
		log.Fatal("Encode BowParams error:", err)
	}

	bowParamsFile.WriteString(output.String())
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
