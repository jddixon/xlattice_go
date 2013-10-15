package lfs

import (
	"os"
)

func PathExists(path string) (whether bool, err error) {
	_, err = os.Stat(path)
	if err == nil {
		whether = true
	} else if os.IsNotExist(err) {
		err = nil
	}
	return
}
