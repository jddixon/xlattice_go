package u256x256

// xlattice_go/sha3_test.go

import (
	"code.google.com/p/go.crypto/sha3"
	// "fmt"
	// "github.com/jddixon/xlattice_go/rnglib"
	. "launchpad.net/gocheck"
	// "testing"
)

func (s *XLSuite) setUp3() {
	dataPath	= "myData"
	uPath		= "myU3"
	uInDir		= "myU3/in"
	uTmpDir		= "myU3/tmp"
	s.setUpHashTest()
}
func (s *XLSuite) TestCopyAndPut(c *C) {
	s.setUp3()
	s.doTestCopyAndPut(c, New(uPath), sha3.NewKeccak256())
}
func (s *XLSuite) TestExists(c *C) {
	s.setUp3()
	s.doTestExists(c, New(uPath), sha3.NewKeccak256())
}
func (s *XLSuite) TestFileLen(c *C) {
	s.setUp3()
	s.doTestFileLen(c, New(uPath), sha3.NewKeccak256())
}
func (s *XLSuite) TestFileHash(c *C) {
	s.setUp3()
	s.doTestFileHash(c, New(uPath), sha3.NewKeccak256())
}
func (s *XLSuite) TestGetPathForKey(c *C) {
	s.setUp3()
	s.doTestGetPathForKey(c, New(uPath), sha3.NewKeccak256())
}
func (s *XLSuite) TestPut(c *C) {
	s.setUp3()
	s.doTestPut(c, New(uPath), sha3.NewKeccak256())
}
func (s *XLSuite) TestPutData(c *C) {
	s.setUp3()
	s.doTestPutData(c, New(uPath), sha3.NewKeccak256())
}
