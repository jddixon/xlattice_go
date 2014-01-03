package overlay

// xlattice_go/overlay/delCallBackI.go

/**
 * This is a callback interface.
 */

type DelCallBackI interface {

	/**
	 * If the delete operation succeeded, returns zero.
	 * Otherwise returns a non-zero application-specific status code.
	 *
	 * @param status application-specific status code
	 */
	FinishedDel(status int)

	CallBackI
}
