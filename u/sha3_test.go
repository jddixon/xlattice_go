package u

// xlattice_go/u/sha3_test.go

import (
	"code.google.com/p/go.crypto/sha3"
	// "fmt"
	// "github.com/jddixon/xlattice_go/rnglib"
	. "launchpad.net/gocheck"
	// "testing"
)

func (s *XLSuite) setUp3() {
	dataPath = "myData"
	uPath = "myU3"
	uInDir = "myU3/in"
	uTmpDir = "myU3/tmp"
	s.setUpHashTest()
	usingSHA1 = false
}
func (s *XLSuite) TestCopyAndPut3(c *C) {
	s.setUp3()
	s.doTestCopyAndPut(c, New(uPath), sha3.NewKeccak256())
}
func (s *XLSuite) TestExists3(c *C) {
	s.setUp3()
	s.doTestExists(c, New(uPath), sha3.NewKeccak256())
}
func (s *XLSuite) TestFileLen3(c *C) {
	s.setUp3()
	s.doTestFileLen(c, New(uPath), sha3.NewKeccak256())
}
func (s *XLSuite) TestFileHash3(c *C) {
	s.setUp3()
	s.doTestFileHash(c, New(uPath), sha3.NewKeccak256())
}
func (s *XLSuite) TestGetPathForKey3(c *C) {
	s.setUp3()
	s.doTestGetPathForKey(c, New(uPath), sha3.NewKeccak256())
}
func (s *XLSuite) TestPut3(c *C) {
	s.setUp3()
	s.doTestPut(c, New(uPath), sha3.NewKeccak256())
}
func (s *XLSuite) TestPutData3(c *C) {
	s.setUp3()
	s.doTestPutData(c, New(uPath), sha3.NewKeccak256())
}
