package reg

// xlattice_go/reg/dump_test.go

import (
	"fmt"
	xn "github.com/jddixon/xlattice_go/node"
	xi "github.com/jddixon/xlattice_go/nodeID"
	"strings"
)

var _ = fmt.Print
var _ = xi.SHA1_LEN

const (
	SPACES   = "          "
	NIL_MAP = "<nil>     "
	SP_COUNT = len(SPACES)
)

func dumpRight(cell *xn.BNIMapCell, indent string) {
	firstCell := cell
	ss := []string{indent}
	var offset int
	var lastCell *xn.BNIMapCell
	for ; cell != nil; cell = cell.NextCol {
		if cell.CellNode == nil {
			ss = append(ss, NIL_MAP)
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
func dumpDown(cell *xn.BNIMapCell, indent string) {
	for ; cell != nil; cell = cell.ThisCol {
		dumpRight(cell, indent)
	}
}
func DumpBNIMap(pm *xn.BNIMap, where string) {
	fmt.Printf( "BEGIN MAP DUMP %s --------\n", where)
	dumpDown(pm.NextCol, "")
	fmt.Println("END MAP MAP -------------------------------")
}
