package crypto

// xlattice_go/crypto/pkcs7_test.go

import (
	"crypto/aes"
	xr "github.com/jddixon/xlattice_go/rnglib"
	. "launchpad.net/gocheck"
)

// TODO: MOVE THIS TO crypto/ =======================================

func (s *XLSuite) TestPKCS7Padding(c *C) {
	rng := xr.MakeSimpleRNG()
	seven := make([]byte, 7)
	rng.NextBytes(&seven)

	fifteen := make([]byte, 15)
	rng.NextBytes(&fifteen)

	sixteen := make([]byte, 16)
	rng.NextBytes(&sixteen)

	seventeen := make([]byte, 17)
	rng.NextBytes(&seventeen)

	padding := PKCS7Padding(seven, aes.BlockSize)
	c.Assert(len(padding), Equals, aes.BlockSize-7)
	c.Assert(padding[0], Equals, byte(aes.BlockSize-7))

	padding = PKCS7Padding(fifteen, aes.BlockSize)
	c.Assert(len(padding), Equals, aes.BlockSize-15)
	c.Assert(padding[0], Equals, byte(aes.BlockSize-15))

	padding = PKCS7Padding(sixteen, aes.BlockSize)
	c.Assert(len(padding), Equals, aes.BlockSize)
	c.Assert(padding[0], Equals, byte(16))

	padding = PKCS7Padding(seventeen, aes.BlockSize)
	expectedLen := 2*aes.BlockSize - 17
	c.Assert(len(padding), Equals, expectedLen)
	c.Assert(padding[0], Equals, byte(expectedLen))

	paddedSeven, err := AddPKCS7Padding(seven, aes.BlockSize)
	c.Assert(err, IsNil)
	unpaddedSeven, err := StripPKCS7Padding(paddedSeven, aes.BlockSize)
	c.Assert(err, IsNil)
	c.Assert(seven, DeepEquals, unpaddedSeven)

	paddedFifteen, err := AddPKCS7Padding(fifteen, aes.BlockSize)
	c.Assert(err, IsNil)
	unpaddedFifteen, err := StripPKCS7Padding(paddedFifteen, aes.BlockSize)
	c.Assert(err, IsNil)
	c.Assert(fifteen, DeepEquals, unpaddedFifteen)

	paddedSixteen, err := AddPKCS7Padding(sixteen, aes.BlockSize)
	c.Assert(err, IsNil)
	unpaddedSixteen, err := StripPKCS7Padding(paddedSixteen, aes.BlockSize)
	c.Assert(err, IsNil)
	c.Assert(sixteen, DeepEquals, unpaddedSixteen)

	paddedSeventeen, err := AddPKCS7Padding(seventeen, aes.BlockSize)
	c.Assert(err, IsNil)
	unpaddedSeventeen, err := StripPKCS7Padding(paddedSeventeen, aes.BlockSize)
	c.Assert(err, IsNil)
	c.Assert(seventeen, DeepEquals, unpaddedSeventeen)
}

// END MOVE THIS TO crypto/ =========================================
