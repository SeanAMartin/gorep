package main

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/SeanAMartin/gorep/internal/display"
	"github.com/SeanAMartin/gorep/internal/search"
)

func main() {
	startTime := time.Now()
	processed := 0
	pattern := os.Args[1]
	dirPath := os.Args[2]

	ch := make(chan display.Display)
	wg := &sync.WaitGroup{}
	paths := search.GetRecursiveFilePaths(dirPath)

	for _, path := range paths {
		wg.Add(1)
		go search.SearchFile(path, pattern, wg, ch)
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
