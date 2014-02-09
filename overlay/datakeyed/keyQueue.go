package datakeyed

// xlattice_go/overlay/datakeyed/cbQueue.go

import (
	xi "github.com/jddixon/xlattice_go/nodeID"
	"sync"
)

/**
 * When a data item is not found in the in-memory cache (MemCache),
 * an instance of this class is queued on the key, a NodeID.  If
 * the data item is found on disk or otherwise fetched, each queue
 * member is called back with a success status code.  If the item
 * cannot be found, each is called with an appropriate failure code.
 */

type KeyQueue struct {
	// protected final static NonBlockingLog debugLog = NonBlockingLog.getInstance("debug.log")

	mx       sync.RWMutex
	memCache MemCacheI
	myID     *xi.NodeID

	status int
}

// CONSTRUCTORS /////////////////////////////////////////////////

/**
 * The queue is always constructed when we have an item but no
 * queue to put it in.  So we make the queue and add that item
 * to it.
 *
 * @param id  NodeID that these CallBacks are waiting for
 * @param cb  CallBack that is being queued up
 * @throws IllegalArgumentException if an argument is nil
 */
func NewKeyQueue(mCache MemCacheI, id *xi.NodeID) (
	kq *KeyQueue, err error) {

	if mCache == nil {
		err = NilMemCache
	} else if id == nil {
		err = NilNodeID
	} else {

		// DEBUG_MSG(" constructor, cbQ on " + StringLib.byteArrayToHex(id.value()))
		kq = &KeyQueue{
			memCache: mCache,
			myID:     id,
			status:   -1,
		}
	}
	return
}

// LOGGING //////////////////////////////////////////////////////

//protected void DEBUG_MSG(String msg) {
//    if (debugLog != nil)
//    debugLog.message("KeyQueue" + msg)
//}

// PROPERTIES ///////////////////////////////////////////////////


func (kq *KeyQueue) GetNodeID() *xi.NodeID {
	// XXX SHOULD COPY
	return kq.myID
}

func (kq *KeyQueue) Size() (n uint) {
	kq.mx.RLock()
	defer kq.mx.RUnlock()
	
	// return uint(len(kq.cbQ))

	// XXX STUB
	return
}


/**
 * Useful only for debugging?
 */
func (kq *KeyQueue) GetStatus() int {
	return kq.status
}
