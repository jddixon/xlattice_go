package nodeID

// xlattice_go/nodeID/idMapHAMT.go

import (
	"fmt"
	gh "github.com/jddixon/hamt_go"
	"sync"
)

var _ = fmt.Print

type IDMapHAMT struct {
	h  *gh.HAMT
	mu sync.RWMutex
}

func NewIDMapHAMT(w, t uint) *IDMapHAMT {
	return &IDMapHAMT{h: gh.NewHAMT(w, t)}
}

// Create an IDMapHAMT with the default parameters w = 4, t = 4
//
func NewNewIDMapHAMT() *IDMapHAMT {
	return NewIDMapHAMT(4, 4)
}

func (m *IDMapHAMT) Delete(key []byte) (err error) {

	bKey, err := gh.NewBytesKey(key)
	if err == nil {
		// XXX A very coarse lock
		m.mu.Lock()
		defer m.mu.Unlock()
		err = m.h.Delete(bKey)
	}
	return
}

// Add an item to the map.  This should be idempotent: adding a key
// that is already in the map should have no effect at all.
//
func (m *IDMapHAMT) Insert(key []byte, value interface{}) (err error) {
	bKey, err := gh.NewBytesKey(key)
	if err == nil {
		// XXX A very coarse lock
		m.mu.Lock()
		defer m.mu.Unlock()
		err = m.h.Insert(bKey, value)
	}
	return
}

// Return the value associated with the key or nil if there is no
// such value.
func (m *IDMapHAMT) Find(key []byte) (value interface{}, err error) {

	bKey, err := gh.NewBytesKey(key)
	if err == nil {
		// XXX A very coarse lock
		m.mu.RLock()
		defer m.mu.RUnlock()
		value, err = m.h.Find(bKey)
		if err == gh.NotFound {
			err = nil
			value = nil
		}
	}
	return
}

// Returns number of entries in the map, the number of mapForDepth
// structures, and the deepest we have gone.
func (m *IDMapHAMT) Size() (x, y, z uint) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	items := m.h.GetLeafCount()
	tables := m.h.GetTableCount()
	depth := uint(0) // XXX NOT YET
	return items, tables, depth
}
