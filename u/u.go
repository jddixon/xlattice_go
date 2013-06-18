package u

// xlattice_go/u/u.go

import (
	"code.google.com/p/go.crypto/sha3"
	"encoding/hex"
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
func FileSHA3 (path string) string {
    if (path == "") || ! PathExists(path) {
        return SHA3_NONE
	}
    d := sha3.NewKeccak256()
    f := io.FileIO(path, "r")
    r := io.BufferedReader(f)
    for ; true ; {
        byteStr := r.read(io.DEFAULT_BUFFER_SIZE)
        if (len(byteStr) == 0) {
            break
		}
        d.Write(byteStr)
	}
	digest := d.Sum(nil)
	return hex.EncodeToString(digest)
}
// - GetData3 --------------------------------------------------------
func GetData3(uPath, key string) ([]byte, error) {
    path := GetPathForKey(uPath, key)
    if ! PathExists(path) {
        return nil, nil
    } else {
		if src, err := os.Open(srcName); err != nil { 
			return nil, err
		}
		defer src.Close()
		// XXX THIS WILL NOT WORK FOR LARGER FILES!  It will ignore
		//     anything over 64 KB
		data := make([]byte, 256*256)
		count, err := src.Read(data)
		if err == nil {
			return data, nil
		} else {
			return nil, err
		}
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
func Put3(inFile, uPath, key string) (int64, string, error){
    hash := FileSHA3(inFile)
    if (hash != key) {
        fmt.Printf("expected %s to have key %s, but the content key is %s\n", 
                inFile, key, hash)
        return 0, errors.NewError("IllegalArgument: key does not match content")
	}
    len = os.Stat(inFile).Size()
    topSubDir = hash[0:2]
    lowerDir  = hash[2:4]
    targetDir = filepath.Join(uPath, topSubDir, lowerDir)
    if ! PathExists(targetDir) {
		// XXX ERROR NOT HANDLED
         _ =os.MkdirAll(targetDir)

	}
    fullishPath = filepath.Join(targetDir, key)
    if (PathExists(fullishPath)) {
		// XXX ERROR NOT HANDLED
        _ = os.Remote(inFile)
    } else {
		// XXX ERROR NOT HANDLED
        _ = os.Rename(inFile, fullishPath)
		// XXX ERROR NOT HANDLED
		_ = os.Chmod(fullishPath, 0444)
	}
    return len, hash, err
}
// - putData3 --------------------------------------------------------
func PutData3(data, uPath, key string) (length int64, hash string, err error){
    s = sha3.NewKeccak256()
    s.Write(data)
	hash = hex.EncodeToString(s.Sum(nil))
    if (hash != key) {
        fmt.Printf("expected data to have key %s, but content key is %s", 
               key, hash)
        err = errors.NewError("content/key mismatch")       
		return
    }
	length    = len(data)          // XXX POINTLESS
    topSubDir = hash[0:2]
    lowerDir  = hash[2:4]
    targetDir = uPath + '/' + topSubDir + '/' + lowerDir + '/'
    if ! PathExists(targetDir) {
		// XXX ERROR NOT HANDLED
         _ =os.MkdirAll(targetDir)
	}
    fullishPath := filepath.Join(targetDir, key)
    if (PathExists(fullishPath)) {
        // print "DEBUG: file is already present"
        pass
    } else {
		if dest, err = os.Create(fullishPath); err != nil { 
			return 
		}
		defer dest.Close()
		// XXX ERROR NOT HANDLED
		length, err := dest.Write(data) 
	}
    return 
}
