package merkletree

// xlattice_go/util/merkletree/merkletree.go

import (
	"code.google.com/p/go.crypto/sha3"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	//xu "github.com/jddixon/xlattice_go/util"
	"hash"
	"io/ioutil"
	"os"
	"path"
	re "regexp"
	"strings"
)

var _ = fmt.Print

type MerkleTree struct {
	bound   bool
	exRE    []*re.Regexp // exclusions
	matchRE []*re.Regexp // must be matched
	nodes   []MerkleNodeI

	path       string
	MerkleNode // so name, hash, usingSHA1
}

func NewNewMerkleTree(name string, usingSHA1 bool) (*MerkleTree, error) {
	return NewMerkleTree(name, usingSHA1, nil, nil)
}

// Create an unbound MerkleTree with a nil hash and an empty nodes list.
// exRE and matchRE must have been validated by the calling code

func NewMerkleTree(name string, usingSHA1 bool, exRE, matchRE []*re.Regexp) (
	mt *MerkleTree, err error) {

	// this validates its parameters
	mn, err := NewMerkleNode(name, nil, usingSHA1)
	if err == nil {
		mt = &MerkleTree{
			exRE:       exRE,
			matchRE:    matchRE,
			MerkleNode: *mn,
		}
	}
	return
}

func (mt *MerkleTree) IsLeaf() bool {
	return false
}

func (mt *MerkleTree) addNode(node MerkleNodeI) (err error) {
	if node == nil {
		err = NilNode
	} else {
		mt.nodes = append(mt.nodes, node)
	}
	return
}
func ParseFirstLine(line string) (
	indent int, treeHash []byte, dirName string, err error) {

	line = strings.TrimRight(line, " \t")

	groups := FIRST_LINE_RE_1.FindStringSubmatch(line)
	if groups == nil {
		groups = FIRST_LINE_RE_3.FindStringSubmatch(line)
	}
	if groups == nil {
		err = CantParseFirstLine
	}
	if err == nil {
		treeHash, err = hex.DecodeString(groups[2])
	}
	if err == nil {
		indent = len(groups[1]) / 2
		dirName = groups[3]
		dirName = dirName[0 : len(dirName)-1] // drop terminating slash
	}
	return
}

func ParseOtherLine(line string) (
	nodeDepth int, nodeHash []byte, nodeName string, isDir bool, err error) {

	line = strings.TrimRight(line, " \t")

	groups := OTHER_LINE_RE_1.FindStringSubmatch(line)
	if groups == nil {
		groups = OTHER_LINE_RE_3.FindStringSubmatch(line)
	}
	if groups == nil {
		err = CantParseOtherLine
	}
	if err == nil {
		nodeHash, err = hex.DecodeString(groups[2])
	}
	if err == nil {
		nodeDepth = len(groups[1]) / 2
		nodeName = groups[3]
		if strings.HasSuffix(nodeName, "/") {
			isDir = true
			nodeName = nodeName[0 : len(nodeName)-1]
		}
	}
	return
}

// The string array is expected to follow conventional indentation
// rules, with zero indentation on the first line and some multiple
// of two spaces on all successive lines.

func ParseMerkleTreeFromStrings(ss *[]string) (mt *MerkleTree, err error) {

	var (
		indent    int
		treeHash  []byte
		dirName   string
		usingSHA1 bool
		stack     []MerkleNodeI
		stkDepth  int
		curTree   *MerkleTree
		// lastWasDir	bool	// not being used
	)
	if len(*ss) == 0 {
		err = EmptySerialization
	}
	if err == nil {
		firstLine := (*ss)[0]
		firstLine = strings.TrimRight(firstLine, " \t")
		indent, treeHash, dirName, err = ParseFirstLine(firstLine)
		if err == nil && indent > 0 {
			err = InitialIndent
		}
	}
	if err == nil {
		usingSHA1 = len(treeHash) == SHA1_LEN
		mt, err = NewNewMerkleTree(dirName, usingSHA1)
		if err != nil {
			return
		}
		mt.SetHash(treeHash)
		_ = indent // NEVER USED?	XXX
		curTree = mt
		stack = append(stack, curTree) // rootTree = mt
		stkDepth++                     // always step after pushing tree
	}
	for i := 1; i < len(*ss); i++ {
		var (
			lineIndent int
			thisHash   []byte
			name       string
			isDir      bool
		)
		line := (*ss)[i]
		line = strings.TrimRight(line, " \t")
		if len(line) == 0 {
			continue
		}
		// Note the hash may not be of the expected type.
		lineIndent, thisHash, name, isDir, err := ParseOtherLine(line)
		if err != nil {
			break
		}
		if lineIndent < stkDepth {
			for lineIndent < stkDepth {
				stkDepth--
				stack = stack[:len(stack)-1]
			}
			curTree = stack[len(stack)-1].(*MerkleTree) // MAY NOT BE!!
		}
		if stkDepth != lineIndent {
			fmt.Printf("INTERNAL ERROR: stkDepth %d, lineIndent, %d\n",
				stkDepth, lineIndent)
		}
		if isDir {
			// create and set attributes of a new node
			var newTree *MerkleTree
			newTree, err = NewNewMerkleTree(name, usingSHA1)
			if err != nil {
				break
			}
			newTree.SetHash(thisHash)

			//  add the new node into the existing tree
			curTree.addNode(newTree)
			stack = append(stack, newTree)
			stkDepth++
			curTree = newTree
		} else {
			var newNode *MerkleLeaf
			// create and set attributes of new node
			newNode, err = NewMerkleLeaf(name, thisHash, usingSHA1)
			if err != nil {
				break
			}
			// add the new node into the existing tree
			curTree.addNode(newNode)
		}
	}
	return
}

func ParseMerkleTree(s string) (mt *MerkleTree, err error) {

	if s == "" {
		err = EmptySerialization
	} else {
		ss := strings.Split(s, "\r\n")
		mt, err = ParseMerkleTreeFromStrings(&ss)
	}
	return
}

//    @staticmethod
//    def createFromFile(pathToFile):
//        if not os.path.exists(pathToFile):
//            raise RuntimeError(
//                "MerkleTree.createFromFile: file "%s" does not exist" % pathToFile)
//        with open(pathToFile, "r") as f:
//            line = f.readline()
//            line = line.rstrip()
//            m = re.match(MerkleTree.FIRST_LINE_PAT_1, line)
//            if m == None:
//                m = re.match(MerkleTree.FIRST_LINE_PAT_3, line)
//                usingSHA1 = false
//            else:
//                usingSHA1 = true
//            if m == None:
//                raise RuntimeError(
//                        "line "%s" does not match expected pattern" %  line)
//            dirName = m.group(3)
//            tree = MerkleTree(dirName, usingSHA1)
//#           if m.group(3) != "bind":
//#               raise RuntimeError(
//#                       "expected "bind" in first line, found %s" % m.group(3))
//            tree.setHash(m.group(2))
//            line = f.readline()
//            while line:
//                line = line.rstrip()
//                if line == "":
//                    continue
//                if mt.usingSHA1:
//                    m = re.match(MerkleTree.OTHER_LINE_PAT_1, line)
//                else:
//                    m = re.match(MerkleTree.OTHER_LINE_PAT_3, line)
//
//                if m == None:
//                    raise RuntimeError(
//                            "line "%s" does not match expected pattern" %  line)
//                tree._add(m.group(3), m.group(2))
//                line = f.readline()
//
//        return tree

//    @staticmethod
//    def CreateMerkleTreeFromFileSystem(pathToDir, usingSHA1 = false,
//                                        exRE = None, matchRE = None):

func CreateMerkleTreeFromFileSystem(pathToDir string, usingSHA1 bool,
	exRE, matchRE []*re.Regexp) (tree *MerkleTree, err error) {

	var (
		dirName string
		files   []os.FileInfo
	)
	found, err := PathExists(pathToDir)
	if err == nil && !found {
		err = FileNotFound
	}
	if err == nil {
		parts := strings.Split(pathToDir, "/")
		if len(parts) == 1 {
			dirName = pathToDir
		} else {
			dirName = parts[len(parts)-1]
		}
		tree, err = NewMerkleTree(dirName, usingSHA1, exRE, matchRE)
	}
	if err == nil {
		var shaX hash.Hash

		// we are promised that this is sorted
		files, err = ioutil.ReadDir(pathToDir)
		if usingSHA1 {
			shaX = sha1.New()
		} else {
			shaX = sha3.NewKeccak256()
		}
		shaXCount := 0
		for i := 0; i < len(files); i++ {
			var node MerkleNodeI
			file := files[i]
			name := file.Name()

			// XXX should continue if any exRE matches
			// XXX should NOT continue if any matchRE match

			pathToFile := path.Join(pathToDir, name)
			mode := file.Mode()
			if mode&os.ModeSymlink != 0 {
				// DEBUG
				fmt.Printf("    LINK: %s, skipping\n", name)
				// END
				continue
			} else if mode.IsDir() {
				node, err = CreateMerkleTreeFromFileSystem(
					pathToFile, usingSHA1, exRE, matchRE)
			} else if mode.IsRegular() {
				// XXX will this ignore symlinks?
				node, err = CreateMerkleLeafFromFileSystem(
					pathToFile, name, usingSHA1)
			}
			if err != nil {
				break
			}
			if node != nil {
				// update tree-level hash
				if node.GetHash() != nil { // IS THIS POSSIBLE?
					shaXCount++
					shaX.Write(node.GetHash())
					tree.nodes = append(tree.nodes, node)
				}
			}
		}
		if err == nil && shaXCount > 0 {
			tree.SetHash(shaX.Sum(nil))
		}
	}
	return
}

// OTHER METHODS AND PROPERTIES =====================================

// Return a pointer to the MerkleTree"s list of component nodes.
// This is a potentially dangerous operation.

func (mt *MerkleTree) Nodes() []MerkleNodeI {
	return mt.nodes
}

func (mt *MerkleTree) AddNode(mn MerkleNodeI) (err error) {

	if mn == nil {
		err = NilMerkleNode
	}
	if err == nil {
		mt.nodes = append(mt.nodes, mn)
	}
	return
}

// SERIALIZATION ////////////////////////////////////////////////

// Called ToString because it returns an error.  XXX Consider dropping
// the indent argument.

func (mt *MerkleTree) ToString(indent string) (str string, err error) {
	var ss []string
	err = mt.ToStrings(indent, &ss) // top level indent
	if err == nil {
		str = strings.Join(ss, "\r\n")
		str += "\r\n"
	}
	return
}

// Serialize a MerkleTree node recursively.
func (mt *MerkleTree) toStringsNotTop(indent string, ss *[]string) (err error) {

	var top string
	topHash := mt.GetHash()
	if len(topHash) == 0 {
		if mt.usingSHA1 {
			top = fmt.Sprintf("%s%s %s/", indent, SHA1_NONE, mt.name)
		} else {
			top = fmt.Sprintf("%s%s %s/", indent, SHA3_NONE, mt.name)
		}
	} else {
		hexHash := hex.EncodeToString(topHash)
		top = fmt.Sprintf("%s%s %s/", indent, hexHash,
			mt.name) // <--- LEVEL 0 NODE
	}
	*ss = append(*ss, top)
	// WORKING HERE - THIS JUST COPIED FROM ToStrings:
	myIndent := indent + "  "
	for i := 0; i < len(mt.nodes); i++ {
		node := mt.nodes[i]
		if node.IsLeaf() {
			mLeaf := node.(*MerkleLeaf)
			err = mLeaf.ToStrings(myIndent, ss)
		} else {
			mTree := node.(*MerkleTree)
			err = mTree.toStringsNotTop(myIndent, ss) // recurses
		}
		if err != nil {
			break
		}
	}
	return
}

//    def toStringNotTop(self, indent):
//        """ indent is the indentation to be used for the top node"""
//        s      = []                             # a list of strings
//        if mt._hash == None:
//            if mt.usingSHA1:
//                top = "%s%s %s/\r\n" % (indent, SHA1_NONE, mt.name)
//            else:
//                top = "%s%s %s/\r\n" % (indent, SHA3_NONE, mt.name)
//        else:
//            top = "%s%s %s/\r\n" % (indent, binascii.b2a_hex(mt._hash),
//                              mt.name)
//        s.append(top)
//        # DEBUG
//        # print "toStringNotTop appends: %s" % top
//        # END
//        indent = indent + "  "              # <--- LEVEL 2+ NODE
//        for node in mt.nodes:
//            if isinstance(node, MerkleLeaf):
//                s.append( node.toString(indent) )
//            else:
//                s.append( node.toStringNotTop(indent) )     # recurses
//
//        return "".join(s)

//        Indent is the initial indentation of the serialized list, NOT the
//        extra indentation added at each recursion, which is fixed at 2 spaces.
//        Using code should take into account that the last line is CR-LF
//        terminated, and so a split on CRLF will generate an extra blank line

func (mt *MerkleTree) ToStrings(indent string, ss *[]string) (err error) {

	var top string
	topHash := mt.GetHash()
	if len(topHash) == 0 {
		if mt.usingSHA1 {
			top = fmt.Sprintf("%s%s %s/", indent, SHA1_NONE, mt.name)
		} else {
			top = fmt.Sprintf("%s%s %s/", indent, SHA3_NONE, mt.name)
		}
	} else {
		hexHash := hex.EncodeToString(topHash)
		top = fmt.Sprintf("%s%s %s/", indent, hexHash,
			mt.name) // <--- LEVEL 0 NODE
	}
	*ss = append(*ss, top)
	myIndent := indent + "  "
	for i := 0; i < len(mt.nodes); i++ {
		node := mt.nodes[i]
		if node.IsLeaf() {
			mLeaf := node.(*MerkleLeaf)
			err = mLeaf.ToStrings(myIndent, ss)
		} else {
			mTree := node.(*MerkleTree)
			err = mTree.toStringsNotTop(myIndent, ss) // recurses
		}
		if err != nil {
			break
		}
	}
	return
}

func (mt *MerkleTree) Equal(any interface{}) bool {
	if any == mt {
		return true
	}
	if any == nil {
		return false
	}
	switch v := any.(type) {
	case *MerkleTree:
		_ = v
	default:
		return false
	}
	other := any.(*MerkleTree) // type assertion
	if other == nil {
		return false
	}

	// compare MerkleNode-level properties (name, hash)
	myNode := &mt.MerkleNode
	otherNode := other.MerkleNode
	if !myNode.Equal(&otherNode) {
		return false
	}
	// compare component nodes
	myLen := len(mt.nodes)
	otherLen := len(other.nodes)
	if myLen != otherLen {
		return false
	}
	for i := 0; i < myLen; i++ {
		if !mt.nodes[i].Equal(other.nodes[i]) { // recurses
			return false
		}
	}
	return true
}
