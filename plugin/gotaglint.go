package plugin

import (
	"golang.org/x/tools/go/analysis"

	"github.com/hrvadl/gotaglint/internal/tags"
)

func New(conf any) ([]*analysis.Analyzer, error) {
	return []*analysis.Analyzer{tags.NewAnalyzer()}, nil
}
