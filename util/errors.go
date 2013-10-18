package util

// xlattice_go/util/errors.go

import (
	"errors"
)

var (
	InvalidName           = errors.New("not a valid entity name")
	TooManyPartsInVersion = errors.New("too many parts in version")
)
