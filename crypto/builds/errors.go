package builds

import (
	e "errors"
)

var (
	CantAddToSignedList = e.New("can't add, list has been signed")
	EmptyHash           = e.New("empty hash slice parameter")
	EmptyPath           = e.New("empty path parameter")
)
