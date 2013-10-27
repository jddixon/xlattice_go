package u16x16

// xlattice_go/u/u16x16/sha1_test.go

import (
	"crypto/sha1"
	// "fmt"
	// "github.com/jddixon/xlattice_go/rnglib"
	. "launchpad.net/gocheck"
	// "testing"
)

func (s *XLSuite) setUp1() {
	dataPath = "myData"
	uPath = "myU1"
	uInDir = "myU1/in"
	uTmpDir = "myU1/tmp"
	s.setUpHashTest()
	usingSHA1 = true
}

func (s *XLSuite) TestCopyAndPut1(c *C) {
	s.setUp1()
	s.doTestCopyAndPut(c, New(uPath), sha1.New())
}
func (s *XLSuite) TestExists1(c *C) {
	s.setUp1()
	s.doTestExists(c, New(uPath), sha1.New())
}
func (s *XLSuite) TestFileLen1(c *C) {
	s.setUp1()
	s.doTestFileLen(c, New(uPath), sha1.New())
}
func (s *XLSuite) TestFileHash1(c *C) {
	s.setUp1()
	s.doTestFileHash(c, New(uPath), sha1.New())
}
func (s *XLSuite) TestGetPathForKey1(c *C) {
	s.setUp1()
	s.doTestGetPathForKey(c, New(uPath), sha1.New())
}
func (s *XLSuite) TestPut1(c *C) {
	s.setUp1()
	s.doTestPut(c, New(uPath), sha1.New())
}
func (s *XLSuite) TestPutData1(c *C) {
	s.setUp1()
	s.doTestPutData(c, New(uPath), sha1.New())
}
