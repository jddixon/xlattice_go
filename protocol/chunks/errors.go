package chunks

import (
	e "errors"
)

var (
	EmptyTitle = e.New("empty title parameter")
	NilData    = e.New("nil data parameter")
	NilDatum   = e.New("nil datum parameter")
)
