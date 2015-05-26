// chris 052515

package steg

import "testing"

func testReadByte(t *testing.T, c chunk) byte {
	r := c.read()
	t.Logf("chunk: %#v; read: %#v", string(c), string(r))
	return r
}
