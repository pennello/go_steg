// chris 061715

// +build !appengine

// Stegserve runs a local HTTP server serving up several endpoints
// providing access to the steganographic embedding package steg of
// which it is a part.  The endpoints are as follows.
//
//	/	Simple web GUI.
//	/api	Header-based API.
//	/mime	Multipart MIME-based API, used by the web GUI.
//		Buffers file uploads into memory.
//
// /api takes the following header arguments.  See the GoDoc
// documentation of the steg command for a fuller explanation of these
// arguments.
//
//	X-Steg-Atom-Size	defaults to 1; can be 1, 2, or 3
//	X-Steg-Box		defaults to false;
//				accepts values recognized by
//				strconv.ParseBool
//	X-Steg-Carrier		optional; valid URL
//	X-Steg-Input		defaults to use the request body;
//				valid URL
//	X-Steg-Offset		defaults to 0; read/write offset
//
// /mime takes the following form-data arguments.  See the GoDoc
// documentation of the steg command for a fuller explanation of these
// arguments.
//
//	atom-size	required; can be 1, 2, or 3
//	box		defaults to false; "on" for true,
//			also accepts values recognized by
//			strconv.ParseBool
//	carrier		optional; valid URL or file upload
//	input		required; valid URL or file upload
//	offset		defaults to 0; read/write offset
//
// This command provides a demonstration of the sort of network
// proxying interface one might implement to provide remote
// steganographic services.  Given the character of steganographic
// embedding, a more practical implementation would go to greater
// lengths to obscure the purpose of the endpoint.  For example, if you
// were to steganographically embed an advertisement in a video, you
// might hit and endpoint and provide it with just the IDs of the video
// and the advertisement.
//
// Options are:
//
//	-http="": host:port address on which to listen
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

	// HTTP handler functions initialized in hande.go.

	log.Fatal(http.ListenAndServe(*addr, nil))
}
