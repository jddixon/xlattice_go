package xlattice_go

// xlattice_go/entityName.go

import (
	"errors"
	"regexp"
)

const NAME_STARTERS = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz_"
const NAME_CHARS = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789_"

var namePat = "^[" + NAME_STARTERS + "]" + "[" + NAME_CHARS + "]*$"
var nameRE = regexp.MustCompile(namePat)

// go won't let me be a constant :-(
var invalidName = errors.New("not a valid xlattice entity name")

func validEntityName(name string) (err error) {
	if !nameRE.MatchString(name) {
		err = invalidName
	}
	return
}

func NAME_PAT() string        { return namePat }
func NAME_RE() *regexp.Regexp { return nameRE }
func INVALID_NAME() error     { return invalidName }
