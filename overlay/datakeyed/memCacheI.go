package datakeyed

// xlattice_go/overlay/datakeyed/memCacheI.go

import (
	xi	"github.com/jddixon/xlattice_go/nodeID"
	xo	"github.com/jddixon/xlattice_go/overlay"
)


// import org.xlattice.CryptoException
// import org.xlattice.NodeID
// import org.xlattice.crypto.SHA1Digest
// import org.xlattice.crypto.SignedList
// import org.xlattice.overlay.CallBack
// import org.xlattice.overlay.DataKeyed
// import org.xlattice.overlay.GetCallBack
// import org.xlattice.overlay.DelCallBack
// import org.xlattice.overlay.PutCallBack
// import org.xlattice.util.NonBlockingLog
// import org.xlattice.util.StringLib;         // DEBUG


type MemCacheI interface {

    // Cannot be part of the interface because 'final static'
    // public final static MemCache getInstance()
    // public final static MemCache getInstance(String s)

    // LOGGING //////////////////////////////////////////////////////
    /** Subclasses should override.  */
    DEBUG_MSG(msg string)
    ERROR_MSG(msg string)

    // PROPERTIES ///////////////////////////////////////////////////
    Add (id *xi.NodeID, b []byte)
    ByteCount () int64
    Clear()
    ItemCount () int64
    GetPathToXLattice () string

	xo.DataKeyedI
}
