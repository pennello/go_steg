// chris 061715

package main

import (
	"io"
	"log"

	"net/http"

	"chrispennello.com/go/steg/cmd"
)

func errorResponse(w http.ResponseWriter, status int, err error) {
	w.WriteHeader(status)
	_, err = io.WriteString(w, err.Error())
	if err != nil {
		log.Print(err)
	}
}

func handle(w http.ResponseWriter, req *http.Request) {
	s, err := parse(req)
	if err != nil {
		errorResponse(w, 400, err)
		return
	}
	err = cmd.Main(w, s)
	if err != nil {
		errorResponse(w, 500, err)
		return
	}
}
