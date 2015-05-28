// chris 052515

package steg

import (
	"errors"
	"io"

	"chrispennello.com/go/swar"
)

var ErrShortRead = errors.New("insufficient carrier data to write requested data")

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

func (c chunk) read() byte {
	ret := byte(0)
	for i := bitIndex(0); i < 8; i++ {
		ret |= c.readBit(i)
	}
	return ret
}

// Read a chunk from an io.Reader.  If there is an error reading, even
// after completely reading the chunk, that error is returned.
func readChunk(c chunk, r io.Reader) (complete bool, err error) {
	// We'll use this as a byte slice here internally.
	p := []byte(c)
	t := 0 // Total number of bytes read.
	for {
		n, err := r.Read(p[t:])
		t += n
		if t == len(p) {
			// We're done reading.  But do check for an
			// error...
			if err != nil {
				return true, err
			}
			break
		}
		if err != nil {
			if err == io.EOF {
				err = ErrShortRead
			}
			return false, err
		}
	}
	return true, nil
}

type Reader struct {
	src io.Reader
}

func NewReader(src io.Reader) Reader {
	return Reader{src: src}
}

func (r Reader) Read(p []byte) (int, error) {
	c := newChunk()
	n := 0
	var complete bool
	var err error
	for ; n < len(p); n++ {
		complete, err = readChunk(c, r.src)
		if !complete {
			return n, err
		}
		p[n] = c.read()
	}
	return n, err
}
