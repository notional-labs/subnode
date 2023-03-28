package rpc

import (
	"fmt"
	"github.com/notional-labs/subnode/cmd"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

func StartRpcServer() {
	hostProxy := make(map[string]*httputil.ReverseProxy)

	cfg := cmd.GetConfig()
	for _, s := range cfg.Upstream {
		target, err := url.Parse(s.Rpc)
		if err != nil {
			panic(err)
		}
		hostProxy[s.Rpc] = httputil.NewSingleHostReverseProxy(target)
	}

	handler := func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			fmt.Print("r.RequestURI=%s\n", r.RequestURI)

			// see `/doc/rpc.md`
			if strings.HasPrefix(r.RequestURI, "/abci_info") {

			}

			r.Host = r.URL.Host
			hostProxy[r.RequestURI].ServeHTTP(w, r)
		} else {
			w.Header().Set("Content-Type", "text/html")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Oops! Something was wrong"))
		}
	}

	// handle all requests to your server using the proxy
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
