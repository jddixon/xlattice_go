package nodeID

import (
	e "errors"
)

var (
	MaxDepthExceeded = e.New("max IDMap depth exceeded")
	MaxDepthTooLarge = e.New("max IDMap depth too large")
	NilID            = e.New("nil ID argument")
)
