// chris 061715

package main

import (
	"fmt"
	"io"
	"log"

	"net/http"
)

func errorResponse(w http.ResponseWriter, status int, err error) {
	w.WriteHeader(status)
	_, err = io.WriteString(w, err.Error())
	if err != nil {
		log.Print(err)
	}
}

func mux(w http.ResponseWriter, r *http.Request) {
	a := parseArgs(r.Header)
	ca, err := parseCommon(a)
	if err != nil {
		errorResponse(w, 400, err)
		return
	}
	ma, err := parseMux(a)
	if err != nil {
		errorResponse(w, 400, err)
		return
	}

	_, err = fmt.Fprintf(w, "%#v %#v", ca, ma)
	if err != nil {
		log.Print(err)
	}
}

func read(w http.ResponseWriter, r *http.Request) {
	a := parseArgs(r.Header)
	ca, err := parseCommon(a)
	if err != nil {
		errorResponse(w, 400, err)
		return
	}
	ra, err := parseRead(a)
	if err != nil {
		errorResponse(w, 400, err)
		return
	}

	_, err = fmt.Fprintf(w, "%#v %#v", ca, ra)
	if err != nil {
		log.Print(err)
	}
}
