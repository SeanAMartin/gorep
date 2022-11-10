package main

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"sync"
	"time"
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

func searchFile(path string, pattern *regexp.Regexp, wg *sync.WaitGroup, ch chan Display) {
	defer wg.Done()

	count := 0

	f, _ := os.Open(path)
	defer f.Close()

	scan := bufio.NewScanner(f)
	for scan.Scan() {
		count += 1
		if pattern.Match(scan.Bytes()) {
			ch <- Display{path, SearchResult{lineNumber: count, line: scan.Text()}}
		}
	}
}

func main() {
	startTime := time.Now()
	processed := 0

	pattern := os.Args[1]
	dirPath := os.Args[2]

	compiledPattern := regexp.MustCompile(pattern)

	ch := make(chan Display)
	wg := &sync.WaitGroup{}
	paths := getRecursiveFilePaths(dirPath)

	for _, path := range paths {
		wg.Add(1)
		go searchFile(path, compiledPattern, wg, ch)
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	for d := range ch {
		d.PrettyPrint()
		processed += 1
	}

	elapsed := time.Since(startTime)
	fmt.Printf("%v results found in %v seconds", processed, elapsed)
}
