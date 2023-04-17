package report

import "errors"

var (
	ErrBeforeLastRun = errors.New("cannot schedule report before last run")
)
