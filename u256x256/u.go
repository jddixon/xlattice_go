package u256x256

// xlattice_go/u256x256/u.go

import (
	"crypto/sha1"
	"code.google.com/p/go.crypto/sha3"
	"encoding/hex"
	"errors"
	"fmt"			// DEBUG
	"io/ioutil"
	"github.com/jddixon/xlattice_go/rnglib"
	"io"
	"os"
	"path/filepath"
	"time"
)

// ....x....1....x....2....x....3....x....4
const SHA1_NONE = "0000000000000000000000000000000000000000"

// ....x....1....x....2....x....3....x....4....x....5....x....6....
const SHA3_NONE = "0000000000000000000000000000000000000000000000000000000000000000"

const DEFAULT_BUFFER_SIZE = 256 * 256

// XXX THIS DOES NOT BELONG HERE.  It is used to make unique
// file names, but this results in a race condition.
func makeSimpleRNG() *rnglib.PRNG {
	t := time.Now().Unix()
	rng := rnglib.NewSimpleRNG(t)
	return rng
}


// PACKAGE-LEVEL FUNCTIONS //////////////////////////////////////////
func CopyFile(destName, srcName string) (written int64, err error) {
	var (
		src, dest *os.File
	)
	if src, err = os.Open(srcName); err != nil {
		return
	}
	defer src.Close()
	if dest, err = os.Create(destName); err != nil {
		return
	}
	defer dest.Close()
	return io.Copy(dest, src)		// returns written, err
}

// - FileSHA1 --------------------------------------------------------
// returns the SHA1 hash of the contents of a file
func FileSHA1(path string) (hash string, err error) {
	hash = SHA1_NONE
	if (path == "") || !PathExists(path) {
		err = errors.New("IllegalArgument: empty path or non-existent file")
		return
	}
	data2, err	:= ioutil.ReadFile(path)
	if err == nil {
		d2		:= sha1.New()
		d2.Write(data2)
		digest2 := d2.Sum(nil)
		hash	= hex.EncodeToString(digest2)
	}
	return
}	// GEEP

// - FileSHA3 --------------------------------------------------------
// returns the SHA3 hash of the contents of a file
func FileSHA3(path string) (hash string, err error) {
	hash = SHA3_NONE
	if (path == "") || !PathExists(path) {
		err = errors.New("IllegalArgument: empty path or non-existent file")
		return
	}

	// THIS METHOD DID NOT RETURN CORRECT RESULTS
//	d := sha3.NewKeccak256()
//	var f *os.File
//	if f, err = os.Open(path); err != nil {
//		return
//	}
//	defer f.Close()
//	buffer := make([]byte, DEFAULT_BUFFER_SIZE)
//	var count int
//	for true {
//		count, err = f.Read(buffer)
//		// DEBUG
//		fmt.Printf("FileSHA3 read loop; count = %v, err = %v\n", count, err)
//		// END
//		if count == 0 || err == io.EOF{
//			break
//		}
//		if err != nil { return }
//		d.Write(buffer)
//	}
//	if err == io.EOF {
//		err = nil
//	}
//	digest := d.Sum(nil)
//	hash = hex.EncodeToString(digest)

	// METHOD 2
	data2, err	:= ioutil.ReadFile(path)
	if err == nil {
		d2		:= sha3.NewKeccak256()
		d2.Write(data2)
		digest2 := d2.Sum(nil)
		hash	= hex.EncodeToString(digest2)
	}
	return
}	// GEEP

func PathExists(fName string) bool {
	if _, err := os.Stat(fName); os.IsNotExist(err) {
		return false
	}
	return true
}

// CLASS, so to speak ///////////////////////////////////////////////
type U256x256 struct {
	path	string					// all parameters are 
	rng		*rnglib.PRNG			//	... private
	inDir	string
	tmpDir	string
}

func New(path string) *U256x256 {
	var u U256x256
	u.path		= path					// XXX validate
	u.inDir		= filepath.Join(path, "in")
	u.tmpDir	= filepath.Join(path, "tmp")
	u.rng		= makeSimpleRNG()
	return &u
}

// - Exists ---------------------------------------------------------
func (u *U256x256) Exists(key string) bool {
	path := u.GetPathForKey(key)
	return PathExists(path)
}
// - FileLen --------------------------------------------------------
func (u *U256x256) FileLen(key string) (length int64, err error) {
	// XXX ERROR IF EMPTY KEY
	path	:= u.GetPathForKey(key)
	if path == "" {
		err = errors.New("IllegalArgument: no key specified")
	}
	info, _ := os.Stat(path)				// ERRORS IGNORED
	length	= info.Size()
	return
}
// - GetPathForKey --------------------------------------------------
// Returns a path to a file with the content key passed.
// XXX NEED TO RESPEC TO RETURN ERROR IF INVALID KEY (blank, wrong length, etc.
func (u *U256x256) GetPathForKey(key string) string {
	if key == "" {
		return key
	}
	topSubDir	:= key[0:2]
	lowerDir	:= key[2:4]
	return filepath.Join(u.path, topSubDir, lowerDir, key[4:])
}
//- copyAndPut3 -------------------------------------------------------
func (u *U256x256) CopyAndPut3(path, key string) (
								written int64, hash string, err error) {
	// the temporary file MUST be created on the same device
	// xxx POSSIBLE RACE CONDITION
	tmpFileName := filepath.Join(u.tmpDir, u.rng.NextFileName(16))
	for ; PathExists(tmpFileName) ; {
		tmpFileName = filepath.Join(u.tmpDir, u.rng.NextFileName(16))
	}
	written, err = CopyFile(tmpFileName, path)		// dest <== src
	if err == nil {
		written, hash, err = u.Put3(tmpFileName, key)
	}
	return
}

// - GetData3 --------------------------------------------------------
func (u *U256x256) GetData3(key string) (data []byte, err error) {
	var path string
	path = u.GetPathForKey(key)
	if !PathExists(path) {
		// XXX SHOULD BE PREDEFINED ERROR
		err = errors.New("IllegalArgument: file does not exist")
		return
	} else {
		var src *os.File
		if src, err = os.Open(path); err != nil {
			return
		}
		defer src.Close()
		var count int
		// XXX THIS WILL NOT WORK FOR LARGER FILES!  It will ignore
		//     anything over 64 KB
		data = make([]byte, DEFAULT_BUFFER_SIZE)
		count, err = src.Read(data)
		// XXX COUNT IS IGNORED
		_ = count
		return
	}
}

// - Put3 ------------------------------------------------------------
// tmp is the path to a local file which will be renamed into U (or deleted
// if it is already present in U)
// u.path is an absolute or relative path to a U directory organized 256x256
// key is an sha3 content hash.
// If the operation succeeds we return the length of the file (which must
// not be zero.  Otherwise we return 0.
// we don't do much checking
func (u *U256x256) Put3(inFile, key string) (length int64, hash string, err error) {
	hash, err = FileSHA3(inFile)
	if err != nil {	
		fmt.Printf("DEBUG: FileSHA3 returned error %v\n", err)
		return 
	}
	if hash != key {
		fmt.Printf("expected %s to have key %s, but the content key is %s\n",
			inFile, key, hash)
		err = errors.New("IllegalArgument: Put3: key does not match content")
		return
	}
	// XXX ERROR NOT HANDLED
	info, _ := os.Stat(inFile)
	length = info.Size()
	topSubDir := hash[0:2]
	lowerDir := hash[2:4]
	targetDir := filepath.Join(u.path, topSubDir, lowerDir)
	if !PathExists(targetDir) {
		// XXX ERROR NOT HANDLED; MODE IS SUSPECT
		_ = os.MkdirAll(targetDir, 0775)

	}
	fullishPath := filepath.Join(targetDir, key[4:])
	if PathExists(fullishPath) {
		// drop the temporary input file
		// XXX ERROR NOT HANDLED
		_ = os.Remove(inFile)
	} else {
		// rename the temporary file into U
		// XXX ERROR NOT HANDLED
		_ = os.Rename(inFile, fullishPath)
		// XXX ERROR NOT HANDLED
		_ = os.Chmod(fullishPath, 0444)
	}
	return
}

// - putData3 --------------------------------------------------------
func (u *U256x256) PutData3(data []byte, key string) (
							length int64, hash string, err error) {
	s := sha3.NewKeccak256()
	s.Write(data)
	hash = hex.EncodeToString(s.Sum(nil))
	if hash != key {
		fmt.Printf("expected data to have key %s, but content key is %s",
			key, hash)
		err = errors.New("content/key mismatch")
		return
	}
	length		= int64(len(data))
	topSubDir	:= hash[0:2]
	lowerDir	:= hash[2:4]
	targetDir	:= filepath.Join(u.path, topSubDir, lowerDir)
	if !PathExists(targetDir) {
		// XXX ERROR NOT HANDLED; MODE QUESTIONABLE
		_ = os.MkdirAll(targetDir, 0775)
	}
	fullishPath := filepath.Join(targetDir, key)
	if !PathExists(fullishPath) {
		var dest *os.File
		if dest, err = os.Create(fullishPath); err != nil {
			return
		}
		var count int
		defer dest.Close()
		// XXX ERROR NOT HANDLED
		count, err = dest.Write(data)
		length = int64(count)
	}
	return
} // GEEP

// SHA1 CODE ========================================================

// CopyAndPut1 ------------------------------------------------------
func (u *U256x256) CopyAndPut1(path, key string) (
								written int64, hash string, err error) {
	// the temporary file MUST be created on the same device
	// xxx POSSIBLE RACE CONDITION
	tmpFileName := filepath.Join(u.tmpDir, u.rng.NextFileName(16))
	for ; PathExists(tmpFileName) ; {
		tmpFileName = filepath.Join(u.tmpDir, u.rng.NextFileName(16))
	}
	written, err = CopyFile(tmpFileName, path)		// dest <== src
	if err == nil {
		written, hash, err = u.Put1(tmpFileName, key)
	}
	return
}

// - GetData1 --------------------------------------------------------
func (u *U256x256) GetData1(key string) (data []byte, err error) {
	var path string
	path = u.GetPathForKey(key)
	if !PathExists(path) {
		// XXX SHOULD BE PREDEFINED ERROR
		err = errors.New("IllegalArgument: file does not exist")
		return
	} else {
		var src *os.File
		if src, err = os.Open(path); err != nil {
			return
		}
		defer src.Close()
		var count int
		// XXX THIS WILL NOT WORK FOR LARGER FILES!  It will ignore
		//     anything over 64 KB
		data = make([]byte, DEFAULT_BUFFER_SIZE)
		count, err = src.Read(data)
		// XXX COUNT IS IGNORED
		_ = count
		return
	}
}

// - Put1 ------------------------------------------------------------
// tmp is the path to a local file which will be renamed into U (or deleted
// if it is already present in U)
// u.path is an absolute or relative path to a U directory organized 256x256
// key is an sha1 content hash.
// If the operation succeeds we return the length of the file (which must
// not be zero.  Otherwise we return 0.
// we don't do much checking
func (u *U256x256) Put1(inFile, key string) (length int64, hash string, err error) {
	hash, err = FileSHA1(inFile)
	if err != nil {	
		fmt.Printf("DEBUG: FileSHA1 returned error %v\n", err)
		return 
	}
	if hash != key {
		fmt.Printf("expected %s to have key %s, but the content key is %s\n",
			inFile, key, hash)
		err = errors.New("IllegalArgument: Put1: key does not match content")
		return
	}
	// XXX ERROR NOT HANDLED
	info, _ := os.Stat(inFile)
	length = info.Size()
	topSubDir := hash[0:2]
	lowerDir := hash[2:4]
	targetDir := filepath.Join(u.path, topSubDir, lowerDir)
	if !PathExists(targetDir) {
		// XXX ERROR NOT HANDLED; MODE IS SUSPECT
		_ = os.MkdirAll(targetDir, 0775)

	}
	fullishPath := filepath.Join(targetDir, key[4:])
	if PathExists(fullishPath) {
		// drop the temporary input file
		// XXX ERROR NOT HANDLED
		_ = os.Remove(inFile)
	} else {
		// rename the temporary file into U
		// XXX ERROR NOT HANDLED
		_ = os.Rename(inFile, fullishPath)
		// XXX ERROR NOT HANDLED
		_ = os.Chmod(fullishPath, 0444)
	}
	return
}

// PutData1 ---------------------------------------------------------
func (u *U256x256) PutData1(data []byte, key string) (
							length int64, hash string, err error) {
	s := sha1.New()
	s.Write(data)
	hash = hex.EncodeToString(s.Sum(nil))
	if hash != key {
		fmt.Printf("expected data to have key %s, but content key is %s",
			key, hash)
		err = errors.New("content/key mismatch")
		return
	}
	length		= int64(len(data))
	topSubDir	:= hash[0:2]
	lowerDir	:= hash[2:4]
	targetDir	:= filepath.Join(u.path, topSubDir, lowerDir)
	if !PathExists(targetDir) {
		// XXX ERROR NOT HANDLED; MODE QUESTIONABLE
		_ = os.MkdirAll(targetDir, 0775)
	}
	fullishPath := filepath.Join(targetDir, key)
	if !PathExists(fullishPath) {
		var dest *os.File
		if dest, err = os.Create(fullishPath); err != nil {
			return
		}
		var count int
		defer dest.Close()
		// XXX ERROR NOT HANDLED
		count, err = dest.Write(data)
		length = int64(count)
	}
	return
} // GEEP
