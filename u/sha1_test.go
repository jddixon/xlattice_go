package u

// xlattice_go/u/sha1_test.go

import (
	"crypto/sha1"
	. "launchpad.net/gocheck"
)

func (s *XLSuite) setUp1() {
	dataPath = "myData"
	uPath = "myU1"
	uInDir = "myU1/in"
	uTmpDir = "myU1/tmp"
	s.setUpHashTest()
	usingSHA1 = true
}

func (s *XLSuite) doTestCopyAndPut1(c *C, ds DirStruc) {
	s.setUp1()
	myU, err := New(uPath, ds, 0)
	c.Assert(err, IsNil)
	s.doTestCopyAndPut(c, myU, sha1.New())
}
func (s *XLSuite) TestCopyAndPut1(c *C) {
	s.setUp1()
	s.doTestCopyAndPut1(c, DIR16x16)
	s.doTestCopyAndPut1(c, DIR256x256)
} // FOO

func (s *XLSuite) doTestExists1(c *C, ds DirStruc) {
	s.setUp1()
	myU, err := New(uPath, ds, 0)
	c.Assert(err, IsNil)
	s.doTestExists(c, myU, sha1.New())
}
func (s *XLSuite) TestExists1(c *C) {
	s.setUp1()
	s.doTestExists1(c, DIR16x16)
	s.doTestExists1(c, DIR256x256)
} // FOO

func (s *XLSuite) doTestFileLen1(c *C, ds DirStruc) {
	s.setUp1()
	myU, err := New(uPath, ds, 0)
	c.Assert(err, IsNil)
	s.doTestFileLen(c, myU, sha1.New())
}
func (s *XLSuite) TestFileLen1(c *C) {
	s.setUp1()
	s.doTestFileLen1(c, DIR16x16)
	s.doTestFileLen1(c, DIR256x256)
} // FOO

func (s *XLSuite) doTestFileHash1(c *C, ds DirStruc) {
	s.setUp1()
	myU, err := New(uPath, ds, 0)
	c.Assert(err, IsNil)
	s.doTestFileHash(c, myU, sha1.New())
}
func (s *XLSuite) TestFileHash1(c *C) {
	s.setUp1()
	s.doTestFileHash1(c, DIR16x16)
	s.doTestFileHash1(c, DIR256x256)
} // FOO

func (s *XLSuite) doTestGetPathForKey1(c *C, ds DirStruc) {
	s.setUp1()
	myU, err := New(uPath, ds, 0)
	c.Assert(err, IsNil)
	s.doTestGetPathForKey(c, myU, sha1.New())
}
func (s *XLSuite) TestGetPathForKey1(c *C) {
	s.setUp1()
	s.doTestGetPathForKey1(c, DIR16x16)
	s.doTestGetPathForKey1(c, DIR256x256)
} // FOO

func (s *XLSuite) doTestPut1(c *C, ds DirStruc) {
	s.setUp1()
	myU, err := New(uPath, ds, 0)
	c.Assert(err, IsNil)
	s.doTestPut(c, myU, sha1.New())
}
func (s *XLSuite) TestPut1(c *C) {
	s.setUp1()
	s.doTestPut1(c, DIR16x16)
	s.doTestPut1(c, DIR256x256)
} // FOO

func (s *XLSuite) doTestPutData1(c *C, ds DirStruc) {
	s.setUp1()
	myU, err := New(uPath, ds, 0)
	c.Assert(err, IsNil)
	s.doTestPutData(c, myU, sha1.New())
}
func (s *XLSuite) TestPutData1(c *C) {
	s.setUp1()
	s.doTestPutData1(c, DIR16x16)
	s.doTestPutData1(c, DIR256x256)
} // FOO
