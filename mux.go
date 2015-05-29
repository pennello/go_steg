// chris 052815

package steg

import "io"

// Mux reads, one byte at a time, from the message reader,
// steganographically embeds its data into the data read from the
// carrier, and then writes the resultant data into the destination
// writer.  If the carrier has more than enough data for the message,
// the rest of the carrier data is simply copied to the writer.
//
// Can return io.ErrUnexpectedEOF if an EOF was encountered before being
// able to read a sufficient number amount of data from the carrier for
// the message.  Can also return other errors as well encountered during
// the writes.  Successful iff err != nil.
func Mux(dst io.Writer, carrier io.Reader, msg io.Reader) (err error) {
	w := NewWriter(dst, carrier)
	// XXX Just reading one byte at a time for simplicity.  The
	// caller can certainly wrap its readers and writers in buffered
	// versions, but this is still internally a lot of function
	// calls for a reasonable amount of data...
	b := make([]byte, 1)
	for {
		_, err = io.ReadFull(msg, b)
		if err != nil {
			if err != io.EOF {
				return err
			}
			break
		}
		_, err = w.Write(b)
		if err != nil {
			return err
		}
	}
	_, err = w.Copy()
	return err
}
