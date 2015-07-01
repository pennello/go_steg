// chris 061715 Command to run the server.

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

	http.HandleFunc("/", handle)
	log.Fatal(http.ListenAndServe(*addr, nil))
}
