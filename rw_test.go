// chris 052715 Read/write tests.

package steg

import (
	"bytes"
	"testing"

	cryptorand "crypto/rand"
	mathrand "math/rand"
)

func testReadWriteByte(t *testing.T, b byte) {
	c := newChunk()
	backup := newChunk()
	copy([]byte(backup), []byte(c))
	c.write(b)
	testChunkDiff(t, backup, c)
	r := c.read()
	if r != b {
		t.Errorf("didn't read back %v after writing to chunk (got %v)", b, r)
	}
}

func TestReadWriteByte(t *testing.T) {
	testReadWriteByte(t, 'a')
	testReadWriteByte(t, '0')
	testReadWriteByte(t, '\x00')

	for i := 0; i < 1000; i++ {
		testReadWriteByte(t, byte(mathrand.Uint32()))
	}
}

func testReaderWriter(t *testing.T, secret []byte) {
	var n int
	var err error

	// First, do write and test.
	carrierBytes := make([]byte, len(secret)*chunkSize+19)
	_, err = cryptorand.Read(carrierBytes)
	if err != nil {
		t.Errorf("failed to read random carrier bytes %v", err)
		return
	}
	carrier := bytes.NewReader(carrierBytes)
	dst := new(bytes.Buffer)
	w := NewWriter(dst, carrier)
	n, err = w.Write(secret)
	if n != len(secret) {
		t.Errorf("incomplete write; n = %v, err = %v", n, err)
		return
	}
	if err != nil {
		t.Errorf("write error %v", err)
		return
	}
	testBytesDiff(t, carrierBytes, dst.Bytes(), len(secret))

	// Now read it back and test that.
	test := make([]byte, len(secret))
	r := NewReader(dst)
	n, err = r.Read(test)
	if n != len(test) {
		t.Errorf("incomplete read; n = %v, err = %v", n, err)
		return
	}
	if err != nil {
		t.Errorf("read error %v", err)
		return
	}
	if !bytes.Equal(test, secret) {
		t.Errorf("failed to read back out %#v (got %#v)", string(secret), string(test))
	}
}

func testReaderWriterRandom(t *testing.T) {
	length := mathrand.Intn(32)
	secret := make([]byte, length)
	_, err := cryptorand.Read(secret)
	if err != nil {
		t.Errorf("failed to generate random data, err = %v", err)
		return
	}
	testReaderWriter(t, secret)
}

func TestReaderWriter(t *testing.T) {
	testReaderWriter(t, []byte("top secret!"))
	testReaderWriter(t, []byte(""))
	testReaderWriter(t, []byte("\x00"))

	for i := 0; i < 1000; i++ {
		testReaderWriterRandom(t)
	}
}
