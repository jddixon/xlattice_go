package merkletree

import (
	"bytes"
	"code.google.com/p/go.crypto/sha3"
	"crypto/sha1"
	"fmt"
	xr "github.com/jddixon/xlattice_go/rnglib"
	"io/ioutil"
	. "launchpad.net/gocheck"
	"strings"
)

func (s *XLSuite) doTestSimpleConstructor(c *C, rng *xr.PRNG, usingSHA1 bool) {
	fileName := rng.NextFileName(8)
	leaf1, err := NewMerkleLeaf(fileName, nil, usingSHA1)
	c.Assert(err, IsNil)
	c.Assert(leaf1.Name(), Equals, fileName)
	c.Assert(len(leaf1.GetHash()), Equals, 0)
	c.Assert(leaf1.UsingSHA1(), Equals, usingSHA1)

	fileName2 := rng.NextFileName(8)
	for fileName2 == fileName {
		fileName2 = rng.NextFileName(8)
	}
	leaf2, err := NewMerkleLeaf(fileName2, nil, usingSHA1)
	c.Assert(err, IsNil)
	c.Assert(leaf2.Name(), Equals, fileName2)

	c.Assert(leaf1.Equal(leaf1), Equals, true)
	c.Assert(leaf1.Equal(leaf2), Equals, false)
}

func (s *XLSuite) doTestSHA(c *C, rng *xr.PRNG, usingSHA1 bool) {

	var hash, fHash []byte
	var sHash string

	// name guaranteed to be unique
	length, pathToFile := rng.NextDataFile("tmp", 1024, 256)
	data, err := ioutil.ReadFile(pathToFile)
	c.Assert(err, IsNil)
	c.Assert(len(data), Equals, length)
	parts := strings.Split(pathToFile, "/")
	c.Assert(len(parts), Equals, 2)
	fileName := parts[1]

	if usingSHA1 {
		sha := sha1.New()
		sha.Write(data)
		hash = sha.Sum(nil)
		fHash, err = SHA1File(pathToFile)
	} else {
		sha := sha3.NewKeccak256()
		sha.Write(data)
		hash = sha.Sum(nil)
		fHash, err = SHA3File(pathToFile)
	}
	c.Assert(err, IsNil)
	c.Assert(bytes.Equal(hash, fHash), Equals, true)

	ml, err := CreateMerkleLeafFromFileSystem(pathToFile, fileName, usingSHA1)
	c.Assert(err, IsNil)
	c.Assert(ml.Name(), Equals, fileName)
	c.Assert(bytes.Equal(ml.GetHash(), hash), Equals, true)
	c.Assert(ml.UsingSHA1(), Equals, usingSHA1)

	// TODO: test ToString
	_ = sHash // TODO

}
func (s *XLSuite) TestMerkleLeaf(c *C) {
	if VERBOSITY > 0 {
		fmt.Println("TEST_MERKLE_LEAF")
	}
	rng := xr.MakeSimpleRNG()
	s.doTestSimpleConstructor(c, rng, true)  // using SHA1
	s.doTestSimpleConstructor(c, rng, false) // not using SHA1

}
