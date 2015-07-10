// chris 061715

package main

import (
	"io"
	"log"

	"net/http"

	"chrispennello.com/go/steg/cmd"
)

// HTTP handler function initialization common to both local servers as
// well as App Engine.
func init() {
	http.HandleFunc("/", index)
	http.HandleFunc("/api", api)
	http.HandleFunc("/form", form)
}

func errorResponse(w http.ResponseWriter, status int, err error) {
	w.WriteHeader(status)
	_, err = io.WriteString(w, err.Error())
	if err != nil {
		log.Print(err)
	}
}

func api(w http.ResponseWriter, req *http.Request) {
	s, err := parseApi(req)
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

func form(w http.ResponseWriter, req *http.Request) {
	s, err := parseForm(req)
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

func index(w http.ResponseWriter, req *http.Request) {
	http.ServeFile(w, req, "static/html/index.html")
}
