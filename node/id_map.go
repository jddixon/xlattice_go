package node

// xlattice_go/node/id_map.go

import (
	"bytes"
	"fmt"
	"sync"
)

var _ = fmt.Print

type Thinger struct {
	Next  *MapForDepth
	Key   *[]byte
	Value interface{} // holds pointer to object of interest
}

type IDMap struct {
	MaxDepth uint // in bytes
	count    int
	mu       sync.RWMutex
	MapForDepth
}
type MapForDepth struct {
	Cells [256]Thinger
}

// Create an IDMap with the depth specified, where depth is the length
// of the longest key that can be inserted.
//
func NewIDMap(maxDepth uint) (m *IDMap, err error) {
	if maxDepth > MAX_MAX_DEPTH {
		err = MaxDepthTooLarge
	} else {
		m = &IDMap{
			MaxDepth: maxDepth,
		}
	}
	return
}

// Create an IDMap with the default depth.
//
func NewNewIDMap() (m *IDMap, err error) {
	return NewIDMap(MAX_MAX_DEPTH)
}

// Add an item to the map.  This should be idempotent: adding a key
// that is already in the map should have no effect at all.  The cell map
// allows us to efficiently return a pointer to the item, given its nodeID.

func (m *IDMap) Insert(key []byte, value interface{}) (err error) {
	if key == nil {
		err = NilID
	} else {
		var depth uint
		// XXX This is a very coarse lock; could lock at the mapForDepth
		// level instead
		m.mu.Lock()
		defer m.mu.Unlock()
		curMap := &m.MapForDepth
		for depth = 0; depth < m.MaxDepth; depth++ {
			nextByte := uint(key[depth])
			cell := &curMap.Cells[nextByte]
			if cell.Next == nil {
				if cell.Key == nil {
					// we are done
					cell.Key = &key
					cell.Value = value
					m.count++
				} else if bytes.Equal(*cell.Key, key) {
					// it's already there
				} else {
					err = m.handleCollision(depth, curMap, cell, key, value)
				}
				break
			} else {
				curMap = cell.Next
			}
		}
		if err == nil && depth >= m.MaxDepth {
			err = MaxDepthExceeded
		}
	}
	return
}

// The cell is occupied by a key with a different value.  We are guaranteed
// that the pointer to the next map is nil
//
func (m *IDMap) handleCollision(curDepth uint, curMap *MapForDepth,
	curCell *Thinger, key []byte, value interface{}) (err error) {

	if curDepth >= m.MaxDepth-1 {
		err = MaxDepthExceeded
	} else {
		keyB := *curCell.Key
		valueB := curCell.Value
		curCell.Key = nil
		curCell.Value = nil
		curCell.Next = new(MapForDepth)
		curMap = curCell.Next
		curDepth++
		aByte := uint(key[curDepth])
		bByte := uint(keyB[curDepth])

		bCell := &curMap.Cells[bByte]
		bCell.Key = &keyB
		bCell.Value = valueB
		if aByte != bByte {
			aCell := &curMap.Cells[aByte]
			aCell.Key = &key
			aCell.Value = value
		} else {
			err = m.handleCollision(curDepth, curMap, bCell, key, value)
		}
	}
	return
}

// Return the value associated with the key or nil if there is no
// such value.
//
func (m *IDMap) Find(key []byte) (value interface{}, err error) {

	if key == nil {
		err = NilID
	} else {
		m.mu.RLock()
		// XXX A very coarse lock; could lock at MapForDepth level instead
		defer m.mu.RUnlock()
		var depth uint
		curMap := &m.MapForDepth
		for depth = 0; depth < m.MaxDepth; depth++ {
			nextByte := uint(key[depth])
			cell := curMap.Cells[nextByte]
			if cell.Next == nil {
				// there is a match at this cell or no match at all
				if cell.Key != nil && bytes.Equal(*cell.Key, key) {
					value = cell.Value
				}
				break
			} else {
				curMap = cell.Next
			}
		}
		if err == nil && depth >= m.MaxDepth {
			err = MaxDepthExceeded
		}
	}
	return
}

func (m *IDMap) Size() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.count
}
