package overlay

// xlattice_go/overlay/getCallBackI.go

/**
 * This is a callback interface.
 */
type GetCallBackI interface {

	/**
	 * If whatever was requested was found, it is returned as the
	 * value of the byte array and the status code is zero; otherwise
	 * the byte array is null and the status code is non-zero.
	 *
	 * @param status application-specific status code
	 * @param data   requested value as byte array
	 */
	FinishedGet(status int, data []byte)

	CallBackI
}
