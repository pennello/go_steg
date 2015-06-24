// chris 061715

// +build !appengine

package main

import (
	"flag"
	"log"
	"os"

	"net/http"
)

func main() {
	addr := flag.String("http", "", "host:port address on which to listen")
	flag.Parse()

	if *addr == "" {
		flag.Usage()
		os.Exit(2)
	}

	http.HandleFunc("/mux", mux)
	http.HandleFunc("/read", read)
	log.Fatal(http.ListenAndServe(*addr, nil))
}
