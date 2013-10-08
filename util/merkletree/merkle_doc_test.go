package merkletree

// xlatttice_go/util/merkletree/merkle_doc_test.go

import (
	//"code.google.com/p/go.crypto/sha3"
	//"crypto/sha1"
	// "encoding/hex"
	"fmt"
	xr "github.com/jddixon/xlattice_go/rnglib"
	//xu "github.com/jddixon/xlattice_go/util"
	//"io/ioutil"
	. "launchpad.net/gocheck"
	re "regexp"
	//"strings"
)

func (s *XLSuite) doTestMerkleDoc(c *C, rng *xr.PRNG, usingSHA1 bool) {
	fileName := rng.NextFileName(8)

	// XXX STUB
	_ = fileName
}

func (s *XLSuite) TestMerkleDoc(c *C) {
	if VERBOSITY > 0 {
		fmt.Println("TEST_MERKLE_DOC")
	}
	rng := xr.MakeSimpleRNG()
	s.doTestMerkleDoc(c, rng, true)  // using SHA1
	s.doTestMerkleDoc(c, rng, false) // not using SHA1
}

// REGEXP TESTS =====================================================
func (s *XLSuite) doTestForExpectedExclusions(c *C, exRE *re.Regexp) {
	// should always match
	c.Assert(exRE.MatchString("."), Equals, true)
	c.Assert(exRE.MatchString(".."), Equals, true)
	c.Assert(exRE.MatchString(".merkle"), Equals, true)
	c.Assert(exRE.MatchString(".svn"), Equals, true)
	c.Assert(exRE.MatchString(".foo.swp"), Equals, true)
	c.Assert(exRE.MatchString("junkEverywhere"), Equals, true)
}
func (s *XLSuite) doTestForExpectedMatches(c *C,
	matchRE *re.Regexp, names []string) {

	for i := 0; i < len(names); i++ {
		name := names[i]
		c.Assert(matchRE.MatchString(name), Equals, true)
	}
}
func (s *XLSuite) doTestForExpectedMatchFailures(c *C,
	matchRE *re.Regexp, names []string) {

	for i := 0; i < len(names); i++ {
		name := names[i]
		m := matchRE.MatchString(name)
		if m {
			fmt.Printf("WE HAVE A MATCH ON '%s'\n", name)
			// self.assertEquals( None, where )
		}
	}
}

// test utility for making excluded file name regexes

func (s *XLSuite) TestMakeExRE(c *C) {
	exRE, err := MakeExRE(nil)
	c.Assert(err, IsNil)
	c.Assert(exRE, NotNil)
	s.doTestForExpectedExclusions(c, exRE)

	// should not be present
	c.Assert(exRE.MatchString("bar"), Equals, false)
	c.Assert(exRE.MatchString("foo"), Equals, false)

	var exc []string
	exc = append(exc, "^foo")
	exc = append(exc, "bar$")
	exc = append(exc, "^junk*")
	exRE, err = MakeExRE(exc)
	c.Assert(err, IsNil)
	s.doTestForExpectedExclusions(c, exRE)

	c.Assert(exRE.MatchString("foobarf"), Equals, true)
	c.Assert(exRE.MatchString(" foobarf"), Equals, false)
	c.Assert(exRE.MatchString(" foobarf"), Equals, false)

	// bear in mind that match must be at the beginning
	c.Assert(exRE.MatchString("ohMybar"), Equals, true)
	c.Assert(exRE.MatchString("ohMybarf"), Equals, false)
	c.Assert(exRE.MatchString("junky"), Equals, true)
	c.Assert(exRE.MatchString(" junk"), Equals, false)
}

// test utility for making matched file name regexes

func (s *XLSuite) TestMakeMatchRE(c *C) {
	matchRE, err := MakeMatchRE(nil)
	c.Assert(err, IsNil)
	c.Assert(matchRE, IsNil)

	var matches []string
	matches = append(matches, "^foo")
	matches = append(matches, "bar$")
	matches = append(matches, "^junk*")
	matchRE, err = MakeMatchRE(matches)
	c.Assert(err, IsNil)
	cases := []string{"foo", "foolish", "roobar", "junky"}
	s.doTestForExpectedMatches(c, matchRE, cases)

	cases = []string{" foo", "roobarf", "myjunk"}
	s.doTestForExpectedMatchFailures(c, matchRE, cases)

	matches = []string{"\\.tgz$"}
	matchRE, err = MakeMatchRE(matches)
	c.Assert(err, IsNil)

	cases = []string{"junk.tgz", "notSoFoolish.tgz"}
	s.doTestForExpectedMatches(c, matchRE, cases)
	cases = []string{"junk.tar.gz", "foolish.tar.gz"}
	s.doTestForExpectedMatchFailures(c, matchRE, cases)

	matches = []string{"\\.tgz$", "\\.tar\\.gz$"}
	matchRE, err = MakeMatchRE(matches)
	c.Assert(err, IsNil)

	cases = []string{
		"junk.tgz", "notSoFoolish.tgz", "junk.tar.gz", "ohHello.tar.gz"}
	s.doTestForExpectedMatches(c, matchRE, cases)

	cases = []string{"junk.gz", "foolish.tar"}
	s.doTestForExpectedMatchFailures(c, matchRE, cases)
}
