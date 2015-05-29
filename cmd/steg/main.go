// chris 052815

// Steg is a command-line interface to the steganographic embedding
// package steg of which it is a part.
//
// With no arguments, steg reads from standard in, extracts
// steganographically-embedded data, and writes it to standard out.
//
// If a path to carrier data is provided as the only positional
// command-line argument, then steg will read a message from standard
// in, embed it in the carrier data, and write the modified data to
// standard out.
package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"../../../steg"
)

func mux(carrier io.Reader, offset int64) (err error) {
	m := steg.NewMux(os.Stdout, carrier, os.Stdin)
	if offset != 0 {
		_, err = m.CopyN(offset)
		if err != nil {
			return err
		}
	}
	return m.Mux()
}

func extract(offset int64) (err error) {
	r := steg.NewReader(os.Stdin)
	if offset != 0 {
		err = r.Discard(offset)
		if err != nil {
			return err
		}
	}
	_, err = io.Copy(os.Stdout, r)
	return err
}

func main() {
	log.SetFlags(0)
	log.SetPrefix(fmt.Sprintf("%s: ", os.Args[0]))
	offset := flag.Int64("offset", 0, "when writing, carrier data "+
		"offset after which to embed message data;\n"+
		"             when reading, offset after which to start "+
		"reading")
	flag.Parse()
	a := flag.Args()

	var carrier io.Reader
	var err error
	if len(a) != 1 {
		carrier = nil
	} else {
		carrier, err = os.Open(a[0])
		if err != nil {
			log.Fatalf("failed to open carrier %v", err)
		}
	}

	var errlabel string
	if carrier == nil {
		err = extract(*offset)
		if err == steg.ErrShortRead {
			// Short reads are ok.  We just got the end of
			// the file!
			err = nil
		}
		errlabel = "extract"
	} else {
		err = mux(carrier, *offset)
		errlabel = "mux"
	}

	if err != nil {
		log.Fatalf("%s error: %v", errlabel, err)
	}
}
