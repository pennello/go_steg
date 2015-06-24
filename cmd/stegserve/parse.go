// chris 0622315 HTTP header argument parsing.

package main

import (
	"errors"
	"strconv"
	"strings"

	"net/http"
	"net/url"
)

type commonArgs struct {
	atomSize int8
	box      bool
	offset   int64
}

type readArgs struct {
	input *url.URL
}

type muxArgs struct {
	carrier *url.URL
	message *url.URL
}

func parseArgs(h http.Header) map[string]string {
	r := make(map[string]string)
	for k, v := range h {
		if strings.HasPrefix(k, "X-Steg-") {
			r[strings.ToLower(k[7:])] = v[0]
		}
	}
	return r
}

func parseCommon(args map[string]string) (*commonArgs, error) {
	atomSizeStr, ok := args["atom-size"]
	if !ok {
		atomSizeStr = "1"
	}
	atomSize, err := strconv.ParseInt(atomSizeStr, 0, 8)
	if err != nil {
		return nil, errors.New("invalid atom size value")
	}
	if atomSize < 1 || atomSize > 3 {
		return nil, errors.New("atom size must be 1, 2, or 3")
	}

	boxStr, ok := args["box"]
	if !ok {
		boxStr = "false"
	}
	box, err := strconv.ParseBool(boxStr)
	if err != nil {
		return nil, errors.New("invalid box value")
	}

	offsetStr, ok := args["offset"]
	if !ok {
		offsetStr = "0"
	}
	offset, err := strconv.ParseInt(offsetStr, 0, 64)
	if err != nil {
		return nil, errors.New("invalid offset value")
	}
	if offset < 0 {
		return nil, errors.New("offset must be positive")
	}

	return &commonArgs{atomSize: int8(atomSize), box: box, offset: offset}, nil
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

func parseRead(args map[string]string) (*readArgs, error) {
	inputStr, ok := args["input"]
	if !ok {
		return nil, errors.New("input required")
	}
	input, err := parseURL(inputStr)
	if err != nil {
		return nil, err
	}

	return &readArgs{input: input}, nil
}

func parseMux(args map[string]string) (*muxArgs, error) {
	carrierStr, ok := args["carrier"]
	if !ok {
		return nil, errors.New("carrier required")
	}
	carrier, err := parseURL(carrierStr)
	if err != nil {
		return nil, err
	}

	messageStr, ok := args["message"]
	if !ok {
		return nil, errors.New("message required")
	}
	message, err := parseURL(messageStr)
	if err != nil {
		return nil, err
	}

	return &muxArgs{carrier: carrier, message: message}, nil
}
