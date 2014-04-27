package overlay

// xlattice_go/overlay/gocheck.go

import (
	. "gopkg.in/check.v1"
	"testing"
)

// gocheck tie-in /////////////////////
func Test(t *testing.T) { TestingT(t) }

type XLSuite struct{}

var _ = Suite(&XLSuite{})

// end gocheck setup //////////////////
