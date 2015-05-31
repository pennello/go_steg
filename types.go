// chris 052515

package steg

import (
	"chrispennello.com/go/swar"
)

type Ctx struct {
	atomSize  uint // in bytes
	chunkSize uint // in bytes
}

type atom struct {
	data []byte
	ctx *Ctx
}

type chunk struct {
	data []byte
	ctx *Ctx
}

func NewCtx(atomSize uint) *Ctx {
	if atomSize < 1 {
		panic("inappropriate atom size")
	}
	if atomSize > 3 {
		// See the chunk.ReadBit implementation.
		panic("unsupported atom size")
	}
	// (2 ^ (atomSize * 8)) / 8
	chunkSize := uint(1) << (atomSize*8 - 3)
	return &Ctx{atomSize: atomSize, chunkSize: chunkSize}
}

func (ctx *Ctx) newAtom() *atom {
	return &atom{data: make([]byte, ctx.atomSize), ctx: ctx}
}

func (ctx *Ctx) newChunk() *chunk {
	return &chunk{data: make([]byte, ctx.chunkSize), ctx: ctx}
}
