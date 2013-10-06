package merkletree

// xlattice_go/util/merkletree/merkle_tree_test.go

import (
	//"code.google.com/p/go.crypto/sha3"
	//"crypto/sha1"
	"encoding/hex"
	"fmt"
	xr "github.com/jddixon/xlattice_go/rnglib"
	xu "github.com/jddixon/xlattice_go/util"
	// "io/ioutil"
	. "launchpad.net/gocheck"
	//"strings"
)

const (
	ONE			= 1
	FOUR		= 4
	MAX_NAME_LEN= 8
)

// UTILITY FUNCTIONS ================================================


// SUBTESTS =========================================================
func (s *XLSuite) doTestParser(c *C, rng *xr.PRNG, usingSHA1 bool) {

	var tHash []byte
	if usingSHA1 {
		tHash = make([]byte, SHA1_LEN)
	} else {
		tHash = make([]byte, SHA3_LEN)
	}
	rng.NextBytes(&tHash)				// not really a hash, of course
	sHash := hex.EncodeToString(tHash)	// string form of tHash

	dirName := rng.NextFileName(8) + "/"
	nameWithoutSlash := dirName[0:len(dirName)-1]

	indent := rng.Intn(4)
	var lSpaces, rSpaces string
	for i := 0; i < indent; i++ {
		lSpaces += "  "		// on the left
		rSpaces += " "		// on the right
	}

	// TEST FIRST LINE PARSER -----------------------------
	line := lSpaces + sHash + " " + dirName + rSpaces

	indent2, treeHash2, dirName2, err := ParseFirstLine(line)
	c.Assert(err, IsNil)
	c.Assert(indent2, Equals, indent)
	c.Assert(xu.SameBytes(treeHash2, tHash), Equals, true)
	c.Assert(dirName2, Equals, nameWithoutSlash)

	// TEST OTHER LINE PARSER -----------------------------
	yesIsDir := rng.NextBoolean()
	if yesIsDir {
		line = lSpaces + sHash + " " + dirName + rSpaces
	} else {
		line = lSpaces + sHash + " " + nameWithoutSlash + rSpaces
	}
	nodeDepth, nodeHash, nodeName, isDir, err := ParseOtherLine(line)
	c.Assert(err, IsNil)
	c.Assert(nodeDepth, Equals, indent)
	c.Assert(xu.SameBytes(nodeHash, tHash), Equals, true)
	c.Assert(nodeName, Equals, nameWithoutSlash)
	c.Assert(isDir, Equals, yesIsDir)	
}
func (s *XLSuite) TestMerkleTree(c *C) {
	if VERBOSITY > 0 {
		fmt.Println("TEST_MERKLE_LEAF")
	}
	rng := xr.MakeSimpleRNG()

	s.doTestParser(c, rng, true)		// usingSHA1
	s.doTestParser(c, rng, false)		// using SHA3 instead
}
