package u

// xlattice_go/u/u.go

import (
	"code.google.com/p/go.crypto/sha3"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	xf "github.com/jddixon/xlattice_go/util/lfs"
	"io"
	"io/ioutil"
	"os"
)

func New(path string, ds DirStruc, perm os.FileMode) (UI, error) {

	switch ds {
	case DIR_FLAT:
		return NewUFlat(path, perm)
	case DIR16x16:
		return NewU16x16(path, perm)
	case DIR256x256:
		return NewU256x256(path, perm)
	default:
		return nil, DirStrucNotRecognized
	}
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
