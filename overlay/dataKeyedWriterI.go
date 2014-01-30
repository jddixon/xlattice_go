package overlay

// xlattice_go/overlay/dataKeyedWriter.go

import (
	xc "github.com/jddixon/xlattice_go/crypto"
	xi "github.com/jddixon/xlattice_go/nodeID"
)

type DataKeyedWriterI interface {
	/**
	 * Delete the data item whose hash or title key is the value of
	 * the nodeID.  This call is synchronous: it blocks.
	 *
	 * XXX This operation may be ambiguous.  It is intended for use
	 * XXX in, for example, LRU caches.
	 *
	 * @param nodeID whose value of the SHA1 or SHA3/256 content hash
	 * @return whether the operation is successful
	 */
	Delete(nodeID *xi.NodeID) error

	/**
	 * Store a data item whose content hash is the value of nodeID.
	 * If an item with this key is already present, return false and
	 * do not write to disk.  If there is no item with this key and
	 * storage is successful, return true.  This is a synchronous call.
	 *
	 * @param nodeID, the SHA1 or SHA3/256 content hash of the buffer
	 * @param b      data being stored
	 */
	Put(nodeID *xi.NodeID, b []byte) error

	/**
	 * Store a SignedList.  The value of the NodeID is the key
	 * calculated from the RSA public key and title of the list.
	 * If a data item with this title key is already present, retrieve
	 * it.  Replace it if the stored item verifies but has a more
	 * recent timestamp.
	 *
	 * @param nodeID whose value is the key of the SignedList
	 * @param lst    the SignedList being stored
	 * @return       whether the list was written to store
	 */
	PutSigned(nodeID *xi.NodeID, lst xc.SignedListI) error
}
