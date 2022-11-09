package main

//--Summary:
//  Create a grep clone that can do simple substring searching
//  within files. It must auto-recurse into subdirectories.
//
//--Requirements:
//* Use goroutines to search through the files for a substring match
//* Display matches to the terminal as they are found
//  * Display the line number, file path, and complete line containing the match
//* Recurse into any subdirectories looking for matches
//* Use any synchronization method to ensure that all files
//  are searched, and all results are displayed before the program
//  terminates.
//
//--Notes:
//* Program invocation should follow the pattern:
//    mgrep search_string search_dir

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
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
	fmt.Printf("Line Number: %v\nFilePath: %v\nLine: %v\n", d.lineNumber, d.filePath, d.line)
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

func main() {
	pattern := os.Args[1]
	dirPath := os.Args[2]

	paths := getRecursiveFilePaths(dirPath)
	for _, path := range paths {
		input := fileFromPath(path)
		var displays []Display
		lines := splitIntoLines(input)

		for index, line := range lines {
			if searchResult, ok := searchLine(pattern, line, index); ok {
				displays = append(displays, Display{path, searchResult})
			}
		}

		for _, display := range displays {
			display.PrettyPrint()
		}
	}
}
