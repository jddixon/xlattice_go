package main

// xlattice_go/cmd/xlReg/xlReg.co

import (
	"encoding/hex"
	"flag"
	"fmt"
	xi "github.com/jddixon/xlattice_go/nodeID"
	"github.com/jddixon/xlattice_go/reg"
	xt "github.com/jddixon/xlattice_go/transport"
	xu "github.com/jddixon/xlattice_go/util"
	"os"
	"path"
)

func Usage() {
	fmt.Printf("Usage: %s [OPTIONS]\n", os.Args[0])
	fmt.Printf("where the options are:\n")
	flag.PrintDefaults()
}

const (
	DEFAULT_ADDR = "127.0.0.1"
	DEFAULT_NAME = "xlReg"
	DEFAULT_LFS  = "/var/app/xlReg"
	DEFAULT_PORT = 44444 // for the registry, not clients
)

var (
	// these need to be referenced as pointers
	address  = flag.String("a", DEFAULT_ADDR, "registry IP address")
	hexID    = flag.String("i", "", "hex reg ID")
	justShow = flag.Bool("j", false, "display option settings and exit")
	lfs      = flag.String("lfs", DEFAULT_LFS, "path to work directory")
	name     = flag.String("n", DEFAULT_NAME, "registry name")
	port     = flag.Int("p", DEFAULT_PORT, "registry listening port")
	testing  = flag.Bool("T", false, "test run")
	verbose  = flag.Bool("v", false, "be talkative")
)

func init() {
	fmt.Println("init() invocation") // DEBUG
}

func main() {
	var err error
	var id *xi.NodeID

	flag.Usage = Usage
	flag.Parse()

	// FIXUPS ///////////////////////////////////////////////////////

	if *hexID == "" {
		id, _ = xi.New(nil)
		*hexID = hex.EncodeToString(id.Value())
	} else {
		var data []byte
		data, err = hex.DecodeString(*hexID)
		if err == nil {
			id, err = xi.New(data)
		}
	}
	if err != nil {
		fmt.Println("error processing NodeID: %s\n", err.Error())
		os.Exit(-1)
	}
	if *testing {
		if *name == DEFAULT_NAME || *name == "" {
			*name = "testReg"
		}
		if *lfs == DEFAULT_LFS || *lfs == "" {
			*lfs = "./myApp/xlReg"
		} else {
			*lfs = path.Join("tmp", *lfs)
		}
		if *port == DEFAULT_PORT || *port == 0 {
			*port = 33333
		}
	}
	addrAndPort := fmt.Sprintf("%s:%d", *address, *port)
	endPoint, err := xt.NewTcpEndPoint(addrAndPort)
	if err != nil {
		fmt.Printf("not a valid endPoint: %s\n", addrAndPort)
		// XXX STUB XXX
	}

	// SANITY CHECKS ////////////////////////////////////////////////
	if err == nil {
		err = xu.CheckLFS(*lfs)		// tries to create if it doesn't exist
	}
	// DISPLAY STUFF ////////////////////////////////////////////////
	if *verbose {
		fmt.Printf("xlReg v%s %s\n", reg.VERSION, reg.VERSION_DATE)
	}
	if *verbose || *justShow {
		fmt.Printf("address      = %v\n", *address)
		fmt.Printf("endPoint     = %v\n", endPoint)
		fmt.Printf("hexID        = %v\n", *hexID)
		fmt.Printf("justShow     = %v\n", *justShow)
		fmt.Printf("lfs          = %s\n", *lfs)
		fmt.Printf("name         = %s\n", *name)
		fmt.Printf("port         = %d\n", *port)
		fmt.Printf("testing      = %v\n", *testing)
		fmt.Printf("verbose      = %v\n", *verbose)
	}
	if *justShow {
		return
	}
	// SET UP OPTIONS ///////////////////////////////////////////////
	var opt reg.RegOptions
	opt.ID = id
	opt.Lfs = *lfs
	opt.Port = *port
	opt.Testing = *testing
	opt.Verbose = *verbose

	r, err := setup(&opt)
	if err == nil {
		err = serve(r)
	}
	_ = err
}
func setup(opt *reg.RegOptions) (r *reg.Registry, err error) {
	// If LFS/.xlattice/config exists, we load that.  Otherwise we
	// create a node.  In either case we force the node to listen on
	// the designated port

	// XXX STUB XXX

	r, err = reg.NewRegistry(nil,		// nil = clusters so far
		opt.Name, opt.ID, opt.Lfs,
		nil, nil,						// opt.CKey, opt.SKey,
		nil,							// overlays
		opt.EndPoint)

	return r, err
}
func serve(r *reg.Registry) (err error) {

	// XXX STUB XXX

	// ORDERLY SHUTDOWN /////////////////////////////////////////////

	return
}
