// chris 052715 Read/write tests.

package steg

import (
	"bytes"
	"testing"

	cryptorand "crypto/rand"
	mathrand "math/rand"
)

func testReadWriteAtom(t *testing.T, a *atom) {
	c := a.ctx.newChunk()
	backup := a.ctx.newChunk()
	copy([]byte(backup.data), []byte(c.data))
	c.write(a)
	testChunkDiff(t, backup, c)
	r := c.readAtom()
	if !bytes.Equal(r.data, a.data) {
		t.Errorf("didn't read back %v after writing to chunk (got %v)", a.data, r.data)
	}
}

func TestReadWriteAtom(t *testing.T) {
	ctx := NewCtx(1)
	a := ctx.newAtom()

	copy(a.data, []byte{'a'})
	testReadWriteAtom(t, a)
	copy(a.data, []byte{'0'})
	testReadWriteAtom(t, a)
	copy(a.data, []byte{'\x00'})
	testReadWriteAtom(t, a)

	for i := 0; i < 100; i++ {
		copy(a.data, []byte{byte(mathrand.Uint32())})
		testReadWriteAtom(t, a)
	}

	ctx2 := NewCtx(2)
	a2 := ctx2.newAtom()
	copy(a2.data, []byte{'a', 'a'})
	testReadWriteAtom(t, a)
	copy(a2.data, []byte{'0', '0'})
	testReadWriteAtom(t, a)
	copy(a2.data, []byte{'\x00', '\x00'})
	testReadWriteAtom(t, a)

	for i := 0; i < 100; i++ {
		r := mathrand.Uint32()
		copy(a.data, []byte{byte(r), byte(r >> 8)})
		testReadWriteAtom(t, a)
	}
}

func testReaderWriter(t *testing.T, secret []byte, atomSize uint) {
	ctx := NewCtx(atomSize)

	var n int
	var err error

	// First, do write and test.
	carrierBytes := make([]byte, uint(len(secret))/ctx.atomSize*ctx.chunkSize+19)
	_, err = cryptorand.Read(carrierBytes)
	if err != nil {
		t.Errorf("failed to read random carrier bytes %v", err)
		return
	}
	carrier := bytes.NewReader(carrierBytes)
	dst := new(bytes.Buffer)
	w := ctx.NewWriter(dst, carrier)
	n, err = w.Write(secret)
	if n != len(secret) {
		t.Errorf("incomplete write; n = %v, err = %v", n, err)
		return
	}
	if err != nil {
		t.Errorf("write error %v", err)
		return
	}
	testBytesDiff(t, carrierBytes, dst.Bytes(), len(secret)/int(ctx.atomSize))

	// Now read it back and test that.
	test := make([]byte, len(secret))
	r := ctx.NewReader(dst)
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

func testReaderWriterRandom(t *testing.T, atomSize uint) {
	length := uint(mathrand.Intn(32))
	rem := length % atomSize
	if rem != 0 {
		length += atomSize - rem
	}
	secret := make([]byte, length)
	_, err := cryptorand.Read(secret)
	if err != nil {
		t.Errorf("failed to generate random data, err = %v", err)
		return
	}
	testReaderWriter(t, secret, atomSize)
}

func TestReaderWriter(t *testing.T) {
	testReaderWriter(t, []byte("top secret!"), 1)
	testReaderWriter(t, []byte(""), 1)
	testReaderWriter(t, []byte("\x00\xab"), 1)

	for i := 0; i < 100; i++ {
		testReaderWriterRandom(t, 1)
	}

	testReaderWriter(t, []byte("top secret!!"), 2)
	testReaderWriter(t, []byte(""), 2)
	testReaderWriter(t, []byte("\x00\xab"), 2)

	for i := 0; i < 100; i++ {
		testReaderWriterRandom(t, 2)
	}
}
