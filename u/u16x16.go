package u16x16

// xlattice_go/u/u16x16/u.go

import (
	"code.google.com/p/go.crypto/sha3"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt" // DEBUG
	"github.com/jddixon/xlattice_go/rnglib"
	xf "github.com/jddixon/xlattice_go/util/lfs"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

const SHA1_LEN = 40 // length of hex version
const SHA3_LEN = 64

//                 ....x....1....x....2....x....3....x....4
const SHA1_NONE = "0000000000000000000000000000000000000000"

//                 ....x....1....x....2....x....3....x....4....x....5....x....6....
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
	return io.Copy(dest, src) // returns written, err
}

// - FileSHA1 --------------------------------------------------------
// returns the SHA1 hash of the contents of a file
func FileSHA1(path string) (hash string, err error) {
	var data2 []byte
	hash = SHA1_NONE
	found, err := xf.PathExists(path)
	if err == nil && !found {
		err = errors.New("IllegalArgument: empty path or non-existent file")
	}
	if err == nil {
		data2, err = ioutil.ReadFile(path)
	}
	if err == nil {
		d2 := sha1.New()
		d2.Write(data2)
		digest2 := d2.Sum(nil)
		hash = hex.EncodeToString(digest2)
	}
	return
}

// - FileSHA3 --------------------------------------------------------
// returns the SHA3 hash of the contents of a file
func FileSHA3(path string) (hash string, err error) {
	var data2 []byte

	hash = SHA3_NONE
	found, err := xf.PathExists(path)
	if err == nil && !found {
		err = errors.New("IllegalArgument: empty path or non-existent file")
	}

	if err == nil {
		data2, err = ioutil.ReadFile(path)
	}
	if err == nil {
		d2 := sha3.NewKeccak256()
		d2.Write(data2)
		digest2 := d2.Sum(nil)
		hash = hex.EncodeToString(digest2)
	}
	return
}

// CLASS, so to speak ///////////////////////////////////////////////
type U16x16 struct {
	path   string       // all parameters are
	rng    *rnglib.PRNG //	... private
	inDir  string
	tmpDir string
}

func New(path string) *U16x16 {
	var u U16x16
	u.path = path // XXX validate
	u.inDir = filepath.Join(path, "in")
	u.tmpDir = filepath.Join(path, "tmp")
	u.rng = makeSimpleRNG()
	return &u
}

// - Exists ---------------------------------------------------------
func (u *U16x16) Exists(key string) bool {
	path := u.GetPathForKey(key)
	found, _ := xf.PathExists(path) // err ignored
	return found
}

// - FileLen --------------------------------------------------------
func (u *U16x16) FileLen(key string) (length int64, err error) {
	// XXX ERROR IF EMPTY KEY
	path := u.GetPathForKey(key)
	if path == "" {
		err = errors.New("IllegalArgument: no key specified")
	}
	info, _ := os.Stat(path) // ERRORS IGNORED
	length = info.Size()
	return
}

// - GetPathForKey --------------------------------------------------
// Returns a path to a file with the content key passed.
// XXX NEED TO RESPEC TO RETURN ERROR IF INVALID KEY (blank, wrong length, etc.
func (u *U16x16) GetPathForKey(key string) string {
	if key == "" {
		return key
	}
	topSubDir := key[0:1]
	lowerDir := key[1:2]
	return filepath.Join(u.path, topSubDir, lowerDir, key[2:])
}

//- copyAndPut3 -------------------------------------------------------
func (u *U16x16) CopyAndPut3(path, key string) (
	written int64, hash string, err error) {
	// the temporary file MUST be created on the same device
	// xxx POSSIBLE RACE CONDITION
	tmpFileName := filepath.Join(u.tmpDir, u.rng.NextFileName(16))
	found, _ := xf.PathExists(tmpFileName) // XXX error ignored
	for found {
		tmpFileName = filepath.Join(u.tmpDir, u.rng.NextFileName(16))
		found, _ = xf.PathExists(tmpFileName)
	}
	written, err = CopyFile(tmpFileName, path) // dest <== src
	if err == nil {
		written, hash, err = u.Put3(tmpFileName, key)
	}
	return
}

// - GetData3 --------------------------------------------------------
func (u *U16x16) GetData3(key string) (data []byte, err error) {
	var path string
	path = u.GetPathForKey(key)
	found, err := xf.PathExists(path)
	if err == nil && !found {
		err = FileNotFound
	}
	if err == nil {
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
	}
	return
}

// - Put3 ------------------------------------------------------------
// tmp is the path to a local file which will be renamed into U (or deleted
// if it is already present in U)
// u.path is an absolute or relative path to a U directory organized 16x16
// key is an sha3 content hash.
// If the operation succeeds we return the length of the file (which must
// not be zero.  Otherwise we return 0.
// we don't do much checking
func (u *U16x16) Put3(inFile, key string) (
	length int64, hash string, err error) {

	var fullishPath string

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
	info, err := os.Stat(inFile)
	if err != nil {
		return
	}
	length = info.Size()
	topSubDir := hash[0:1]
	lowerDir := hash[1:2]
	targetDir := filepath.Join(u.path, topSubDir, lowerDir)
	found, err := xf.PathExists(targetDir)
	if err == nil && !found {
		// XXX MODE IS SUSPECT
		err = os.MkdirAll(targetDir, 0775)
	}
	if err == nil {
		var found bool

		fullishPath = filepath.Join(targetDir, key[2:])
		found, err = xf.PathExists(fullishPath)
		if err == nil {
			if found {
				// drop the temporary input file
				err = os.Remove(inFile)
			} else {
				// rename the temporary file into U
				err = os.Rename(inFile, fullishPath)
			}
		}
	}
	if err == nil {
		err = os.Chmod(fullishPath, 0444)
	}
	return
}

// - putData3 --------------------------------------------------------
func (u *U16x16) PutData3(data []byte, key string) (length int64, hash string, err error) {
	s := sha3.NewKeccak256()
	s.Write(data)
	hash = hex.EncodeToString(s.Sum(nil))
	if hash != key {
		fmt.Printf("expected data to have key %s, but content key is %s",
			key, hash)
		err = errors.New("content/key mismatch")
		return
	}
	length = int64(len(data))
	topSubDir := hash[0:1]
	lowerDir := hash[1:2]
	targetDir := filepath.Join(u.path, topSubDir, lowerDir)
	found, err := xf.PathExists(targetDir)
	if err == nil && !found {
		err = os.MkdirAll(targetDir, 0775)
	}
	fullishPath := filepath.Join(targetDir, key[2:])
	found, err = xf.PathExists(fullishPath)
	if !found {
		var dest *os.File
		dest, err = os.Create(fullishPath)
		if err == nil {
			var count int
			defer dest.Close()
			count, err = dest.Write(data)
			if err == nil {
				length = int64(count)
			}
		}
	}
	return
}

// SHA1 CODE ========================================================

// CopyAndPut1 ------------------------------------------------------
func (u *U16x16) CopyAndPut1(path, key string) (
	written int64, hash string, err error) {
	// the temporary file MUST be created on the same device
	// xxx POSSIBLE RACE CONDITION
	tmpFileName := filepath.Join(u.tmpDir, u.rng.NextFileName(16))
	found, err := xf.PathExists(tmpFileName)
	for found {
		tmpFileName = filepath.Join(u.tmpDir, u.rng.NextFileName(16))
		found, err = xf.PathExists(tmpFileName)
	}
	written, err = CopyFile(tmpFileName, path) // dest <== src
	if err == nil {
		written, hash, err = u.Put1(tmpFileName, key)
	}
	return
}

// - GetData1 --------------------------------------------------------
func (u *U16x16) GetData1(key string) (data []byte, err error) {

	var (
		path string
		src  *os.File
	)
	path = u.GetPathForKey(key)
	found, err := xf.PathExists(path)
	if err == nil && !found {
		err = FileNotFound
	}
	if err == nil {
		src, err = os.Open(path)
	}
	if err == nil {
		defer src.Close()
		var count int
		// XXX THIS WILL NOT WORK FOR LARGER FILES!  It will ignore
		//     anything over 64 KB
		data = make([]byte, DEFAULT_BUFFER_SIZE)
		count, err = src.Read(data)
		// XXX COUNT IS IGNORED
		_ = count
	}
	return
}

// - Put1 ------------------------------------------------------------
// tmp is the path to a local file which will be renamed into U (or deleted
// if it is already present in U)
// u.path is an absolute or relative path to a U directory organized 16x16
// key is an sha1 content hash.
// If the operation succeeds we return the length of the file (which must
// not be zero.  Otherwise we return 0.
// we don't do much checking
func (u *U16x16) Put1(inFile, key string) (
	length int64, hash string, err error) {

	var (
		found                          bool
		fullishPath                    string
		topSubDir, lowerDir, targetDir string
	)
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
	info, err := os.Stat(inFile)
	if err != nil {
		return
	}
	length = info.Size()
	topSubDir = hash[0:1]
	lowerDir = hash[1:2]
	targetDir = filepath.Join(u.path, topSubDir, lowerDir)
	found, err = xf.PathExists(targetDir)
	if err == nil && !found {
		// XXX MODE IS SUSPECT
		err = os.MkdirAll(targetDir, 0775)

	}
	if err == nil {
		fullishPath = filepath.Join(targetDir, key[2:])
		found, err = xf.PathExists(fullishPath)
	}
	if err == nil {
		if found {
			// drop the temporary input file
			err = os.Remove(inFile)
		} else {
			// rename the temporary file into U
			err = os.Rename(inFile, fullishPath)
			if err == nil {
				err = os.Chmod(fullishPath, 0444)
			}
		}
	}
	return
}

// PutData1 ---------------------------------------------------------
func (u *U16x16) PutData1(data []byte, key string) (
	length int64, hash string, err error) {

	var fullishPath string
	var found bool

	s := sha1.New()
	s.Write(data)
	hash = hex.EncodeToString(s.Sum(nil))
	if hash != key {
		fmt.Printf("expected data to have key %s, but content key is %s",
			key, hash)
		err = errors.New("content/key mismatch")
		return
	}
	length = int64(len(data))
	topSubDir := hash[0:1]
	lowerDir := hash[1:2]
	targetDir := filepath.Join(u.path, topSubDir, lowerDir)
	found, err = xf.PathExists(targetDir)
	if err == nil && !found {
		// MODE QUESTIONABLE
		err = os.MkdirAll(targetDir, 0775)
	}
	if err == nil {
		fullishPath = filepath.Join(targetDir, key[2:])
		found, err = xf.PathExists(fullishPath)
		if err == nil && !found {
			var dest *os.File
			dest, err = os.Create(fullishPath)
			if err == nil {
				var count int
				defer dest.Close()
				count, err = dest.Write(data)
				if err == nil {
					length = int64(count)
				}
			}
		}
	}
	return
}
