// chris 052515

package steg

import (
	"bytes"
	"strings"
	"testing"

	"math/rand"
)

func testReadAtomHello1(t *testing.T) {
	ctx := NewCtx(1)
	c := ctx.newChunk()
	copy(c.data, []byte(helloString))
	a := c.readAtom()
	if a.data[0] != helloByte {
		t.Fail()
	}
}

func testReadAtomHello2(t *testing.T) {
	c := testHelloChunk(1)
	a := c.readAtom()
	if a.data[0] != helloByte {
		t.Fail()
	}
}

func testReadAtomRepeat(t *testing.T, atomSize uint) {
	ctx := NewCtx(atomSize)
	c := ctx.newChunk()
	s := string(rand.Intn(256))
	copy(c.data, []byte(strings.Repeat(s, len(c.data))))
	a := c.readAtom()
	if !bytes.Equal(a.data, make([]byte, len(a.data))) {
		t.Fail()
	}
}

func TestReadAtom(t *testing.T) {
	testReadAtomHello1(t)
	testReadAtomHello2(t)
	testReadAtomRepeat(t, 1)
	testReadAtomRepeat(t, 2)
}

//func testReadChunkHello(t *testing.T) {
//	const repeat = 3
//	r := bytes.NewReader([]byte(strings.Repeat(string(testHelloChunk()), repeat)))
//	c := newChunk()
//	err := readChunk(c, r)
//	if err != nil {
//		t.Errorf("read error %v", err)
//		return
//	}
//	if !bytes.Equal([]byte(c), []byte(testHelloChunk())) {
//		t.Errorf("didn't read hello back, got %#v", string(c))
//	}
//}
//
//func testReadChunkShortRead(t *testing.T) {
//	r := bytes.NewReader([]byte(testHelloChunk())[:30])
//	c := newChunk()
//	err := readChunk(c, r)
//	if err == nil {
//		t.Errorf("no error, c = %#v", string(c))
//	}
//}
//
//func TestReadChunk(t *testing.T) {
//	testReadChunkHello(t)
//	testReadChunkShortRead(t)
//}
//
//func testReaderHello(t *testing.T) {
//	const repeat = 3
//	b := bytes.NewBuffer([]byte(strings.Repeat(string(testHelloChunk()), repeat)))
//	r := NewReader(b)
//	out := make([]byte, repeat)
//	n, err := r.Read(out)
//	if n != repeat {
//		t.Errorf("didn't read enough bytes")
//		return
//	}
//	if err != nil {
//		t.Errorf("read error %v", err)
//		return
//	}
//	expect := []byte(strings.Repeat(string([]byte{helloByte}), repeat))
//	if !bytes.Equal(out, expect) {
//		t.Errorf("unexpected string read %#v (expected %#v)", out, expect)
//	}
//}
//
//func testReaderShortRead1(t *testing.T) {
//	const repeat = 3
//	b := bytes.NewBuffer([]byte(strings.Repeat(string(testHelloChunk()), repeat)))
//	r := NewReader(b)
//	requestedBytes := repeat + 1
//	out := make([]byte, requestedBytes)
//	n, err := r.Read(out)
//	if n == requestedBytes {
//		t.Errorf("read too many bytes, n = %v, out = %#v", n, string(out))
//		return
//	}
//	if err == nil {
//		t.Errorf("no error, n = %v, out = %#v", n, string(out))
//		return
//	}
//}
//
//func testReaderShortRead2(t *testing.T) {
//	const repeat = 2
//	b := bytes.NewBuffer([]byte(strings.Repeat(string(testHelloChunk()), repeat))[:48])
//	r := NewReader(b)
//	requestedBytes := repeat
//	out := make([]byte, requestedBytes)
//	n, err := r.Read(out)
//	if n == requestedBytes {
//		t.Errorf("read too many bytes, n = %v, out = %#v", n, string(out))
//		return
//	}
//	if err == nil {
//		t.Errorf("no error, n = %v, out = %#v", n, string(out))
//		return
//	}
//}
//
//func TestReader(t *testing.T) {
//	testReaderHello(t)
//	testReaderShortRead1(t)
//	testReaderShortRead2(t)
//}
