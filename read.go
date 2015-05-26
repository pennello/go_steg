// chris 052515

package steg

import "chrispennello.com/go/swar"

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
