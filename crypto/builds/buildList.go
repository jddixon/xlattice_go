package builds

// xlattice_go/crypto/builds/buildList.go

import (
	"bufio"
	"crypto/rsa"
	xc "github.com/jddixon/xlattice_go/crypto"
	"io"
)

/**
 * Serialized, a build list is a list of files and their extended hashes.
 * Each content line starts with base64-encoded extended hash which is
 * followed by a single space and then the file name, including the
 * path.  Lines end with CRLF.
 *
 * The hash for a serialized BuildList, its title key, is the 20-byte
 * SignedList hash, an SHA1-based function of the BuildList's title and
 * RSA public key.
 *
 * The digital signature in the last line is calculated from the
 * SHA1 digest of the header lines (public key, title, and timestamp
 * lines, each CRLF-terminated) and the content lines.
 */
type BuildList struct {
	content []*Item
	xc.SignedList
}

const (
     SEPARATOR   = "/"
     SEPARATOR_CHAR = '/'
)

func NewBuildList(pubkey *rsa.PublicKey, title string) (
	bl *BuildList, err error) {

	sl, err := xc.NewSignedList(pubkey, title)
	if err == nil {
		bl = &BuildList{ SignedList: *sl }
	} 
	return
}

// SignedList ABSTRACT METHODS //////////////////////////////////

/**
 * Read a series of content lines, each consisting of a hash
 * followed by a space followed by a file name.  The hash is
 * base-64 encoded.
 *
 * The text of the line, excluding the line terminator, is
 * included in the digest.
 */
func (bl *BuildList) ReadContents(in *bufio.Reader) (err error) {

	// XXX STUB

	return
}

/** 
 * Return the number of content lines 
 */
func (bl *BuildList) Size () int {
    return len(bl.content)
}

/**
 * Return the Nth content item in string form, without any CRLF.
 */
func (bl *BuildList) Get(n int) (s string, err error) {
	if n < 0 || bl.Size() <= n {
		err = xc.NdxOutOfRange
	} else {
		s = bl.content[n].String()
	}
	return
}

/**
 * Add a content line to the BuildList.  In string form, the
 * content line begins with the extended hash of the Item
 * (the content hash if it is a data file) followed by a space
 * followed by the name of the Item.  If the name is a path,
 * the SEPARATOR character is a UNIX/Linux-style forward slash,
 * BuildList.SEPARATOR.
 *
 * @param hash  extended hash of Item, its file key
 * @param name  file or path name of Item
 * @return      reference to this BuildList, to ease chaining
 */
func (bl *BuildList) Add (hash []byte, name string) (err error) {
	
	if bl.IsSigned() {
		err = CantAddToSignedList
	} else {
		var item *Item
		item, err = NewItem(hash, name)
		if err == nil {
			bl.content = append(bl.content, item)
		}
	}
    return
}
/** 
 * Return the SHA1 hash for the Nth Item.  
 * XXX Should be modified to return a copy.
 */
func (bl *BuildList) GetItemHash(n int) []byte {
    return bl.content[n].ehash
}
/**
 * Returns the path + fileName for the Nth content line, in
 * a form usable with the operating system.  That is, the
 * SEPARATOR is File.SEPARATOR instead of BuildList.SEPARATOR,
 * if there is a difference.
 *
 * @param n content line
 * @return the path + file name for the Nth Item
 */
func (bl *BuildList) GetPath(n uint) string {

	// XXX NEEDS VALIDATION
    return bl.content[n].path
}

func ParseBuildList (rd io.Reader) (bl *BuildList, err error) {
    // super (reader)

	// XXX STUB

	return
}

