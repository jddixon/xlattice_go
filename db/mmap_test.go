package db

// xlattice_go/db/mmap_test.go

import (
	"fmt"
	xr "github.com/jddixon/xlattice_go/rnglib"
	"io/ioutil"
	.  "launchpad.net/gocheck"
	gm "launchpad.net/gommap"
	"os"
	"strings"
)

// XXX should be a utility routine
func (s *XLSuite) ScratchFileName(c *C, rng *xr.PRNG, dirName string) (fileName string) {
	length := len(dirName)
	c.Assert(length > 0, Equals, true)
	if strings.HasSuffix(dirName, "/") {
		dirName = dirName[:length - 1]
	}
	err := os.MkdirAll(dirName, 0755)
	c.Assert(err, IsNil)

	fileName = fmt.Sprintf("%s/%s", dirName, rng.NextFileName(8))
	for {
		_, err = os.Stat(fileName)
		if os.IsNotExist(err) {
			break
		}
		c.Assert(err, IsNil)
		fileName = fmt.Sprintf("%s/%s", dirName, rng.NextFileName(8))
	}
	return

}
func (s *XLSuite) TestMmap(c *C) {
	if VERBOSITY > 0 {
		fmt.Println("TEST_MMAP")
	}

	rng := xr.MakeSimpleRNG()
	pathToFile := s.ScratchFileName(c, rng, "tmp")
	fmt.Printf("FILE: %s\n", pathToFile)

	_ = pathToFile

	data := make([]byte, BLOCK_SIZE)
	rng.NextBytes(&data)
	lastByte := data[BLOCK_SIZE - 1]
	err :=  ioutil.WriteFile(pathToFile, data, 0644)
	c.Assert(err, IsNil)

	f, err := os.Open(pathToFile)
	c.Assert(err, IsNil)
	inCore, err := gm.MapAt(0, f.Fd(), 0, 2 * BLOCK_SIZE,
							gm.PROT_READ | gm.PROT_WRITE,
							gm.MAP_PRIVATE)

	const (
		ASCII_A = byte(64)
		ASCII_B = byte(65)
	)
	inCore[ BLOCK_SIZE - 1] = ASCII_A
	inCore.Sync(gm.MS_SYNC)								// blocks
	fmt.Println("first write succeeds\n")				// DEBUG

	// XXX This faults; ie we get a "fatal error: fault"
	// inCore[ 2 * BLOCK_SIZE - 1] = ASCII_B
	// c.Assert(err, IsNil)

	err = inCore.UnsafeUnmap()
	c.Assert(err, IsNil)

	f.Close()

	data2, err := ioutil.ReadFile(pathToFile)
	c.Assert(err, IsNil)
	// this succeeds, but shouldn't
	c.Assert( data2[BLOCK_SIZE - 1], Equals, lastByte )
	// this fails, meaning that the memory-mapping doesn't work
	c.Assert( data2[BLOCK_SIZE - 1], Equals, ASCII_A )
}

