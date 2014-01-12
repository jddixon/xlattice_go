package datakeyed

// xlattice_go/overlay/datakeyed/cbQueue.go

import (
	xi "github.com/jddixon/xlattice_go/nodeID"
	xo "github.com/jddixon/xlattice_go/overlay"
	"sync"
)

/**
 * When a data item is not found in the in-memory cache (MemCache),
 * an instance of this class is queued on the key, a NodeID.  If
 * the data item is found on disk or otherwise fetched, each queue
 * member is called back with a success status code.  If the item
 * cannot be found, each is called with an appropriate failure code.
 */

// KeyQueue must implement GetCallBackI

type KeyQueue struct {
	// protected final static NonBlockingLog debugLog = NonBlockingLog.getInstance("debug.log")

	cbQ      []xo.CallBackI
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
func NewKeyQueue(mCache MemCacheI, id *xi.NodeID, cb xo.GetCallBackI) (
	kq *KeyQueue, err error) {

	if mCache == nil {
		err = NilMemCache
	} else if id == nil {
		err = NilNodeID
	} else if cb == nil {
		err = NilCallBack
	} else {

		// DEBUG_MSG(" constructor, cbQ on " + StringLib.byteArrayToHex(id.value()))
		kq = &KeyQueue{
			cbQ:      []xo.CallBackI{cb},
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

/**
 * This method must be externally synchronized.
 *
 * @param  cb  CallBack that is being queued up
 */
func (kq *KeyQueue) Add(cb xo.GetCallBackI) (err error) {
	if cb == nil {
		err = NilCallBack
	} else {
		// DEBUG_MSG(".add(callback), cbQ on " + StringLib.byteArrayToHex(myID.value()))
		kq.cbQ = append(kq.cbQ, (cb))
	}
	return
}

func (kq *KeyQueue) GetNodeID() *xi.NodeID {
	// XXX SHOULD COPY
	return kq.myID
}

func (kq *KeyQueue) Size() uint {
	kq.mx.RLock()
	defer kq.mx.RUnlock()
	return uint(len(kq.cbQ))
}

// INTERFACE GetCallBack ////////////////////////////////////////

/**
 * If whatever was requested was found, it is returned as the
 * value of the byte array and the status code is zero; otherwise
 * the byte array is nil and the status code is non-zero.
 *
 * @param status application-specific status code
 * @param data   requested value as byte array or nil if failure
 */
func (kq *KeyQueue) FinishedGet(status int, data []byte) (err error) {

	// DEBUG_MSG(".finishedGet, " + StringLib.byteArrayToHex(myID.value()) + ":\n    status = " + status + ", " + cbQ.size() + " callbacks pending")

	kq.status = status
	var myData []byte
	if status == xo.OK {
		myData = data
		if myData != nil { // let's be paranoid
			kq.memCache.Add(kq.myID, myData)
		}
	}

	kq.mx.Lock()
	defer kq.mx.Unlock()

	count := uint(len(kq.cbQ))
	for i := uint(0); i < count; i++ {
		cb := kq.cbQ[i].(xo.GetCallBackI)
		cb.FinishedGet(status, myData)
	}
	kq.cbQ = kq.cbQ[:0]
	return
}

/**
 * Useful only for debugging?
 */
func (kq *KeyQueue) GetStatus() int {
	return kq.status
}
