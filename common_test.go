// chris 052515

package steg

import "log"

const helloString = "hello, there, how are you? fine."

// The byte embedded in the above string.
const helloByte = 0xdb

func testHelloChunk() chunk {
	b := []byte(helloString)
	if len(b) != chunkSize {
		log.Panicf("test hello chunk wrong size (%v != %v)", len(b), chunkSize)
	}
	return chunk(b)
}
