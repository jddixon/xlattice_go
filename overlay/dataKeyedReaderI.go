package overlay

// xlattice_go/overlay/dataKeyedReader.go

import (
	xi "github.com/jddixon/xlattice_go/nodeID"
)

type DataKeyedReaderI interface {

	/**
	 * Retrieve data by content key (content hash).
	 */
	Get(nodeID *xi.NodeID, listener GetCallBackI)

	/**
	 * Retrieve a serialized SignedList, given its key, calculated
	 * from the RSA public key and title of the list.
	 */
	GetSigned(nodeID *xi.NodeID, listener GetCallBackI)
}
