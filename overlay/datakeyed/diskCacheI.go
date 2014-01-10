package datakeyed

// xlattice_go/overlay/datakeyed/diskCacheI.go

import (
	xi "github.com/jddixon/xlattice_go/nodeID"
	xo "github.com/jddixon/xlattice_go/overlay"
)

type DiskCacheI interface {

	/** initializes (clears) the underlying Bloom filter */
	// Init ()

	/** badly named */
	GetPathToXLattice() string

	/** For testing, should be deprecated ASAP. */
	// GetThreadPool() *DiskIOThreadPool

	// BLOOM FILTER INTERFACE ///////////////////////////////////////
	/**
	 * The theoretical false positive rate is
	 *   (1 - e(-kN/M))^k
	 * where k is the number of key functions, N is the number
	 * of keys in the filter, and there are 2^M bits in the
	 * filter.
	 *
	 * @param  n number of keys in the filter
	 * @return an approximation to the false positive rate
	 */
	FalsePositives(n int) float64

	/**
	 * Add a key to the filter.
	 */
	Insert(id *xi.NodeID) error

	/**
	 * This method has a non-zero false positive rate.  The filter
	 * can normally be set up to make this a very small number.
	 *
	 * @return whether a key is in the filter
	 * @see #FalsePositives
	 */
	IsMember(id *xi.NodeID) bool

	/**
	 * Remove a key from the filter.
	 */
	Remove(id *xi.NodeID) error

	/**
	 * A counter is incremented and decremented whenever a key is
	 * inserted or removed.  The value of this counter is returned
	 * as a crude estimate of the number of keys in the filter.
	 *
	 * @return the approximate number of keys in the filter
	 */
	Size() uint

	// DISK I/O QUEUES //////////////////////////////////////////////
	AcceptReadJob(id *xi.NodeID, cb xo.GetCallBackI) error

	AcceptWriteJob(id *xi.NodeID, data []byte, cb xo.PutCallBackI) error
}
