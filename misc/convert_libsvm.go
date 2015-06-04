package main

import (
	"io/ioutil"

	"github.com/lazywei/mockingbird"
)

func ConvertLibsvm(samplePath string) {
	langsIdx := map[string]int{}
	tokensIdx := map[string]int{}

	langDirs := getSubEntries(samplePath, true, nil)

	for _, langDir := range langDirs {

		codeFiles := getSubEntries(langDir, false, nil)
		langName := getLastSegInPath(langDir)

		if _, ok := langsIdx[langName]; !ok {
			langsIdx[langName] = len(langsIdx)
		}

		for _, codeFile := range codeFiles {

			fileContent, err := ioutil.ReadFile(codeFile)
			if err != nil {
				panic(err)
			}

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
