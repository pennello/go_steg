// chris 052815

package steg

import "io"

// Mux reads, one byte at a time, from the message reader,
// steganographically embeds its data into the data read from the
// carrier, and then writes the resultant data into the destination
// writer.  If the carrier has more than enough data for the message,
// the rest of the carrier data is simply copied to the writer.
//
// Can return ErrShortCarrier if an EOF was encountered before being
// able to read a sufficient number amount of data from the carrier for
// the message.  Can return other errors as well encountered during the
// writes.
//
// Can return ErrInsufficientData if the reader does not contain
// sufficient data to read an integral number of atoms.
//
// Successful iff err != nil.
func (m Mux) Mux() (err error) {
	// Would be nice if we could just call io.Copy, but the buffer
	// size isn't ensured to be a multiple of the atom size.

	// XXX Just reading one atom at a time for simplicity.  The
	// caller can certainly wrap its readers and writers in buffered
	// versions, but this is still internally a lot of function
	// calls for a reasonable amount of data...

	var n int
	a := m.ctx.newAtom()
	for {
		n, err = io.ReadFull(m.msg, a.data)
		if err != nil {
			if err == io.EOF {
				if n == 0 {
					break
				}
				return ErrInsufficientData
			}
			return err
		}
		_, err = m.w.Write(a.data)
		if err != nil {
			return err
		}
	}
	_, err = m.w.Copy()
	return err
}

// CopyN copies n bytes from the carrier to the destination.
//
// The idea is that you'd call this before sending through any of your
// message data to get past critical headers in your carrier before
// embedding your data.
//
// Counterpart to Reader.Discard.
func (m Mux) CopyN(n int64) (written int64, err error) {
	return m.w.CopyN(n)
}
