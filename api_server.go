package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
)

func StartApiServer() {
	handler := func(w http.ResponseWriter, r *http.Request) {
		prunedNode := SelectPrunedNodeApi()
		selectedHost := prunedNode.Backend.Api // default to pruned node

		if r.Method == "GET" {
			fmt.Printf("r.RequestURI=%s\n", r.RequestURI)

			xCosmosBlockHeight := r.Header.Get("x-cosmos-block-height")
			if xCosmosBlockHeight != "" {
				height, err := strconv.ParseInt(xCosmosBlockHeight, 10, 64)
				if err != nil {
					SendError(w)
				}

				node, err := SelectMatchedNodeApi(height)
				if err != nil {
					SendError(w)
				}

				selectedHost = node.Backend.Api
			} else {
				selectedHost = prunedNode.Backend.Api
			}

			r.Host = r.URL.Host
			ProxyMapApi[selectedHost].ServeHTTP(w, r)
		} else if r.Method == "POST" {
			selectedHost = prunedNode.Backend.Api
			r.Host = r.URL.Host
			ProxyMapApi[selectedHost].ServeHTTP(w, r)
		} else {
			SendError(w)
		}
	}
	// handle all requests to your server using the proxy
	//http.HandleFunc("/", handler)
	serverMux := http.NewServeMux()
	serverMux.HandleFunc("/", handler)
	go func() {
		log.Fatal(http.ListenAndServe(":1337", serverMux))
	}()
}
