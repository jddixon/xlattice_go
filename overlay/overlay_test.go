package overlay

// xlattice_go/overlay/overlay_test.go

import (
	"github.com/bmizerany/assert"
	x "github.com/jddixon/xlattice_go"
	"testing"
)

func TestCtor(t *testing.T) {
	rng := x.MakeRNG()
	name := rng.NextFileName(8)

	o, err := NewOverlay(name, nil, "tcpip", 0.42)
	assert.Equal(t, nil, err)
	assert.NotEqual(t, nil, o)
}
