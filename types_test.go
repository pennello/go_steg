// chris 053115

package steg

import (
	"math"
	"testing"
)

func testNewCtxPanic(t *testing.T, atomSize uint) {
	defer func() {
		if r := recover(); r == nil {
			t.Error(r)
		}
	}()
	NewCtx(atomSize)
}

func testNewCtx(t *testing.T, atomSize uint) {
	ctx := NewCtx(atomSize)
	a := ctx.newAtom()
	if atomSize != uint(len(a.data)) {
		t.Fail()
	}
	expect := int(math.Pow(2, float64(atomSize*8)) / 8)
	c := ctx.newChunk()
	if expect != len(c.data) {
		t.Fail()
	}
}

func TestNewCtx(t *testing.T) {
	testNewCtxPanic(t, 0)
	testNewCtxPanic(t, 4)
	testNewCtx(t, 1)
	testNewCtx(t, 2)
	testNewCtx(t, 3)
}
