package display

import (
	"fmt"

	"github.com/SeanAMartin/gorep/internal/search"
)

type Display struct {
	filePath string
	search.SearchResult
}

func (d Display) PrettyPrint() {
	fmt.Printf("Line Number: %v\nFilePath: %v\nLine: %v\n\n", d.lineNumber, d.filePath, d.line)
}
