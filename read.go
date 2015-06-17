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

func (a *atom) asUint32() uint32 {
	r := uint32(0)
	for i := uint8(0); i < uint8(a.ctx.atomSize); i++ {
		r |= uint32(a.data[i]) << (i * 8)
	}
	return r
}

// readBitMask xors the bits specified out of the byte B by applying the
// mask and then stores the value in the atom a at the specified atom
// bit index abi.
func (c *chunk) readBitMask(a *atom, abi uint8, mask byte, B byte) {
	// First, extract the desired bits from the chunk byte by using
	// the mask.
	x := mask & B
	// We then want to XOR together the bits specified by the mask.
	// The key is to recognize that this is the same as taking the
	// 8-bit population count (ones count, or Hamming weight) and
	// then examining the parity.  If even, then 0; if odd, then 1.
	x = swar.Ones8(x) % 2
	// XOR the bit into the atom.
	a.xorBit(x, abi)
	// Done!
}

// readBit xors the bits specified out of the byte B by computing a mask
// and invoking readBitMask.
func (c *chunk) readBit(a *atom, abi uint8, cBi uint32, B byte) {
	// We compute the power of two threshold for this bit index.  If
	// the byte index is above this value (mod power), then we will
	// include these bytes in the xor for this bit.  If it's below,
	// then we will not.  This is expressed as using a mask of
	// either 0xff (for inclusion) or 0x00 (for exclusion).  This
	// threshold exponentially increases with each atom bit index
	// abi until ultimately we are excluding the entirety of the
	// bottom half of the bytes and including the entirety of the
	// top half.  See the table in the notes.

	// c.ctx.atomSize won't be bigger than 3, so abi will be no
	// larger than 23, so power or thresh won't overflow an int32.
	thresh := int32(1) << (uint8(abi) - 3)
	power := thresh << 1
	// c.ctx.atomSize won't be bigger than 3, so cBi will be less
	// than 2Mi., so cBi won't overflow an int32.
	value := int32(cBi) % power
	mask := byte(swar.IntegerSelect32(value, thresh, 0x00, 0xff))
	c.readBitMask(a, abi, mask, B)
}

// readAtom creates a new atom and reads its contents out of the chunk.
func (c *chunk) readAtom() *atom {
	a := c.ctx.newAtom()
	atomBits := c.ctx.atomSize * 8
	// cBi: chunk byte index
	for cBi := uint32(0); cBi < c.ctx.chunkSize; cBi++ {
		B := c.data[cBi]

		// Bit indexes 0, 1, and 2 are special.  This is because
		// for these index values, the bits to be selected from
		// the chunk data are all within a byte.  Above these
		// bit indexes, we simply select particular *bytes* _en
		// masse_, but below them, we must select subsets of
		// bits within the bytes.
		c.readBitMask(a, 0, 0xaa, B)
		c.readBitMask(a, 1, 0xcc, B)
		c.readBitMask(a, 2, 0xf0, B)

		// abi: atom bit index
		for abi := uint8(3); abi < atomBits; abi++ {
			c.readBit(a, abi, cBi, B)
		}
	}
	return a
}

// Read reads steganographically-embedded bytes from the underlying
// source io.Reader.  Returns the number of bytes read as well as an
// error, if one occurred.
//
// Can return io.EOF or io.ErrUnexpectedEOF if an EOF was encountered
// before being able to read a sufficient number of bytes to extract the
// requested amount of data.
//
// n == len(p) iff err != nil
//
// The current implementation is somewhat naive.  Each chunk is read
// completely into memory from the underlying source reader.  In
// particular, for atom size 3, this means that 2MiB at a time will be
// read into memory.
func (r *Reader) Read(p []byte) (n int, err error) {
	c := r.ctx.newChunk()
	for n < len(p) {
		if r.cur == nil {
			_, err = io.ReadFull(r.src, c.data)
			if err != nil {
				return n, err
			}
			r.cur = c.readAtom()
			r.cn = int(c.ctx.atomSize)
		}
		nn := copy(p[n:], r.cur.data[int(c.ctx.atomSize)-r.cn:])
		n += nn
		r.cn -= nn
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
