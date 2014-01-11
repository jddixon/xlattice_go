package builds

// xlattice_go/crypto/builds/buildList_test.go

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"fmt"
	xc "github.com/jddixon/xlattice_go/crypto"
	xr "github.com/jddixon/xlattice_go/rnglib"
	. "launchpad.net/gocheck"
)

var _ = fmt.Print

const (
	// a folded RSA public key
	docPK1    = "rsa AL0zGtdGkuJdH1vd4TaUMmRvdEBepnGfAbvZXPkdsVq367VUevbfzNL4W6u+Ks8+BksZzZPc"
	docPK2    = "yLJsnDZr7mE/rHSwQ7la1HlSWwNDlhQtCnKTlSoqffVhofhtak/SqBOJVLkWrouaK60uCiZV0Hw"
	docPK3    = "YTM6Pqo8sqYinA3W8mvK2tsW/ 65537"
	docPubKey = docPK1 + FOLD + docPK2 + FOLD + docPK3

	docTitle      = "document 1"
	docTime       = "2004-11-18 20:03:34"
	docEncodedSig = "tIQJ+7Y27eIyQCb3esTgU/AdBfPDAGEOhU/KShAo5N5dfxtjkH04N5IwvyftEJd5jM0kHB1LD1TtavoxZ0gx4eADizHcDjEpZOiO+wUHIcbGsuvLUvZvBttPPBRuRfZgZXkvvSMBX0KIwRVgFqwaRB5gzQyD2skcP2kGFBWrFdM="
	testDoc       = docPubKey + CRLF +
		docTitle + CRLF +
		docTime + CRLF +
		docEncodedSig
)

func (s *XLSuite) TestEmptyBuildList(c *C) {
	var (
		err    error
		myList *BuildList
		key    *rsa.PrivateKey
		pubKey *rsa.PublicKey
	)
	key, err = rsa.GenerateKey(rand.Reader, 1024)
	c.Assert(err, IsNil)
	pubKey = &key.PublicKey
	myList, err = NewBuildList(pubKey, "document 1")
	c.Assert(err, IsNil)
	c.Assert(myList, NotNil)
	c.Assert(myList.Size(), Equals, uint(0))
	c.Assert(myList.IsSigned(), Equals, false)

	myList.Sign(key)
	c.Assert(myList.IsSigned(), Equals, true)
	c.Assert(myList.Verify(), IsNil)

	err = myList.Sign(key)
	c.Assert(err, Equals, xc.ListAlreadySigned)
}

func (s *XLSuite) TestGeneratedBuildList(c *C) {
	var (
		err    error
		myList *BuildList
		key    *rsa.PrivateKey
		pubKey *rsa.PublicKey
	)
	rng := xr.MakeSimpleRNG()

	hash0 := make([]byte, xc.SHA1_LEN)
	hash1 := make([]byte, xc.SHA1_LEN)
	hash2 := make([]byte, xc.SHA1_LEN)
	hash3 := make([]byte, xc.SHA1_LEN)
	rng.NextBytes(hash0)
	rng.NextBytes(hash1)
	rng.NextBytes(hash2)
	rng.NextBytes(hash3)

	key, err = rsa.GenerateKey(rand.Reader, 1024)
	c.Assert(err, IsNil)
	pubKey = &key.PublicKey
	myList, err = NewBuildList(pubKey, "document 1")
	c.Assert(err, IsNil)
	c.Assert(myList, NotNil)
	c.Assert(myList.Size(), Equals, uint(0))
	c.Assert(myList.IsSigned(), Equals, false)

	// XXX NOTE WE CAN ADD DUPLICATE OR CONFLICTING ITEMS !! XXX
	err = myList.Add(hash0, "fileForHash0")
	c.Assert(err, IsNil)
	c.Assert(myList.Size(), Equals, uint(1))

	err = myList.Add(hash1, "fileForHash1")
	c.Assert(err, IsNil)
	c.Assert(myList.Size(), Equals, uint(2))

	err = myList.Add(hash2, "fileForHash2")
	c.Assert(err, IsNil)
	c.Assert(myList.Size(), Equals, uint(3))

	err = myList.Add(hash3, "fileForHash3")
	c.Assert(err, IsNil)
	c.Assert(myList.Size(), Equals, uint(4))

	// check (arbitrarily) second content line
	expected1 := base64.StdEncoding.EncodeToString(hash1) + " fileForHash1"
	actual1, err := myList.Get(1)
	c.Assert(err, IsNil)
	c.Assert(expected1, Equals, actual1)
	err = myList.Sign(key)
	c.Assert(err, IsNil)
	c.Assert(myList.IsSigned(), Equals, true)
	err = myList.Verify()
	c.Assert(err, IsNil)

	myDoc := myList.String()
	_ = myDoc

	//    BuildList list2 = ParseBuildList (new StringReader(myDoc))
	//    c.AssertNotNull(list2)
	//    c.AssertEquals(4, list2.Size())
	//    c.AssertTrue (list2.isSigned())
	//    c.AssertEquals(myDoc, list2.toString())
	//    c.AssertTrue (list2.verify())
	//
	//    // test item gets - sloppy naming, so can't loop :-(
	//    byte[] b = myList.getHash(0)
	//    c.AssertEquals (hash0.length, b.length)
	//    for (int i = 0; i < hash0.length; i++)
	//        c.AssertEquals (hash0[i], b[i])
	//    c.AssertEquals ("fileForHash0", myList.getPath(0))
	//
	//    b = myList.getHash(1)
	//    c.AssertEquals (hash1.length, b.length)
	//    for (int i = 0; i < hash1.length; i++)
	//        c.AssertEquals (hash1[i], b[i])
	//    c.AssertEquals ("fileForHash1", myList.getPath(1))
}
