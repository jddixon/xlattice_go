package lfs

// xlattice_go/util/fs/lfs.go

import (
	"os"
)

// If the directory named does not exist, create it, restricting
// visibility to the owner.  If the directory name is empty, call it
// "lfs", that is, ./lfs/

func CheckLFS(lfs string) (err error) {
	if lfs == "" {
		lfs = "lfs"
	}
	return os.MkdirAll(lfs, 0700)
}
