package rnglib

import (
	. "launchpad.net/gocheck"
	"time"
)

// utility functions ////////////////////////////////////////////////
func makeSystemRNG() *PRNG {
	t := time.Now().Unix()
	rng := NewSystemRNG(t)
	return rng
}

// unit tests ///////////////////////////////////////////////////////
func (s *XLSuite) TestSystemConstuctor(c *C) {
	s.doTestConstructor(c, makeSystemRNG())
}
func (s *XLSuite) TestSystemNextBoolean(c *C) {
	s.doTestNextBoolean(c, makeSystemRNG())
}
func (s *XLSuite) TestSystemNextByte(c *C) {
	s.doTestNextByte(c, makeSystemRNG())
}
func (s *XLSuite) TestSystemNextBytes(c *C) {
	s.doTestNextBytes(c, makeSystemRNG())
}
func (s *XLSuite) TestSystemNextFileName(c *C) {
	s.doTestNextFileName(c, makeSystemRNG())
}
func (s *XLSuite) TestSystemNextDataFile(c *C) {
	s.doTestNextDataFile(c, makeSystemRNG())
}
func (s *XLSuite) TestSystemNextDataDir(c *C) {
	s.doTestNextDataDir(c, makeSystemRNG())
}
