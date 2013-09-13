package reg

import (
	"errors"
)

var (
	ClusterMustHaveTwo = errors.New("cluster must have at least two members")
)
