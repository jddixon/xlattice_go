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
type XLSuite struct {}
var _ = Suite(&XLSuite{})
// end gocheck setUp //////////////////

var (
	dataPath	string
	uPath		string
	uInDir		string
	uTmpDir 	string
)

// SETUP AND TEARDOWN ///////////////////////////////////////////////
func (s *XLSuite) setUpHashTest() {
	var err error
	if ! PathExists(dataPath) {
		// MODE SUSPECT
		if err = os.MkdirAll(dataPath, 0775) ; err != nil {
			fmt.Printf("error creating %s: %v\n", dataPath, err)
		}
	}
	if ! PathExists(uPath) {
		// MODE SUSPECT
		if err = os.MkdirAll(uPath, 0775) ; err != nil {
			fmt.Printf("error creating %s: %v\n", uPath, err)
		}
	}
	if ! PathExists(uInDir) {
		// MODE SUSPECT
		if err = os.MkdirAll(uInDir, 0775) ; err != nil {
			fmt.Printf("error creating %s: %v\n", uInDir, err)
		}
	}
	if ! PathExists(uTmpDir) {
		// MODE SUSPECT
		if err = os.MkdirAll(uTmpDir, 0775) ; err != nil {
			fmt.Printf("error creating %s: %v\n", uTmpDir, err)
		}
	}
}
func (s *XLSuite) teardownHashTest() {
	// arguably should remove the two directories
}

// UNIT TESTS ///////////////////////////////////////////////////////
func (s *XLSuite) doTestCopyAndPut(c *C, u *U256x256, digest hash.Hash) {
    //we are testing sha1Key = u.CopyAndPut3(path, key)

    // create a random file                   maxLen   minLen
    dLen, dPath	:= u.rng.NextDataFile(dataPath, 16*1024,    1)
    dKey, err   := FileSHA3(dPath)
	c.Assert(err, Equals, nil)			// actual, Equals, expected

    // invoke function
    uLen, uKey, err	:= u.CopyAndPut3(dPath, dKey)
	c.Assert(err, Equals, nil)
    c.Assert(dLen, Equals, uLen)
    c.Assert(dKey, Equals, uKey)

    // verify that original and copy both exist
    c.Assert(  PathExists(dPath), Equals, true )
    xPath := u.GetPathForKey( uKey  )
    c.Assert(  PathExists(xPath), Equals, true )

	// HACK - SIMPLEST Keccak TEST VECTOR
	dKey, err	= FileSHA3("abc")
	fmt.Printf("SHA3-256 for 'abc' is %s\n", dKey)		
	// END HACK
}
func (s *XLSuite) doTestExists(c *C, u *U256x256, digest hash.Hash) {
    //we are testing whether = u.Exists( key)

    dLen, dPath := u.rng.NextDataFile(dataPath, 16*1024,    1)
    dKey, err   := FileSHA3(dPath)			
	c.Assert(err, Equals, nil)
    uLen, uKey,err:= u.CopyAndPut3(dPath, dKey)	
	c.Assert(err, Equals, nil)
	c.Assert(dLen, Equals, uLen)
    kPath		:= u.GetPathForKey( uKey  )
    c.Assert( true, Equals,  PathExists(kPath) )
    c.Assert( true, Equals, u.Exists( uKey) )
    os.Remove(kPath)
    c.Assert(false, Equals, PathExists(kPath) )
    c.Assert(false, Equals,  u.Exists( uKey) )
}
func (s *XLSuite) doTestFileLen(c *C, u *U256x256, digest hash.Hash) {
    //we are testing len = u.fileLen(key)

    dLen, dPath := u.rng.NextDataFile(dataPath, 16*1024,    1)
    dKey, err   := FileSHA3(dPath)		
	c.Assert(err, Equals, nil)
    uLen,uKey,err	:= u.CopyAndPut3(dPath, dKey)
	c.Assert(err, Equals, nil)
    c.Assert(dLen, Equals, uLen)
    kPath		:= u.GetPathForKey( uKey  )
	_ = kPath									// NOT USED
	length,err	:= u.FileLen(uKey)			
	c.Assert(err, Equals, nil)
    c.Assert(dLen, Equals, length)
}

func (s *XLSuite) doTestFileHash(c *C, u *U256x256, digest hash.Hash) {
    // we are testing sha1Key = fileSHA3(path)
    dLen, dPath := u.rng.NextDataFile(dataPath, 16*1024,    1)
	data, err	:= ioutil.ReadFile(dPath)
	c.Assert(err, Equals, nil)
	c.Assert(dLen, Equals, int64(len(data)))			
    digest.Write(data)
	hash		:= digest.Sum(nil)
    dKey		:= hex.EncodeToString(hash)		// 'expected'
    fsha1,err	:= FileSHA3(dPath)				// 'actual'
	c.Assert(err,   Equals, nil)
    c.Assert(fsha1, Equals, dKey)				// * FAILS *
}

func (s *XLSuite) doTestGetPathForKey(c *C, u *U256x256, digest hash.Hash) {
    // we are testing path = GetPathForKey(key)

    dLen, dPath	:= u.rng.NextDataFile(dataPath, 16*1024,    1)
    dKey, _     := FileSHA3(dPath)				// ERRORS IGNORED
    uLen,uKey,_	:= u.CopyAndPut3(dPath, dKey)
	c.Assert(uLen, Equals, dLen)
    kPath		:= u.GetPathForKey( uKey  )

    // XXX implementation-dependent test
    expectedPath := fmt.Sprintf("%s/%s/%s/%s", 
							u.path, uKey[0:2], uKey[2:4], uKey [4:])
    c.Assert(expectedPath, Equals, kPath)
}

func (s *XLSuite) doTestPut(c *C, u *U256x256, digest hash.Hash) {
    //we are testing (len,hash)  = put(inFile, key)

    dLen, dPath := u.rng.NextDataFile(dataPath, 16*1024,    1)
    dKey, _     := FileSHA3(dPath)				// ERRORS IGNORED
	data, _		:= ioutil.ReadFile(dPath)		// ERRORS IGNORED
    dupePath	:= filepath.Join(dataPath, dKey)
	_ = ioutil.WriteFile(dupePath, data, 0664)	// ERRORS IGNORED

    uLen,uKey,_ := u.Put3(dPath, dKey)			// ERRORS IGNORED
	c.Assert( dLen, Equals, uLen)
    kPath		:= u.GetPathForKey( uKey  )
	_ = kPath									// NOT USED

    // inFile is renamed
    c.Assert(false, Equals, PathExists(dPath)  )
    c.Assert( true, Equals,  u.Exists( uKey) )

    dupeLen, dupeKey,_	:= u.Put3(dupePath, dKey)	// ERRORS IGNORED
	c.Assert(uLen, Equals, dupeLen)
    // dupe file is deleted'
    c.Assert( uKey, Equals, dupeKey )
    c.Assert(false, Equals, PathExists(dupePath)  )
    c.Assert( true, Equals,  u.Exists( uKey) )
}

func (s *XLSuite) doTestPutData(c *C, u *U256x256, digest hash.Hash) {
    // we are testing (len,hash)  = putData3(data, key)

    dLen, dPath := u.rng.NextDataFile(dataPath, 16*1024,    1)
    dKey,err    := FileSHA3(dPath)					
	c.Assert( err, Equals, nil)
	data,err	:= ioutil.ReadFile(dPath)			
	c.Assert( err, Equals, nil)
	c.Assert( int64(len(data)), Equals, dLen)

    uLen, uKey, err	:= u.PutData3(data, dKey)			
	c.Assert( err, Equals, nil)				// FAILS, wrong key
	c.Assert(dLen, Equals, uLen)
    c.Assert(dKey, Equals, uKey)
    c.Assert( true, Equals,  u.Exists( dKey))
    xPath := u.GetPathForKey( uKey  )
	_ = xPath										// NOT USED
}
