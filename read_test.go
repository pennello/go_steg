// chris 052515

package steg

import (
	"bytes"
	"strings"
	"testing"

	"encoding/binary"

	cryptorand "crypto/rand"
	mathrand "math/rand"
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

func testReadAtomRepeat(t *testing.T, atomSize uint8) {
	ctx := NewCtx(atomSize)
	c := ctx.newChunk()
	s := string(mathrand.Intn(256))
	copy(c.data, []byte(strings.Repeat(s, int(ctx.chunkSize))))
	a := c.readAtom()
	// Since everything is just a repeated single byte, should get
	// all 0s back.
	if !bytes.Equal(a.data, make([]byte, ctx.atomSize)) {
		t.Fail()
	}
}

func TestReadAtom(t *testing.T) {
	testReadAtomHello1(t)
	testReadAtomHello2(t)
	testReadAtomRepeat(t, 1)
	testReadAtomRepeat(t, 2)
}

func testReaderHello(t *testing.T) {
	const repeat = 3
	const atomSize = 1
	ctx := NewCtx(atomSize)
	s := string(testHelloChunk(atomSize).data)
	src := bytes.NewBuffer([]byte(strings.Repeat(s, repeat)))
	r := ctx.NewReader(src)
	out := make([]byte, atomSize*repeat)
	n, err := r.Read(out)
	if n != len(out) {
		t.Errorf("didn't read enough bytes, atomSize = %v", atomSize)
		return
	}
	if err != nil {
		t.Errorf("read error %v, atomSize = %v", err, atomSize)
		return
	}
	expect := []byte(strings.Repeat(string([]byte{helloByte}), int(atomSize)*repeat))
	if !bytes.Equal(out, expect) {
		t.Errorf("unexpected string read %#v (expected %#v)", out, expect)
	}
}

func testReaderLorem(t *testing.T) {
	const repeat = 3
	const atomSize = 2
	ctx := NewCtx(atomSize)
	s := string(testLoremChunk().data)
	src := bytes.NewBuffer([]byte(strings.Repeat(s, repeat)))
	r := ctx.NewReader(src)
	out := make([]byte, atomSize*repeat)
	n, err := r.Read(out)
	if n != len(out) {
		t.Errorf("didn't read enough bytes, atomSize = %v", atomSize)
		return
	}
	if err != nil {
		t.Errorf("read error %v, atomSize = %v", err, atomSize)
		return
	}
	expect := []byte(strings.Repeat(string(loremBytes), repeat))
	if !bytes.Equal(out, expect) {
		t.Errorf("unexpected string read %#v (expected %#v)", out, expect)
	}
}

func testReaderShortRead1(t *testing.T, atomSize uint8) {
	const repeat = 3
	ctx := NewCtx(atomSize)
	s := string(testHelloChunk(atomSize).data)
	b := bytes.NewBuffer([]byte(strings.Repeat(s, repeat)))
	r := ctx.NewReader(b)
	requestedBytes := int(atomSize)*repeat + 1
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

func testReaderShortRead2(t *testing.T, atomSize uint8) {
	const repeat = 2
	ctx := NewCtx(atomSize)
	s := string(testHelloChunk(atomSize).data)
	b := bytes.NewBuffer([]byte(strings.Repeat(s, repeat))[:ctx.chunkSize*2/3])
	r := ctx.NewReader(b)
	requestedBytes := int(atomSize) * repeat
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
	testReaderLorem(t)
	testReaderShortRead1(t, 1)
	testReaderShortRead1(t, 2)
	testReaderShortRead2(t, 1)
	testReaderShortRead2(t, 2)
}

func testAsUint(t *testing.T, a *atom) {
	p := make([]byte, 4)
	copy(p, a.data)
	out := a.asUint32()
	expect := binary.LittleEndian.Uint32(p)
	if out != expect {
		t.Errorf("(%v).asUint() != %v (was %v)", a.data, expect, out)
	}
}

func testAsUintCtx(t *testing.T, atomSize uint8) {
	ctx := NewCtx(atomSize)
	for i := 0; i < 1000; i++ {
		a := ctx.newAtom()
		_, err := cryptorand.Read(a.data)
		if err != nil {
			t.Error(err)
			return
		}
		testAsUint(t, a)
	}
}

func TestAsUint(t *testing.T) {
	testAsUintCtx(t, 1)
	testAsUintCtx(t, 2)
	testAsUintCtx(t, 3)
}
