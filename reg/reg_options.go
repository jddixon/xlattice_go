package reg

// xlattice_go/reg/reg_options.go

import (
	xt "github.com/jddixon/xlattice_go/transport"
	"log"
)

// Options normally set from the command line or derived from those.
// Not used in this package but used by xlReg
type RegOptions struct {
	Address  string
	EndPoint xt.EndPointI // derived from Address, Port
	K        uint
	Lfs      string
	Logger   *log.Logger
	M        uint
	Name     string
	Port     string
	Testing  bool
	Verbose  bool
}
