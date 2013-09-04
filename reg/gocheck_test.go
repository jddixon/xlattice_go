package reg

import (
	. "launchpad.net/gocheck"
	"testing"
)

// IF USING gocheck, need a file like this in each package=directory.

func Test(t *testing.T) { TestingT(t) }

type XLSuite struct{}

var _ = Suite(&XLSuite{})

const (
	VERBOSITY = 1
	SHA1_LEN  = 20
	SHA3_LEN  = 32
)
