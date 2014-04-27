package util

// xlattice_go/util/timestamp_test.go

import (
	"github.com/jddixon/xlattice_go/rnglib"
	. "gopkg.in/check.v1"
	"time"
)

func (s *XLSuite) TestGoodTimes(c *C) {
	rng := rnglib.MakeSimpleRNG()
	_ = rng

	stdLayout := "Mon Jan 2 15:04:05 -0700 MST 2006"
	t, err := time.Parse(stdLayout, stdLayout)
	c.Assert(err, IsNil)
	ts := Timestamp(t.UnixNano())
	_ = ts

	myLayout := "2006-01-02 15:04:05" // a UTC time
	myUTC := "2006-01-02 22:04:05"
	t2, err := time.Parse(myLayout, myUTC)
	c.Assert(err, IsNil)

	c.Assert(t2.Unix(), Equals, t.Unix())

	// This finally tests Timestamp.String()
	c.Assert(ts.String(), Equals, myUTC)

	utc, err := ParseTimestamp(myUTC)
	c.Assert(err, IsNil)
	c.Assert(utc, Equals, Timestamp(t2.UnixNano()))
}
