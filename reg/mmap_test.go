package reg

// xlattice_go/reg/mmap_test.go

import (
	"bytes"
	"fmt"
	xr "github.com/jddixon/xlattice_go/rnglib"
	"io/ioutil"
	. "launchpad.net/gocheck"
	gm "launchpad.net/gommap"
	"os"
	"reflect"
	"strings"
	"unsafe"
)

// XXX should be a utility routine
func (s *XLSuite) scratchFileName(c *C, rng *xr.PRNG, dirName string) (fileName string) {
	length := len(dirName)
	c.Assert(length > 0, Equals, true)
	if strings.HasSuffix(dirName, "/") {
		dirName = dirName[:length-1]
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
	pathToFile := s.scratchFileName(c, rng, "tmp")

	// XXX take care: this is ONE block
	data := make([]byte, BLOCK_SIZE)
	rng.NextBytes(&data)
	err := ioutil.WriteFile(pathToFile, data, 0644)
	c.Assert(err, IsNil)

	f, err := os.OpenFile(pathToFile, os.O_CREATE|os.O_RDWR, 0640)
	c.Assert(err, IsNil)

	// XXX Changing this from gm.MAP_PRIVATE to gm.MAP_SHARED made
	// the tests at the bottom succeed.  That is, changes made to
	// memory were written to disk by the Sync.
	inCore, err := gm.MapAt(0, f.Fd(), 0, 2*BLOCK_SIZE,
		gm.PROT_READ|gm.PROT_WRITE, gm.MAP_SHARED)
	c.Assert(err, IsNil)
	c.Assert(inCore, Not(IsNil))
	// The next succeeds, so it has grabbed that much memory ...
	c.Assert(len(inCore), Equals, 2*BLOCK_SIZE)

	// these are per-block flags
	boolz, err := inCore.IsResident()
	c.Assert(err, IsNil)
	c.Assert(boolz[0], Equals, true)

	// HACKING ABOUT ////////////////////////////////////////////////
	//for i := 0; i < BLOCK_SIZE; i++ {
	//	inCore[i] = byte(0)
	//}

	// inCore is an MMap
	mh := (*reflect.SliceHeader)(unsafe.Pointer(&inCore))

	const SIZEOF_UINT64 = 8 // bytes

	type MMap64 []uint64
	var inCore64 MMap64
	ih := (*reflect.SliceHeader)(unsafe.Pointer(&inCore64))

	ih.Data = mh.Data
	ih.Len = mh.Len / SIZEOF_UINT64
	ih.Cap = mh.Cap / SIZEOF_UINT64

	inCore[0] = 0x7f
	// END HACKING //////////////////////////////////////////////////

	// This succeeds, so the mapping from disk succeeded.
	c.Assert(bytes.Equal(inCore[1:BLOCK_SIZE], data[1:]), Equals, true)

	const (
		ASCII_A = byte(64)
	)
	inCore[BLOCK_SIZE-1] = ASCII_A
	inCore.Sync(gm.MS_SYNC) // should block

	// With the change to gm.MAP_SHARED, this does not seem to be
	// necessary:
	//
	// if the Sync didn't flush the ASCII_A to disk, this should do it.
	//err = inCore.UnsafeUnmap()
	// c.Assert(err, IsNil)

	f.Close()

	data2, err := ioutil.ReadFile(pathToFile)
	c.Assert(err, IsNil)

	// if this succeeds, then the flush to disk succeeded
	c.Assert(data2[BLOCK_SIZE-1], Equals, ASCII_A)
}
