package reg

// xlattice_go/reg/admin_client.go

import (
	"fmt"

	xi "github.com/jddixon/xlattice_go/nodeID"
)

var _ = fmt.Print

type AdminClient struct {
	clusterName string
	clusterID   *xi.NodeID
	clusterSize uint32 // this is a FIXED size, aka MaxSize

	members []ClusterMember

	ClientNode
}
