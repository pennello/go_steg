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

func TestReader(t *testing.T) {
	testReaderHello(t)
}
