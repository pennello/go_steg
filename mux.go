// chris 052815

package steg

import "io"

// Mux multiplexes a message on a carrier into a destination.  It
// steganographically embeds data from the message into the carrier and
// then writes the modified data into the destination.
type Mux struct {
	w   Writer
	msg io.Reader
}

// NewMux returns a fresh Mux, ready to multiplex a message on a carrier
// into a destination.
func NewMux(dst io.Writer, carrier, msg io.Reader) Mux {
	w := NewWriter(dst, carrier)
	return Mux{w: w, msg: msg}
}

// Mux reads, one byte at a time, from the message reader,
// steganographically embeds its data into the data read from the
// carrier, and then writes the resultant data into the destination
// writer.  If the carrier has more than enough data for the message,
// the rest of the carrier data is simply copied to the writer.
//
// Can return ErrShortCarrier if an EOF was encountered before being
// able to read a sufficient number amount of data from the carrier for
// the message.  Can return other errors as well encountered during the
// writes.  Successful iff err != nil.
func (m Mux) Mux() (err error) {
	// XXX Just reading one byte at a time for simplicity.  The
	// caller can certainly wrap its readers and writers in buffered
	// versions, but this is still internally a lot of function
	// calls for a reasonable amount of data...
	b := make([]byte, 1)
	for {
		_, err = io.ReadFull(m.msg, b)
		if err != nil {
			if err != io.EOF {
				return err
			}
			break
		}
		_, err = m.w.Write(b)
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
