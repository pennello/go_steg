// chris 052815

// Steg is a command-line interface to the steganographic embedding
// package steg of which it is a part.
//
// The atom size may be specified as 1, 2, or 3.  The default is 1.
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
// avoid overwriting sensitive headers in the carrier data.  Note that
// specifying an offset effectivly reduces the size of the carrier
// available to embed your message.
//
// Frequently, the data to be embedded will be less than the capacity
// provided by the carrier.  In this case, on extraction, you'll want
// some way to know not to read more than was embedded.  A mechanism for
// this is provided with the box flag.  This will enable the use of a
// simple size-checking encapsulation format.  If you use it on write,
// you'll want to use it on read as well.  Note that using the box flag
// effectively increases the size of your input data.
//
// If you are embedding input data into a carrier with the box flag, and
// the input data is large, you may want to specify a path to the input
// data explicitly instead of using standard in.  Otherwise, steg will
// buffer all of the input data into memory so it can determine the size
// for use with the size-checking encapsulation format.  Of course, if
// the input data is small, then this isn't an issue.
//
// When embedding, steg will check the effective input data size against
// the capacity of the effective carrier size.  If it's insufficient,
// steg will error out early with an informative message.
//
// Options are:
//
//	-atomsize=1: atom size (1, 2, or 3)
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

var state struct {
	ctx         *steg.Ctx
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

func init() {
	atomSizeUsage := "atom size (1, 2, or 3)"
	atomSize := flag.Uint("atomsize", 1, atomSizeUsage)

	carrierUsage := "path to message carrier"
	carrier := flag.String("carrier", "", carrierUsage)

	inputUsage := "path to input; can be - for standard in"
	input := flag.String("input", "-", inputUsage)

	boxUsage := "use size-checking encapsulation format"
	box := flag.Bool("box", false, boxUsage)

	offsetUsage := "read/write offset"
	offset := flag.Int64("offset", 0, offsetUsage)

	flag.Parse()

	if *atomSize < 1 || *atomSize > 3 {
		log.Fatalf("atom size must be 1, 2, or 3")
	}

	state.ctx = steg.NewCtx(uint8(*atomSize))
	state.carrier, state.carrierSize = getCarrier(*carrier)
	state.input, state.inputSize = getInput(*input)
	state.box = *box
	state.offset = *offset

	if state.offset < 0 {
		log.Fatalf("offset must be positive")
	}
}

func extract() {
	var err error
	sr := state.ctx.NewReader(state.input)
	if state.offset != 0 {
		err = sr.Discard(state.offset)
		if err != nil {
			log.Fatalf("extract error: %v", err)
		}
	}
	r := io.Reader(sr)
	if state.box {
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

func mux() {
	var err error
	inputSize := state.inputSize
	carrierSize := state.carrierSize
	message := state.input
	if state.box {
		message = databox.NewMarshaller(state.input, state.inputSize)
		inputSize += databox.HeaderSize
	}
	m := state.ctx.NewMux(os.Stdout, state.carrier, message)
	if state.offset != 0 {
		_, err = m.CopyN(state.offset)
		if err != nil {
			log.Fatalf("mux error: %v", err)
		}
		carrierSize -= state.offset
	}
	capacity := state.ctx.Capacity(carrierSize)
	if capacity < inputSize {
		log.Fatalf("mux error: input size %v > capacity %v", inputSize, capacity)
	}
	err = m.Mux()
	if err != nil {
		log.Fatalf("mux error: %v", err)
	}
}

func main() {
	log.SetFlags(0)
	log.SetPrefix(fmt.Sprintf("%s: ", os.Args[0]))
	if state.carrier == nil {
		extract()
	} else {
		mux()
	}
}
