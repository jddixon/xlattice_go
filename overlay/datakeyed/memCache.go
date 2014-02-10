package datakeyed

// xlattice_go/overlay/datakeyed/memCache.go

import (
	xi "github.com/jddixon/xlattice_go/nodeID"
	// xo "github.com/jddixon/xlattice_go/overlay"
	"sync"
	// "time"
)

type MemCache struct {
	byteCount uint64
	itemCount uint
	maxBytes  uint64 // bytes in cache; should treat as const
	maxItems  uint   // items in cache; should treat as const

	idMap *xi.IDMap
	mu    sync.RWMutex
}

func NewMemCache(maxBytes uint64, maxItems uint) (mc *MemCache, err error) {
	idMap, err := xi.NewNewIDMap()
	if err == nil {
		mc = &MemCache{
			maxBytes: maxBytes,
			maxItems: maxItems,
			idMap:    idMap,
		}
	}
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
