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
	"io/ioutil"
	"log"
	"os"
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

func searchLine(pattern string, line string, lineNumber int) SearchResult {
	var searchResult SearchResult
	if strings.Contains(line, pattern) {
		searchResult.lineNumber = lineNumber + 1
		searchResult.line = line
	}
	return searchResult
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

func main() {
	pattern := os.Args[1]
	filePath := os.Args[2]
	input := fileFromPath(filePath)
	var displays []Display
	lines := splitIntoLines(input)
	fmt.Println(len(lines) - 1)

	for index, line := range lines {
		searchResult := searchLine(pattern, line, index)
		if searchResult.lineNumber != 0 {
			displays = append(displays, Display{filePath, searchResult})
		}
	}

	for _, display := range displays {
		display.PrettyPrint()
	}
}
