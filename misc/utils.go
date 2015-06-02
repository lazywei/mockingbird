package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
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
