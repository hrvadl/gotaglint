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
	match      *string
	tagsJoined *string
)

func init() {
	match = fs.String("match", "", "regex which will determine matched files to check")
	tagsJoined = fs.String("buildtags", "", "tags mandatory for the matched file separated by ','")
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
	if match == nil {
		return nil, ErrMatchPatterIsNotDefined
	}

	r, err := regexp.Compile(*match)
	if err != nil {
		return nil, errors.Join(ErrInvalidPattern, err)
	}

	if tagsJoined == nil {
		return nil, nil
	}

	tags := strings.Split(*tagsJoined, ",")
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
