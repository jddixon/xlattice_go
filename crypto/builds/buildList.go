package builds

// xlattice_go/crypto/builds/buildList.go

import (
	"bufio"
	"bytes"
	"crypto/rsa"
	"encoding/base64"
	xc "github.com/jddixon/xlattice_go/crypto"
	"io"
	"strings"
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

func NewBuildList(pubkey *rsa.PublicKey, title string) (
	bl *BuildList, err error) {

	sl, err := xc.NewSignedList(pubkey, title)
	if err == nil {
		bl = &BuildList{SignedList: *sl}
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

	for err == nil {
		var (
			hash, line []byte
			path       string
			item       *Item
		)
		line, err = xc.NextLineWithoutCRLF(in)
		if err == nil || err == io.EOF {
			if bytes.Equal(line, xc.CONTENT_END) {
				break
			} else {
				// Parse the line.  We expect it to consist of a base64-
				// encoded hash followed by a space followed by a POSIX
				// path.
				line = bytes.Trim(line, " \t")
				if len(line) == 0 {
					err = EmptyContentLine
				} else {
					parts := bytes.Split(line, SPACE_SEP)
					if len(parts) != 2 {
						err = IllFormedContentLine
					} else {
						_, err = base64.StdEncoding.Decode(hash, parts[0])
						if err == nil {
							path = string(parts[1])
						}
					}
				}
				if err == nil {
					item, err = NewItem(hash, path)
					if err == nil {
						bl.content = append(bl.content, item)
					}
				}
			}
		}
	}
	return
}

/**
 * Return the number of content lines
 */
func (bl *BuildList) Size() uint {
	return uint(len(bl.content))
}

/**
 * Return the Nth content item in string form, without any CRLF.
 */
func (bl *BuildList) Get(n uint) (s string, err error) {
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
func (bl *BuildList) Add(hash []byte, name string) (err error) {

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
func (bl *BuildList) GetItemHash(n uint) []byte {
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

func (bl *BuildList) String() (s string) {

	var (
		err error
	)
	pubKey, title, timestamp := bl.Strings()

	// we leave out pubKey because it is newline-terminated
	ss := []string{title, timestamp}
	ss = append(ss, string(xc.CONTENT_START))
	for i := uint(0); err == nil && i < bl.Size(); i++ {
		var line string
		line, err = bl.Get(i)
		if err == nil || err == io.EOF {
			ss = append(ss, line)
			if err == io.EOF {
				err = nil
				break
			}
		}
	}
	if err == nil {
		ss = append(ss, string(xc.CONTENT_END))
		myDigSig := base64.StdEncoding.EncodeToString(bl.GetDigSig())
		ss = append(ss, myDigSig)
		s = string(pubKey) + strings.Join(ss, CRLF) + CRLF
	}
	return
}
func ParseBuildList(rd io.Reader) (bl *BuildList, err error) {
	// super (reader)

	// XXX STUB

	return
}
