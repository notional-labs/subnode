package rpc

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func StartRpcServer() {
	hostProxy := make(map[string]*httputil.ReverseProxy)

	target, err := url.Parse("https://google.com:443")
	if err != nil {
		panic(err)
	}
	hostProxy["/google"] = httputil.NewSingleHostReverseProxy(target)

	target, err = url.Parse("https://rpc-osmosis-ia.cosmosia.notional.ventures:443")
	if err != nil {
		panic(err)
	}
	hostProxy["/osmosis"] = httputil.NewSingleHostReverseProxy(target)

	handler := func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			fmt.Println("%v", r.RequestURI)

			r.Host = r.URL.Host
			hostProxy[r.RequestURI].ServeHTTP(w, r)
		} else {
			w.Header().Set("Content-Type", "text/html")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Oops! Something was wrong"))
		}
	}

	if err != nil {
		panic(err)
	}

	// handle all requests to your server using the proxy
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
