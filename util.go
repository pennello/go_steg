// chris 061115

package steg

// Capacity returns the largest message a carrier of the given size can
// embed, in bytes.
func (ctx *Ctx) Capacity(carrierSize int64) (messageSize int64) {
	return int64(uint(carrierSize) / ctx.chunkSize * ctx.atomSize)
}
