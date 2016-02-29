package db

import (
	. "gopkg.in/check.v1"
	"testing"
)

// IF USING test framework, need a file like this in each package=directory.

func Test(t *testing.T) { TestingT(t) }

type XLSuite struct{}

var _ = Suite(&XLSuite{})

const (
	BLOCK_SIZE = 4096
	SHA1_LEN   = 20
	SHA3_LEN   = 32
	VERBOSITY  = 1
)
