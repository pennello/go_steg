// chris 052515

package steg

// Fixed chunk size for reading/writing bytes.  For every chunk of carrier
// data, we can embed 1 byte of message data.
const chunkSize = 32

// These should range from 0 to 7.
type bitIndex uint

// We _could_ define a chunk as a [32]byte, but then we'd just end up
// implementing slices all over again...
type chunk []byte

// The masks for the chunk.readBit implementation end up being the same as a
// chunk.
type byteMasks chunk

func newChunk() chunk {
	return chunk(make([]byte, chunkSize))
}
