// chris 052515

package steg

import "io"

// A Ctx is a context that encapsulates the desired atom size.  Create
// atoms, chunks, Readers, Writers, and Muxes from a context.
type Ctx struct {
	// Atom size at most 3, so it will fit in a uint8.
	atomSize uint8 // in bytes
	// Chunk size at most 2Mi, so it will fit in a uint32.
	chunkSize uint32 // in bytes
	// In general, it's safe to cast either of these to an int.

	// Even a chunk bit index will be at most 2Mi * 8 - 1 =
	// 16Mi - 1, which will fit in a uint32.
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
	cn int
}

// A Writer enables you to write steganographically-embedded bytes into
// a destination io.Writer by using the data read from a carrier
// io.Reader.  Implements io.Writer.
type Writer struct {
	ctx *Ctx

	dst     io.Writer
	carrier io.Reader
}

// Mux multiplexes a message on a carrier into a destination.  It
// steganographically embeds data from the message into the carrier and
// then writes the modified data into the destination.
type Mux struct {
	ctx *Ctx

	w   *Writer
	msg io.Reader
}

// NewCtx returns a fresh Ctx, ready to create the other types.  Panics
// if atomSize is not 1, 2, or 3.
func NewCtx(atomSize uint8) *Ctx {
	if atomSize < 1 {
		panic("inappropriate atom size")
	}
	// NB: atomSize <= 3 depended on elsewhere in the code for type
	// safety.
	if atomSize > 3 {
		// See the chunk.ReadBit implementation.
		panic("unsupported atom size")
	}
	// (2 ^ (atomSize * 8)) / 8
	chunkSize := uint32(1) << (atomSize*8 - 3)
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

// NewWriter returns a fresh Writer, ready to write
// steganographically-embedded bytes to the destination io.Writer using
// the data from the carrier io.Reader.
func (ctx *Ctx) NewWriter(dst io.Writer, carrier io.Reader) *Writer {
	return &Writer{ctx: ctx, dst: dst, carrier: carrier}
}

// NewMux returns a fresh Mux, ready to multiplex a message on a carrier
// into a destination.
func (ctx *Ctx) NewMux(dst io.Writer, carrier, msg io.Reader) Mux {
	w := ctx.NewWriter(dst, carrier)
	return Mux{ctx: ctx, w: w, msg: msg}
}
