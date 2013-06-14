package rnglib

import (
	. "launchpad.net/gocheck"
)

// unit tests ///////////////////////////////////////////////////////
func (s *XLSuite) TestConstuctor(c *C) {
	s.doTestConstructor	(c, makeSimpleRNG())
}
func (s *XLSuite) TestNextBoolean(c *C) {
	s.doTestNextBoolean(c, makeSimpleRNG())
}
func (s *XLSuite) TestNextByte(c *C) {
	s.doTestNextByte(c, makeSimpleRNG())
}
func (s *XLSuite) TestNextBytes(c *C) {
	s.doTestNextBytes(c, makeSimpleRNG())
}
func (s *XLSuite) TestNextFileName(c *C) {
	s.doTestNextFileName(c, makeSimpleRNG())
}
func (s *XLSuite) TestNextDataFile(c *C) {
	s.doTestNextDataFile(c, makeSimpleRNG())
}
func (s *XLSuite) TestNextDataDir(c *C) {
	s.doTestNextDataDir(c, makeSimpleRNG())
}

