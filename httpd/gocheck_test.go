package httpd

import (
	. "gopkg.in/check.v1"
	"testing"
)

// test framework tie-in /////////////////////
func Test(t *testing.T) { TestingT(t) }

type XLSuite struct{}

var _ = Suite(&XLSuite{})

// end test framework setup //////////////////
