package main

import (
	"log"
	"net/http"
)

func StartApiServer() {
	handler := func(w http.ResponseWriter, r *http.Request) {

	}
	// handle all requests to your server using the proxy
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":1337", nil))
}
