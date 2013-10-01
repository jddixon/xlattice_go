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

func dumpRight(cell *BNIMapCell, indent string) {
	firstCell := cell
	ss := []string{indent}
	var offset int
	var lastCell *BNIMapCell
	for ; cell != nil; cell = cell.NextCol {
		if cell.CellNode == nil {
			ss = append(ss, NIL_PEER)
		} else {
			ss = append(ss, fmt.Sprintf("%-10s", cell.CellNode.GetName()))
		}
		lastCell = cell
		offset++
	}
	line := strings.Join(ss, "")
	fmt.Printf("%s\n", line)
	if firstCell != lastCell && lastCell.ThisCol != nil {
		for i := 0; i < offset-1; i++ {
			indent += SPACES
		}
		dumpDown(lastCell.ThisCol, indent)
	}
}
func dumpDown(cell *BNIMapCell, indent string) {
	for ; cell != nil; cell = cell.ThisCol {
		dumpRight(cell, indent)
	}
}
func DumpBNIMap(pm *BNIMap, where string) {
	fmt.Printf("PEER MAP DUMP %s --------\n", where)
	dumpDown(pm.NextCol, "")
	fmt.Println("END PEER MAP DUMP -------------------------------")
}
