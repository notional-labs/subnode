package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func StartRpcServer() {
	handler := func(w http.ResponseWriter, r *http.Request) {
		prunedNode := SelectPrunedNode(ProtocolTypeRpc)
		selectedHost := prunedNode.Backend.Rpc // default to pruned node

		if r.Method == "GET" { // URI over HTTP
			// see `/doc/rpc.md` to see the logic

			fmt.Printf("r.RequestURI=%s\n", r.RequestURI)

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
						return
					}

					node, err := SelectMatchedBackend(height, ProtocolTypeRpc)
					if err != nil {
						SendError(w)
						return
					}

					selectedHost = node.Backend.Rpc
				} else {
					selectedHost = prunedNode.Backend.Rpc
				}
			}

			r.Host = r.URL.Host
			ProxyMapRpc[selectedHost].ServeHTTP(w, r)
		} else if r.Method == "POST" { // JSONRPC over HTTP
			body, err := io.ReadAll(r.Body)
			if err != nil {
				SendError(w)
				return
			}

			fmt.Printf("body=%s\n", string(body))

			var j0 interface{}
			err = json.Unmarshal(body, &j0)
			if err != nil {
				SendError(w)
				return
			}

			m0 := j0.(map[string]interface{})
			method := m0["method"].(string)
			//params := m0["params"].([]interface{})

			fmt.Printf("method=%s, params=%+v\n", method, m0["params"])

			// note that params could be positional parameters or named parameters
			// eg.,  "params": [ "9045128", "1", "30" ]
			//   or  "params": { "height": "9045128", "page": "1", "per_page": "30" }

			if method == "abci_info" ||
				strings.HasPrefix(method, "broadcast_") ||
				method == "check_tx" ||
				method == "consensus_state" ||
				method == "dump_consensus_state" ||
				method == "genesis" ||
				method == "genesis_chunked" ||
				method == "health" ||
				method == "net_info" ||
				method == "num_unconfirmed_txs" ||
				method == "status" ||
				method == "subscribe" ||
				method == "unconfirmed_txs" ||
				method == "unsubscribe" ||
				method == "unsubscribe_all" {
				selectedHost = prunedNode.Backend.Rpc
			} else if method == "block" ||
				method == "block_results" ||
				method == "commit" ||
				method == "consensus_params" ||
				method == "validators" {

				height := int64(-1)

				if positionalParams, ok := m0["params"].([]interface{}); ok { // positional parameters
					// height is 1st param
					if len(positionalParams) < 1 {
						SendError(w)
						return
					}

					heightParam := positionalParams[0].(string)
					height, err = strconv.ParseInt(heightParam, 10, 64)
					if err != nil {
						SendError(w)
						return
					}
				} else if namedParams, ok := m0["params"].(map[string]interface{}); ok { // named parameters
					heightParam := namedParams["height"].(string)
					height, err = strconv.ParseInt(heightParam, 10, 64)
					if err != nil {
						SendError(w)
						return
					}
				}

				if height >= 0 {
					node, err := SelectMatchedBackend(height, ProtocolTypeRpc)
					if err != nil {
						SendError(w)
						return
					}

					selectedHost = node.Backend.Rpc
				}
			} else if method == "abci_query" {
				height := int64(-1)

				if positionalParams, ok := m0["params"].([]interface{}); ok { // positional parameters
					// height is 3rd param
					if len(positionalParams) < 3 {
						SendError(w)
						return
					}

					heightParam := positionalParams[2].(string)
					height, err = strconv.ParseInt(heightParam, 10, 64)
					if err != nil {
						SendError(w)
						return
					}
				} else if namedParams, ok := m0["params"].(map[string]interface{}); ok { // named parameters
					if heightParam, ok := namedParams["height"].(string); ok {
						height, err = strconv.ParseInt(heightParam, 10, 64)
						if err != nil {
							SendError(w)
							return
						}
					}
				}

				if height >= 0 {
					node, err := SelectMatchedBackend(height, ProtocolTypeRpc)
					if err != nil {
						SendError(w)
						return
					}
					selectedHost = node.Backend.Rpc
				}
			}

			r.Body = io.NopCloser(bytes.NewBuffer(body)) // assign a new body with previous byte slice
			r.Host = r.URL.Host
			ProxyMapRpc[selectedHost].ServeHTTP(w, r)
		} else {
			SendError(w)
			return
		}
	}

	// handle all requests to your server using the proxy
	//http.HandleFunc("/", handler)
	serverMux := http.NewServeMux()
	serverMux.HandleFunc("/", handler)
	go func() {
		log.Fatal(http.ListenAndServe(":26657", serverMux))
	}()
}

func SendError(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte("Oops! Something was wrong"))
}
