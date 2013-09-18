package reg

// xlattice_go/reg/regData.go

// This file contains functions and structures used to describe
// and manage the cluster data managed by the registry.

import (
	"fmt"
)

var _ = fmt.Print

type RegData struct {
	clusters []*RegCluster
	members  []*ClusterMember
}
