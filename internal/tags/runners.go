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
	if !r.MatchString(filename) {
		return nil
	}

	if containsBuildTag(f, tags) {
		return nil
	}

	missing := MissingTag{Pos: f.Name.NamePos, Token: f.Name, RequiredTags: tags}
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

	if !r.MatchString(filename) {
		return nil
	}

	node, err := parser.ParseFile(pass.Fset, filename, nil, parser.ParseComments)
	if err != nil {
		return fmt.Errorf("%w %s: %w", ErrFailedToReadFile, filename, err)
	}

	if containsBuildTag(node, tags) {
		return nil
	}

	missing := MissingTag{Pos: node.Name.NamePos, Token: node.Name, RequiredTags: tags}
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
	Pos          token.Pos
	Token        ast.Node
	RequiredTags []string
}

func report(pass *analysis.Pass, mt MissingTag) error {
	var buf bytes.Buffer
	if err := printer.Fprint(&buf, pass.Fset, mt.Token); err != nil {
		return errors.Join(ErrFailedToGenerateReport, err)
	}
	pass.Reportf(mt.Pos, "required build tag (%s) is not found", strings.Join(mt.RequiredTags, " "))
	return nil
}
