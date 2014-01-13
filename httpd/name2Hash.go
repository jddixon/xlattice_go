package httpd

// xlattice_go/httpd/name2Hash.go

import (
	xd "github.com/jddixon/xlattice_go/overlay/datakeyed"
	"sync"
)

/**
 * Maintains data structures mapping path names to NodeIDs, which
 * are used to retrieve data from a MemCache, an in-memory cache of
 * byte slices.
 */
type Name2Hash struct { // must implement xo.NameKeyedI

	buildLists []string // should be specific object
	hashCache  *xd.MemCache
	hashMap    map[string][]byte
	siteNames  []string
	mx         sync.Mutex
}
