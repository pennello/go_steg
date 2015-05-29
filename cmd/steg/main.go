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

func mux(carrier io.Reader) error {
	return steg.Mux(os.Stdout, carrier, os.Stdin)
}

func read() error {
	_, err := io.Copy(os.Stdout, steg.NewReader(os.Stdin))
	return err
}

func main() {
	log.SetFlags(0)
	log.SetPrefix(fmt.Sprintf("%s: ", os.Args[0]))
	flag.Parse()
	a := flag.Args()

	var carrier io.Reader
	var err error
	if len(a) == 1 {
		carrier, err = os.Open(a[0])
		if err != nil {
			log.Fatalf("failed to open carrier %v", err)
		}
	} else {
		carrier = nil
	}

	var errlabel string
	if carrier == nil {
		err = read()
		if err == steg.ErrShortRead {
			// Short reads are ok.  We just got the end of
			// the file!
			err = nil
		}
		errlabel = "read"
	} else {
		err = mux(carrier)
		errlabel = "mux"
	}

	if err != nil {
		log.Fatalf("%s error: %v", errlabel, err)
	}
}
