// chris 0622315 API parameter parsing.

package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"mime"
	"strconv"

	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"

	"chrispennello.com/go/steg"
	"chrispennello.com/go/steg/cmd"
)

// getHeader returns "" if the key is not present in the request header.
func getHeader(req *http.Request, key string) string {
	key = fmt.Sprintf("X-Steg-%s", key)
	values, ok := req.Header[key]
	if !ok {
		return ""
	}
	return values[0]
}

func parseURL(rawurl string) (u *url.URL, err error) {
	// http://stackoverflow.com/a/417184
	if len(rawurl) > 2048 {
		return nil, errors.New("url too long")
	}
	u, err = url.Parse(rawurl)
	if err != nil {
		return nil, errors.New("invalid url")
	}
	if u.Scheme != "http" {
		return nil, errors.New("http urls only")
	}
	return u, nil
}

// Note that size can be -1 if the content length is unknown.
func getURL(u *url.URL) (body io.ReadCloser, size int64, err error) {
	resp, err := http.Get(u.String())
	if err != nil {
		return nil, -2, err
	}
	return resp.Body, resp.ContentLength, nil
}

func getCarrier(u *url.URL) (carrier io.ReadCloser, size int64, err error) {
	if u == nil {
		return nil, -2, nil
	}
	return getURL(u)
}

func getInput(req *http.Request, u *url.URL) (input io.ReadCloser, size int64, err error) {
	if u == nil {
		return req.Body, req.ContentLength, nil
	}
	return getURL(u)
}

func parseAtomSize(atomSizeStr string) (uint8, error) {
	atomSize, err := strconv.ParseInt(atomSizeStr, 0, 8)
	if err != nil {
		return 0, errors.New("invalid atom size value")
	}
	if atomSize < 1 || atomSize > 3 {
		return 0, errors.New("atom size must be 1, 2, or 3")
	}
	return uint8(atomSize), nil
}

func parseBox(boxStr string) (bool, error) {
	box, err := strconv.ParseBool(boxStr)
	if err != nil {
		return false, errors.New("invalid box value")
	}
	return box, nil
}

func parseOffset(offsetStr string) (int64, error) {
	offset, err := strconv.ParseInt(offsetStr, 0, 64)
	if err != nil {
		return -2, errors.New("invalid offset value")
	}
	if offset < 0 {
		return -2, errors.New("offset must be positive")
	}
	return offset, nil
}

func parseApi(req *http.Request) (s *cmd.State, err error) {
	atomSizeStr := getHeader(req, "Atom-Size")
	if atomSizeStr == "" {
		atomSizeStr = "1"
	}
	atomSize, err := parseAtomSize(atomSizeStr)
	if err != nil {
		return nil, err
	}

	var carrier *url.URL
	carrierStr := getHeader(req, "Carrier")
	if carrierStr == "" {
		// No muxing.
		carrier = nil
	} else {
		carrier, err = parseURL(carrierStr)
		if err != nil {
			return nil, err
		}
	}

	var input *url.URL
	inputStr := getHeader(req, "Input")
	if inputStr == "" {
		// We'll just use the request body.
		input = nil
	} else {
		input, err = parseURL(inputStr)
		if err != nil {
			return nil, err
		}
	}

	boxStr := getHeader(req, "Box")
	if boxStr == "" {
		boxStr = "false"
	}
	box, err := parseBox(boxStr)
	if err != nil {
		return nil, err
	}

	offsetStr := getHeader(req, "Offset")
	if offsetStr == "" {
		offsetStr = "0"
	}
	offset, err := parseOffset(offsetStr)
	if err != nil {
		return nil, err
	}

	s = new(cmd.State)
	s.Ctx = steg.NewCtx(atomSize)
	s.Carrier, s.CarrierSize, err = getCarrier(carrier)
	if err != nil {
		return nil, err
	}
	s.Input, s.InputSize, err = getInput(req, input)
	if err != nil {
		err2 := s.Carrier.Close()
		if err2 != nil {
			log.Print(err2)
		}
		return nil, err
	}
	s.Box = box
	s.Offset = offset

	return s, nil
}

type bytesBufferCloser struct {
	*bytes.Buffer
}

func (bytesBufferCloser) Close() error {
	return nil
}

// parsePart parses a multipart mime part for an input or carrier form
// entry.  File uploads will be buffered entirely into memory.
func parsePart(part *multipart.Part) (u *url.URL, rc io.ReadCloser, err error) {
	filename := part.FileName()
	if filename != "" {
		// We assume entire file.
		// Parts are being read from multipart mime, so we
		// buffer the entire thing into memory. :[

		buf := new(bytes.Buffer)
		_, err := buf.ReadFrom(part)
		if err != nil {
			return nil, nil, err
		}
		return nil, bytesBufferCloser{buf}, nil
	}
	// We assume URL.
	inputBytes, err := ioutil.ReadAll(part)
	if err != nil {
		return nil, nil, err
	}
	rawurl := string(inputBytes)
	if rawurl == "" {
		// Form may have just been
		// submited with no URL, but
		// with file input in another
		// field of the same name.
		// We'll check for missing
		// values after we go through
		// all the parts.
		return nil, nil, nil
	}
	u, err = parseURL(rawurl)
	if err != nil {
		return nil, nil, err
	}
	return u, nil, nil
}

func parseForm(req *http.Request) (s *cmd.State, err error) {
	contenttype, ok := req.Header["Content-Type"]
	if !ok {
		return nil, errors.New("content-type required")
	}
	mediatype, params, err := mime.ParseMediaType(contenttype[0])
	if err != nil {
		return nil, err
	}
	if mediatype != "multipart/form-data" {
		return nil, errors.New("multipart/form-data required")
	}
	boundary, ok := params["boundary"]
	if !ok {
		return nil, errors.New("form boundary required")
	}

	// XXX This can behave poorly if the input is weird.  For
	// example, if the input provides two parts with carrier files,
	// one will be orphaned and never closed.

	mpr := multipart.NewReader(req.Body, boundary)
	s = new(cmd.State)

	var atomSize uint8
	var carrier *url.URL
	var carrierReader io.ReadCloser
	var input *url.URL
	var inputReader io.ReadCloser

	for {
		part, err := mpr.NextPart()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}

		switch part.FormName() {

		// This includes the empty string, returned when there
		// is no form name present.
		default:
			continue

		case "atom-size":
			{
				atomSizeBytes, err := ioutil.ReadAll(part)
				if err != nil {
					return nil, err
				}
				// No default handling--form should provide the
				// default of 1.
				atomSize, err = parseAtomSize(string(atomSizeBytes))
				if err != nil {
					return nil, err
				}
			}

		case "carrier":
			{
				carrier, carrierReader, err = parsePart(part)
				if err != nil {
					return nil, err
				}
			}

		case "input":
			{
				input, inputReader, err = parsePart(part)
				if err != nil {
					return nil, err
				}
			}

		case "box":
			{
				boxBytes, err := ioutil.ReadAll(part)
				if err != nil {
					return nil, err
				}
				boxStr := string(boxBytes)
				if boxStr == "on" {
					boxStr = "1"
				}
				box, err := parseBox(boxStr)
				if err != nil {
					return nil, err
				}
				s.Box = box
			}

		case "offset":
			{
				offsetBytes, err := ioutil.ReadAll(part)
				if err != nil {
					return nil, err
				}
				offset, err := parseOffset(string(offsetBytes))
				if err != nil {
					return nil, err
				}
				s.Offset = offset
			}
		}
	}

	if atomSize == 0 {
		return nil, errors.New("atom-size required")
	}
	s.Ctx = steg.NewCtx(atomSize)

	if carrierReader != nil {
		s.Carrier, s.CarrierSize = carrierReader, -1
	} else if carrier != nil {
		s.Carrier, s.CarrierSize, err = getURL(carrier)
		if err != nil {
			return nil, err
		}
	} else {
		// No carrier, no problem--we just won't do any muxing.
	}

	if inputReader != nil {
		s.Input, s.InputSize = inputReader, -1
	} else if input != nil {
		s.Input, s.InputSize, err = getURL(input)
		if err != nil {
			err2 := s.Carrier.Close()
			if err2 != nil {
				log.Print(err2)
			}
			return nil, err
		}
	} else {
		return nil, errors.New("input required")
	}

	return s, nil
}
