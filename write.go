// chris 052515

package steg

import (
	//"errors"
	//"io"
)

//// ErrShortCarrier is similar to ErrShortRead, but is specialized for
//// errors reading from the carrier io.Reader in Writer.Write.
//var ErrShortCarrier = errors.New("not enough carrier data")

func (a *atom) xorBit(bit uint8, bitIndex uint) {
	// The bits in bitIndex above 3 tell us in which atom Byte to
	// put the bit, and the low 3 bits tell us which bit in that
	// byte this is for.
	// aBi: atom byte index
	aBi := bitIndex >> 3
	// absi: atom bit sub-index
	absi := bitIndex & 0x7
	// XOR the bit.
	a.data[aBi] ^= bit << absi
	// Done!
}

//// A Writer enables you to write steganographically-embedded bytes into
//// a destination io.Writer by using the data read from a carrier
//// io.Reader.  Implements io.Writer.
//type Writer struct {
//	dst     io.Writer
//	carrier io.Reader
//}
//
//// NewWriter returns a fresh Writer, ready to write
//// steganographically-embedded bytes to the destination io.Writer using
//// the data from the carrier io.Reader.
//func NewWriter(dst io.Writer, carrier io.Reader) Writer {
//	return Writer{dst: dst, carrier: carrier}
//}
//
//// Write chunk into destination io.Reader.
//func (w Writer) write(c chunk) error {
//	// XXX Can io.Writer.Write return an error even if n = len(p)?
//	_, err := w.dst.Write([]byte(c))
//	return err
//}
//
//// Write steganographically-embedded bytes to the destination io.Writer
//// using data from the carrier io.Reader.  Returns the number of bytes
//// written, as well as an error, if one occurred.
////
//// Can return ErrShortCarrier if an EOF was encountered before being
//// able to read a sufficient number of bytes from the carrier to embed
//// the requested data.  Note that in this case, you're sort of sunk--we
//// couldn't read enough data from the carrier to embed some byte, so the
//// carrier data was therefore thrown away before being written to the
//// destination.
////
//// n == len(p) iff err != nil
//func (w Writer) Write(p []byte) (n int, err error) {
//	c := newChunk()
//	for _, b := range p {
//		err = readChunk(c, w.carrier)
//		if err != nil {
//			if err == ErrShortRead {
//				err = ErrShortCarrier
//			}
//			return n, err
//		}
//		c.write(b)
//		err = w.write(c)
//		if err != nil {
//			// We may have written _some_ of the bytes of c,
//			// but won't have written all of them.  We
//			// consider this to be b _not_ having been
//			// written, so n remains unincremented.
//			return n, err
//		}
//		n++
//	}
//	return n, nil
//}
//
//// Copy copies from the carrier to the destination without doing any
//// steganographic embedding.  It's implemented by a simple call to
//// io.Copy.
////
//// The idea is that you'd call this to send through the rest of your
//// carrier data after you've finished successfully with any Writes.
//func (w Writer) Copy() (written int64, err error) {
//	return io.Copy(w.dst, w.carrier)
//}
//
//// CopyN copies n bytes from the carrier to the destination without
//// doing any steganographic embedding.  It's implemented by a simple
//// call to io.CopyN.
////
//// The idea is that you'd call this before sending through any of your
//// message data to get past critical headers in your carrier before
//// embedding your data.
////
//// Counterpart to Reader.Discard.
//func (w Writer) CopyN(n int64) (written int64, err error) {
//	return io.CopyN(w.dst, w.carrier, n)
//}
