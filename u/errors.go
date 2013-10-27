package u

import "errors"

var (
	DirStrucNotRecognized = errors.New("DirStruc not recognized")
	EmptyKey              = errors.New("empty key parameter")
	FileNotFound          = errors.New("file not found")
)
