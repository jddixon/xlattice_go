package merkletree

// xlattice_go/util/merkletree/merkletree.go

import (
	"code.google.com/p/go.crypto/sha3"
	"crypto/sha1"
	//"encoding/hex"
	"fmt"
	xu "github.com/jddixon/xlattice_go/util"
	"hash"
	"os"
	"path"
	re "regexp"
	"strings"
)

var _ = fmt.Print

// The path to a tree, and the SHA hash of the path and the treehash.

type MerkleDoc struct {
	bound   bool
	exRE    []*re.Regexp // exclusions
	matchRE []*re.Regexp // must be matched
	tree    *MerkleTree

	path      string
	hash      []byte
	usingSHA1 bool
}

// This belongs in ../, in the utilities directory.

func PathExists(path string) (whether bool, err error) {
	_, err = os.Stat(path)
	if err == nil {
		whether = true
	} else if os.IsNotExist(err) {
		err = nil
	}
	return
}

// XXX "MUST ADD matchRE and exRE and test on their values at this level."

func NewMerkleDoc(pathToDir string, usingSHA1, binding bool, tree *MerkleTree,
	exRE, matchRE []*re.Regexp) (m *MerkleDoc, err error) {

	if pathToDir == "" {
		err = EmptyPath
	}
	if err == nil {
		if strings.HasSuffix(pathToDir, "/") {
			pathToDir += "/"
		}
		self := MerkleDoc{
			exRE:      exRE,
			matchRE:   matchRE,
			path:      pathToDir,
			tree:      tree,
			usingSHA1: usingSHA1,
		}
		if tree != nil {
			var digest hash.Hash
			if usingSHA1 {
				digest = sha1.New()
			} else {
				digest = sha3.NewKeccak256()
			}
			digest.Write(tree.hash)
			digest.Write([]byte(pathToDir))
			self.hash = digest.Sum(nil)
		} else if !binding {
			err = NilTreeButNotBinding
			if err == nil && binding {
				var whether bool
				fullerPath := path.Join(pathToDir, tree.name)
				whether, err = PathExists(fullerPath)
				if err == nil && !whether {
					err = DirectoryNotFound
				}
			}
			if err == nil {
				m = &self
			}
		}
	}
	return
}
func (md *MerkleDoc) Equal(any interface{}) bool {
	if any == md {
		return true
	}
	if any == nil {
		return false
	}
	switch v := any.(type) {
	case *MerkleDoc:
		_ = v
	default:
		return false
	}
	other := any.(*MerkleDoc) // type assertion

	return md.path == other.path &&
		xu.SameBytes(md.hash, other.hash) &&
		md.tree.Equal(other.tree)
}

func (md *MerkleDoc) Hash() []byte {
	return md.hash
}
func (md *MerkleDoc) GetPath() string {
	return md.path
}
func (md *MerkleDoc) SetPath(value string) (err error) {
	// XXX STUB: MUST CHECK VALUE
	md.path = value
	return
}
func (md *MerkleDoc) GetExRE() []*re.Regexp {
	return md.exRE
}
func (md *MerkleDoc) GetMatchRE() []*re.Regexp {
	return md.matchRE
}

//# -------------------------------------------------------------------
//class MerkleDoc():
//
//
//    @property
//    def tree(self):
//        return self._tree
//
//    @tree.setter
//    def tree(self, value):
//        # XXX CHECKS
//        self._tree = value
//
//    @property
//    def bound(self):
//        return self._bound
//
//    @property
//    def usingSHA1(self):
//        return self._usingSHA1
//
//    # QUASI-CONSTRUCTORS ############################################
//    @staticmethod
//    def createFromFileSystem(pathToDir, usingSHA1 = False,
//                             exclusions = None, matches = None):
//        """
//        Create a MerkleDoc based on the information in the directory
//        at pathToDir.  The name of the directory will be the last component
//        of pathToDir.  Return the MerkleTree.
//        """
//        if not pathToDir:
//            raise RuntimeError("cannot create a MerkleTree, no path set")
//        if not os.path.exists(pathToDir):
//            raise RuntimeError(
//                "MerkleTree: directory '%s' does not exist" % self._path)
//        (path, delim, name) = pathToDir.rpartition('/')
//        if path == '':
//            raise RuntimeError("cannot parse inclusive path " + pathToDir)
//        path += '/'
//        exRE = None
//        if exclusions:
//            exRE    = MerkleDoc.makeExRE(exclusions)
//        matchRE = None
//        if matches:
//            matchRE = MerkleDoc.makeMatchRE(matches)
//        tree = MerkleTree.createFromFileSystem(pathToDir, usingSHA1,
//                                            exRE, matchRE)
//        # creates the hash
//        doc  = MerkleDoc(path, usingSHA1, False, tree, exRE, matchRE)
//        doc.bound = True
//        return doc
//
//    @staticmethod
//    def createFromSerialization(s):
//        if s == None:
//            raise RuntimeError ("MerkleDoc.createFromSerialization: no input")
//        sArray = s.split('\r\n')                # note CR-LF
//        return MerkleDoc.createFromStringArray(sArray)
//
//    @staticmethod
//    def createFromStringArray(s):
//        """
//        The string array is expected to follow conventional indentation
//        rules, with zero indentation on the first line and some multiple
//        of two spaces on all successive lines.
//        """
//        if s == None:
//            raise RuntimeError('null argument')
//        # XXX check TYPE - must be array of strings
//        if len(s) == 0:
//            raise RuntimeError("empty string array")
//
//        (docHash, docPath) = \
//                            MerkleDoc.parseFirstLine(s[0].rstrip())
//#       print "DEBUG: doc first line: hash = %s, path = %s" % (
//#                               docHash, docPath)
//        usingSHA1 = (40 == len(docHash))
//
//        tree = MerkleTree.createFromStringArray( s[1:] )
//
//        #def __init__ (self, path, binding = False, tree = None,
//        #    exRE    = None,    # exclusions, which are Regular Expressions
//        #    matchRE = None):   # matches, also Regular Expressions
//        doc = MerkleDoc( docPath, usingSHA1, False, tree )
//        doc.hash = docHash
//        return doc
//
//    # CLASS METHODS #################################################
//    @staticmethod
//    def parseFirstLine(line):
//        line = line.rstrip()
//        m = re.match(MerkleDoc.FIRST_LINE_PAT_1, line)
//        if m == None:
//            m = re.match(MerkleDoc.FIRST_LINE_PAT_3, line)
//        if m == None:
//            raise RuntimeError(
//                    "MerkleDoc first line <%s> does not match expected pattern" %  line)
//        docHash  = m.group(1)
//        docPath  = m.group(2)          # includes terminating slash
//        return (docHash, docPath)

// Given a string array of regular expressions, append a list of standard
// exclusions, and then return the compiled regexp.

func MakeExRE(excl []string) (exRE *re.Regexp, err error) {
	excl = append(excl, "^\\.$")       // .
	excl = append(excl, "^\\.\\.$")    // ..
	excl = append(excl, "^\\.merkle$") // merkletree hidden files
	excl = append(excl, "^\\.git$")    // git control data
	excl = append(excl, "^\\.svn$")    // subversion control data
	// some might disagree with these:
	excl = append(excl, "^junk")
	excl = append(excl, "^\\..*\\.swp$") // vi editor files

	exPat := strings.Join(excl, "|")
	exRE, err = re.Compile(exPat)
	return
}

// Given a possibly empty list of expressions to be matched, return
// nil if the list is empty or a regular expression which matches
// any of the patterns in the list.
func MakeMatchRE(matches []string) (matchRE *re.Regexp, err error) {
	if len(matches) > 0 {
		matchPat := strings.Join(matches, "|")
		matchRE, err = re.Compile(matchPat)
	}
	return
}

//    # SERIALIZATION #################################################
//    def __str__(self):
//        return self.toString()
//
//    def toString(self):
//        return ''.join([
//            "%s %s\r\n" % ( self.hash, self.path),
//            self._tree.toString('')
//            ])	// GEEP
