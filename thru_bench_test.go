// chris 061515 Throughput benchmark testing.

package steg

import (
	"bytes"
	"testing"

	"crypto/rand"
	"io/ioutil"
)

// benchSize is the number of bytes in the carrier for a benchmarking
// test.
const benchSize = 3 * 1000 * 1000

func benchmarkSetup(b *testing.B, atomSize uint8) (ctx *Ctx, carrierBytes, msgBytes []byte) {
	ctx = NewCtx(atomSize)

	carrierBytes = make([]byte, benchSize)
	if _, err := rand.Read(carrierBytes); err != nil {
		b.Error(err)
		return
	}
	capacity := ctx.Capacity(benchSize)
	msgBytes = make([]byte, capacity)
	if _, err := rand.Read(msgBytes); err != nil {
		b.Error(err)
		return
	}

	return ctx, carrierBytes, msgBytes
}

func benchmarkN(b *testing.B, atomSize uint8) {
	ctx, carrierBytes, msgBytes := benchmarkSetup(b, atomSize)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		carrier := bytes.NewBuffer(carrierBytes)
		msg := bytes.NewBuffer(msgBytes)
		ctx.NewMux(ioutil.Discard, carrier, msg).Mux()
	}
}

// Benchmark by atom size.
func Benchmark1(b *testing.B) { benchmarkN(b, 1) }
func Benchmark2(b *testing.B) { benchmarkN(b, 2) }
func Benchmark3(b *testing.B) { benchmarkN(b, 3) }
