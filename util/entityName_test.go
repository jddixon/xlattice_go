package util

import (
	"github.com/jddixon/xlattice_go/rnglib"
	. "gopkg.in/check.v1"
	"strconv"
	"strings"
)

func (s *XLSuite) noDotsOrDashes(rng *rnglib.PRNG) string {
	var length int = 3 + rng.Intn(16)
	var name = rng.NextFileName(length)
	for len(name) < 3 || strings.ContainsAny(name, ".-") ||
		strings.ContainsAny(name[0:1], "0123456789") {
		name = rng.NextFileName(length)
	}
	return name
}

func (s *XLSuite) TestGoodNames(c *C) {
	rng := rnglib.MakeSimpleRNG()
	var count int = 3 + rng.Intn(16)
	for i := 0; i < count; i++ {
		s := s.noDotsOrDashes(rng)
		c.Assert(ValidEntityName(s), IsNil)
	}
}
func (s *XLSuite) TestBadNames(c *C) {
	rng := rnglib.MakeSimpleRNG()
	var count int = 3 + rng.Intn(16)
	for i := 0; i < count; i++ {
		s := s.noDotsOrDashes(rng)
		length := len(s)
		c.Assert(length > 2, Equals, true)
		offset := 1 + rng.Intn(length-2)
		switch length % 3 {
		case 0: // error: starts with digit
			s = strconv.Itoa(rng.Intn(10)) + s[1:]
		case 1: // error: contains dot
			s = s[0:offset] + "." + s[offset+1:]
		case 2: // error: contains dash
			s = s[0:offset] + "-" + s[offset+1:]
		}
		c.Assert(ValidEntityName(s), Not(IsNil))
		c.Assert(ValidEntityName(s), Equals, INVALID_NAME())
	}
}
