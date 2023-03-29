package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"strings"
)

func StartRpcServer() {
	hostProxy := make(map[string]*httputil.ReverseProxy)

	InitPool()

	for _, s := range Pool {
		target, err := url.Parse(s.Name)
		if err != nil {
			panic(err)
		}
		hostProxy[s.Name] = httputil.NewSingleHostReverseProxy(target)
	}

	handler := func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			// see `/doc/rpc.md` to see the logic

			fmt.Printf("r.RequestURI=%s\n", r.RequestURI)

			prunedNode := SelectPrunedNode()
			selectedHost := prunedNode.Backend.Rpc // default to pruned node

			if strings.HasPrefix(r.RequestURI, "/abci_info") ||
				strings.HasPrefix(r.RequestURI, "/broadcast_") ||
				strings.HasPrefix(r.RequestURI, "/check_tx") {
				selectedHost = prunedNode.Backend.Rpc
			} else if strings.HasPrefix(r.RequestURI, "/abci_query") ||
				strings.HasPrefix(r.RequestURI, "/block") ||
				strings.HasPrefix(r.RequestURI, "/block_results") ||
				strings.HasPrefix(r.RequestURI, "/commit") {
				heightParam := r.URL.Query().Get("height")
				if heightParam != "" {
					height, err := strconv.ParseInt(heightParam, 10, 64)
					if err != nil {
						SendError(w)
					}

					node, err := SelectMatchedNode(height)
					if err != nil {
						SendError(w)
					}

					selectedHost = node.Backend.Rpc
				} else {
					selectedHost = prunedNode.Backend.Rpc
				}
			}

			r.Host = r.URL.Host
			hostProxy[selectedHost].ServeHTTP(w, r)
		} else {
			SendError(w)
		}
	}

	// handle all requests to your server using the proxy
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func SendError(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte("Oops! Something was wrong"))
}
