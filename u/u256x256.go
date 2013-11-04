package u

// xlattice_go/u/u256x256.go

import (
	"code.google.com/p/go.crypto/sha3"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt" // DEBUG
	xr "github.com/jddixon/xlattice_go/rnglib"
	xf "github.com/jddixon/xlattice_go/util/lfs"
	//"io/ioutil"
	"os"
	"path/filepath"
	//"time"
)

// CLASS, so to speak ///////////////////////////////////////////////
type U256x256 struct {
	path   string // all parameters are
	inDir  string // ... private
	tmpDir string
	rng    *xr.PRNG
}

// Create a new 256x256 file system, ORing perm into the default permissions.
// If perm is 0, the default is to allow user and group access.
// If the root is U, then this creates U/, U/tmp, U/in, and the top-level
// hex directories U/xx
func NewU256x256(path string, perm os.FileMode) (udir *U256x256, err error) {
	// TODO: validate path
	var (
		inDir, tmpDir string
	)
	err = os.MkdirAll(path, 0750|perm)
	if err == nil {
		inDir = filepath.Join(path, "in")
		err = os.MkdirAll(inDir, 0770|perm)
		if err == nil {
			tmpDir = filepath.Join(path, "tmp")
			err = os.MkdirAll(tmpDir, 0700)
			if err == nil {
				for i := 0; i < 256; i++ {
					hexDir := fmt.Sprintf("%02x", i)
					hexPath := filepath.Join(path, hexDir)
					err = os.MkdirAll(hexPath, 0750|perm)
					if err != nil {
						break
					}
				}
			}
		}
	}
	udir = &U256x256{
		path:   path,
		rng:    xr.MakeSimpleRNG(),
		inDir:  inDir,
		tmpDir: tmpDir,
	}
	return
}

func (u *U256x256) GetDirStruc() DirStruc { return DIR256x256 }
func (u2 *U256x256) GetRNG() *xr.PRNG {
	return u2.rng
}
func (u2 *U256x256) GetPath() string {
	return u2.path
}

// - Exists ---------------------------------------------------------
func (u2 *U256x256) Exists(key string) (found bool, err error) {
	path, err := u2.GetPathForKey(key)
	if err == nil {
		found, err = xf.PathExists(path)
	}
	return
}

// - FileLen --------------------------------------------------------
func (u2 *U256x256) FileLen(key string) (length int64, err error) {
	path, err := u2.GetPathForKey(key)
	if err == nil {
		var info os.FileInfo
		info, err = os.Stat(path)
		if err == nil {
			length = info.Size()
		}
	}
	return
}

// - GetPathForKey --------------------------------------------------
// Returns a path to a file with the content key passed.
func (u2 *U256x256) GetPathForKey(key string) (path string, err error) {
	if key == "" {
		err = EmptyKey
	} else {
		topSubDir := key[0:2]
		lowerDir := key[2:4]
		path = filepath.Join(u2.path, topSubDir, lowerDir, key[4:])
	}
	return
}

//- copyAndPut3 -------------------------------------------------------
func (u2 *U256x256) CopyAndPut3(path, key string) (
	written int64, hash string, err error) {
	// the temporary file MUST be created on the same device
	tmpFileName := filepath.Join(u2.tmpDir, u2.rng.NextFileName(16))
	found, _ := xf.PathExists(tmpFileName) // XXX error ignored
	for found {
		tmpFileName = filepath.Join(u2.tmpDir, u2.rng.NextFileName(16))
		found, _ = xf.PathExists(tmpFileName)
	}
	written, err = CopyFile(tmpFileName, path) // dest <== src
	if err == nil {
		written, hash, err = u2.Put3(tmpFileName, key)
	}
	return
}

// - GetData3 --------------------------------------------------------
func (u2 *U256x256) GetData3(key string) (data []byte, err error) {
	var (
		found bool
		path  string
	)
	path, err = u2.GetPathForKey(key)
	if err == nil {
		found, err = xf.PathExists(path)
	}
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
// u2.path is an absolute or relative path to a U directory organized 256x256
// key is an sha3 content hash.
// If the operation succeeds we return the length of the file (which must
// not be zero.  Otherwise we return 0.
// we don't do much checking
func (u2 *U256x256) Put3(inFile, key string) (
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
	topSubDir := hash[0:2]
	lowerDir := hash[2:4]
	targetDir := filepath.Join(u2.path, topSubDir, lowerDir)
	found, err := xf.PathExists(targetDir)
	if err == nil && !found {
		// XXX MODE IS SUSPECT
		err = os.MkdirAll(targetDir, 0775)
	}
	if err == nil {
		var found bool

		fullishPath = filepath.Join(targetDir, key[4:])
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
func (u2 *U256x256) PutData3(data []byte, key string) (length int64, hash string, err error) {
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
	topSubDir := hash[0:2]
	lowerDir := hash[2:4]
	targetDir := filepath.Join(u2.path, topSubDir, lowerDir)
	found, err := xf.PathExists(targetDir)
	if err == nil && !found {
		err = os.MkdirAll(targetDir, 0775)
	}
	fullishPath := filepath.Join(targetDir, key[4:])
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
// XXX SHOULD RETURN ERROR
func (u2 *U256x256) CopyAndPut1(path, key string) (
	written int64, hash string, err error) {
	// the temporary file MUST be created on the same device
	// xxx POSSIBLE RACE CONDITION
	tmpFileName := filepath.Join(u2.tmpDir, u2.rng.NextFileName(16))
	found, err := xf.PathExists(tmpFileName)
	for found {
		tmpFileName = filepath.Join(u2.tmpDir, u2.rng.NextFileName(16))
		found, err = xf.PathExists(tmpFileName)
	}
	written, err = CopyFile(tmpFileName, path) // dest <== src
	if err == nil {
		written, hash, err = u2.Put1(tmpFileName, key)
	}
	return
}

// - GetData1 --------------------------------------------------------
func (u2 *U256x256) GetData1(key string) (data []byte, err error) {

	var (
		found bool
		path  string
		src   *os.File
	)
	path, err = u2.GetPathForKey(key)
	if err == nil {
		found, err = xf.PathExists(path)
	}
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
// u2.path is an absolute or relative path to a U directory organized 256x256
// key is an sha1 content hash.
// If the operation succeeds we return the length of the file (which must
// not be zero.  Otherwise we return 0.
// we don't do much checking
func (u2 *U256x256) Put1(inFile, key string) (
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
	topSubDir = hash[0:2]
	lowerDir = hash[2:4]
	targetDir = filepath.Join(u2.path, topSubDir, lowerDir)
	found, err = xf.PathExists(targetDir)
	if err == nil && !found {
		// XXX MODE IS SUSPECT
		err = os.MkdirAll(targetDir, 0775)

	}
	if err == nil {
		fullishPath = filepath.Join(targetDir, key[4:])
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
// XXX SHOULD RETURN ERROR
func (u2 *U256x256) PutData1(data []byte, key string) (
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
	topSubDir := hash[0:2]
	lowerDir := hash[2:4]
	targetDir := filepath.Join(u2.path, topSubDir, lowerDir)
	found, err = xf.PathExists(targetDir)
	if err == nil && !found {
		// MODE QUESTIONABLE
		err = os.MkdirAll(targetDir, 0775)
	}
	if err == nil {
		fullishPath = filepath.Join(targetDir, key[4:])
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
