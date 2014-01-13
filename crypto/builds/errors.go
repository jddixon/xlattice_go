package builds

import (
	e "errors"
)

var (
	EmptyContentLine     = e.New("content line empty after trim")
	EmptyHash            = e.New("empty hash slice parameter")
	EmptyPath            = e.New("empty path parameter")
	IllFormedContentLine = e.New("content line not correctly formed")
)
