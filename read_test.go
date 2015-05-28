// chris 052515

package steg

import (
	"bytes"
	"strings"
	"testing"
)

func testReadByte(t *testing.T, c chunk) byte {
	r := c.read()
	//t.Logf("chunk: %#v; read: %#v", string(c), r)
	return r
}

func testReadByteHello(t *testing.T) {
	out := testReadByte(t, testHelloChunk())
	if out != helloByte {
		t.Errorf("didn't get back u with hat (got %#v)", out)
	}
}

func TestReadByte(t *testing.T) {
	testReadByteHello(t)
	for i := bitIndex(0); i < 8; i++ {
		if testReadByte(t, chunk(masksByIndex[i])) != 0 {
			t.Errorf("mask at index %v didn't yield 0", i)
		}
	}
}

func testReadChunkHello(t *testing.T) {
	const repeat = 3
	r := bytes.NewReader([]byte(strings.Repeat(string(testHelloChunk()), repeat)))
	c := newChunk()
	complete, err := readChunk(c, r)
	if !complete {
		t.Errorf("incomplete read, err = %v", err)
		return
	}
	if err != nil {
		t.Errorf("read error %v", err)
		return
	}
	if !bytes.Equal([]byte(c), []byte(testHelloChunk())) {
		t.Errorf("didn't read hello back, got %#v", string(c))
	}
}

func testReadChunkShortRead(t *testing.T) {
	r := bytes.NewReader([]byte(testHelloChunk())[:30])
	c := newChunk()
	complete, err := readChunk(c, r)
	if complete {
		t.Errorf("complete read, c = %#v", string(c))
		return
	}
	if err == nil {
		t.Errorf("no error, c = %#v", string(c))
	}
}

func TestReadChunk(t *testing.T) {
	testReadChunkHello(t)
	testReadChunkShortRead(t)
}

func testReaderHello(t *testing.T) {
	const repeat = 3
	b := bytes.NewBuffer([]byte(strings.Repeat(string(testHelloChunk()), repeat)))
	r := NewReader(b)
	out := make([]byte, repeat)
	n, err := r.Read(out)
	if n != repeat {
		t.Errorf("didn't read enough bytes")
		return
	}
	if err != nil {
		t.Errorf("read error %v", err)
		return
	}
	expect := []byte(strings.Repeat(string([]byte{helloByte}), repeat))
	if !bytes.Equal(out, expect) {
		t.Errorf("unexpected string read %#v (expected %#v)", out, expect)
	}
}

func testReaderShortRead1(t *testing.T) {
	const repeat = 3
	b := bytes.NewBuffer([]byte(strings.Repeat(string(testHelloChunk()), repeat)))
	r := NewReader(b)
	requestedBytes := repeat + 1
	out := make([]byte, requestedBytes)
	n, err := r.Read(out)
	if n == requestedBytes {
		t.Errorf("read too many bytes, n = %v, out = %#v", n, string(out))
		return
	}
	if err == nil {
		t.Errorf("no error, n = %v, out = %#v", n, string(out))
		return
	}
}

func testReaderShortRead2(t *testing.T) {
	const repeat = 2
	b := bytes.NewBuffer([]byte(strings.Repeat(string(testHelloChunk()), repeat))[:48])
	r := NewReader(b)
	requestedBytes := repeat
	out := make([]byte, requestedBytes)
	n, err := r.Read(out)
	if n == requestedBytes {
		t.Errorf("read too many bytes, n = %v, out = %#v", n, string(out))
		return
	}
	if err == nil {
		t.Errorf("no error, n = %v, out = %#v", n, string(out))
		return
	}
}

func TestReader(t *testing.T) {
	testReaderHello(t)
	testReaderShortRead1(t)
	testReaderShortRead2(t)
}
