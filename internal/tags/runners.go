package tags

import (
	"bytes"
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"regexp"
	"slices"
	"strings"

	"golang.org/x/tools/go/analysis"
)

func checkGoFile(
	pass *analysis.Pass,
	f *ast.File,
	r *regexp.Regexp,
	tags []string,
) error {
	filename := pass.Fset.Position(f.Pos()).Filename
	if !r.Match([]byte(filename)) {
		return nil
	}

	if containsBuildTag(f, tags) {
		return nil
	}

	missing := MissingTag{Pos: f.Name.NamePos, Token: f.Name}
	return report(pass, missing)
}

func checkOtherFile(
	pass *analysis.Pass,
	filename string,
	r *regexp.Regexp,
	tags []string,
) error {
	if !strings.HasSuffix(filename, ".go") {
		return nil
	}

	if !r.Match([]byte(filename)) {
		return nil
	}

	node, err := parser.ParseFile(pass.Fset, filename, nil, parser.ParseComments)
	if err != nil {
		return fmt.Errorf("%w %s: %w", ErrFailedToReadFile, filename, err)
	}

	if containsBuildTag(node, tags) {
		return nil
	}

	missing := MissingTag{Pos: node.Name.NamePos, Token: node.Name}
	return report(pass, missing)
}

func containsBuildTag(f *ast.File, tags []string) bool {
	for _, cg := range f.Comments {
		for _, c := range cg.List {
			if slices.Contains(tags, c.Text) {
				return true
			}
		}
	}

	return false
}

type MissingTag struct {
	Pos   token.Pos
	Token ast.Node
}

func report(pass *analysis.Pass, mt MissingTag) error {
	var buf bytes.Buffer
	if err := printer.Fprint(&buf, pass.Fset, mt.Token); err != nil {
		return errors.Join(ErrFailedToGenerateReport, err)
	}
	pass.Reportf(mt.Pos, "matching build tag is not found")
	return nil
}