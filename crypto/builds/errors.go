package builds

import (
	e "errors"
)

var (
	EmptyHash = e.New("empty hash slice parameter")
	EmptyPath = e.New("empty path parameter")
)
