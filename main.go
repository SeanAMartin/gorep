package main

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type SearchResult struct {
	line       string
	lineNumber int
}

type Display struct {
	filePath string
	SearchResult
}

func (d Display) PrettyPrint() {
	fmt.Printf("Line Number: %v\nFilePath: %v\nLine: %v\n\n", d.lineNumber, d.filePath, d.line)
}

func searchLine(pattern string, line string, lineNumber int) (SearchResult, bool) {
	if strings.Contains(line, pattern) {
		return SearchResult{lineNumber: lineNumber + 1, line: line}, true
	}
	return SearchResult{}, false
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

func getRecursiveFilePaths(inputDir string) []string {
	var paths []string
	err := filepath.Walk(inputDir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
			return err
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

func searchPath(path string, pattern string, wg *sync.WaitGroup, ch chan Display) {
	defer wg.Done()
	for _, display := range searchFile(path, pattern) {
		ch <- display
	}
}

func searchFile(path string, pattern string) []Display {
	var out []Display
	input := fileFromPath(path)
	lines := splitIntoLines(input)
	for index, line := range lines {
		if searchResult, ok := searchLine(pattern, line, index); ok {
			out = append(out, Display{path, searchResult})
		}
	}
	return out
}

func main() {
	pattern := os.Args[1]
	dirPath := os.Args[2]
	ch := make(chan Display)
	wg := &sync.WaitGroup{}
	paths := getRecursiveFilePaths(dirPath)

	for _, path := range paths {
		wg.Add(1)
		go searchPath(path, pattern, wg, ch)
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	for d := range ch {
		d.PrettyPrint()
	}

}
