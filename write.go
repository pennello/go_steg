// chris 052515

package steg

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
