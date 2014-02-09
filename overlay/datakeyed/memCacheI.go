package datakeyed

// xlattice_go/overlay/datakeyed/memCacheI.go

import (
	xi "github.com/jddixon/xlattice_go/nodeID"
	xo "github.com/jddixon/xlattice_go/overlay"
)

type MemCacheI interface {

	// Cannot be part of the interface because 'final static'
	// public final static MemCache getInstance()
	// public final static MemCache getInstance(String s)

	// LOGGING //////////////////////////////////////////////////////
	/** Subclasses should override.  */
	DEBUG_MSG(msg string)
	ERROR_MSG(msg string)

	// PROPERTIES ///////////////////////////////////////////////////
	Add(id *xi.NodeID, b []byte) error
	ByteCount() int64
	Clear()
	ItemCount() int64
	GetPathToXLattice() string

	xo.DataKeyedI
}
