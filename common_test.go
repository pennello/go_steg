// chris 052515

package steg

import (
	"log"
	"testing"
)

func testHelloChunk() chunk {
	b := []byte("hello, there, how are you? fine.")
	if len(b) != chunkSize {
		log.Panicf("test hello chunk wrong size (%v != %v)", len(b), chunkSize)
	}
	return chunk(b)
}

func testReadByte(t *testing.T, c chunk) byte {
	r := c.read()
	//t.Logf("chunk: %#v; read: %#v", string(c), string(r))
	return r
}
