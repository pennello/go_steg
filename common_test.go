// chris 052515 Common code for test routines.

package steg

import (
	"os"
	"strings"
	"testing"
	"time"

	"math/rand"
)

const helloString = "hello, there, how are you? fine."

// The byte embedded in the above string.
const helloByte = 0xdb

func testHelloChunk(atomSize uint) *chunk {
	ctx := NewCtx(atomSize)
	c := ctx.newChunk()
	n := len(c.data) / len(helloString)
	if n*len(helloString) != len(c.data) {
		panic("non-integral chunk multiple")
	}
	copy(c.data, []byte(strings.Repeat(helloString, n)))
	return c
}

func TestMain(m *testing.M) {
	rand.Seed(time.Now().UTC().UnixNano())
	os.Exit(m.Run())
}
