// chris 052815

// Steg is a command-line interface to the steganographic embedding
// package steg of which it is a part.
package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"chrispennello.com/go/util/databox"

	"../../../steg"
)

type argSpec struct {
	carrier   io.Reader
	input     io.Reader
	inputSize int64
	box       bool
	offset    int64
}

func getCarrier(path string) io.Reader {
	if path == "" {
		return nil
	}
	carrier, err := os.Open(path)
	if err != nil {
		log.Fatalf("failed to open carrier %v", err)
	}
	return carrier
}

func getInput(path string) (input io.Reader, inputSize int64) {
	if path == "-" {
		return os.Stdin, -1
	}
	inputFile, err := os.Open(path)
	if err != nil {
		log.Fatalf("failed to open input %v", err)
	}
	fi, err := inputFile.Stat()
	if err != nil {
		log.Fatalf("failed to stat input %v", err)
	}
	return inputFile, fi.Size()
}

func getArgs() argSpec {
	var args argSpec

	carrierUsage := "path to message carrier"
	carrier := flag.String("carrier", "", carrierUsage)

	inputUsage := "path to input; can be - for standard in"
	input := flag.String("input", "-", inputUsage)

	boxUsage := "use length-checking encapsulation format"
	box := flag.Bool("box", false, boxUsage)

	offsetUsage := "read/write offset"
	offset := flag.Int64("offset", 0, offsetUsage)

	flag.Parse()

	args.carrier = getCarrier(*carrier)
	args.input, args.inputSize = getInput(*input)
	args.box = *box
	args.offset = *offset

	return args
}

func extract(args argSpec) {
	var err error
	sr := steg.NewReader(args.input)
	if args.offset != 0 {
		err = sr.Discard(args.offset)
		if err != nil {
			log.Fatalf("extract error: %v", err)
		}
	}
	r := io.Reader(sr)
	if args.box {
		r = databox.NewUnmarshaller(r)
	}
	_, err = io.Copy(os.Stdout, r)
	if err == steg.ErrShortRead {
		// Short reads are ok on extract.  We just got the end
		// of the file!
		err = nil
	}
	if err != nil {
		log.Fatalf("extract error: %v", err)
	}
}

func mux(args argSpec) {
	var err error
	message := args.input
	if args.box {
		message = databox.NewMarshaller(args.input, args.inputSize)
	}
	m := steg.NewMux(os.Stdout, args.carrier, message)
	if args.offset != 0 {
		_, err = m.CopyN(args.offset)
		if err != nil {
			log.Fatalf("mux error: %v", err)
		}
	}
	err = m.Mux()
	if err != nil {
		log.Fatalf("mux error: %v", err)
	}
}

func main() {
	log.SetFlags(0)
	log.SetPrefix(fmt.Sprintf("%s: ", os.Args[0]))
	args := getArgs()
	if args.carrier == nil {
		extract(args)
	} else {
		mux(args)
	}
}
