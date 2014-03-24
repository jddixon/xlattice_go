package chunks

// xlattice_go/protocol/chunks/chunkListAssyDisassy_test.go

import (
	"bytes"
	"code.google.com/p/go.crypto/sha3"
	"crypto/rand"
	"crypto/rsa"
	"encoding/hex"
	"fmt"
	xi "github.com/jddixon/xlattice_go/nodeID"
	xr "github.com/jddixon/xlattice_go/rnglib"
	"github.com/jddixon/xlattice_go/u"
	xu "github.com/jddixon/xlattice_go/util"
	xf "github.com/jddixon/xlattice_go/util/lfs"
	. "launchpad.net/gocheck"
	"os"
	"path"
	"time"
)

var _ = fmt.Print

func (s *XLSuite) TestChunkListAssyDisassy(c *C) {
	if VERBOSITY > 0 {
		fmt.Println("TEST_CHUNK_LIST_ASSY_DISASSY")
	}
	rng := xr.MakeSimpleRNG()

	// make a slice 3 to 7 chunks long, fill with random-ish data ---
	chunkCount := 3 + rng.Intn(5) // so 3 to 7, inclusive
	lastChunkLen := 1 + rng.Intn(MAX_DATA_BYTES-1)
	dataLen := (chunkCount-1)*MAX_DATA_BYTES + lastChunkLen
	data := make([]byte, dataLen)
	rng.NextBytes(data)

	// calculate datum, the SHA3 hash of the data -------------------
	d := sha3.NewKeccak256()
	d.Write(data)
	hash := d.Sum(nil)
	datum, err := xi.NewNodeID(hash)
	c.Assert(err, IsNil)

	// DEBUG
	datumStr := hex.EncodeToString(datum.Value())
	fmt.Printf("DATUM is %s\n", datumStr)
	// END

	// create tmp if it doesn't exist -------------------------------
	found, err := xf.PathExists("tmp")
	c.Assert(err, IsNil)
	if !found {
		err = os.MkdirAll("tmp", 0755)
		c.Assert(err, IsNil)
	}

	// create scratch subdir with unique name -----------------------
	var pathToU string
	for {
		dirName := rng.NextFileName(8)
		pathToU = path.Join("tmp", dirName)
		found, err = xf.PathExists(pathToU)
		c.Assert(err, IsNil)
		if !found {
			break
		}
	}

	// create a FLAT uDir at that point -----------------------------
	myU, err := u.New(pathToU, u.DIR_FLAT, 0) // 0 means default perm
	c.Assert(err, IsNil)

	// write the test data into uDir --------------------------------
	bytesWritten, key, err := myU.PutData(data, datum.Value())
	c.Assert(err, IsNil)
	c.Assert(bytes.Equal(datum.Value(), key), Equals, true)
	c.Assert(bytesWritten, Equals, int64(dataLen))

	skPriv, err := rsa.GenerateKey(rand.Reader, 1024) // cheap key
	sk := &skPriv.PublicKey
	c.Assert(err, IsNil)
	c.Assert(skPriv, NotNil)

	// Verify the file is present in uDir ---------------------------
	// (yes this is a test of uDir logic but these are early days ---
	// XXX uDir.Exist(arg) - arg should be []byte, no string
	keyStr := hex.EncodeToString(key)
	found, err = myU.Exists(keyStr)
	c.Assert(err, IsNil)
	c.Assert(found, Equals, true)

	// DEBUG
	keyPath, err := myU.GetPathForKey(keyStr)
	c.Assert(err, IsNil)
	fmt.Printf("%s is present in uDir at %s\n", keyStr, keyPath)
	// END

	// use the data file to build a chunkList, writing the chunks ---
	title := rng.NextFileName(8)
	now := xu.Timestamp(time.Now().UnixNano())
	// DEBUG
	fmt.Printf("the UTC time is %s\n", now.String())
	// END

	// make a reader --------------------------------------
	pathToData, err := myU.GetPathForKey(keyStr)
	c.Assert(err, IsNil)
	reader, err := os.Open(pathToData) // open for read only
	c.Assert(err, IsNil)
	defer reader.Close()

	chunkList, err := NewChunkList(sk, title, now, reader, int64(dataLen), key, myU)
	c.Assert(err, IsNil)
	err = chunkList.Sign(skPriv)
	c.Assert(err, IsNil)

	// XXX STUB
	_ = sk // DEBUG
	_ = title

	// rebuild the complete file from the chunkList and files present
	// in myU

	// verify that the rebuilt file is identical to the original ----
}
