package util

import (
	. "launchpad.net/gocheck"
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
