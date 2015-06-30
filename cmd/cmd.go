// chris 062915

// Package cmd implements common code for steg commands.
package cmd

import (
	"fmt"
	"io"

	"chrispennello.com/go/steg"
	"chrispennello.com/go/util/databox"
)

// State represents the state of your command.  Fill it in by parsing
// arguments, input, etc., and then pass it to Main to execute the
// command.  If there is no carrier (i.e., you're extracting), then set
// it to nil; CarrierSize is not inspected in this case.  InputSize may
// be set to -1 to indicate that the input is being streamed.  Note that
// in this case, if Box is being used, then all of the input data will
// be read into memory by the databox library.
type State struct {
	Ctx         *steg.Ctx
	Carrier     io.Reader
	CarrierSize int64
	Input       io.Reader
	InputSize   int64
	Box         bool
	Offset      int64
}

func extract(dst io.Writer, s *State) error {
	var err error
	sr := s.Ctx.NewReader(s.Input)
	if s.Offset != 0 {
		err = sr.Discard(s.Offset)
		if err != nil {
			return fmt.Errorf("extract error: %v", err)
		}
	}
	r := io.Reader(sr)
	if s.Box {
		r = databox.NewUnmarshaller(r)
	}
	_, err = io.Copy(dst, r)
	if err == steg.ErrShortRead {
		// Short reads are ok on extract.  We just got to the
		// end of the file!
		err = nil
	}
	if err != nil {
		return fmt.Errorf("extract error: %v", err)
	}
	return nil
}

func mux(dst io.Writer, s *State) error {
	stream := s.InputSize == -1
	inputSize := s.InputSize
	carrierSize := s.CarrierSize
	message := s.Input
	if s.Box {
		message = databox.NewMarshaller(s.Input, s.InputSize)
		inputSize += databox.HeaderSize
	}
	m := s.Ctx.NewMux(dst, s.Carrier, message)
	if s.Offset != 0 {
		_, err := m.CopyN(s.Offset)
		if err != nil {
			return fmt.Errorf("mux error: %v", err)
		}
		carrierSize -= s.Offset
	}
	capacity := s.Ctx.Capacity(carrierSize)
	if !stream && capacity < inputSize {
		return fmt.Errorf("mux error: input size %v > capacity %v", inputSize, capacity)
	}
	err := m.Mux()
	if err != nil {
		return fmt.Errorf("mux error: %v", err)
	}
	return nil
}

// Main is the entry point for common command logic.  Pass in a
// destination writer and a pointer to a state struct you've prepared.
func Main(dst io.Writer, s *State) error {
	if s.Carrier == nil {
		return extract(dst, s)
	} else {
		return mux(dst, s)
	}
}
