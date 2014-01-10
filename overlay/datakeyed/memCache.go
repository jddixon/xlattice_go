package datakeyed

// xlattice_go/overlay/datakeyed/memCache.go

import (
	xi "github.com/jddixon/xlattice_go/nodeID"
	// xo	"github.com/jddixon/xlattice_go/overlay"
)

type MemCache struct {
	maxCount uint   // items in cache
	maxBytes uint64 // bytes in cache

}

// LOGGING //////////////////////////////////////////////////////
/** Subclasses should override.  */
func (mc *MemCache) DEBUG_MSG(msg string) {

	// XXX STUB

	return
}
func (mc *MemCache) ERROR_MSG(msg string) {

	// XXX STUB

	return
}

// PROPERTIES ///////////////////////////////////////////////////

func (mc *MemCache) Add(id *xi.NodeID, b []byte) {

	// XXX STUB

	return
}
func (mc *MemCache) ByteCount() (count uint64) {

	// XXX STUB

	return
}
func (mc *MemCache) Clear() {

	// XXX STUB

	return
}
func (mc *MemCache) ItemCount() (count uint64) {

	// XXX STUB

	return
}
func (mc *MemCache) GetPathToXLattice() (path string) {

	// XXX STUB

	return
}
