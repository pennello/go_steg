// chris 061715

package main

import (
	"io"
	"log"

	"html/template"
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
	tmpl, err := template.ParseFiles("index.tmpl")
	if err != nil {
		errorResponse(w, 500, err)
		return
	}
	err = tmpl.Execute(w, nil)
	if err != nil {
		errorResponse(w, 500, err)
		return
	}
}
