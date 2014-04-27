package crypto

// xlattice_go/crypto/signed_list_test.go

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	xr "github.com/jddixon/xlattice_go/rnglib"
	. "gopkg.in/check.v1"
	//. "launchpad.net/gocheck"
	"strings"
)

var _ = fmt.Print

/**
 * Generate a few random RSA keys, create MyLists, test.
 */
func (s *XLSuite) TestGenerateSignedList(c *C) {
	rng := xr.MakeSimpleRNG()
	_ = rng

	for i := 0; i < 8; i++ {

		// create keys
		skPriv, err := rsa.GenerateKey(rand.Reader, 1024)
		c.Assert(err, IsNil)
		c.Assert(skPriv, NotNil)
		pubKey := skPriv.PublicKey
		c.Assert(pubKey, NotNil)

		// create and test signed list
		myList, err := NewMockSignedList(&pubKey, "document 1")
		c.Assert(err, IsNil)
		c.Assert(myList, NotNil)

		// add a few lines
		count := uint(3 + rng.Intn(10))
		for i := uint(0); i < count; i++ {
			s := rng.NextFileName(16)
			n := myList.AddItem(s)
			c.Assert(n, Equals, i)
		}
		c.Assert(myList.Size(), Equals, count)

		err = myList.Sign(skPriv)
		c.Assert(myList.IsSigned(), Equals, true)
		c.Assert(myList.Verify(), IsNil)

		// Generate a new SignedList from the serialization of the
		// current one, use it to test Reader constructor.
		myDoc := myList.String()
		c.Assert(myDoc, Not(Equals), "")

		// deserialize = parse it
		reader := strings.NewReader(myDoc)
		myList2, err := ParseMockSignedList(reader)
		c.Assert(err, IsNil)
		c.Assert(myList2, NotNil)

		// assert that it's signed
		c.Assert(myList2.IsSigned(), Equals, true)

		// verify the digSig
		c.Assert(myList2.Verify(), IsNil)
	}
}

func (s *XLSuite) TestListHash(c *C) {
	rng := xr.MakeSimpleRNG()
	_ = rng

	for i := 0; i < 8; i++ {
		skPriv, err := rsa.GenerateKey(rand.Reader, 1024)
		c.Assert(err, IsNil)
		c.Assert(skPriv, NotNil)
		pubKey := skPriv.PublicKey
		c.Assert(pubKey, NotNil)

		myList, err := NewMockSignedList(&pubKey, "document 1")
		c.Assert(err, IsNil)
		c.Assert(myList, NotNil)

		// add a few lines
		count := uint(3 + rng.Intn(10))
		for i := uint(0); i < count; i++ {
			s := rng.NextFileName(16)
			n := myList.AddItem(s)
			c.Assert(n, Equals, i)
		}
		c.Assert(myList.Size(), Equals, count)

		myHash := myList.GetHash()
		list2, err := NewMockSignedList(&pubKey, "document 1")
		c.Assert(err, IsNil)
		hash2 := list2.GetHash()
		// pubkey and title the same so hashes are the same
		c.Assert(bytes.Equal(myHash, hash2), Equals, true)

		list2, err = NewMockSignedList(&pubKey, "document 2")
		c.Assert(err, IsNil)
		hash2 = list2.GetHash()
		// titles differ so hashes differ
		c.Assert(bytes.Equal(myHash, hash2), Equals, false)

		//      // a build list with the same key and title has same hash
		//      BuildList buildList = new BuildList(pubKey, "document 1")
		//      bHash = buildList.GetHash()
		//      c.AssertEquals (20, bHash.length)
		//      checkSameHash (bHash, myHash)
	}
}
