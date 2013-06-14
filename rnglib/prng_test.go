package rnglib

import (
	"fmt"
	. "launchpad.net/gocheck"
	"os"
	"strings"
	"testing"
	"time"
)

// SETUP FOR THIS PACKAGE (rnglib) //////////////////////////////////
const TMP_DIR = "tmp"

// gocheck tie-in /////////////////////
func Test(t *testing.T) { TestingT(t) }

type XLSuite struct{}

var _ = Suite(&XLSuite{})

// end gocheck setup //////////////////

// utility functions ////////////////////////////////////////////////
func (s *XLSuite) buildData(count uint32) *[]byte {
	p := make([]byte, count)
	return &p
}
func (s *XLSuite) makeSimpleRNG() *PRNG {
	t := time.Now().Unix() // int64 sec
	rng := NewSimpleRNG(t)
	return rng
}

// parameterized unit tests /////////////////////////////////////////
func (s *XLSuite) doTestConstructor(c *C, rng *PRNG) {
	c.Assert(rng, Not(IsNil)) // NOT
}
func (s *XLSuite) doTestNextBoolean(c *C, rng *PRNG) {
	val := rng.NextBoolean()
	c.Assert(val, Not(IsNil)) // NOT

	valAsIface := interface{}(val)
	switch v := valAsIface.(type) {
	default:
		fmt.Printf("expected type bool, found %T", v)
		// assert.Fail("whatever NextBoolean() returns is not a bool")
	case bool:
		/* empty statement */
	}
}
func (s *XLSuite) doTestNextByte(c *C, rng *PRNG) {
	_ = c
	_ = rng
}
func (s *XLSuite) doTestNextBytes(c *C, rng *PRNG) {
	count := uint32(1)          // minimum length of buffer
	count += rng.NextInt32(256) // maximum
	data := s.buildData(count)  // so 1 .. 256 bytes
	rng.NextBytes(data)
	actualLen := uint32(len(*data))
	c.Assert(0, Not(Equals), actualLen) // NOT
	c.Assert(actualLen, Equals, count)

}
func (s *XLSuite) doTestNextFileName(c *C, rng *PRNG) {
	maxLen := uint32(1)         // minimum length of name
	maxLen += rng.NextInt32(16) // maximum
	name := rng.NextFileName(int(maxLen))
	// DEBUG
	fmt.Printf("next file name is %s\n", name)
	// END
	actualLen := len(name)
	c.Assert(0, Not(Equals), actualLen) // NOT
	// assert.True( t, actualLen < maxLen)
}
func (s *XLSuite) doTestNextDataFile(c *C, rng *PRNG) {
	minLen := int(rng.NextInt32(4))            // minimum length of file
	maxLen := minLen + int(rng.NextInt32(256)) // maximum

	// XXX should return err, which should be nil
	fileLen, pathToFile := rng.NextDataFile(TMP_DIR, maxLen, minLen)
	// DEBUG
	fmt.Printf("data file is %s; size is %d\n", pathToFile, fileLen)
	// END

	stats, err := os.Stat(pathToFile)
	c.Assert(err, IsNil)
	fileName := stats.Name()
	c.Assert(TMP_DIR+"/"+fileName, Equals, pathToFile)
	c.Assert(stats.Size(), Equals, int64(fileLen))

}
func (s *XLSuite) doNextDataDirTest(c *C, rng *PRNG, width int, depth int) {
	dirName := rng.NextFileName(8)
	dirPath := TMP_DIR + "/" + dirName
	pathExists, err := PathExists(dirPath)
	if err != nil {
		panic("error invoking PathExists on " + dirPath)
	}
	if pathExists {
		if strings.HasPrefix(dirPath, "/") {
			panic("attempt to remove absolute path " + dirPath)
		}
		if strings.Contains(dirPath, "..") {
			panic("attempt to remove path containing ..: " + dirPath)
		}
		os.RemoveAll(dirPath)
	}
	rng.NextDataDir(dirPath, width, depth, 32, 0)
}
func (s *XLSuite) doTestNextDataDir(c *C, rng *PRNG) {
	s.doNextDataDirTest(c, rng, 1, 1)
	s.doNextDataDirTest(c, rng, 1, 4)
	s.doNextDataDirTest(c, rng, 4, 1)
	s.doNextDataDirTest(c, rng, 4, 4)
}
