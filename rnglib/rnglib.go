package rnglib

import (
	"strings"
)

func Version() (string, string) {
	return VERSION, VERSION_DATE
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
