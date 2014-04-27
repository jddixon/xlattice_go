package crypto

// xlattice_go/crypto/gocheck.go

import (
	//. "launchpad.net/gocheck"
	. "gopkg.in/check.v1"
	"testing"
)

func Test(t *testing.T) { TestingT(t) }

type XLSuite struct{}

var _ = Suite(&XLSuite{})
