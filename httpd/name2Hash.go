package httpd

// xlattice_go/httpd/name2Hash.go

import (
	xd "github.com/jddixon/xlattice_go/overlay/datakeyed"
)
type Name2Hash struct {		// must implement xo.NameKeyedI


	// map
	hashCache	xd.MemCacheI
	siteNames []string
	buildLists	[]string		// should be specific object
}
