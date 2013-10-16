package lfs

// xlattice_go/util/fs/lfs.go

import (
	"os"
	"strings"
)

// If the directory named does not exist, create it, restricting
// visibility to the owner.  If the directory name is empty, call it
// "lfs", that is, ./lfs/

func CheckLFS(lfs string) (err error) {
	// XXX should verify that LFS is a directory
	if lfs == "" {
		lfs = "lfs"
	}
	return os.MkdirAll(lfs, 0700)
}

// Given a path to a file, create any missing intermediate directories.
func MkdirsToFile(pathToFile string, perm os.FileMode) (err error) {

	parts := strings.Split(pathToFile, "/")
	if len(parts) > 1 {
		var pathToDir string
		if len(parts) == 2 {
			// just drop the file name
			pathToDir = parts[0]
		} else {
			pathToDir = strings.Join( parts[:len(parts)-1], "/")
		}
		err = os.MkdirAll(pathToDir, perm)
	}
	return
}
