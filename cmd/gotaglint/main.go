package main

import (
	"golang.org/x/tools/go/analysis/singlechecker"

	"github.com/hrvadl/gotaglint/internal/tags"
)

func main() {
	singlechecker.Main(tags.NewAnalyzer())
}
