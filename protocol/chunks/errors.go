package chunks

import (
	e "errors"
)

var (
	NilData  = e.New("nil data parameter")
	NilDatum = e.New("nil datum parameter")
)
