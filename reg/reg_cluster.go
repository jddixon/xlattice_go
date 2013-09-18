package reg

// xlattice_go/reg/reg_cluster.go

// This file contains functions and structures used to describe
// and manage the clusters managed by the registry.

import (
	"crypto/rsa"
	"encoding/hex"
	"fmt"
	xn "github.com/jddixon/xlattice_go/node"
	xi "github.com/jddixon/xlattice_go/nodeID"
	xu "github.com/jddixon/xlattice_go/util"
	"strconv"
	"strings"
)

var _ = fmt.Print

// cluster bit flags
const (
	CLUSTER_DELETED = 1 << iota
)

type RegCluster struct {
	Attrs         uint64 // a field of bit flags
	Name          string // must be unique
	ID            []byte // must be unique
	MaxSize       int    // a maximum > 1
	Members       []*ClusterMember
	MembersByName map[string]*ClusterMember
	MembersByID   *xn.BaseNodeMap
}

func NewRegCluster(attrs uint64, name string, id *xi.NodeID, maxSize int) (
	rc *RegCluster, err error) {

	if name == "" {
		name = "xlCluster"
	}
	nameMap := make(map[string]*ClusterMember)
	if maxSize < 2 {
		err = ClusterMustHaveTwo
	} else {
		var bnm xn.BaseNodeMap // empty map
		rc = &RegCluster{
			Attrs:         attrs,
			Name:          name,
			ID:            id.Value(),
			MaxSize:       maxSize,
			MembersByName: nameMap,
			MembersByID:   &bnm,
		}
	}
	return
}

func (rc *RegCluster) Size() int {
	return len(rc.MembersByName)
}

func (rc *RegCluster) AddToCluster(name string, id *xi.NodeID,
	commsPubKey, sigPubKey *rsa.PublicKey, attrs uint64) (err error) {

	if _, ok := rc.MembersByName[name]; ok {
		// XXX surely something more complicated is called for!
		return
	}
	member, err := NewClusterMember(name, id, commsPubKey, sigPubKey, attrs)
	if err == nil {
		rc.MembersByName[name] = member

		// XXX add to MembersByID

	}
	return
}

func (rc *RegCluster) AddMember(member *ClusterMember) (err error) {

	name := member.GetName()
	if _, ok := rc.MembersByName[name]; ok {
		// XXX surely something more complicated is called for!
		return
	}
	// no existing member has the same name
	rc.MembersByName[name] = member
	rc.Members = append(rc.Members, member)

	// XXX ADD TO MembersByID

	return
}

// EQUAL ////////////////////////////////////////////////////////////
func (rc *RegCluster) Equal(any interface{}) bool {

	if any == rc {
		return true
	}
	if any == nil {
		return false
	}
	switch v := any.(type) {
	case *RegCluster:
		_ = v
	default:
		return false
	}
	other := any.(*RegCluster) // type assertion
	if rc.Attrs != other.Attrs {
		// DEBUG
		fmt.Printf("rc.Equal: ATTRS DIFFER %s vs %s\n", rc.Attrs, other.Attrs)
		// END
		return false
	}
	if rc.Name != other.Name {
		// DEBUG
		fmt.Printf("rc.Equal: NAMES DIFFER %s vs %s\n", rc.Name, other.Name)
		// END
		return false
	}
	if !xu.SameBytes(rc.ID, other.ID) {
		// DEBUG
		rcHexID := hex.EncodeToString(rc.ID)
		otherHexID := hex.EncodeToString(other.ID)
		fmt.Printf("rc.Equal: IDs DIFFER %s vs %s\n", rcHexID, otherHexID)
		// END
		return false
	}
	if rc.MaxSize != other.MaxSize {
		// DEBUG
		fmt.Printf("rc.Equal: MAXSIZES DIFFER %d vs %d\n",
			rc.MaxSize, other.MaxSize)
		// END
		return false
	}
	if rc.Size() != other.Size() {
		// DEBUG
		fmt.Printf("rc.Equal:ACTUAL SIZES DIFFER %d vs %d\n",
			rc.Size(), other.Size())
		// END
		return false
	}
	// Members			[]*ClusterMember
	for i := 0; i < rc.Size(); i++ {
		rcMember := rc.Members[i]
		otherMember := other.Members[i]
		if !rcMember.Equal(otherMember) {
			return false
		}
	}
	return true
}

// SERIALIZATION ////////////////////////////////////////////////////

func (rc *RegCluster) Strings() (ss []string) {

	ss = []string{"regCluster {"}

	ss = append(ss, fmt.Sprintf("    Attrs: 0x%016x", rc.Attrs))
	ss = append(ss, "    Name: "+rc.Name)
	ss = append(ss, "    ID: "+hex.EncodeToString(rc.ID))
	ss = append(ss, fmt.Sprintf("    MaxSize: %d", rc.MaxSize))

	ss = append(ss, "    Members {")
	for i := 0; i < len(rc.Members); i++ {
		mem := rc.Members[i].Strings()
		for i := 0; i < len(mem); i++ {
			ss = append(ss, "        "+mem[i])
		}
	}
	ss = append(ss, "    }")
	ss = append(ss, "}")

	// DEBUG
	fmt.Println("SERIALIZED CLUSTER: ==================================")
	for i := 0; i < len(ss); i++ {
		fmt.Printf("%s\n", ss[i])
	}
	fmt.Println("END SERIALIZED CLUSTER: ==============================")
	// END
	return
}

func (rc *RegCluster) String() string {
	return strings.Join(rc.Strings(), "\n")
}
func ParseRegCluster(s string) (rc *RegCluster, rest []string, err error) {
	ss := strings.Split(s, "\n")
	return ParseRegClusterFromStrings(ss)
}
func ParseRegClusterFromStrings(ss []string) (
	rc *RegCluster, rest []string, err error) {

	var (
		attrs   uint64
		name    string
		id      *xi.NodeID
		maxSize int
	)
	rest = ss

	line := xn.NextNBLine(&rest) // the line is trimmed
	if line != "regCluster {" {
		fmt.Println("MISSING regCluster {")
		err = IllFormedCluster
	} else {
		line = xn.NextNBLine(&rest)
		if strings.HasPrefix(line, "Attrs: ") {
			var i int64
			i, err = strconv.ParseInt(line[7:], 0, 64)
			if err == nil {
				attrs = uint64(i)
			}
		} else {
			fmt.Printf("BAD ATTRS in line '%s'", line)
			err = IllFormedCluster
		}
	}
	if err == nil {
		line = xn.NextNBLine(&rest)
		if strings.HasPrefix(line, "Name: ") {
			name = line[6:]
		} else {
			fmt.Printf("BAD NAME in line '%s'", line)
			err = IllFormedCluster
		}
	}
	if err == nil {
		// collect ID
		line = xn.NextNBLine(&rest)
		if strings.HasPrefix(line, "ID: ") {
			var val []byte
			val, err = hex.DecodeString(line[4:])
			if err == nil {
				id, err = xi.New(val)
			}
		} else {
			fmt.Println("BAD ID")
			err = IllFormedCluster
		}
	}
	if err == nil {
		line = xn.NextNBLine(&rest)
		if strings.HasPrefix(line, "MaxSize: ") {
			maxSize, err = strconv.Atoi(line[9:])
		} else {
			fmt.Println("BAD MAX_SIZE")
			err = IllFormedCluster
		}
	}
	if err == nil {
		rc, err = NewRegCluster(attrs, name, id, maxSize)
	}
	if err == nil {
		line = xn.NextNBLine(&rest)
		if line == "Members {" {
			for {
				line = strings.TrimSpace(rest[0]) // peek
				if line == "}" {
					break
				}
				var member *ClusterMember
				member, rest, err = ParseClusterMemberFromStrings(rest)
				if err != nil {
					break
				}
				err = rc.AddMember(member)
				if err != nil {
					break
				}
			}
		} else {
			err = MissingMembersList
		}
	}

	// expect closing brace for Members list
	if err == nil {
		line = xn.NextNBLine(&rest)
		if line != "}" {
			err = MissingClosingBrace
		}
	}
	// expect closing brace  for cluster
	if err == nil {
		line = xn.NextNBLine(&rest)
		if line != "}" {
			err = MissingClosingBrace
		}
	}

	return
}
