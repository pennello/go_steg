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

func testChunkDiff(t *testing.T, a, b chunk) {
	bitsDifferent := 0
	for i := 0; i < chunkSize; i++ {
		bitsDifferent += int(swar.Ones8(a[i] ^ b[i]))
	}
	if bitsDifferent != 1 {
		t.Errorf("%#v and %#v differ by more than 1 bit (by %v)", a, b, bitsDifferent)
	}
}

func testWrite(t *testing.T, c chunk, b byte) {
	testReadByte(t, c)
	// Make backup.
	cb := newChunk()
	copy([]byte(cb), []byte(c))
	c.write(b)
	r := testReadByte(t, c)
	if r != b {
		t.Errorf("failed to write %#v and read back (got %#v); original was %#v, after writing: %#v", b, r, cb, c)
	}
	testChunkDiff(t, cb, c)
}

func testWriteHello(t *testing.T, b byte) {
	testWrite(t, testHelloChunk(), b)
}

func testWriteRandom(t *testing.T) {
	buf := make([]byte, chunkSize+1)
	_, err := cryptorand.Read(buf)
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

func testBytesDiff(t *testing.T, a, b []byte, expectBitsDifferent int) {
	var m int
	if len(a) < len(b) {
		m = len(a)
	} else {
		m = len(b)
	}
	bitsDifferent := 0
	for i := 0; i < m; i++ {
		bitsDifferent += int(swar.Ones8(a[i] ^ b[i]))
	}
	if bitsDifferent != expectBitsDifferent {
		t.Errorf("more than %v bit difference (is %v); %#v %#v",
			expectBitsDifferent,
			bitsDifferent,
			string(a),
			string(b))
	}
}

func testWriterHello(t *testing.T) {
	var n int
	var err error
	testBytes := []byte("secret message")
	dst := new(bytes.Buffer)
	carrierBytes := []byte(strings.Repeat(string(testHelloChunk()), len(testBytes)+17))
	carrier := bytes.NewReader(carrierBytes)
	w := NewWriter(dst, carrier)
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

func testWriterRandom(t *testing.T) {
	var n int
	var err error
	testBytes := make([]byte, mathrand.Intn(10)+1)
	testSize := len(testBytes) * (3 / 2 * chunkSize)
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
	w := NewWriter(dst, carrier)
	n, err = w.Write(testBytes)
	if n != len(testBytes) {
		t.Errorf("incomplete write; n = %v, err = %v", n, err)
		return
	}
	if err != nil {
		t.Errorf("write error %v", err)
		return
	}
	testBytesDiff(t, carrierBytes, dst.Bytes(), len(testBytes))
}

func TestWriter(t *testing.T) {
	testWriterHello(t)
	for i := 0; i < 1000; i++ {
		testWriterRandom(t)
	}
}
