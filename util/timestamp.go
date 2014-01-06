package util

// xlattice_go/util/timestamp.go

import (
	"time"
)

/**
 * Convenience class handling YYYY-MM-DD HH:MM:SS formatted dates.
 */
type Timestamp int64

const (
	layout = "2006-01-02 15:04:05" // construed as UTC
)

func (t Timestamp) String() (x string) {
	utc := time.Unix(0, int64(t)) // a time value
	x = utc.UTC().Format(layout)
	return
}
