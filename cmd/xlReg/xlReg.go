package main

// xlattice_go/cmd/xlReg/xlReg.co

import (
	"crypto/rand"
	"crypto/rsa"
	"flag"
	"fmt"
	xn "github.com/jddixon/xlattice_go/node"
	xi "github.com/jddixon/xlattice_go/nodeID"
	"github.com/jddixon/xlattice_go/reg"
	xt "github.com/jddixon/xlattice_go/transport"
	xf "github.com/jddixon/xlattice_go/util/lfs"
	"io/ioutil"
	"log"
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
	address = flag.String("a", DEFAULT_ADDR,
		"registry IP address")
	clearFilter = flag.Bool("c", false,
		"clear Bloom filer at beginning of run")
	ephemeral = flag.Bool("e", false,
		"server is ephemeral, does not persist data")
	justShow = flag.Bool("j", false,
		"display option settings and exit")
	k = flag.Int("k", int(reg.DEFAULT_K),
		"number of hash functions in Bloom filter")
	lfs = flag.String("lfs", DEFAULT_LFS,
		"path to work directory")
	logFile = flag.String("l", "log",
		"path to log file")
	m = flag.Int("m", int(reg.DEFAULT_M),
		"exponent in Bloom filter")
	name = flag.String("n", DEFAULT_NAME,
		"registry name")
	port = flag.Int("p", DEFAULT_PORT,
		"registry listening port")
	testing = flag.Bool("T", false,
		"this is a test run")
	verbose = flag.Bool("v", false,
		"be talkative")
)

func init() {
	fmt.Println("init() invocation") // DEBUG
}

func main() {
	var err error

	flag.Usage = Usage
	flag.Parse()

	// FIXUPS ///////////////////////////////////////////////////////

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
	var backingFile string
	if !*ephemeral {
		backingFile = path.Join(*lfs, "idFilter.dat")
	}
	endPoint, err := xt.NewTcpEndPoint(addrAndPort)
	if err != nil {
		fmt.Printf("not a valid endPoint: %s\n", addrAndPort)
		// XXX STUB XXX
	}
	// SANITY CHECKS ////////////////////////////////////////////////
	if err == nil {
		if *m < 2 {
			*m = 20
		}
		if *k < 2 {
			*k = 8
		}
		err = xf.CheckLFS(*lfs) // tries to create if it doesn't exist
		if err == nil {
			if *logFile != "" {
				*logFile = path.Join(*lfs, *logFile)
			}
		}
	}
	// DISPLAY STUFF ////////////////////////////////////////////////
	if *verbose || *justShow {
		fmt.Printf("address      = %v\n", *address)
		fmt.Printf("backingFile  = %v\n", backingFile)
		fmt.Printf("clearFilter  = %v\n", *clearFilter)
		fmt.Printf("endPoint     = %v\n", endPoint)
		fmt.Printf("ephemeral    = %v\n", *ephemeral)
		fmt.Printf("justShow     = %v\n", *justShow)
		fmt.Printf("k            = %d\n", *k)
		fmt.Printf("lfs          = %s\n", *lfs)
		fmt.Printf("logFile      = %s\n", *logFile)
		fmt.Printf("m            = %d\n", *m)
		fmt.Printf("name         = %s\n", *name)
		fmt.Printf("port         = %d\n", *port)
		fmt.Printf("testing      = %v\n", *testing)
		fmt.Printf("verbose      = %v\n", *verbose)
	}
	if *justShow {
		return
	}
	// SET UP OPTIONS ///////////////////////////////////////////////
	var (
		f      *os.File
		logger *log.Logger
		opt    reg.RegOptions
		rs     *reg.RegServer
	)
	if *logFile != "" {
		f, err = os.OpenFile(*logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0660)
		if err == nil {
			logger = log.New(f, "", log.Ldate|log.Ltime)
		}
	}
	if f != nil {
		defer f.Close()
	}
	if err == nil {
		opt.Address = *address
		opt.BackingFile = backingFile
		opt.ClearFilter = *clearFilter
		opt.Ephemeral = *ephemeral
		opt.K = uint(*k)
		opt.Lfs = *lfs
		opt.Logger = logger
		opt.M = uint(*m)
		opt.Lfs = *lfs
		opt.Port = fmt.Sprintf("%d", *port)
		opt.Testing = *testing
		opt.Verbose = *verbose

		rs, err = setup(&opt)
		if err == nil {
			err = serve(rs)
		}
	}
	_ = logger // NOT YET
	_ = err
}
func setup(opt *reg.RegOptions) (rs *reg.RegServer, err error) {
	// If LFS/.xlattice/node.config exists, we load that.  Otherwise we
	// create a node.  In either case we force the node to listen on
	// the designated port

	var (
		e                []xt.EndPointI
		pathToConfigFile string
		node             *xn.Node
		ckPriv, skPriv   *rsa.PrivateKey
	)

	greetings := fmt.Sprintf("xlReg v%s %s start run\n",
		reg.VERSION, reg.VERSION_DATE)
	// fmt.Print(greetings)
	opt.Logger.Print(greetings)

	pathToConfigFile = path.Join(path.Join(opt.Lfs, ".xlattice"), "node.config")
	found, err := xf.PathExists(pathToConfigFile)
	if err == nil {
		if found {
			// The registry node already exists.  Parse it and we are done.
			var data []byte
			data, err = ioutil.ReadFile(pathToConfigFile)
			if err == nil {
				node, _, err = xn.Parse(string(data))
			}
		} else {
			// We need to create a registry node from scratch.
			nodeID, _ := xi.New(nil)
			ep, err := xt.NewTcpEndPoint(opt.Address + ":" + opt.Port)
			if err == nil {
				e = []xt.EndPointI{ep}
				ckPriv, err = rsa.GenerateKey(rand.Reader, 2048)
				if err == nil {
					skPriv, err = rsa.GenerateKey(rand.Reader, 2048)
				}
			}
			if err == nil {
				node, err = xn.New("xlReg", nodeID, opt.Lfs, ckPriv, skPriv,
					nil, e, nil)
			}
			if err == nil {
				err = xf.MkdirsToFile(pathToConfigFile, 0700)
				if err == nil {
					err = ioutil.WriteFile(pathToConfigFile,
						[]byte(node.String()), 0400)
				}
			}
		}
	}
	if err == nil {
		var r *reg.Registry
		r, err = reg.NewRegistry(nil, // nil = clusters so far
			node, ckPriv, skPriv, opt)
		if err == nil {
			// DEBUG
			fmt.Printf("Registry name: %s\n", node.GetName())
			fmt.Printf("         ID:   %s\n", node.GetNodeID().String())
			// END
		}
		if err == nil {
			var verbosity int
			if opt.Verbose {
				verbosity++
			}
			rs, err = reg.NewRegServer(r, opt.Testing, verbosity)
		}
	}
	if err != nil {
		fmt.Printf("ERROR: %s\n", err.Error())
	}
	return
}
func serve(rs *reg.RegServer) (err error) {

	err = rs.Run() // non-blocking
	if err == nil {
		<-rs.DoneCh
	}

	// XXX STUB XXX

	// ORDERLY SHUTDOWN /////////////////////////////////////////////

	return
}
