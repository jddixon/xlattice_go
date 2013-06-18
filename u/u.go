package u

// xlattice_go/u/u.go

import (
	"code.google.com/p/go.crypto/sha3"
	"encoding/hex"
	"fmt"				// DEBUG
	"errors"
	"github.com/jddixon/xlattice_go/rnglib"
	"io"
	"os"
	"path/filepath"
	"time"
)

                    // ....x....1....x....2....x....3....x....4
const SHA1_NONE     = "0000000000000000000000000000000000000000"
      // ....x....1....x....2....x....3....x....4....x....5....x....6....
const SHA3_NONE     = 
        "0000000000000000000000000000000000000000000000000000000000000000"

const DEFAULT_BUFFER_SIZE = 256*256

// XXX THIS DOES NOT BELONG HERE.  It is used to make unique
// file names, but this results in a race condition.
func MakeSimpleRNG() *rnglib.PRNG {
	t := time.Now().Unix()
	rng := rnglib.NewSimpleRNG(t)
	return rng
}
var RNG = MakeSimpleRNG()

func CopyFile(destName, srcName string) (written int64, err error) {
	var (src, dest *os.File)
	if src, err = os.Open(srcName);		err != nil { return }
	defer src.Close()
	if dest, err = os.Open(destName);	err != nil { return }
	defer dest.Close()
	return io.Copy(dest, src)
}

func PathExists(fName string) bool {
	if _, err := os.Stat(fName); os.IsNotExist(err) {
		return false
	}
	return true
}
//- copyAndPut3 -------------------------------------------------------
func CopyAndPut3(path, uPath, key string) (int64, string, error) {
    // the temporary file MUST be created on the same device
    tmpDir := filepath.Join(uPath, "tmp")
    // xxx POSSIBLE RACE CONDITION
    tmpFileName := filepath.Join(tmpDir, RNG.NextFileName(16))
    for ;  PathExists(tmpFileName) ; {
        tmpFileName = filepath.Join(tmpDir, RNG.NextFileName(16))
	}
    CopyFile(tmpFileName, path)
    return Put3(tmpFileName, uPath, key)
}
// - FileSHA3 --------------------------------------------------------
// returns the SHA1 hash of the contents of a file
func FileSHA3 (path string) (hash string, err error) {
	hash = SHA3_NONE
    if (path == "") || ! PathExists(path) {
		err = errors.New("IllegalArgument: empty or missing path")
		return
	}
    d := sha3.NewKeccak256()
	var f *os.File
	if f, err = os.Open(path); err != nil { 
		return
	}
	defer f.Close()
	buffer := make([]byte, DEFAULT_BUFFER_SIZE)
	var count int
    for ; true ; {
		// err IS IGNORED
        count, err = f.Read(buffer)
        if count == 0 {
            break
		}
        d.Write(buffer)
	}
	digest := d.Sum(nil)
	hash = hex.EncodeToString(digest)
	return
}
// - GetData3 --------------------------------------------------------
func GetData3(uPath, key string) (data []byte, err error) {
	var path string
    path, err = GetPathForKey(uPath, key)
    if ! PathExists(path) {
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
// uPath is an absolute or relative path to a U directory organized 256x256
// key is an sha3 content hash.
// If the operation succeeds we return the length of the file (which must
// not be zero.  Otherwise we return 0.
// we don't do much checking
func Put3(inFile, uPath, key string) (length int64, hash string, err error){
	length = 0
	// XXX IGNORING ERRORS
    hash, err = FileSHA3(inFile)
    if (hash != key) {
        fmt.Printf("expected %s to have key %s, but the content key is %s\n", 
                inFile, key, hash)
        err = errors.New("IllegalArgument: key does not match content")
		return
	}
	// XXX ERROR NOT HANDLED
	info,_		:= os.Stat(inFile)
    length		= info.Size()
    topSubDir	:= hash[0:2]
    lowerDir	:= hash[2:4]
    targetDir	:= filepath.Join(uPath, topSubDir, lowerDir)
    if ! PathExists(targetDir) {
		// XXX ERROR NOT HANDLED; MODE IS SUSPECT
         _ =os.MkdirAll(targetDir, 0775)

	}
    fullishPath := filepath.Join(targetDir, key)
    if (PathExists(fullishPath)) {
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
func PutData3(data []byte, uPath, key string) (
								length uint64, hash string, err error){
    s := sha3.NewKeccak256()
    s.Write(data)
	hash = hex.EncodeToString(s.Sum(nil))
    if (hash != key) {
        fmt.Printf("expected data to have key %s, but content key is %s", 
               key, hash)
        err = errors.New("content/key mismatch")       
		return
    }
	length    = uint64(len(data))
    topSubDir := hash[0:2]
    lowerDir  := hash[2:4]
    targetDir := filepath.Join(uPath, topSubDir, lowerDir)
    if ! PathExists(targetDir) {
		// XXX ERROR NOT HANDLED; MODE QUESTIONABLE
         _ =os.MkdirAll(targetDir, 0775)
	}
    fullishPath := filepath.Join(targetDir, key)
    if (!PathExists(fullishPath)) {
		var dest *os.File
		if dest, err = os.Create(fullishPath); err != nil { 
			return 
		}
		var count int
		defer dest.Close()
		// XXX ERROR NOT HANDLED
		count, err = dest.Write(data) 
		length = uint64(count)
	}
    return 
}
// - GetPathForKey ---------------------------------------------------
// returns a path to a file with the content key passed, or None if there
// is no such file
func GetPathForKey(uPath, key string) (path string, err error) {
	path = ""
    if !PathExists(uPath) {
        fmt.Printf ("HASH %s: UDIR DOES NOT EXIST: %s", key, uPath)
        err = errors.New("IllegalArgument: path does not exist")
		return
    }
    topSubDir := key[0:2]
    lowerDir  := key[2:4]
    path = filepath.Join(uPath, topSubDir, lowerDir, key[4:])
	return
}
