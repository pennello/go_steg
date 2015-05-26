// chris 052515

package steg

import "testing"

func TestReadByte(t *testing.T) {
	testReadByte(t, chunk("hello, there, how are you? fine."))
	//for i := bitIndex(0); i < 8; i++ {
	//	testReadByte(t, chunk(masksByIndex[i]))
	//}
}
