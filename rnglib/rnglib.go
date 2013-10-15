package rnglib

import (
	"strings"
)

// Version number tracked in ../CHANGES
func Version() (string, string) {
	return "0.3.1", "2013-08-24"
}

// a crude attempt at properties
var _FILE_NAME_STARTERS = strings.Split(
	"ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789_", "")

func FILE_NAME_STARTERS() []string {
	return _FILE_NAME_STARTERS
}

var _FILE_NAME_CHARS = strings.Split(
	"ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789_-.", "")

func FILE_NAME_CHARS() []string {
	return _FILE_NAME_CHARS
}

