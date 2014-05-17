package httpd

// xlattice_go/httpd/siteList.go

import (
	"bufio"
	"bytes"
	"crypto/rsa"
	"encoding/base64"
	xc "github.com/jddixon/xlCrypto_go"
	"io"
	"strings"
)

/**
 * Serialized, a site list is a list of Web site names.  The names
 * must end with a File.separator. Lines end with CRLF.
 *
 * XXX TODO: Change to sort entries, using this to eliminate
 * XXX any duplicates.
 *
 * XXX Also ADD PORT NUMBERS with default to 80.
 */

type SiteList struct {
	content []string
	xc.SignedList
}

func NewSiteList(pubkey *rsa.PublicKey, title string) (
	stl *SiteList, err error) {

	sl, err := xc.NewSignedList(pubkey, title)
	if err == nil {
		stl = &SiteList{SignedList: *sl}
	}
	return
}

// SignedList METHODS ///////////////////////////////////////////

// Return the Nth content item in string form, without any CRLF.
func (stl *SiteList) Get(n uint) (s string, err error) {
	if n < 0 || stl.Size() <= n {
		err = xc.NdxOutOfRange
	} else {
		s = stl.content[n]
	}
	return
}

/**
 * Read a series of content lines, each consisting of a simple
 * name.  This is an Internet domain name and so may contain no
 * spaces or other delimiters.  The name must end with the
 * file separator (File.separator).
 *
 * The text of the line, excluding the line terminator, is
 * included in the digest.
 */
func (stl *SiteList) ReadContents(in *bufio.Reader) (err error) {

	for err == nil {
		var line []byte
		line, err = xc.NextLineWithoutCRLF(in)
		if err == nil || err == io.EOF {
			if bytes.Equal(line, xc.CONTENT_END) {
				break
			} else {
				stl.content = append(stl.content, string(line))
			}
		}
	}
	return
}

/** @return the number of content lines */
func (stl *SiteList) Size() uint {
	return uint(len(stl.content))
}

// SiteList-SPECIFIC METHODS ///////////////////////////////////

/**
 * Add a content line to the SiteList.  The line is just a
 * domain name terminated by a File.separator.
 *
 * @param name  file or path name of item
 */

func (stl *SiteList) AddItem(s string) (err error) {

	if s == "" {
		err = EmptySiteDomainName
	} else if stl.IsSigned() {
		err = xc.CantAddToSignedList
	} else {
		stl.content = append(stl.content, s)
	}
	return
}

/**
 * Serialize the entire document.  If any error is encountered, this
 * function silently returns an empty string.
 */
func (stl *SiteList) String() (s string) {

	var (
		err error
	)
	pubKey, title, timestamp := stl.SignedList.Strings()

	// we leave out pubKey because it is newline-terminated
	ss := []string{title, timestamp}

	// content lines --------------------------------------
	ss = append(ss, string(xc.CONTENT_START))
	for i := uint(0); err == nil && i < stl.Size(); i++ {
		var line string
		line, err = stl.Get(i)
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

		myDigSig := base64.StdEncoding.EncodeToString(stl.GetDigSig())
		ss = append(ss, myDigSig)
		s = string(pubKey) + strings.Join(ss, xc.CRLF) + xc.CRLF
	}
	return
}
func ParseSiteList(in io.Reader) (stl *SiteList, err error) {
	var (
		digSig, line []byte
	)
	bin := bufio.NewReader(in)
	sl, err := xc.ParseSignedList(bin)
	if err == nil {
		stl = &SiteList{SignedList: *sl}
		err = stl.ReadContents(bin)
		if err == nil {
			// try to read the digital signature line
			line, err = xc.NextLineWithoutCRLF(bin)
			if err == nil || err == io.EOF {
				digSig, err = base64.StdEncoding.DecodeString(string(line))
			}
			if err == nil || err == io.EOF {
				stl.SetDigSig(digSig)
			}
		}
	}
	if err == io.EOF {
		err = nil
	}
	return
}
