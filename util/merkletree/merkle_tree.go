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

func ParseMerkleTreeFromStrings(ss []string) (mt *MerkleTree, err error) {

	// XXX STUB

	return
}

//class MerkleTree(MerkleNode):
// ...

//    @staticmethod
//    def createFromStringArray(s):
//        """
//        The string array is expected to follow conventional indentation
//        rules, with zero indentation on the first line and some multiple
//        of two spaces on all successive lines.
//        """
//        if s == None:
//            raise RuntimeError("null argument")
//
//        # XXX should check TYPE - must be array of strings
//
//        if len(s) == 0:
//            raise RuntimeError("empty string array")
//        (indent, treeHash, dirName) = \
//                            MerkleTree.parseFirstLine(s[0].rstrip())
//        usingSHA1   = (40 == len(treeHash))
//        rootTree    = MerkleTree(dirName, usingSHA1)    # an empty tree
//        rootTree.setHash(treeHash)
//
//        if indent != 0:
//            print "INTERNAL ERROR: initial line indent %d" % indent
//
//        stack      = []
//        stkDepth   = 0
//        curTree    = rootTree
//        stack.append(curTree)           # rootTree
//        stkDepth  += 1                  # always step after pushing tree
//        lastWasDir = False
//
//        # REMEMBER THAT PYTHON HANDLES LARGE RANGES BADLY
//        for n in range(1, len(s)):
//            line = s[n].rstrip()
//#           print "LINE: " + line       # DEBUG
//            if len(line) == 0:
//                n += 1
//                continue
//            # XXX SHOULD/COULD CHECK THAT HASHES ARE OF THE RIGHT TYPE
//            (lineIndent, hash, name, isDir) = MerkleTree.parseOtherLine(line)
//#           print "DEBUG: item %d, lineIndent %d, stkDepth %d, name %s" % (
//#                           n, lineIndent, stkDepth, name)    # DEBUG
//            if lineIndent < stkDepth:
//                while lineIndent < stkDepth:
//                    stkDepth -= 1
//                    stack.pop()
//                curTree = stack[-1]
//
//#               print "DEBUG: item %d, lineIndent %d, stkDepth %d BEYOND LOOP curTree is %s" % (
//#                           n, lineIndent, stkDepth, curTree.name)    # DEBUG
//
//#               MerkleTree.showStack(stack)         # DEBUG
//
//                if not stkDepth == lineIndent:
//                    print "ERROR: stkDepth != lineIndent"
//
//            if isDir:
//                # create and set attributes of new node
//                newTree = MerkleTree(name, usingSHA1)  # , curTree)
//                newTree.setHash(hash)
//                # add the new node into the existing tree
//                curTree.addNode(newTree)
//                stack.append(newTree)
//                stkDepth += 1
//                curTree   = newTree
//#               # DEBUG
//#               MerkleTree.showStack( stack )
//#               # END
//            else:
//                # create and set attributes of new node
//                newNode = MerkleLeaf(name, usingSHA1, hash)
//                # add the new node into the existing tree
//                curTree.addNode(newNode)
//#               print "DEBUG: added node %s to tree %s" % (newNode.name,
//#                                                          curTree.name)
//            n += 1
//        return rootTree         # BAR
//
//    @staticmethod
//    def createFromSerialization(s):
//        if s == None:
//            raise RuntimeError ("MerkleTree.createFromSerialization: no input")
//        sArray = s.split("\r\n")                # note CR-LF
//        return MerkleTree.createFromStringArray(sArray)
//
func ParseMerkleTree(s string) (mt *MerkleTree, err error) {

	if s == "" {
		err = EmptySerialization
	} else {
		ss := strings.Split(s, "\r\n")
		mt, err = ParseMerkleTreeFromStrings(ss)
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
//                usingSHA1 = False
//            else:
//                usingSHA1 = True
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
//                if self._usingSHA1:
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
//    def CreateMerkleTreeFromFileSystem(pathToDir, usingSHA1 = False,
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
			// DEBUG
			fmt.Printf("FILE: %s\n", name)
			// END

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

//    # OTHER METHODS AND PROPERTIES ##################################

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
func (mt *MerkleTree) ToString(indent string) (s string) {

	// XXX STUB
	return
}

//    @property
//    def __str__(self):
//        return self.toString("")

//    def toStringNotTop(self, indent):
//        """ indent is the indentation to be used for the top node"""
//        s      = []                             # a list of strings
//        if self._hash == None:
//            if self._usingSHA1:
//                top = "%s%s %s/\r\n" % (indent, SHA1_NONE, self.name)
//            else:
//                top = "%s%s %s/\r\n" % (indent, SHA3_NONE, self.name)
//        else:
//            top = "%s%s %s/\r\n" % (indent, binascii.b2a_hex(self._hash),
//                              self.name)
//        s.append(top)
//        # DEBUG
//        # print "toStringNotTop appends: %s" % top
//        # END
//        indent = indent + "  "              # <--- LEVEL 2+ NODE
//        for node in self.nodes:
//            if isinstance(node, MerkleLeaf):
//                s.append( node.toString(indent) )
//            else:
//                s.append( node.toStringNotTop(indent) )     # recurses
//
//        return "".join(s)
//
//    def toString(self, indent):
//        """
//        indent is the initial indentation of the serialized list, NOT the
//        extra indentation added at each recursion, which is fixed at 2 spaces.
//        Using code should take into account that the last line is CR-LF
//        terminated, and so a split on CRLF will generate an extra blank line
//        """
//        s      = []                             # a list of strings
//        if self._hash == None:
//            if self._usingSHA1:
//                top = "%s%s %s/\r\n" % (indent, SHA1_NONE, self.name)
//            else:
//                top = "%s%s %s/\r\n" % (indent, SHA3_NONE, self.name)
//        else:
//            top = "%s%s %s/\r\n" % (indent, binascii.b2a_hex(self._hash),
//                              self.name)    # <--- LEVEL 0 NODE
//        s.append(top)
//        myIndent = indent + "  "            # <--- LEVEL 1 NODE
//        for node in self.nodes:
//            if isinstance (node, MerkleLeaf):
//                s.append(node.toString(myIndent))
//            else:
//                s.append( node.toStringNotTop(myIndent) )     # recurses
//
//        return "".join(s)

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

	// compare MerkleNode-level properties (name, hash)
	myNode := mt.MerkleNode
	otherNode := other.MerkleNode
	if !myNode.Equal(otherNode) {
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
