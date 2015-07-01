// chris 0622315 HTTP header argument parsing.

package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"strconv"

	"net/http"
	"net/url"

	"chrispennello.com/go/steg"
	"chrispennello.com/go/steg/cmd"
)

// get returns "" if the key is not present in the request header.
func get(h http.Header, key string) string {
	key = fmt.Sprintf("X-Steg-%s", key)
	values, ok := h[key]
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

func parse(req *http.Request) (s *cmd.State, err error) {
	h := req.Header

	atomSizeStr := get(h, "Atom-Size")
	if atomSizeStr == "" {
		atomSizeStr = "1"
	}
	atomSize, err := strconv.ParseInt(atomSizeStr, 0, 8)
	if err != nil {
		return nil, errors.New("invalid atom size value")
	}
	if atomSize < 1 || atomSize > 3 {
		return nil, errors.New("atom size must be 1, 2, or 3")
	}

	var carrier *url.URL
	carrierStr := get(h, "Carrier")
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
	inputStr := get(h, "Input")
	if inputStr == "" {
		// We'll just use the request body.
		input = nil
	} else {
		input, err = parseURL(inputStr)
		if err != nil {
			return nil, err
		}
	}

	boxStr := get(h, "Box")
	if boxStr == "" {
		boxStr = "false"
	}
	box, err := strconv.ParseBool(boxStr)
	if err != nil {
		return nil, errors.New("invalid box value")
	}

	offsetStr := get(h, "Offset")
	if offsetStr == "" {
		offsetStr = "0"
	}
	offset, err := strconv.ParseInt(offsetStr, 0, 64)
	if err != nil {
		return nil, errors.New("invalid offset value")
	}
	if offset < 0 {
		return nil, errors.New("offset must be positive")
	}

	s = new(cmd.State)
	s.Ctx = steg.NewCtx(uint8(atomSize))
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
