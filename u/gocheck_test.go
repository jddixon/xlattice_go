package u

import (
	. "launchpad.net/gocheck" // for Suite
	"testing"
)

func Test(t *testing.T) { TestingT(t) }

type XLSuite struct{}

var _ = Suite(&XLSuite{})
