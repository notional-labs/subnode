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
				strings.HasPrefix(r.RequestURI, "/check_tx") ||
				strings.HasPrefix(r.RequestURI, "/consensus_state") ||
				strings.HasPrefix(r.RequestURI, "/dump_consensus_state") ||
				strings.HasPrefix(r.RequestURI, "/genesis") ||
				strings.HasPrefix(r.RequestURI, "/genesis_chunked") ||
				strings.HasPrefix(r.RequestURI, "/health") ||
				strings.HasPrefix(r.RequestURI, "/net_info") ||
				strings.HasPrefix(r.RequestURI, "/num_unconfirmed_txs") ||
				strings.HasPrefix(r.RequestURI, "/status") ||
				strings.HasPrefix(r.RequestURI, "/subscribe") ||
				strings.HasPrefix(r.RequestURI, "/unconfirmed_txs") ||
				strings.HasPrefix(r.RequestURI, "/unsubscribe") ||
				strings.HasPrefix(r.RequestURI, "/unsubscribe_all") {
				selectedHost = prunedNode.Backend.Rpc
			} else if strings.HasPrefix(r.RequestURI, "/abci_query") ||
				strings.HasPrefix(r.RequestURI, "/block") ||
				strings.HasPrefix(r.RequestURI, "/block_results") ||
				strings.HasPrefix(r.RequestURI, "/commit") ||
				strings.HasPrefix(r.RequestURI, "/consensus_params") ||
				strings.HasPrefix(r.RequestURI, "/validators") {
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
