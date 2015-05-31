// chris 052815

// Steg is a command-line interface to the steganographic embedding
// package steg of which it is a part.
//
// Input can be provided either as a path, or from the default, standard
// in.
//
// If a carrier is provided, then the input will be interpreted as a
// message to embed within the carrier.  The modified output will be
// written to standard out.  Steg refers to this as "muxing".
//
// Sans carrier, the input will be interpreted as a source from which to
// extract steganographically-embedded data.  The extracted data will be
// written to standard out.
//
// An offset may be specified on both read and write.  The idea is to
// avoid overwriting sensitive headers in the carrier data.
//
// Frequently, the data to be embedded will be less than the capacity
// provided by the carrier.  In this case, on extraction, you'll want
// some way to know not to read more than was embedded.  A mechanism for
// this is provided with the box flag.  This will enable the use of a
// simple size-checking encapsulation format.  If you use it on write,
// you'll want to use it on read as well.
//
// If you are embedding input data into a carrier with the box flag, and
// the input data is large, you may want to specify a path to the input
// data explicitly instead of using standard in.  Otherwise, steg will
// buffer all of the input data into memory so it can determine the size
// for use with the size-checking encapsulation format.  Of course, if
// the input data is small, then this isn't an issue.
//
// Options are:
//
//	-box=false:  use size-checking encapsulation format
//	-carrier="": path to message carrier
//	-input="-":  path to input; can be - for standard in
//	-offset=0:   read/write offset
//
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
	carrier     io.Reader
	carrierSize int64
	input       io.Reader
	inputSize   int64
	box         bool
	offset      int64
}

func getFile(path string) (f *os.File, size int64) {
	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	fi, err := f.Stat()
	if err != nil {
		log.Fatal(err)
	}
	return f, fi.Size()
}

func getCarrier(path string) (carrier io.Reader, size int64) {
	if path == "" {
		return nil, -2
	}
	return getFile(path)
}

func getInput(path string) (input io.Reader, size int64) {
	if path == "-" {
		return os.Stdin, -1
	}
	return getFile(path)
}

func getArgs() argSpec {
	var args argSpec

	carrierUsage := "path to message carrier"
	carrier := flag.String("carrier", "", carrierUsage)

	inputUsage := "path to input; can be - for standard in"
	input := flag.String("input", "-", inputUsage)

	boxUsage := "use size-checking encapsulation format"
	box := flag.Bool("box", false, boxUsage)

	offsetUsage := "read/write offset"
	offset := flag.Int64("offset", 0, offsetUsage)

	flag.Parse()

	args.carrier, args.carrierSize = getCarrier(*carrier)
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
