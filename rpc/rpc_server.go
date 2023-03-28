package rpc

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func StartRpcServer() {
	target, err := url.Parse("https://rpc-osmosis-ia.cosmosia.notional.ventures:443")
	if err != nil {
		panic(err)
	}

	proxy := httputil.NewSingleHostReverseProxy(target)

	handler := func(w http.ResponseWriter, r *http.Request) {
		r.Host = r.URL.Host
		proxy.ServeHTTP(w, r)
	}

	if err != nil {
		panic(err)
	}

	// handle all requests to your server using the proxy
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
