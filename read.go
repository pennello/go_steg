// chris 052515

package steg

import (
	"errors"
	"io"
	"io/ioutil"

	"chrispennello.com/go/swar"
)

// ErrShortRead will be returned from Read and Reader.Read when an EOF
// is encountered before being able to read sufficient data.
var ErrShortRead = errors.New("short read")

func (a *atom) asUint() uint {
	r := uint(0)
	for i := uint(0); i < a.ctx.atomSize; i++ {
		r |= uint(a.data[i]) << (i * 8)
	}
	return r
}

func (c *chunk) readBitMask(a *atom, bitIndex uint, mask byte, B byte) {
	// First, extract the desired bits from the chunk byte by using
	// the mask.
	x := mask & B
	// We then want to XOR together the bits specified by the mask.
	// The key is to recognize that this is the same as taking the
	// 8-bit population count (ones count, or Hamming weight) and
	// then examining the parity.  If even, then 0; if odd, then 1.
	x = swar.Ones8(x) % 2
	// XOR the bit into the atom.
	a.xorBit(x, bitIndex)
	// Done!
}

func (c *chunk) readBit(a *atom, bitIndex uint, cBi uint, B byte) {
	// c.ctx.atomSize won't be bigger than 3, so bitIndex will be no
	// larger than 23, so power or thresh won't overflow an int32.
	thresh := int32(1) << (bitIndex - 3)
	power := thresh << 1
	value := int32(cBi) % power
	mask := byte(swar.IntegerSelect32(value, thresh, 0x00, 0xff))
	c.readBitMask(a, bitIndex, mask, B)
}

func (c *chunk) readAtom() *atom {
	a := c.ctx.newAtom()
	bits := c.ctx.atomSize * 8
	// cBi: chunk byte index
	for cBi := uint(0); cBi < c.ctx.chunkSize; cBi++ {
		B := c.data[cBi]

		c.readBitMask(a, 0, 0xaa, B)
		c.readBitMask(a, 1, 0xcc, B)
		c.readBitMask(a, 2, 0xf0, B)

		for bitIndex := uint(3); bitIndex < bits; bitIndex++ {
			c.readBit(a, bitIndex, cBi, B)
		}
	}
	return a
}

// Read into a chunk from an io.Reader.  If there is an error reading,
// even after completely reading the chunk, that error is returned.
// Returns an error iff the chunk was not completely read.
func readChunk(c *chunk, r io.Reader) error {
	_, err := io.ReadFull(r, c.data)
	if err == io.EOF || err == io.ErrUnexpectedEOF {
		err = ErrShortRead
	}
	return err
}

// Read steganographically-embedded bytes from the underlying source
// io.Reader.  Returns the number of bytes read as well as an error, if
// one occurred.
//
// Can return io.EOF or io.ErrUnexpectedEOF if an EOF was encountered
// before being able to read a sufficient number of bytes to extract the
// requested amount of data.
//
// n == len(p) iff err != nil
func (r *Reader) Read(p []byte) (n int, err error) {
	c := r.ctx.newChunk()
	for n < len(p) {
		if r.cur == nil {
			err := readChunk(c, r.src)
			if err != nil {
				return n, err
			}
			r.cur = c.readAtom()
			r.cn = c.ctx.atomSize
		}
		nn := copy(p[n:], r.cur.data[c.ctx.atomSize-r.cn:])
		n += nn
		r.cn -= uint(nn)
		if r.cn == 0 {
			r.cur = nil
		}
	}
	return n, err
}

// Discard reads n bytes into ioutil.Discard, throwing them away.
//
// The idea is that you'd call this to jump ahead by some offset in the
// carrier data before you start reading your
// steganographically-embedded data.
//
// Counterpart to Writer.CopyN and Mux.CopyN.
func (r *Reader) Discard(n int64) error {
	_, err := io.CopyN(ioutil.Discard, r.src, n)
	return err
}
