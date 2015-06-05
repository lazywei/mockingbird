package main

import (
	"fmt"
	"io/ioutil"

	"github.com/lazywei/mockingbird"
)

func ConvertLibsvm(samplePath string) {
	// var start time.Time
	// var elapsed time.Duration
	langsIdx := map[string]int{}
	tokensIdx := map[string]int{}

	langDirs := getSubEntries(samplePath, true, nil)

	for _, langDir := range langDirs {
		fmt.Println(langDir)

		codeFiles := getSubEntries(langDir, false, nil)
		langName := getLastSegInPath(langDir)

		if _, ok := langsIdx[langName]; !ok {
			langsIdx[langName] = len(langsIdx)
		}

		for _, codeFile := range codeFiles {

			// start = time.Now()
			fileContent, err := ioutil.ReadFile(codeFile)
			if err != nil {
				panic(err)
			}
			// elapsed = time.Since(start)
			// fmt.Println("Reading took %s", elapsed)

			tokens := mockingbird.ExtractTokens(string(fileContent))

			for _, token := range tokens {

				if _, ok := tokensIdx[token]; !ok {
					tokensIdx[token] = len(tokensIdx)
				}
			}

		}
	}

	// fmt.Println(langsIdx)
}
