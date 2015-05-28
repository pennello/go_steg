// chris 052515

package steg

import "io"

func (c chunk) write(b byte) {
	cur := c.read()
	// Compare current value with what we need to write.
	x := cur ^ b
	// XXX Don't do anything if x == 0?
	// The 5 high bits are the offset to the byte containing the bit
	// to flip.
	byteIndex := (x >> 3) & 0x1f
	// The 3 low bits are the index of the bit to flip.
	mask := byte(1) << (x & 0x7)
	// Flip the bit.
	c[byteIndex] ^= mask
	// Done!
}

type Writer struct {
	dst     io.Writer
	carrier io.Reader
}

func NewWriter(dst io.Writer, carrier io.Reader) Writer {
	return Writer{dst: dst, carrier: carrier}
}

// Write chunk into destination io.Reader.
func (w Writer) write(c chunk) error {
	_, err := w.dst.Write([]byte(c))
	return err
}

func (w Writer) Write(p []byte) (int, error) {
	c := chunk(make([]byte, chunkSize))
	n := 0 // Total bytes written.
	for _, b := range p {
		var err error
		err = readChunk(c, w.carrier)
		if err != nil {
			return n, err
		}
		c.write(b)
		err = w.write(c)
		if err != nil {
			// We may have written _some_ of the bytes of c,
			// but won't have written all of them.  We
			// consider this to be b _not_ having been
			// written, so n remains unincremented.
			return n, err
		}
		n++
	}
	return n, nil
}

func (w Writer) Copy() (written int64, err error) {
	return io.Copy(w.dst, w.carrier)
}
