package u

import (
	xr "github.com/jddixon/xlattice_go/rnglib"
)

type UI interface {
	Exists(key string) bool
	FileLen(key string) (length int64, err error)
	GetPathForKey(key string) string
	CopyAndPut3(path, key string) (int64, string, error)
	GetData3(key string) (data []byte, err error)
	Put3(inFile, key string) (length int64, hash string, err error)
	PutData3(data []byte, key string) (length int64, hash string, err error)
	CopyAndPut1(path, key string) (int64, string, error)
	GetData1(key string) (data []byte, err error)
	Put1(inFile, key string) (length int64, hash string, err error)
	PutData1(data []byte, key string) (length int64, hash string, err error)

	GetDirStruc() DirStruc
	GetPath() string
	// presumably temporary
	GetRNG() *xr.PRNG
}
