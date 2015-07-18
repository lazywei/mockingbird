package main

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/alecthomas/kingpin.v2"
)

func CollectRosetta(rootPath, destPath string) {
	kingpin.Parse()

	taskRoot := filepath.Join(rootPath, "Task")
	fmt.Println("Rosetta Task Root is:", taskRoot)

	taskDirs := getSubEntries(taskRoot, true, nil)

	for _, taskDir := range taskDirs {
		langDirs := getSubEntries(taskDir, true, nil)

		fmt.Println(taskDir)

		for _, langDir := range langDirs {
			langName := getLastSegInPath(langDir)

			fmt.Println(langDir)

			codeFiles := getSubEntries(langDir, false, map[string]bool{
				"00DESCRIPTION": true,
				"00META.yaml":   true,
			})

			langDestPath := filepath.Join(destPath, langName)

			err := os.MkdirAll(langDestPath, 0755)

			if err != nil {
				fmt.Println("Got error when trying to create dir for each languages")
				panic(err)
			}

			for _, codeFile := range codeFiles {
				codeFileName := getLastSegInPath(codeFile)

				if len(getSubEntries(langDestPath, false, nil)) > 100 {
					// We don't want too many files in one lang.
					break
				}

				cp(filepath.Join(langDestPath, codeFileName), codeFile)
			}
		}
	}
}
