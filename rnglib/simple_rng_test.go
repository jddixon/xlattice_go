package rnglib

import (
	. "gopkg.in/check.v1"
	"time"
)

// utility functions ////////////////////////////////////////////////
func makeSimpleRNG() *PRNG {
	t := time.Now().Unix()
	rng := NewSimpleRNG(t)
	return rng
}

// unit tests ///////////////////////////////////////////////////////
func (s *XLSuite) TestConstuctor(c *C) {
	s.doTestConstructor(c, makeSimpleRNG())
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
