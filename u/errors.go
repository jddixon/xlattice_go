package u

import "errors"

var (
	DirStrucNotRecognized = errors.New("DirStruc not recognized")
	FileNotFound          = errors.New("file not found")
)
