package rnglib

import "fmt"
import "github.com/bmizerany/assert"
import "os"
import "strings"
import "testing"
import "time"

const TMP_DIR = "tmp"

func buildData(count uint32) *[]byte {
	p := make([]byte, count)
	return &p
}
func makeRNG() *SimpleRNG {
	t := time.Now().Unix() // int64 sec
	rng := NewSimpleRNG(t)
	return rng
}
func TestConstuctor(t *testing.T) {
	rng := makeRNG()
	assert.NotEqual(t, rng, nil)
}
func TestNextBoolean(t *testing.T) {
	rng := makeRNG()
	val := rng.NextBoolean()
	assert.NotEqual(t, val, nil)

	valAsIface := interface{}(val)
	switch v := valAsIface.(type) {
	default:
		fmt.Printf("expected type bool, found %T", v)
		// assert.Fail("whatever NextBoolean() returns is not a bool")
	case bool:
		/* empty statement */
	}
}
func TestNextByte(t *testing.T) {
	// rng := makeRNG()
}
func TestNextBytes(t *testing.T) {
	rng := makeRNG()
	count := uint32(1)          // minimum length of buffer
	count += rng.NextInt32(256) // maximum
	data := buildData(count)    // so 1 .. 256 bytes
	rng.NextBytes(data)
	actualLen := uint32(len(*data))
	assert.NotEqual(t, 0, actualLen)
	assert.Equal(t, actualLen, count)

}
func TestNextFileName(t *testing.T) {
	rng := makeRNG()
	maxLen := uint32(1)         // minimum length of name
	maxLen += rng.NextInt32(16) // maximum
	name := rng.NextFileName(int(maxLen))
	// DEBUG
	fmt.Printf("next file name is %s\n", name)
	// END
	actualLen := len(name)
	assert.NotEqual(t, 0, actualLen)
	// assert.True( t, actualLen < maxLen)
}
func TestNextDataFile(t *testing.T) {
	rng := makeRNG()
	minLen := int(rng.NextInt32(4))            // minimum length of file
	maxLen := minLen + int(rng.NextInt32(256)) // maximum

	// XXX should return err, which should be nil
	fileLen, pathToFile := rng.NextDataFile(TMP_DIR, maxLen, minLen)
	// DEBUG
	fmt.Printf("data file is %s; size is %d\n", pathToFile, fileLen)
	// END

	stats, err := os.Stat(pathToFile)
	assert.Equal(t, nil, err)
	fileName := stats.Name()
	assert.Equal(t, TMP_DIR+"/"+fileName, pathToFile)
	assert.Equal(t, stats.Size(), int64(fileLen))

}
func doNextDataDirTest(t *testing.T, rng *SimpleRNG, width int, depth int) {
	dirName := rng.NextFileName(8)
	dirPath := TMP_DIR + "/" + dirName
	pathExists, err := PathExists(dirPath)
	if err != nil {
		panic("error invoking PathExists on " + dirPath)
	}
	if pathExists {
		if strings.HasPrefix(dirPath, "/") {
			panic("attempt to remove absolute path " + dirPath)
		}
		if strings.Contains(dirPath, "..") {
			panic("attempt to remove path containing ..: " + dirPath)
		}
		os.RemoveAll(dirPath)
	}
	rng.NextDataDir(dirPath, width, depth, 32, 0)
}
func TestNextDataDir(t *testing.T) {
	rng := makeRNG()
	doNextDataDirTest(t, rng, 1, 1)
	doNextDataDirTest(t, rng, 1, 4)
	doNextDataDirTest(t, rng, 4, 1)
	doNextDataDirTest(t, rng, 4, 4)
}
