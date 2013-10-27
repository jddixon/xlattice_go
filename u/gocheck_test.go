package u

import (
    "testing"
    . "launchpad.net/gocheck"	// for Suite
)

func Test(t *testing.T) { TestingT(t) }
type XLSuite struct{}
var _ = Suite(&XLSuite{})
