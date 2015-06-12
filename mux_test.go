// chris 061015

package steg

import (
	"bytes"
	"testing"

	cryptorand "crypto/rand"
	mathrand "math/rand"
)

func testMuxPad(t *testing.T, atomSize uint8) {
	ctx := NewCtx(atomSize)

	msgLen := mathrand.Intn(4) + 1
	msgBytes := make([]byte, msgLen)
	if _, err := cryptorand.Read(msgBytes); err != nil {
		t.Error(err)
	}
	msg := bytes.NewBuffer(msgBytes)
	carrierLen := msgLen / int(ctx.atomSize) * int(ctx.chunkSize)
	carrierBytes := make([]byte, carrierLen)
	if _, err := cryptorand.Read(carrierBytes); err != nil {
		t.Error(err)
	}
	carrier := bytes.NewBuffer(carrierBytes)

	dst := new(bytes.Buffer)
	m := ctx.NewMux(dst, carrier, msg)
	if err := m.Mux(); err != nil {
		t.Error(err)
	}
}

func TestMuxPad(t *testing.T) {
	for i := 0; i < 100; i++ {
		testMuxPad(t, 1)
		testMuxPad(t, 2)
	}
}
