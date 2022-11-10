package search

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"
	"sync"

	"github.com/SeanAMartin/gorep/internal/display"
)

type SearchResult struct {
	line       string
	lineNumber int
}

func SearchLine(pattern string, line string, lineNumber int) (SearchResult, bool) {
	if strings.Contains(line, pattern) {
		return SearchResult{lineNumber: lineNumber + 1, line: line}, true
	}
	return SearchResult{}, false
}

func SearchFile(path string, pattern string, wg *sync.WaitGroup, ch chan display.Display) {
	defer wg.Done()
	input := fileFromPath(path)
	lines := splitIntoLines(input)
	for index, line := range lines {
		if searchResult, ok := SearchLine(pattern, line, index); ok {
			ch <- display.Display{path, searchResult}
		}
	}
}

func GetRecursiveFilePaths(inputDir string) []string {
	var paths []string
	err := filepath.Walk(inputDir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
		}
		if !info.IsDir() {
			paths = append(paths, path)
		}
		return nil
	})
	if err != nil {
		fmt.Printf("Error walking the path %q: %v\n", inputDir, err)
	}
	return paths
}

func splitIntoLines(file string) []string {
	lines := strings.Split(file, "\n")
	return lines
}

func fileFromPath(path string) string {
	fileContent, err := ioutil.ReadFile(path)

	if err != nil {
		log.Fatal(err)
	}

	return string(fileContent)
}
