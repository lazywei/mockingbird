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
	totalTokens, totalLangs := getTotalTokensAndLangs(samplePath)

	// Construct BowParams
	langsMapping := map[string]int{}
	tokensMapping := map[string]int{}

	for i, tokens := range totalTokens {
		for _, token := range tokens {
			if _, ok := tokensMapping[token]; !ok {
				tokensMapping[token] = len(tokensMapping)
			}
		}

		langName := totalLangs[i]
		if _, ok := langsMapping[langName]; !ok {
			langsMapping[langName] = len(langsMapping)
		}
	}
	// ----

	bowParams := BowParams{langsMapping, tokensMapping}
	bagOfWords, labels := getSparseFeatures(totalTokens, totalLangs, bowParams)

	saveLibsvm(bagOfWords, labels, filepath.Join(outputDirPath, "samples.libsvm"))
	saveBowParams(bowParams, filepath.Join(outputDirPath, "bow.gob"))
}

func ConvertLibsvmWithBow(samplePath string, outputDirPath string, bowParamsPath string) {
	var bowParams BowParams
	totalTokens, totalLangs := getTotalTokensAndLangs(samplePath)

	input, err := os.Open(bowParamsPath)
	dec := gob.NewDecoder(input)

	err = dec.Decode(&bowParams)
	if err != nil {
		log.Fatal(err)
	}

	bagOfWords, labels := getSparseFeatures(totalTokens, totalLangs, bowParams)

	saveLibsvm(bagOfWords, labels, filepath.Join(outputDirPath, "samples.libsvm"))
}

func saveLibsvm(bagOfWords []SparseFeatures, labels []int, filepath string) {
	outputFile, err := os.Create(filepath)
	if err != nil {
		fmt.Println("Got error when trying to create libsvm format output file")
		panic(err)
	}
	defer outputFile.Close()

	for i, label := range labels {
		outputFile.WriteString(libsvmFmt(label, bagOfWords[i]) + "\n")
	}
}

func saveBowParams(bowParams BowParams, filepath string) {
	var output bytes.Buffer
	enc := gob.NewEncoder(&output)
	err := enc.Encode(bowParams)
	if err != nil {
		log.Fatal("Encode BowParams error:", err)
	}

	bowParamsFile, err := os.Create(filepath)
	if err != nil {
		fmt.Println("Got error when trying to create libsvm format output file")
		panic(err)
	}
	defer bowParamsFile.Close()

	bowParamsFile.WriteString(output.String())
}

func getTotalTokensAndLangs(samplePath string) (totalTokens [][]string, totalLangs []string) {
	totalTokens = [][]string{}
	totalLangs = []string{}

	langDirs := getSubEntries(samplePath, true, nil)

	for _, langDir := range langDirs {
		fmt.Println(langDir)

		codeFiles := getSubEntries(langDir, false, nil)
		langName := getLastSegInPath(langDir)

		// early stop, remove this when ExtractTokens is efficient enough
		earlyStopCnt := 0

		for _, codeFile := range codeFiles {
			fileContent, err := ioutil.ReadFile(codeFile)
			if err != nil {
				panic(err)
			}

			tokens := mockingbird.ExtractTokens(string(fileContent))

			if len(tokens) == 0 {
				fmt.Println("Tokens are empty:", codeFile)
				break
			}

			totalTokens = append(totalTokens, tokens)
			totalLangs = append(totalLangs, langName)

			// early stop, remove this when ExtractTokens is efficient enough
			earlyStopCnt += 1
			if earlyStopCnt >= 10 {
				break
			}
		}
	}

	return totalTokens, totalLangs
}

func getSparseFeatures(totalTokens [][]string, totalLangs []string, bowParams BowParams) (bagOfWords []SparseFeatures, labels []int) {
	bagOfWords = []SparseFeatures{}
	labels = []int{}

	for i, tokens := range totalTokens {
		sparseFeatures := SparseFeatures{}
		for _, token := range tokens {
			sparseFeatures[bowParams.TokensMapping[token]] += 1
		}

		labels = append(labels, bowParams.LangsMapping[totalLangs[i]])

		bagOfWords = append(bagOfWords, sparseFeatures)
	}

	return bagOfWords, labels
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
