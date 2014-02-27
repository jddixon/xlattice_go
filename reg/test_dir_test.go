package reg

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"fmt"
	xc "github.com/jddixon/xlattice_go/crypto"
	xu "github.com/jddixon/xlattice_go/util"
	"io/ioutil"
	. "launchpad.net/gocheck"
	"path"
)

// AAAA makes it run first.
func (s *XLSuite) TestAAAATestDir(c *C) {
	if VERBOSITY > 0 {
		fmt.Println("TEST_TEST_DIR")
	}

	// 00 READ AND INTERPRET test_dir/regCred.dat ///////////////////

	rcFile := path.Join("test_dir", "regCred.dat")
	rcData, err := ioutil.ReadFile(rcFile)
	c.Assert(err, IsNil)
	c.Assert(len(rcData) > 0, Equals, true)
	rc, err := ParseRegCred(string(rcData))
	c.Assert(err, IsNil)
	c.Assert(rc, NotNil)

	// 00 HELLO - REPLY TESTS ///////////////////////////////////////

	// 1. Read key_rsa as key *rsa.PrivateKey
	keyFile := path.Join("test_dir", "key-rsa")
	kd, err := ioutil.ReadFile(keyFile)
	c.Assert(err, IsNil)
	c.Assert(len(kd) > 0, Equals, true)
	key, err := xc.RSAPrivateKeyFromDisk(kd)
	c.Assert(err, IsNil)
	c.Assert(key, NotNil)

	// 2. Extract public key as pubkey *rsa.PublicKey
	pubKey := key.PublicKey

	// 3. Read key_rsa.pub as pubkey2 *rsa.PublicKey
	pubKeyFile := path.Join("test_dir", "key-rsa.pub")
	pkd, err := ioutil.ReadFile(pubKeyFile)
	c.Assert(err, IsNil)
	c.Assert(len(pkd) > 0, Equals, true)
	pubKey2, err := xc.RSAPubKeyFromDisk(pkd)
	c.Assert(err, IsNil)
	c.Assert(pubKey2, NotNil)

	// 4. Verify pubkey == pubkey2
	c.Assert(&pubKey, DeepEquals, pubKey2)

	// 5. Read version1.str as v1Str
	v1File := path.Join("test_dir", "version1.str")
	v, err := ioutil.ReadFile(v1File)
	c.Assert(err, IsNil)
	c.Assert(len(v) > 0, Equals, true)
	v1Str := string(v)

	// 6. Read version1 as []byte
	v1File = path.Join("test_dir", "version1")
	version1, err := ioutil.ReadFile(v1File)
	c.Assert(err, IsNil)
	c.Assert(len(version1), Equals, 4) // length of DecimalVersion

	// 7. Convert version1 to dv1 DecimalVersion
	dv1, err := xu.VersionFromBytes(version1)
	c.Assert(err, IsNil)

	// 8. Verify v1Str == dv1.String()
	c.Assert(v1Str, Equals, dv1.String())

	// 9, 10, 11, 12 same as 5-8 for version2 -----------------------

	// 9. Read version2.str as v2Str
	v2File := path.Join("test_dir", "version2.str")
	v, err = ioutil.ReadFile(v2File)
	c.Assert(err, IsNil)
	c.Assert(len(v) > 0, Equals, true)
	v2Str := string(v)

	// 10. Read version2 as []byte
	v2File = path.Join("test_dir", "version2")
	version2, err := ioutil.ReadFile(v2File)
	c.Assert(err, IsNil)
	c.Assert(len(version2), Equals, 4) // length of DecimalVersion

	// 11. Convert version2 to dv2 DecimalVersion
	dv2, err := xu.VersionFromBytes(version2)
	c.Assert(err, IsNil)

	// 12. Verify v2Str == dv2.String()
	c.Assert(v2Str, Equals, dv2.String())

	// 13, 14, 15, 16 read iv1, key1, salt1, hello-data as []byte ---

	iv1, err := ioutil.ReadFile(path.Join("test_dir", "iv1"))
	c.Assert(err, IsNil)

	key1, err := ioutil.ReadFile(path.Join("test_dir", "key1"))
	c.Assert(err, IsNil)

	salt1, err := ioutil.ReadFile(path.Join("test_dir", "salt1"))
	c.Assert(err, IsNil)

	helloData, err := ioutil.ReadFile(path.Join("test_dir", "hello-data"))
	c.Assert(err, IsNil)

	// 17. helloPlain = iv1 + key1 + salt1 + version1
	var helloPlain []byte
	helloPlain = append(helloPlain, iv1...)
	helloPlain = append(helloPlain, key1...)
	helloPlain = append(helloPlain, salt1...)
	helloPlain = append(helloPlain, version1...)

	// 18. Verify helloPlain == helloData
	bytes.Equal(helloPlain, helloData)

	// 19. Read hello-encrypted as []byte
	helloEncrypted, err := ioutil.ReadFile(
		path.Join("test_dir", "hello-encrypted"))
	c.Assert(err, IsNil)

	// 20. Decrypt helloEncrypted using key => helloDecrypted
	helloDecrypted, err := rsa.DecryptOAEP(sha1.New(), rand.Reader,
		key, helloEncrypted, nil)
	c.Assert(err, IsNil)
	c.Assert(len(helloDecrypted) == 0, Equals, false)

	// 21. Verify helloDecrypted == helloData
	c.Assert(bytes.Equal(helloDecrypted, helloData), Equals, true)

	// 22, 23, 24, 25, 26 read iv2, key2, salt2, padding, reply-data as []byte

	iv2, err := ioutil.ReadFile(path.Join("test_dir", "iv2"))
	c.Assert(err, IsNil)

	key2, err := ioutil.ReadFile(path.Join("test_dir", "key2"))
	c.Assert(err, IsNil)

	salt2, err := ioutil.ReadFile(path.Join("test_dir", "salt2"))
	c.Assert(err, IsNil)

	padding, err := ioutil.ReadFile(path.Join("test_dir", "padding"))
	c.Assert(err, IsNil)

	replyData, err := ioutil.ReadFile(path.Join("test_dir", "reply-data"))
	c.Assert(err, IsNil)

	// 27. helloReply = concat iv2, key2, salt2, version2, salt1, padding
	var helloReply []byte
	helloReply = append(helloReply, iv2...)
	helloReply = append(helloReply, key2...)
	helloReply = append(helloReply, salt2...)
	helloReply = append(helloReply, version2...)
	helloReply = append(helloReply, salt1...)
	helloReply = append(helloReply, padding...)

	// 28. Verify helloReply == replyData
	c.Assert(bytes.Equal(replyData, helloReply), Equals, true)

	// 29. Create aesEngineS1 from iv1, key1

	// 30. helloReplyMsg = aesEngineS1.encrypt(helloReply)

	// 31. Read reply-encrypted as replyEncrypted []byte

	// 32. Verify helloReplyMsg == replyEncrypted

	// 33. Create aesEngineC1 from iv1, key1

	// 34. Use aesEngineC1.decrypt(replyEncrypted) => replyDecrypted

	// 35. Verify replyDecrypted == replyData
}
