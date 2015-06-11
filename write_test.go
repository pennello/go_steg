// chris 052515

package steg

import (
	"bytes"
	"strings"
	"testing"

	cryptorand "crypto/rand"
	mathrand "math/rand"

	"chrispennello.com/go/swar"
)

func testXorBit(t *testing.T, p []byte, bit uint8, bitIndex uint, expect []byte) {
	xorBit(p, bit, bitIndex)
	if !bytes.Equal(p, expect) {
		t.Fail()
	}
}

func TestXorBit(t *testing.T) {
	testXorBit(t, []byte{0}, 1, 0, []byte{1})
	testXorBit(t, []byte{0}, 1, 7, []byte{0x80})
	testXorBit(t, []byte{0, 0, 0}, 1, 7, []byte{0x80, 0, 0})
	testXorBit(t, []byte{0, 0, 0}, 1, 8, []byte{0, 1, 0})
	testXorBit(t, []byte{0, 0, 0x80}, 1, 23, []byte{0, 0, 0})
}

func testAtomXorBit(t *testing.T, atomSize uint, bit uint8, bitIndex uint, expect []byte) {
	ctx := NewCtx(atomSize)
	a := ctx.newAtom()
	a.xorBit(bit, bitIndex)
	if !bytes.Equal(expect, a.data) {
		t.Fail()
	}
}

func TestAtomXorBit(t *testing.T) {
	var atomSize uint
	var e []byte

	atomSize = 3
	e = make([]byte, atomSize)
	e[0] = 1
	testAtomXorBit(t, atomSize, 1, 0, e)

	atomSize = 2
	e = make([]byte, atomSize)
	e[1] = 1
	testAtomXorBit(t, atomSize, 1, 8, e)
}

func testChunkDiff(t *testing.T, a, b *chunk) {
	if a.ctx != b.ctx {
		panic("chunks with different contexts")
	}
	bitsDiff := 0
	for i := uint(0); i < a.ctx.chunkSize; i++ {
		bitsDiff += int(swar.Ones8(a.data[i] ^ b.data[i]))
	}
	if bitsDiff != 1 {
		t.Errorf("%#v and %#v differ by more than 1 bit (by %v)", a.data, b.data, bitsDiff)
	}
}

func testWriteAtom(t *testing.T, c *chunk, a *atom) {
	// Make backup.
	cb := c.ctx.newChunk()
	copy([]byte(cb.data), []byte(c.data))
	c.write(a)
	aa := c.readAtom()
	if !bytes.Equal(a.data, aa.data) {
		t.Errorf("failed to write %#v and read back (got %#v); original was %#v, after writing: %#v",
			a.data, aa.data, cb.data, c.data)
	}
	testChunkDiff(t, cb, c)
}

func testWriteByteHello(t *testing.T, b byte) {
	ctx := NewCtx(1)
	a := ctx.newAtom()
	copy(a.data, []byte{b})
	testWriteAtom(t, testHelloChunk(1), a)
}

func testWriteAtomLorem(t *testing.T, a *atom) {
	ctx := NewCtx(2)
	aa := ctx.newAtom()
	copy(aa.data, a.data)
	testWriteAtom(t, testLoremChunk(), a)
}

func TestWriteAtom(t *testing.T) {
	for x := byte(0x20); x < 0x80; x++ {
		testWriteByteHello(t, x)
	}

	ctx := NewCtx(2)
	for i := 0; i < 100; i++ {
		a := ctx.newAtom()
		_, err := cryptorand.Read(a.data)
		if err != nil {
			t.Error("failed to generate random data for test")
		}
		testWriteAtomLorem(t, a)
	}
}

func testBytesDiff(t *testing.T, a, b []byte, expectBits int) {
	var m int
	if len(a) < len(b) {
		m = len(a)
	} else {
		m = len(b)
	}
	bitsDiff := 0
	for i := 0; i < m; i++ {
		bitsDiff += int(swar.Ones8(a[i] ^ b[i]))
	}
	if bitsDiff != expectBits {
		t.Errorf("not %v bit difference (is %v)", expectBits, bitsDiff)
	}
}

func testWriterHelloExtra(t *testing.T, extra int) {
	ctx := NewCtx(1)
	var n int
	var err error
	testBytes := []byte("secret message")
	dst := new(bytes.Buffer)
	carrierBytes := []byte(strings.Repeat(string(testHelloChunk(1).data), len(testBytes)+extra))
	carrier := bytes.NewReader(carrierBytes)
	w := ctx.NewWriter(dst, carrier)
	n, err = w.Write(testBytes)
	if n != len(testBytes) {
		t.Errorf("incomplete write; n = %v, err = %v", n, err)
		return
	}
	if err != nil {
		t.Errorf("write error %v", err)
		return
	}
	_, err = w.Copy()
	if err != nil {
		t.Errorf("remaining data copy error %v", err)
	}
	testBytesDiff(t, carrierBytes, dst.Bytes(), len(testBytes))
}

func testWriterLoremExtra(t *testing.T, extra int) {
	ctx := NewCtx(2)
	var n int
	var err error
	testBytes := []byte("secret message")
	dst := new(bytes.Buffer)
	carrierBytes := []byte(strings.Repeat(string(testLoremChunk().data), len(testBytes)/2+extra))
	carrier := bytes.NewReader(carrierBytes)
	w := ctx.NewWriter(dst, carrier)
	n, err = w.Write(testBytes)
	if n != len(testBytes) {
		t.Errorf("incomplete write; n = %v, err = %v", n, err)
		return
	}
	if err != nil {
		t.Errorf("write error %v", err)
		return
	}
	_, err = w.Copy()
	if err != nil {
		t.Errorf("remaining data copy error %v", err)
	}
	testBytesDiff(t, carrierBytes, dst.Bytes(), len(testBytes)/2)
}

func testWriterRandom(t *testing.T, atomSize uint) {
	ctx := NewCtx(atomSize)
	var n int
	var err error
	nTestBytes := uint(mathrand.Intn(10) + 1)
	rem := nTestBytes % atomSize
	if rem != 0 {
		nTestBytes += atomSize - nTestBytes
	}
	testBytes := make([]byte, nTestBytes)
	testSize := len(testBytes) * int(3*ctx.chunkSize/2)
	dst := new(bytes.Buffer)
	carrierBytes := make([]byte, testSize)
	_, err = cryptorand.Read(carrierBytes)
	if err != nil {
		t.Errorf("failed to generate random carrier data for test; %v", err)
		return
	}
	carrier := bytes.NewReader(carrierBytes)
	_, err = cryptorand.Read(testBytes)
	if err != nil {
		t.Errorf("failed to generate random test data for test; %v", err)
		return
	}
	w := ctx.NewWriter(dst, carrier)
	n, err = w.Write(testBytes)
	if n != len(testBytes) {
		t.Errorf("incomplete write; n = %v, err = %v", n, err)
		return
	}
	if err != nil {
		t.Errorf("write error %v", err)
		return
	}
	testBytesDiff(t, carrierBytes, dst.Bytes(), len(testBytes)/int(atomSize))
}

func testWriterShortReadRandom(t *testing.T, atomSize uint) {
	ctx := NewCtx(atomSize)
	var n int
	var err error
	nTestBytes := uint(mathrand.Intn(10) + 1)
	rem := nTestBytes % atomSize
	if rem != 0 {
		nTestBytes += atomSize - nTestBytes
	}
	testBytes := make([]byte, nTestBytes)
	testSize := uint(len(testBytes) / int(atomSize) * int(2*ctx.chunkSize/3))
	dst := new(bytes.Buffer)
	carrierBytes := make([]byte, testSize)
	_, err = cryptorand.Read(carrierBytes)
	if err != nil {
		t.Errorf("failed to generate random carrier data for test; %v", err)
		return
	}
	carrier := bytes.NewReader(carrierBytes)
	_, err = cryptorand.Read(testBytes)
	if err != nil {
		t.Errorf("failed to generate random test data for test; %v", err)
		return
	}
	w := ctx.NewWriter(dst, carrier)
	expect := testSize / ctx.chunkSize
	n, err = w.Write(testBytes)
	if n == len(testBytes) {
		t.Errorf("wrote too much; n = %v (expected %v)", n, expect)
		return
	}
	if err == nil {
		t.Errorf("no error; n =  %v", n)
		return
	}
}

func TestWriter(t *testing.T) {
	testWriterHelloExtra(t, 0)
	testWriterHelloExtra(t, 17)
	testWriterLoremExtra(t, 0)
	testWriterLoremExtra(t, 17)
	for i := 0; i < 100; i++ {
		testWriterRandom(t, 1)
		testWriterRandom(t, 2)
		testWriterShortReadRandom(t, 1)
		testWriterShortReadRandom(t, 2)
	}
}
