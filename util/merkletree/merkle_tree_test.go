package merkletree

// xlattice_go/util/merkletree/merkle_tree_test.go

import (
	"bytes"
	"code.google.com/p/go.crypto/sha3"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	xr "github.com/jddixon/xlattice_go/rnglib"
	xf "github.com/jddixon/xlattice_go/util/lfs"
	"hash"
	"io/ioutil"
	. "launchpad.net/gocheck"
	"os"
	"path"
	"strings"
)

const (
	ONE          = 1
	FOUR         = 4
	MAX_NAME_LEN = 8
)

// UTILITY FUNCTIONS ================================================

func (s *XLSuite) getTwoUniqueDirectoryNames(c *C, rng *xr.PRNG) (
	string, string) {

	dirName1 := rng.NextFileName(MAX_NAME_LEN)
	dirName2 := rng.NextFileName(MAX_NAME_LEN)
	for dirName2 == dirName1 {
		dirName2 = rng.NextFileName(MAX_NAME_LEN)
	}
	return dirName1, dirName2
}

// Return a populated test directory with the name etc specified.
// If a directory of that name exists, it is deleted.

func (s *XLSuite) makeOneNamedTestDirectory(c *C, rng *xr.PRNG,
	name string, depth, width int) string {

	name = strings.TrimSpace(name)
	dirPath := fmt.Sprintf("tmp/%s", name)
	if strings.Contains(dirPath, "..") {
		msg := fmt.Sprintf("directory name '%s' contains a double-dot\n", name)
		panic(msg)
	}
	err := os.RemoveAll(dirPath)
	c.Assert(err, IsNil)
	//                                     max/minLen
	rng.NextDataDir(dirPath, depth, width, 32, 1)
	return dirPath
}

// Create and populate two test diretories.

func (s *XLSuite) makeTwoTestDirectories(c *C, rng *xr.PRNG,
	depth, width int) (string, string, string, string) {

	dirName1 := rng.NextFileName(MAX_NAME_LEN)
	dirPath1 := s.makeOneNamedTestDirectory(c, rng, dirName1, depth, width)

	dirName2 := rng.NextFileName(MAX_NAME_LEN)
	for dirName2 == dirName1 {
		dirName2 = rng.NextFileName(MAX_NAME_LEN)
	}
	dirPath2 := s.makeOneNamedTestDirectory(c, rng, dirName2, depth, width)

	return dirName1, dirPath1, dirName2, dirPath2
}

func (s *XLSuite) verifyLeafSHA(c *C, rng *xr.PRNG,
	node MerkleNodeI, pathToFile string, usingSHA1 bool) {

	c.Assert(node.IsLeaf(), Equals, true)
	found, err := xf.PathExists(pathToFile)
	c.Assert(err, IsNil)
	c.Assert(found, Equals, true)
	data, err := ioutil.ReadFile(pathToFile)
	c.Assert(err, IsNil)
	c.Assert(data, NotNil)

	var sha hash.Hash
	if usingSHA1 {
		sha = sha1.New()
	} else {
		sha = sha3.NewKeccak256()
	}
	sha.Write(data)
	sum := sha.Sum(nil)
	c.Assert(node.GetHash(), DeepEquals, sum)
}
func (s *XLSuite) verifyTreeSHA(c *C, rng *xr.PRNG,
	n MerkleNodeI, pathToNode string, usingSHA1 bool) {

	c.Assert(n.IsLeaf(), Equals, false)
	node := n.(*MerkleTree)

	if node.nodes == nil {
		c.Assert(node.GetHash(), Equals, nil)
	} else {
		hashCount := 0
		var sha hash.Hash
		if usingSHA1 {
			sha = sha1.New()
		} else {
			sha = sha3.NewKeccak256()
		}
		for i := 0; i < len(node.nodes); i++ {
			n := node.nodes[i]
			pathToFile := path.Join(pathToNode, n.Name())
			if n.IsLeaf() {
				s.verifyLeafSHA(c, rng, n, pathToFile, usingSHA1)
			} else if !n.IsLeaf() {
				s.verifyTreeSHA(c, rng, n, pathToFile, usingSHA1)
			} else {
				c.Error("unknown node type!")
			}
			if n.GetHash() != nil {
				hashCount += 1
				sha.Write(n.GetHash())
			}
		}
		if hashCount == 0 {
			c.Assert(node.GetHash(), IsNil)
		} else {
			c.Assert(node.GetHash(), DeepEquals, sha.Sum(nil))
		}
	}
}

// PARSER TESTS =====================================================
func (s *XLSuite) doTestParser(c *C, rng *xr.PRNG, usingSHA1 bool) {

	var tHash []byte
	if usingSHA1 {
		tHash = make([]byte, SHA1_LEN)
	} else {
		tHash = make([]byte, SHA3_LEN)
	}
	rng.NextBytes(&tHash)              // not really a hash, of course
	sHash := hex.EncodeToString(tHash) // string form of tHash

	dirName := rng.NextFileName(8) + "/"
	nameWithoutSlash := dirName[0 : len(dirName)-1]

	indent := rng.Intn(4)
	var lSpaces, rSpaces string
	for i := 0; i < indent; i++ {
		lSpaces += "  " // on the left
		rSpaces += " "  // on the right
	}

	// TEST FIRST LINE PARSER -----------------------------
	line := lSpaces + sHash + " " + dirName + rSpaces

	indent2, treeHash2, dirName2, err := ParseFirstLine(line)
	c.Assert(err, IsNil)
	c.Assert(indent2, Equals, indent)
	c.Assert(bytes.Equal(treeHash2, tHash), Equals, true)
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
	c.Assert(bytes.Equal(nodeHash, tHash), Equals, true)
	c.Assert(nodeName, Equals, nameWithoutSlash)
	c.Assert(isDir, Equals, yesIsDir)
}

func (s *XLSuite) TestParser(c *C) {
	if VERBOSITY > 0 {
		fmt.Println("TEST_PARSER")
	}
	rng := xr.MakeSimpleRNG()

	s.doTestParser(c, rng, true)  // usingSHA1
	s.doTestParser(c, rng, false) // using SHA3 instead
} // GEEP

// OTHER UNIT TESTS /////////////////////////////////////////////////

func (s *XLSuite) TestPathlessUnboundConstructor(c *C) {
	if VERBOSITY > 0 {
		fmt.Println("TEST_PATHLESS_UNBOUND_CONSTRUCTOR")
	}
	rng := xr.MakeSimpleRNG()
	s.doTestPathlessUnboundConstructor(c, rng, true)  // usingSHA1
	s.doTestPathlessUnboundConstructor(c, rng, false) // not
}

func (s *XLSuite) doTestPathlessUnboundConstructor(c *C, rng *xr.PRNG, usingSHA1 bool) {

	dirName1, dirName2 := s.getTwoUniqueDirectoryNames(c, rng)

	tree1, err := NewNewMerkleTree(dirName1, true)
	c.Assert(err, IsNil)
	c.Assert(tree1, NotNil)
	c.Assert(tree1.Name(), Equals, dirName1)
	c.Assert(len(tree1.GetHash()), Equals, 0)

	// DEBUG -- think this through before removing this code
	// XXX WITHOUT THIS FIX MerkleNode.Equal() sees the serializations
	// as different because one has a hash length of zero and the
	// other a hash length of 20.
	null := make([]byte, SHA1_LEN)
	tree1.SetHash(null)
	// END

	tree2, err := NewNewMerkleTree(dirName2, true)
	c.Assert(err, IsNil)
	c.Assert(tree2, NotNil)
	c.Assert(tree2.Name(), Equals, dirName2)

	// these tests remain skimpy
	c.Assert(tree1.Equal(tree1), Equals, true)
	c.Assert(tree1.Equal(tree2), Equals, false)

	tree1Str, err := tree1.ToString("")
	c.Assert(len(tree1Str) > 0, Equals, true)

	// there should be no indent on the first line
	c.Assert(tree1Str[0:1], Not(Equals), " ")

	// no extra lines should be added
	lines := strings.Split(tree1Str, "\r\n")
	// this split generates an extra blank line, because the serialization
	// ends with CR-LF
	lineCount := len(lines)
	if lineCount > 0 {
		if lines[lineCount-1] == "" {
			lines = lines[:lineCount-1]
		}
	}
	// there should only be a header line
	c.Assert(len(lines), Equals, 1)

	tree1Rebuilt, err := ParseMerkleTree(tree1Str)
	c.Assert(err, IsNil)

	// compare at the string level
	t1RStr, err := tree1Rebuilt.ToString("")
	c.Assert(err, IsNil)
	c.Assert(t1RStr, Equals, tree1Str)
	c.Assert(tree1.Equal(tree1Rebuilt), Equals, true)
} // GEEP

// ------------------------------------------------------------------
func (s *XLSuite) TestBoundFlatDirs(c *C) {
	if VERBOSITY > 0 {
		fmt.Println("TEST_BOUND_FLAT_DIRS")
	}
	rng := xr.MakeSimpleRNG()
	s.doTestBoundFlatDirs(c, rng, true)
	s.doTestBoundFlatDirs(c, rng, false)
} // GEEP

func (s *XLSuite) doTestBoundFlatDirs(c *C, rng *xr.PRNG, usingSHA1 bool) {

	// test directory is single level, with four data files"""
	dirName1, dirPath1, dirName2, dirPath2 := s.makeTwoTestDirectories(
		c, rng, ONE, FOUR)
	tree1, err := CreateMerkleTreeFromFileSystem(dirPath1, usingSHA1, nil, nil)
	c.Assert(err, IsNil)
	c.Assert(tree1.Name(), Equals, dirName1)
	nodes1 := tree1.nodes
	c.Assert(nodes1, NotNil)
	c.Assert(len(nodes1), Equals, FOUR)
	s.verifyTreeSHA(c, rng, tree1, dirPath1, usingSHA1)

	tree2, err := CreateMerkleTreeFromFileSystem(dirPath2, usingSHA1, nil, nil)
	c.Assert(err, IsNil)
	c.Assert(tree2.Name(), Equals, dirName2)
	nodes2 := tree2.nodes
	c.Assert(nodes2, NotNil)
	c.Assert(len(nodes2), Equals, FOUR)
	s.verifyTreeSHA(c, rng, tree2, dirPath2, usingSHA1)

	c.Assert(tree1.Equal(tree1), Equals, true)
	c.Assert(tree1.Equal(tree2), Equals, false)
	c.Assert(tree1.Equal(""), Equals, false)

	tree1Str, err := tree1.ToString("")
	c.Assert(err, IsNil)
	c.Assert(len(tree1Str) > 0, Equals, true)

	tree1Rebuilt, err := ParseMerkleTree(tree1Str)
	c.Assert(err, IsNil)

	// compare at the string level
	t1RStr, err := tree1Rebuilt.ToString("")
	c.Assert(err, IsNil)
	c.Assert(t1RStr, Equals, tree1Str)

	c.Assert(tree1.Equal(tree1Rebuilt), Equals, true)
} // FOO

// ------------------------------------------------------------------
func (s *XLSuite) TestBoundNeedleDirs(c *C) {
	if VERBOSITY > 0 {
		fmt.Println("TEST_BOUND_NEEDLE_DIRS")
	}
	rng := xr.MakeSimpleRNG()
	s.doTestBoundNeedleDirs(c, rng, true)
	s.doTestBoundNeedleDirs(c, rng, false)
} // GEEP
func (s *XLSuite) doTestBoundNeedleDirs(c *C, rng *xr.PRNG, usingSHA1 bool) {

	//test directories four deep with one data file at the lowest level"""
	dirName1, dirPath1, dirName2, dirPath2 := s.makeTwoTestDirectories(
		c, rng, FOUR, ONE)

	tree1, err := CreateMerkleTreeFromFileSystem(dirPath1, usingSHA1, nil, nil)
	c.Assert(err, IsNil)
	c.Assert(tree1.Name(), Equals, dirName1)
	nodes1 := tree1.nodes
	c.Assert(nodes1, NotNil)
	c.Assert(len(nodes1), Equals, ONE)
	s.verifyTreeSHA(c, rng, tree1, dirPath1, usingSHA1)

	tree2, err := CreateMerkleTreeFromFileSystem(dirPath2, usingSHA1, nil, nil)
	c.Assert(err, IsNil)
	c.Assert(tree2.Name(), Equals, dirName2)
	nodes2 := tree2.nodes
	c.Assert(nodes2, NotNil)
	c.Assert(len(nodes2), Equals, ONE)
	s.verifyTreeSHA(c, rng, tree2, dirPath2, usingSHA1)

	tree1Str, err := tree1.ToString("")
	c.Assert(err, IsNil)
	c.Assert(len(tree1Str) > 0, Equals, true)

	tree1Rebuilt, err := ParseMerkleTree(tree1Str)
	c.Assert(err, IsNil)

	//       # DEBUG
	//       print "NEEDLEDIR TREE1:\n" + tree1Str
	//       print "REBUILT TREE1:\n" + tree1Rebuilt.ToString()
	//       # END
	c.Assert(tree1.Equal(tree1Rebuilt), Equals, true) // FAILS
}

// ==================================================================
// BUGS IN THE PYTHON IMPLEMENTATION
// ==================================================================

func (s *XLSuite) TestGrayBoxesBug(c *C) {
	if VERBOSITY > 0 {
		fmt.Println("TEST_GRAY_BOXES_BUG")
	}

	serialization := "721a08022dd26e7be98b723f26131786fd2c0dc3 grayboxes.com/\r\n" +
		"  fcd3973c66230b9078a86a5642b4c359fe72d7da images/\r\n" +
		"    15e47f4eb55197e1bfffae897e9d5ce4cba49623 grayboxes.gif\r\n" +
		"  2477b9ea649f3f30c6ed0aebacfa32cb8250f3df index.html\r\n"

	// create from string array ----------------------------------
	ss := strings.Split(serialization, "\r\n")
	lineCount := len(ss)
	c.Assert(lineCount > 1, Equals, true)
	ss = ss[:lineCount-1]
	c.Assert(len(ss), Equals, 4)

	tree2, err := ParseMerkleTreeFromStrings(&ss)
	c.Assert(err, IsNil)
	ser2, err := tree2.ToString("")
	c.Assert(err, IsNil)
	c.Assert(ser2, Equals, serialization) // XXX FAILS

	// create from serialization ---------------------------------
	tree1, err := ParseMerkleTree(serialization)
	c.Assert(err, IsNil)
	c.Assert(tree1, NotNil)

	ser1, err := tree1.ToString("")
	c.Assert(err, IsNil)
	c.Assert(ser1, Equals, serialization)

	c.Assert(tree1.Equal(tree2), Equals, true)
}

// ------------------------------------------------------------------

func (s *XLSuite) TestXLatticeBug(c *C) {
	if VERBOSITY > 0 {
		fmt.Println("TEST_XLATTICE_BUG")
	}

	// This test relies on dat.xlattice.org being locally present
	// and an internally consistent merkleization.
	serialization, err := ioutil.ReadFile("./dat.xlattice.org")
	c.Assert(err, IsNil)
	lines := string(serialization)

	// create from serialization ---------------------------------
	tree1, err := ParseMerkleTree(lines)
	c.Assert(err, IsNil)

	//       # DEBUG
	//       print "tree1 has %d nodes" % len(tree1.nodes)
	//       with open("junk.tree1", "w") as t {
	//           t.write( tree1.String() )
	//       # END

	ser1, err := tree1.ToString("")
	c.Assert(err, IsNil)
	c.Assert(ser1, Equals, lines)

	// create from string array ----------------------------------
	ss := strings.Split(lines, "\r\n")
	lineCount := len(ss)
	c.Assert(lineCount, Equals, 2512) // one too many
	ss = ss[:lineCount-1]             // so we discard the last
	c.Assert(len(ss), Equals, 2511)

	tree2, err := ParseMerkleTreeFromStrings(&ss)
	c.Assert(err, IsNil)

	ser2, err := tree2.ToString("") // no extra indent
	c.Assert(err, IsNil)
	c.Assert(ser2, Equals, lines)

	c.Assert(tree1.Equal(tree2), Equals, true)
}
