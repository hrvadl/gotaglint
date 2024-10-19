package tags

import (
	"errors"
	"flag"
	"regexp"
	"slices"
	"strings"

	"golang.org/x/tools/go/analysis"
)

var (
	fs         flag.FlagSet
	match      string
	tagsJoined string
)

func init() {
	fs.StringVar(&match, "match", "", "regex which will determine matched files to check")
	fs.StringVar(
		&tagsJoined,
		"buildtags",
		"",
		"tags mandatory for the matched file separated by ','",
	)
}

func NewAnalyzer() *analysis.Analyzer {
	return &analysis.Analyzer{
		Name:             "gotaglint",
		Doc:              "reports missing build tags",
		Run:              run,
		Flags:            fs,
		RunDespiteErrors: true,
	}
}

func run(pass *analysis.Pass) (any, error) {
	r, err := regexp.Compile(match)
	if err != nil {
		return nil, errors.Join(ErrInvalidPattern, err)
	}

	tags := strings.Split(tagsJoined, ",")
	if len(tags) == 0 {
		return nil, nil
	}

	for i, t := range tags {
		tags[i] = "//go:build " + t
	}

	for _, f := range pass.Files {
		if err = checkGoFile(pass, f, r, tags); err != nil {
			return nil, err
		}
	}

	otherFiles := slices.Concat(pass.OtherFiles, pass.IgnoredFiles)
	for _, filename := range otherFiles {
		if err = checkOtherFile(pass, filename, r, tags); err != nil {
			return nil, err
		}
	}

	return nil, nil
}
