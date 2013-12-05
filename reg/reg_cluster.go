package reg

// xlattice_go/reg/reg_cluster.go

// This file contains functions and structures used to describe
// and manage the clusters managed by the registry.

import (
	"bytes"
	"crypto/rsa"
	"encoding/hex"
	"fmt"
	xn "github.com/jddixon/xlattice_go/node"
	xi "github.com/jddixon/xlattice_go/nodeID"
	xo "github.com/jddixon/xlattice_go/overlay"
	"strconv"
	"strings"
	"sync"
)

var _ = fmt.Print

// cluster bit flags (Attrs)
const (
	CLUSTER_EPHEMERAL = 1 << iota
	CLUSTER_DELETED
)

type RegCluster struct {
	Name          string // must be unique
	ID            []byte // must be unique
	Attrs         uint64 // a field of bit flags
	maxSize       uint   // a maximum; must be > 0
	epCount       uint   // a positive integer, for now is 1 or 2
	Members       []*MemberInfo
	MembersByName map[string]*MemberInfo
	MembersByID   *xn.BNIMap
	mu            sync.RWMutex
}

func NewRegCluster(name string, id *xi.NodeID, attrs uint64,
	maxSize, epCount uint) (rc *RegCluster, err error) {

	if name == "" {
		name = "xlCluster"
	}
	nameMap := make(map[string]*MemberInfo)
	if epCount < 1 {
		err = ClusterMembersMustHaveEndPoint
	}
	if err == nil && maxSize < 1 {
		//err = ClusterMustHaveTwo
		err = ClusterMustHaveMember
	} else {
		var bnm xn.BNIMap // empty map
		rc = &RegCluster{
			Attrs:         attrs,
			Name:          name,
			ID:            id.Value(),
			epCount:       epCount,
			maxSize:       maxSize,
			MembersByName: nameMap,
			MembersByID:   &bnm,
		}
	}
	return
}

func (rc *RegCluster) AddToCluster(name string, id *xi.NodeID,
	commsPubKey, sigPubKey *rsa.PublicKey, attrs uint64, myEnds []string) (
	err error) {

	member, err := NewMemberInfo(
		name, id, commsPubKey, sigPubKey, attrs, myEnds)
	if err == nil {
		err = rc.AddMember(member)
	}
	return
}

func (rc *RegCluster) AddMember(member *MemberInfo) (err error) {

	// verify no existing member has the same name
	name := member.GetName()

	rc.mu.RLock() // <------------------------------------
	_, ok := rc.MembersByName[name]
	rc.mu.RUnlock() // <------------------------------------

	if ok {
		// XXX surely something more complicated is called for!

		fmt.Printf("AddMember: ATTEMPT TO ADD EXISTING MEMBER %s\n", name)
		return
	}
	// XXX CHECK FOR ENTRY IN BNIMap
	// XXX STUB

	rc.mu.Lock()             // <------------------------------------
	index := len(rc.Members) // DEBUG
	_ = index                // we might want to use this
	rc.Members = append(rc.Members, member)
	rc.MembersByName[name] = member
	err = rc.MembersByID.AddToBNIMap(member)
	rc.mu.Unlock() // <------------------------------------

	return
}

func (rc *RegCluster) EndPointCount() uint {
	return rc.epCount
}
func (rc *RegCluster) MaxSize() uint {
	return rc.maxSize
}
func (rc *RegCluster) Size() uint {
	var curSize uint
	rc.mu.RLock() // <------------------------------------
	curSize = uint(len(rc.Members))
	rc.mu.RUnlock() // <------------------------------------
	return curSize
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
	if !bytes.Equal(rc.ID, other.ID) {
		// DEBUG
		rcHexID := hex.EncodeToString(rc.ID)
		otherHexID := hex.EncodeToString(other.ID)
		fmt.Printf("rc.Equal: IDs DIFFER %s vs %s\n", rcHexID, otherHexID)
		// END
		return false
	}
	if rc.epCount != other.epCount {
		// DEBUG
		fmt.Printf("rc.Equal: EPCOUNTS DIFFER %d vs %d\n",
			rc.epCount, other.epCount)
		// END
		return false
	}
	if rc.maxSize != other.maxSize {
		// DEBUG
		fmt.Printf("rc.Equal: MAXSIZES DIFFER %d vs %d\n",
			rc.maxSize, other.maxSize)
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
	// Members			[]*MemberInfo
	for i := uint(0); i < rc.Size(); i++ {
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
	ss = append(ss, fmt.Sprintf("    epCount: %d", rc.epCount))
	ss = append(ss, fmt.Sprintf("    maxSize: %d", rc.maxSize))

	ss = append(ss, "    Members {")
	for i := 0; i < len(rc.Members); i++ {
		mem := rc.Members[i].Strings()
		for i := 0; i < len(mem); i++ {
			ss = append(ss, "        "+mem[i])
		}
	}
	ss = append(ss, "    }")
	ss = append(ss, "}")

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
		attrs            uint64
		name             string
		id               *xi.NodeID
		epCount, maxSize uint
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
		if strings.HasPrefix(line, "epCount: ") {
			var count int
			count, err = strconv.Atoi(line[9:])
			if err == nil {
				epCount = uint(count)
			}
		} else {
			fmt.Println("BAD END POINT COUNT")
			err = IllFormedCluster
		}
	}
	if err == nil {
		line = xn.NextNBLine(&rest)
		if strings.HasPrefix(line, "maxSize: ") {
			var size int
			size, err = strconv.Atoi(line[9:])
			if err == nil {
				maxSize = uint(size)
			}
		} else {
			fmt.Println("BAD MAX_SIZE")
			err = IllFormedCluster
		}
	}
	if err == nil {
		rc, err = NewRegCluster(name, id, attrs, maxSize, epCount)
	}
	if err == nil {
		line = xn.NextNBLine(&rest)
		if line == "Members {" {
			for {
				line = strings.TrimSpace(rest[0]) // peek
				if line == "}" {
					break
				}
				var member *MemberInfo
				member, rest, err = ParseMemberInfoFromStrings(rest)
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

// BaseNodeI INTERFACE //////////////////////////////////////////////

func (rc *RegCluster) GetName() string {
	return rc.Name
}
func (rc *RegCluster) GetNodeID() (id *xi.NodeID) {
	id, _ = xi.New(rc.ID)
	return
}

// Dummy functions to make this compliant with the interface

func (rc *RegCluster) AddOverlay(o xo.OverlayI) (ndx int, err error) {
	return
}
func (rc *RegCluster) SizeOverlays() (size int) {
	return
}
func (rc *RegCluster) GetOverlay(n int) (o xo.OverlayI) {
	return
}

func (rc *RegCluster) GetCommsPublicKey() (ck *rsa.PublicKey) {
	return
}
func (rc *RegCluster) GetSSHCommsPublicKey() (s string) {
	return
}
func (rc *RegCluster) GetSigPublicKey() (sk *rsa.PublicKey) {
	return
}
