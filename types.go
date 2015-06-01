// chris 052515

package steg

import "io"

type Ctx struct {
	atomSize  uint // in bytes
	chunkSize uint // in bytes
}

type atom struct {
	ctx  *Ctx
	data []byte
}

type chunk struct {
	ctx  *Ctx
	data []byte
}

// A Reader wraps an io.Reader and reads steganographically-embedded
// bytes from it.  Implements io.Reader.
type Reader struct {
	ctx *Ctx
	src io.Reader

	// Current atom whose bytes we're returning when Read calls are
	// made.
	cur *atom
	// Remaining bytes before the current atom is exhausted and we
	// need to get another one.
	cn uint
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
	return &atom{ctx: ctx, data: make([]byte, ctx.atomSize)}
}

func (ctx *Ctx) newChunk() *chunk {
	return &chunk{ctx: ctx, data: make([]byte, ctx.chunkSize)}
}

// NewReader returns a fresh Reader, ready to read
// steganographically-embedded bytes from the source io.Reader.
func (ctx *Ctx) NewReader(src io.Reader) *Reader {
	return &Reader{ctx: ctx, src: src, cur: nil, cn: 0}
}
