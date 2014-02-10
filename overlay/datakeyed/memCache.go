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

// Add a file to the collection.  This operation must be idempotent.
func (mc *MemCache) Add(id *xi.NodeID, b []byte) (err error) {

	key := id.Value()

	// XXX POSSIBLE DEADLOCK
	mc.mu.Lock()
	defer mc.mu.Unlock()

	value, err := mc.idMap.Find(key)
	if err == nil {
		if value == nil {
			mc.idMap.Insert(key, b)
			mc.itemCount++
			mc.byteCount += uint64(len(b))
		}
	}
	return
}

func (mc *MemCache) ByteCount() uint64 {
	mc.mu.RLock()
	defer mc.mu.RUnlock()
	return mc.byteCount
}
func (mc *MemCache) Clear() {

	// XXX STUB -- need a function of this name in IDMap !

	return
}
func (mc *MemCache) ItemCount() uint {
	mc.mu.RLock()
	defer mc.mu.RUnlock()
	return mc.itemCount
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
