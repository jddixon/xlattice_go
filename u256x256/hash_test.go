package u256x256

// xlattice_go/hash_test.go

import (
	"encoding/hex"
	"fmt"
	"hash"
	"io/ioutil"
	. "launchpad.net/gocheck"
	"os"
	"path/filepath"
	"testing"
)

// gocheck tie-in /////////////////////
func Test(t *testing.T) { TestingT(t) }

type XLSuite struct{}

var _ = Suite(&XLSuite{})

// end gocheck setUp //////////////////

var (
	dataPath string
	uPath    string
	uInDir   string
	uTmpDir  string
)

// SETUP AND TEARDOWN ///////////////////////////////////////////////
func (s *XLSuite) setUpHashTest() {
	var err error
	if !PathExists(dataPath) {
		// MODE SUSPECT
		if err = os.MkdirAll(dataPath, 0775); err != nil {
			fmt.Printf("error creating %s: %v\n", dataPath, err)
		}
	}
	if !PathExists(uPath) {
		// MODE SUSPECT
		if err = os.MkdirAll(uPath, 0775); err != nil {
			fmt.Printf("error creating %s: %v\n", uPath, err)
		}
	}
	if !PathExists(uInDir) {
		// MODE SUSPECT
		if err = os.MkdirAll(uInDir, 0775); err != nil {
			fmt.Printf("error creating %s: %v\n", uInDir, err)
		}
	}
	if !PathExists(uTmpDir) {
		// MODE SUSPECT
		if err = os.MkdirAll(uTmpDir, 0775); err != nil {
			fmt.Printf("error creating %s: %v\n", uTmpDir, err)
		}
	}
}
func (s *XLSuite) teardownHashTest() {
	// arguably should remove the two directories
}

// UNIT TESTS ///////////////////////////////////////////////////////
func (s *XLSuite) doTestCopyAndPut(
	c *C, u *U256x256, digest hash.Hash, usingSHA1 bool) {
	//we are testing uLen, uKey, err = u.CopyAndPut3(path, key)

	// create a random file                   maxLen   minLen
	dLen, dPath := u.rng.NextDataFile(dataPath, 16*1024, 1)
	dKey, err := FileSHA3(dPath)
	c.Assert(err, Equals, nil) // actual, Equals, expected

	// invoke function
	uLen, uKey, err := u.CopyAndPut3(dPath, dKey)
	c.Assert(err, Equals, nil)
	c.Assert(dLen, Equals, uLen)
	c.Assert(dKey, Equals, uKey)

	// verify that original and copy both exist
	c.Assert(PathExists(dPath), Equals, true)
	xPath := u.GetPathForKey(uKey)
	c.Assert(PathExists(xPath), Equals, true)

	// HACK - SIMPLEST Keccak TEST VECTOR
	if !usingSHA1 {
		dKey, err = FileSHA3("abc")
		fmt.Printf("SHA3-256 for 'abc' is %s\n", dKey)
	}
	// END HACK
}
func (s *XLSuite) doTestExists(c *C, u *U256x256, digest hash.Hash) {
	//we are testing whether = u.Exists( key)

	dLen, dPath := u.rng.NextDataFile(dataPath, 16*1024, 1)
	dKey, err := FileSHA3(dPath)
	c.Assert(err, Equals, nil)
	uLen, uKey, err := u.CopyAndPut3(dPath, dKey)
	c.Assert(err, Equals, nil)
	c.Assert(dLen, Equals, uLen)
	kPath := u.GetPathForKey(uKey)
	c.Assert(true, Equals, PathExists(kPath))
	c.Assert(true, Equals, u.Exists(uKey))
	os.Remove(kPath)
	c.Assert(false, Equals, PathExists(kPath))
	c.Assert(false, Equals, u.Exists(uKey))
}
func (s *XLSuite) doTestFileLen(c *C, u *U256x256, digest hash.Hash) {
	//we are testing len = u.fileLen(key)

	dLen, dPath := u.rng.NextDataFile(dataPath, 16*1024, 1)
	dKey, err := FileSHA3(dPath)
	c.Assert(err, Equals, nil)
	uLen, uKey, err := u.CopyAndPut3(dPath, dKey)
	c.Assert(err, Equals, nil)
	c.Assert(dLen, Equals, uLen)
	kPath := u.GetPathForKey(uKey)
	_ = kPath // NOT USED
	length, err := u.FileLen(uKey)
	c.Assert(err, Equals, nil)
	c.Assert(dLen, Equals, length)
}

func (s *XLSuite) doTestFileHash(c *C, u *U256x256, digest hash.Hash) {
	// we are testing sha1Key = fileSHA3(path)
	dLen, dPath := u.rng.NextDataFile(dataPath, 16*1024, 1)
	data, err := ioutil.ReadFile(dPath)
	c.Assert(err, Equals, nil)
	c.Assert(dLen, Equals, int64(len(data)))
	digest.Write(data)
	hash := digest.Sum(nil)
	dKey := hex.EncodeToString(hash) // 'expected'
	var fKey string
	if len(dKey) == SHA1_LEN {
		fKey, err = FileSHA1(dPath) // 'actual'
	} else {
		c.Assert(len(dKey), Equals, SHA3_LEN)
		fKey, err = FileSHA3(dPath) // 'actual'
	}
	c.Assert(err, Equals, nil)
	c.Assert(fKey, Equals, dKey)
}

func (s *XLSuite) doTestGetPathForKey(
	c *C, u *U256x256, digest hash.Hash, usingSHA1 bool) {
	// we are testing path = GetPathForKey(key)

	dLen, dPath := u.rng.NextDataFile(dataPath, 16*1024, 1)
	var dKey, uKey string
	var uLen int64
	if usingSHA1 {
		dKey, _ = FileSHA1(dPath) // ERRORS IGNORED
		uLen, uKey, _ = u.CopyAndPut1(dPath, dKey)
	} else {
		dKey, _ = FileSHA3(dPath) // ERRORS IGNORED
		uLen, uKey, _ = u.CopyAndPut3(dPath, dKey)
	}
	c.Assert(uLen, Equals, dLen)
	kPath := u.GetPathForKey(uKey)

	// XXX implementation-dependent test
	expectedPath := fmt.Sprintf("%s/%s/%s/%s",
		u.path, uKey[0:2], uKey[2:4], uKey[4:])
	c.Assert(expectedPath, Equals, kPath)
}

func (s *XLSuite) doTestPut(
	c *C, u *U256x256, digest hash.Hash, usingSHA1 bool) {
	//we are testing (len,hash)  = put(inFile, key)

	var dLen, uLen int64
	var dPath, dKey, uKey string
	dLen, dPath = u.rng.NextDataFile(dataPath, 16*1024, 1)
	if usingSHA1 {
		dKey, _ = FileSHA1(dPath) // ERRORS IGNORED
	} else {
		dKey, _ = FileSHA3(dPath) // ERRORS IGNORED
	}
	data, _ := ioutil.ReadFile(dPath) // ERRORS IGNORED
	dupePath := filepath.Join(dataPath, dKey)
	_ = ioutil.WriteFile(dupePath, data, 0664) // ERRORS IGNORED
	if usingSHA1 {
		uLen, uKey, _ = u.Put1(dPath, dKey) // ERRORS IGNORED
	} else {
		uLen, uKey, _ = u.Put3(dPath, dKey) // ERRORS IGNORED
	}
	c.Assert(dLen, Equals, uLen)
	kPath := u.GetPathForKey(uKey)
	_ = kPath // NOT USED

	// inFile is renamed
	c.Assert(false, Equals, PathExists(dPath))
	c.Assert(true, Equals, u.Exists(uKey))
	var dupeLen int64
	var dupeKey string
	if usingSHA1 {
		dupeLen, dupeKey, _ = u.Put1(dupePath, dKey) // ERRORS IGNORED
	} else {
		dupeLen, dupeKey, _ = u.Put3(dupePath, dKey) // ERRORS IGNORED
	}
	c.Assert(uLen, Equals, dupeLen)
	// dupe file is deleted'
	c.Assert(uKey, Equals, dupeKey)
	c.Assert(false, Equals, PathExists(dupePath))
	c.Assert(true, Equals, u.Exists(uKey))
}

func (s *XLSuite) doTestPutData(
	c *C, u *U256x256, digest hash.Hash, usingSHA1 bool) {
	// we are testing (len,hash)  = putData3(data, key)

	var dPath, dKey, uKey string
	var dLen, uLen int64
	var err error

	dLen, dPath = u.rng.NextDataFile(dataPath, 16*1024, 1)
	if usingSHA1 {
		dKey, err = FileSHA1(dPath)
	} else {
		dKey, err = FileSHA3(dPath)
	}
	c.Assert(err, Equals, nil)
	data, err := ioutil.ReadFile(dPath)
	c.Assert(err, Equals, nil)
	c.Assert(int64(len(data)), Equals, dLen)

	if usingSHA1 {
		uLen, uKey, err = u.PutData1(data, dKey)
	} else {
		uLen, uKey, err = u.PutData3(data, dKey)
	}
	c.Assert(err, Equals, nil)
	c.Assert(dLen, Equals, uLen)
	c.Assert(dKey, Equals, uKey)
	c.Assert(true, Equals, u.Exists(uKey))
	xPath := u.GetPathForKey(uKey)
	c.Assert(true, Equals, PathExists(xPath))
}
