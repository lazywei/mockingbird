package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	rosettaPath = kingpin.
			Arg("rosettaPath", "Path to RosettaCodeData root").
			Required().String()

	destPath = kingpin.
			Arg("destPath", "Path for storing converted RosettaCodeData").
			Required().String()
)

// Ref: https://gist.github.com/elazarl/5507969
func cp(dst, src string) error {
	s, err := os.Open(src)
	if err != nil {
		return err
	}
	// no need to check errors on read only file, we already got everything
	// we need from the filesystem, so nothing can go wrong now.
	defer s.Close()
	d, err := os.Create(dst)
	if err != nil {
		return err
	}
	if _, err := io.Copy(d, s); err != nil {
		d.Close()
		return err
	}
	return d.Close()
}

func getSubEntries(path string, dirMode bool, ignoreList map[string]bool) []string {
	subDirs, err := ioutil.ReadDir(path)

	if err != nil {
		fmt.Println("Got error when reading Sub Entries")
		panic(err)
	}

	subDirPaths := []string{}

	for _, subDir := range subDirs {
		subDirPath := filepath.Join(path, subDir.Name())

		if subDir.IsDir() != dirMode {
			continue
		}

		if ignoreList != nil && ignoreList[subDir.Name()] {
			continue
		}

		subDirPaths = append(subDirPaths, subDirPath)
	}

	return subDirPaths
}

func getLastSegInPath(path string) string {
	segs := strings.Split(path, "/")
	return segs[len(segs)-1]
}

func main() {
	kingpin.Parse()

	taskRoot := filepath.Join(*rosettaPath, "Task")
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

			langDestPath := filepath.Join(*destPath, langName)

			err := os.MkdirAll(langDestPath, 0755)

			if err != nil {
				fmt.Println("Got error when trying to create dir for each languages")
				panic(err)
			}

			for _, codeFile := range codeFiles {
				codeFileName := getLastSegInPath(codeFile)

				cp(filepath.Join(langDestPath, codeFileName), codeFile)
			}
		}
	}

}
