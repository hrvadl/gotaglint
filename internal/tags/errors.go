package tags

import "errors"

var (
	ErrMatchPatterIsNotDefined = errors.New("match pattern is not defined")
	ErrInvalidPattern          = errors.New("invalid match pattern")
	ErrFailedToGenerateReport  = errors.New("failed to generate report")
	ErrFailedToReadFile        = errors.New("failed to read file")
)
