// chris 052515

package steg

import (
	"errors"
	"io"
)

// ErrShortCarrier is similar to ErrShortRead, but is specialized for
// errors reading from the carrier io.Reader in Writer.Write.
var ErrShortCarrier = errors.New("not enough carrier data")

// ErrInsufficientData is returned when the number of bytes to write
// passed into a Writer.Write call is not a multiple of the atom size
// being used.
var ErrInsufficientData = errors.New("data size not a multiple of atom size")

// XOR the bit into the byte slice p given the specified bitIndex.
func xorBit(p []byte, bit uint8, bitIndex uint) {
	// The bits in bitIndex above 3 tell us in which slice byte to
	// xor the bit, and the low 3 bits tell us which bit in that
	// byte this is for.
	Bi := bitIndex >> 3   // byte index
	bsi := bitIndex & 0x7 // bit sub-index
	// XOR the bit.
	p[Bi] ^= bit << bsi
	// Done!

}

func (a *atom) xorBit(bit uint8, bitIndex uint) {
	xorBit(a.data, bit, bitIndex)
}

// In-place xor of a with b.  Alters a.
func (a *atom) xor(b *atom) {
	for i := uint(0); i < a.ctx.atomSize; i++ {
		a.data[i] ^= b.data[i]
	}
}

// Zero out the bytes in the atom's data starting at the given offset.
func (a *atom) zero(off int) {
	for ; off < int(a.ctx.atomSize); off++ {
		a.data[off] = 0
	}
}

// Copy data into a.data.
func (a *atom) copy(data []byte) {
	if uint(len(data)) != a.ctx.atomSize {
		panic("mis-matched atom copy")
	}
	copy(a.data, data)
}

func (c *chunk) write(a *atom) {
	// Compare current value with what we need to write.
	x := c.readAtom().asUint() ^ a.asUint()
	// x is now a bit index to which bit in c we need to flip.
	xorBit(c.data, 1, x)
}

// Write chunk into destination io.Reader.
func (w *Writer) write(c *chunk) error {
	// XXX Can io.Writer.Write return an error even if n = len(p)?
	_, err := w.dst.Write(c.data)
	return err
}

// Write steganographically-embedded bytes to the destination io.Writer
// using data from the carrier io.Reader.  Returns the number of bytes
// written, as well as an error, if one occurred.
//
// The number of bytes to write must be a multiple of the atom size
// being used.  If it is not, ErrInsufficientData will be returned
// immediately with zero bytes written.
//
// Can return ErrShortCarrier if an EOF was encountered before being
// able to read a sufficient number of bytes from the carrier to embed
// the requested data.  Note that in this case, you're sort of sunk--we
// couldn't read enough data from the carrier to embed some atom, so the
// carrier data was therefore thrown away before being written to the
// destination.
//
// n == len(p) iff err != nil
func (w *Writer) Write(p []byte) (n int, err error) {
	if uint(len(p))%w.ctx.atomSize != 0 {
		return 0, ErrInsufficientData
	}
	c := w.ctx.newChunk()
	for n < len(p) {
		err = readChunk(c, w.carrier)
		if err != nil {
			if err == ErrShortRead {
				err = ErrShortCarrier
			}
			return n, err
		}
		a := w.ctx.newAtom()
		a.copy(p[n : n+int(w.ctx.atomSize)])
		c.write(a)
		err = w.write(c)
		if err != nil {
			// We may have written _some_ of the bytes of c,
			// but won't have written all of them.  We
			// consider this to be a _not_ having been
			// written, so n remains unincremented.
			return n, err
		}
		n += int(w.ctx.atomSize)
	}
	return n, nil
}

// Copy copies from the carrier to the destination without doing any
// steganographic embedding.  It's implemented by a simple call to
// io.Copy.
//
// The idea is that you'd call this to send through the rest of your
// carrier data after you've finished successfully with any Writes.
func (w *Writer) Copy() (written int64, err error) {
	return io.Copy(w.dst, w.carrier)
}

// CopyN copies n bytes from the carrier to the destination without
// doing any steganographic embedding.  It's implemented by a simple
// call to io.CopyN.
//
// The idea is that you'd call this before sending through any of your
// message data to get past critical headers in your carrier before
// embedding your data.
//
// Counterpart to Reader.Discard.
func (w *Writer) CopyN(n int64) (written int64, err error) {
	return io.CopyN(w.dst, w.carrier, n)
}
