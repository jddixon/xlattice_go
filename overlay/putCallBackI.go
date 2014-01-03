package overlay

// xlattice_go/overlay/putCallBack.go

/**
 * This is a callback interface.
 */

type PutCallBackI interface {

	/**
	 * The put has completed.  If the status is zero, it was
	 * successful.  Otherwise it was unsuccessful.
	 *
	 * @param status application-specific status code.
	 */
	FinishedPut(status int)

	CallBackI
}
