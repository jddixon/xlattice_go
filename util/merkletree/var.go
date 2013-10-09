package merkletree

import (
	re "regexp"
)

var (
	// Notice the terminating forward slash and lack of newlines or CR-LF
	// THIS PATTERN WON"T CATCH SOME ERRORS; eg it permits '///' in paths
	// (?i:RE) is Go's way of saying "ignore lower case".

	// FOUND IN the Python equivalent of merkle_tree.go.  This pattern
	// accepts leading spaces in the first line.

	FIRST_LINE_RE_1 = re.MustCompile("^(?i:( *)([0-9a-f]{40}) ([a-z0-9_\\-\\.]+/))$")
	OTHER_LINE_RE_1 = re.MustCompile("^(?i:([ XYZ]*)([0-9a-f]{40}) ([a-z0-9_\\-\\.]+/?))$")
	FIRST_LINE_RE_3 = re.MustCompile("^(?i:( *)([0-9a-f]{64}) ([a-z0-9_\\-\\.]+/))$")
	OTHER_LINE_RE_3 = re.MustCompile("^(?i:([ XYZ]*)([0-9a-f]{64}) ([a-z0-9_\\-\\.]+/?))$")

	// FOUND IN the Python equivalent of merkle_doc.go.  No leading spaces
	// in the first line.
	FIRST_LINE_RE_1d = re.MustCompile("^(?i:([0-9a-f]{40}) ([a-z0-9_\\-\\./]+/))$")
	FIRST_LINE_RE_3d = re.MustCompile("^(?i:([0-9a-f]{64}) ([a-z0-9_\\-\\./]+/))$")
)
