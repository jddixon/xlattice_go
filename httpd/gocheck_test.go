package httpd

// xlattice_go/httpd/gocheck.go

import (
	. "gopkg.in/check.v1"
	"testing"
)

// gocheck tie-in /////////////////////
func Test(t *testing.T) { TestingT(t) }

type XLSuite struct{}

var _ = Suite(&XLSuite{})

// end gocheck setup //////////////////
