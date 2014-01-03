package overlay

// xlattice_go/overlay/callBackI.go

const (
	/** reports success */
	OK = iota
	/** already present; not usually an error status */
	EXISTS
	/** error found checking parameter list */
	BAD_ARGS
	/** IOException occurred */
	IO_EXCEPTION
	/** operation not successful because item is a directory */
	IS_DIRECTORY
	/** item not found */
	NOT_FOUND
	/** operation is not implemented */
	NOT_IMPLEMENTED
	/** operation failed because item is too large */
	TOO_BIG
	/** crypto verification fails */
	VERIFY_FAILS
)

var STATUS_CODES = []string{
	"OK", "EXISTS", "BAD_ARGS", "IO_EXCEPTION",
	"IS_DIRECTORY", "NOT_FOUND", "NOT_IMPLEMENTED", "TOO_BIG",
	"VERIFY_FAILS",
}

type CallBackI interface {
	GetStatus() int
}
