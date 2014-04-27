package util

import (
	. "gopkg.in/check.v1"
	"testing"
)

// IF USING gocheck, need a file like this in each package=directory.

func Test(t *testing.T) { TestingT(t) }

type XLSuite struct{}

var _ = Suite(&XLSuite{})

// LOCAL VARIATIONS -------------------------------------------------
const (
	VERBOSITY = 1
)
