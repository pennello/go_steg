// chris 052515

package steg

import (
	"testing"

	"crypto/rand"

	"chrispennello.com/go/swar"
)

func testDiff(t *testing.T, a chunk, b chunk) {
	bitsDifferent := 0
	for i := 0; i < chunkSize; i++ {
		bitsDifferent += int(swar.Ones8(a[i] ^ b[i]))
	}
	if bitsDifferent != 1 {
		t.Errorf("%#v and %#v differ by more than 1 bit", a, b)
	}
}

func testWrite(t *testing.T, c chunk, b byte) {
	t.Log("original chunk begin")
	testReadByte(t, c)
	// Make backup.
	cb := chunk(make([]byte, chunkSize))
	copy([]byte(cb), []byte(c))
	t.Logf("writing byte %#v", string(b))
	c.write(b)
	if testReadByte(t, c) != b {
		t.Errorf("failed to write %#v and read back", b)
	}
	testDiff(t, cb, c)
}

func testWriteHello(t *testing.T, b byte) {
	c := chunk("hello, there, how are you? fine.")
	testWrite(t, c, b)
}

func testWriteRandom(t *testing.T) {
	buf := make([]byte, chunkSize+1)
	_, err := rand.Read(buf)
	if err != nil {
		t.Errorf("failed to generate random data for test; %v", err)
		return
	}
	c := chunk(buf[1:])
	b := byte(buf[0])
	testWrite(t, c, b)
}

func TestWrite(t *testing.T) {
	for b := byte(0x20); b < 0x80; b++ {
		testWriteHello(t, b)
	}
	for i := 0; i < 1000; i++ {
		testWriteRandom(t)
	}
}
