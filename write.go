// chris 052515

package steg

import (
	"errors"
	"io"
)

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

var ErrShortRead = errors.New("insufficient carrier data to write requested data")

type Writer struct {
	dst     io.Writer
	carrier io.Reader
}

func NewWriter(dst io.Writer, carrier io.Reader) Writer {
	return Writer{dst: dst, carrier: carrier}
}

// Read a chunk from the carrier.  If there is an error reading from the
// carrier, even after completely reading the chunk from the carrier, that
// error is returned.
func (w Writer) read(c chunk) error {
	// We'll use this as a byte slice here internally.
	p := []byte(c)
	t := 0 // Total number of bytes read.
	for {
		n, err := w.carrier.Read(p[t:])
		t += n
		if t == len(p) {
			// We're done reading from the carrier.  But do
			// check for an error...
			if err != nil && err != io.EOF {
				return err
			}
			break
		}
		if err != nil {
			if err == io.EOF {
				return ErrShortRead
			} else {
				return err
			}
		}
	}
	return nil
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
		err = w.read(c)
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
