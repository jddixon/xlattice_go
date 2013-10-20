package filters

import "errors"

var (
	KeySelectorArgOutOfRange = errors.New("KeySelector arg out of range")
	KeyTooShort              = errors.New("key too short")
	MappingFileTooSmall		 = errors.New("mapping file too small")
	MOutOfRange              = errors.New("m out of range")
	NilKey                   = errors.New("nil key parameter")
	TooManyHashFunctions     = errors.New("too many hash functions")
)
