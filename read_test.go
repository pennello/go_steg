// chris 052515

package steg

import "testing"

func TestReadByte(t *testing.T) {
	testReadByte(t, testHelloChunk())
	for i := bitIndex(0); i < 8; i++ {
		if testReadByte(t, chunk(masksByIndex[i])) != 0 {
			t.Errorf("mask at index %v didn't yield 0", i)
		}
	}
}
