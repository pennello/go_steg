// chris 052515

package steg

import (
	"errors"
	"io"

	"chrispennello.com/go/swar"
)

// ErrShortRead will be returned from Read and Reader.Read when an EOF
// is encountered before being able to read sufficient data.
var ErrShortRead = errors.New("short read")

// Read a single bit with index i from the chunk c.  If you iterate over
// i from 0 to 7, you'll get the bits you need to reconstruct a whole
// byte from a chunk.
func (c chunk) readBit(i bitIndex) byte {
	// Byte we'll return.  Will have the output bit set at index i.
	ret := byte(0)
	// Byte masks specific to this bit index.
	masks := masksByIndex[i]
	// Iterate through the bytes in the chunk.
	for j := 0; j < chunkSize; j++ {
		// First, extract the desired bits from the chunk byte
		// by using the mask.  We then want to XOR together the
		// bits specified by the mask.  The key is to recognize
		// that this is the same as taking the 8-bit population
		// count (ones count, or Hamming weight) of b, and then
		// examining the parity.  If even, then 0; if odd, then
		// 1.
		r := swar.Ones8(masks[j] & c[j])
		r = (r % 2) << i
		// This bit result is only for this chunk byte.  XOR the
		// bit into the return byte.
		ret ^= r
	}
	return ret
}

// Read a byte from a chunk c.
func (c chunk) read() byte {
	ret := byte(0)
	for i := bitIndex(0); i < 8; i++ {
		ret |= c.readBit(i)
	}
	return ret
}

// Read a chunk from an io.Reader.  If there is an error reading, even
// after completely reading the chunk, that error is returned.  Sort of
// similar to io.Reader.Read, returns a boolean complete--whether we
// completely read the chunk.  Returns an error iff no data was read.
func readChunk(c chunk, r io.Reader) (err error) {
	_, err = io.ReadFull(r, []byte(c))
	if err == io.EOF || err == io.ErrUnexpectedEOF {
		err = ErrShortRead
	}
	return err
}

// A Reader wraps an io.Reader and reads steganographically-embedded
// bytes from it.  Implements io.Reader.
type Reader struct {
	src io.Reader
}

// NewReader returns a fresh Reader, ready to read
// steganographically-embedded bytes from the source io.Reader.
func NewReader(src io.Reader) Reader {
	return Reader{src: src}
}

// Read steganographically-embedded bytes from the underlying source
// io.Reader.  Returns the number of bytes read as well as an error, if
// one occurred.
//
// Can return io.EOF or io.ErrUnexpectedEOF if an EOF was encountered
// before being able to read a sufficient number of bytes to extract the
// requested data.
//
// n == len(p) iff err != nil
func (r Reader) Read(p []byte) (n int, err error) {
	c := newChunk()
	for ; n < len(p); n++ {
		err = readChunk(c, r.src)
		if err != nil {
			return n, err
		}
		p[n] = c.read()
	}
	return n, err
}
