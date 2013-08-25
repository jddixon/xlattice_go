package node

// xlattice_go/node/dump_test.go

import (
	"fmt"
	xi "github.com/jddixon/xlattice_go/nodeID"
	"strings"
)

var _ = fmt.Print
var _ = xi.SHA1_LEN

const (
	SPACES   = "          "
	NIL_PEER = "<nil>     "
	SP_COUNT = len(SPACES)
)

func dumpRight(cell *PeerMapCell, indent string) {
	firstCell := cell
	ss := []string{indent}
	var offset int
	var lastCell *PeerMapCell
	for ; cell != nil; cell = cell.nextCol {
		if cell.peer == nil {
			ss = append(ss, NIL_PEER)
		} else {
			ss = append(ss, fmt.Sprintf("%-10s", cell.peer.GetName()))
		}
		lastCell = cell
		offset++
	}
	line := strings.Join(ss, "")
	fmt.Printf("%s\n", line)
	if firstCell != lastCell && lastCell.thisCol != nil {
		for i := 0; i < offset-1; i++ {
			indent += SPACES
		}
		dumpDown(lastCell.thisCol, indent)
	}
}
func dumpDown(cell *PeerMapCell, indent string) {
	for ; cell != nil; cell = cell.thisCol {
		dumpRight(cell, indent)
	}
}
func DumpPeerMap(pm *PeerMap, where string) {
	fmt.Printf("PEER MAP DUMP %s --------\n", where)
	dumpDown(pm.nextCol, "")
	fmt.Println("END PEER MAP DUMP -------------------------------")
}
