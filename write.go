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
	w   io.Writer
	carrier io.Reader
}

func NewWriter(w io.Writer, carrier io.Reader) Writer {
	return Writer{w: w, carrier: carrier}
}

func (w Writer) write(buf []byte, b byte) error {
	nt := 0
	for {
		n, err := w.carrier.Read(buf[nt:])
		nt += n
		if nt == len(buf) {
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

func (w Writer) Write(p []byte) (int, error) {
	buf := make([]byte, chunkSize)
	n := 0
	for _, b := range p {
		err := w.write(buf, b)
		if err != nil {
			return n, err
		}
		n++
	}
	return n, nil
}
