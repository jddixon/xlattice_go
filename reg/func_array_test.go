package reg

// xlattice_go/reg/func_array_test.go

import (
	"fmt"
	//xc "github.com/jddixon/xlattice_go/crypto"
	//xn "github.com/jddixon/xlattice_go/node"
	. "launchpad.net/gocheck"
)

func aa(in string) (out string) {
	return in + "AA"
}
func bb(in string) (out string) {
	return in + "BB"
}
func cc(in string) (out string) {
	return in + "CC"
}
func (s *XLSuite) TestFuncArray(c *C) {
	if VERBOSITY > 0 {
		fmt.Println("TEST_FUNC_ARRAY")
	}
	sweetFA := make([]interface{}, 3)
	sweetFA[0] = aa
	sweetFA[1] = bb
	sweetFA[2] = cc

	f := sweetFA[1].(func(string) string)
	foo := f("foo")
	fmt.Printf("the middle function yields '%s'\n", foo)

	c.Assert(sweetFA[0].(func(string) string)("foo"), Equals, "fooAA")
	c.Assert(sweetFA[1].(func(string) string)("foo"), Equals, "fooBB")
	c.Assert(sweetFA[2].(func(string) string)("foo"), Equals, "fooCC")
}
