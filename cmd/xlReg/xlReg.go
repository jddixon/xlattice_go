package main

import (
	"flag"
	"fmt"
	"github.com/jddixon/xlattice_go/reg"
	"os"
	"path"
)

func Usage() {
	fmt.Printf("Usage: %s [OPTIONS]\n", os.Args[0])
	fmt.Printf("where the options are:\n")
	flag.PrintDefaults()
}

const (
	DEFAULT_LFS  = "/var/XLReg"
	DEFAULT_PORT = 55555
)

var (
	// these need to be referenced as pointers
	justShow = flag.Bool("j", false, "display option settings and exit")
	lfs      = flag.String("lfs", DEFAULT_LFS, "path to work directory")
	port     = flag.Int("p", DEFAULT_PORT, "listening port")
	testing  = flag.Bool("T", false, "test run")
	verbose  = flag.Bool("v", false, "be talkative")
)

func init() {
	fmt.Println("init() invocation") // DEBUG
}

func main() {
	flag.Usage = Usage
	flag.Parse()

	// FIXUPS ///////////////////////////////////////////////////////
	if *testing {
		if *lfs == DEFAULT_LFS || *lfs == "" {
			*lfs = "./tmp"
		} else {
			*lfs = path.Join("tmp", *lfs)
		}
	}
	// SANITY CHECKS ////////////////////////////////////////////////

	// DISPLAY FLAGS ////////////////////////////////////////////////
	if *verbose || *justShow {
		fmt.Printf("justShow               = %v\n", *justShow)
		fmt.Printf("lfs                    = %s\n", *lfs)
		fmt.Printf("port                   = %d\n", *port)
		fmt.Printf("testing                = %v\n", *testing)
		fmt.Printf("verbose                = %v\n", *verbose)
	}
	if *justShow {
		return
	}
	// SET UP OPTIONS ///////////////////////////////////////////////
	var opt reg.RegOptions
	opt.Lfs = *lfs
	opt.Port = *port
	opt.Testing = *testing
	opt.Verbose = *verbose

	// LOAD XLATTICE NODE CONFIGURATION /////////////////////////////
	// If LFS/.xlattice/config exists, we load that.  Otherwise we
	// create a node.  In either case we force the node to listen on
	// the designated port

	// XXX STUB XXX

	// LISTEN AND SERVE /////////////////////////////////////////////

	// ORDERLY SHUTDOWN /////////////////////////////////////////////

}
