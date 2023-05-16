package internal

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
)

var rpcServer *http.Server

func uriOverHttp(w http.ResponseWriter, r *http.Request) {
	prunedNode := SelectPrunedNode(ProtocolTypeRpc)
	selectedHost := prunedNode.Backend.Rpc // default to pruned node

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
		//strings.HasPrefix(r.RequestURI, "/status") ||
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
	} else if strings.HasPrefix(r.RequestURI, "/blockchain") { // base on maxHeight
		heightParam := r.URL.Query().Get("maxHeight")
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
	} else { // try to support partially for other methods
		strQuery := r.URL.Query().Encode()
		//fmt.Printf("query=%s", strQuery)

		if strings.HasPrefix(r.RequestURI, "/status") {
			DoAggeratorUriOverHttp_status(w, strQuery)
			return
		} else if strings.HasPrefix(r.RequestURI, "/tx_search") {
			DoAggeratorUriOverHttp_tx_search(w, strQuery)
			return
		} else if strings.HasPrefix(r.RequestURI, "/block_by_hash") {
			DoAggeratorUriOverHttp_block_by_hash(w, strQuery)
			return
		} else if strings.HasPrefix(r.RequestURI, "/tx") {
			DoAggeratorUriOverHttp_tx(w, strQuery)
			return
		} else if strings.HasPrefix(r.RequestURI, "/block_search") {
			DoAggeratorUriOverHttp_block_search(w, strQuery)
			return
		}
	}

	r.Host = r.URL.Host
	ProxyMapRpc[selectedHost].ServeHTTP(w, r)
}

func jsonRpcOverHttp(w http.ResponseWriter, r *http.Request) {
	prunedNode := SelectPrunedNode(ProtocolTypeRpc)
	selectedHost := prunedNode.Backend.Rpc // default to pruned node

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

	if m0, ok := j0.(map[string]interface{}); ok {
		if method, ok := m0["method"].(string); ok {
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

					if heightParam, ok := positionalParams[0].(string); ok {
						height, err = strconv.ParseInt(heightParam, 10, 64)
						if err != nil {
							SendError(w)
							return
						}
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
			} else if method == "blockchain" { // base on maxHeight
				height := int64(-1)

				if positionalParams, ok := m0["params"].([]interface{}); ok { // positional parameters
					// maxHeight is 2nd param
					if len(positionalParams) < 2 {
						SendError(w)
						return
					}

					if heightParam, ok := positionalParams[1].(string); ok {
						height, err = strconv.ParseInt(heightParam, 10, 64)
						if err != nil {
							SendError(w)
							return
						}
					}
				} else if namedParams, ok := m0["params"].(map[string]interface{}); ok { // named parameters
					if heightParam, ok := namedParams["maxHeight"].(string); ok {
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
			} else if method == "abci_query" {
				height := int64(-1)

				if positionalParams, ok := m0["params"].([]interface{}); ok { // positional parameters
					// height is 3rd param
					if len(positionalParams) < 3 {
						SendError(w)
						return
					}

					if heightParam, ok := positionalParams[2].(string); ok {
						height, err = strconv.ParseInt(heightParam, 10, 64)
						if err != nil {
							SendError(w)
							return
						}
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
			} else { // try to support partially for other methods
				if method == "tx_search" {
					DoAggeratorJsonRpcOverHttp_tx_search(w, body)
					return
				} else if method == "block_by_hash" {
					DoAggeratorJsonRpcOverHttp_block_by_hash(w, body)
					return
				} else if method == "tx" {
					DoAggeratorJsonRpcOverHttp_tx(w, body)
					return
				} else if method == "block_search" {
					DoAggeratorJsonRpcOverHttp_block_search(w, body)
					return
				}
			}
		}
	}

	r.Body = io.NopCloser(bytes.NewBuffer(body)) // assign a new body with previous byte slice
	r.Host = r.URL.Host
	ProxyMapRpc[selectedHost].ServeHTTP(w, r)
}

func StartRpcServer() {
	handler := func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" { // URI over HTTP
			// see `/doc/rpc.md` and `https://github.com/tendermint/tendermint/blob/main/light/proxy/routes.go` to see the logic
			uriOverHttp(w, r)

		} else if r.Method == "POST" { // JSONRPC over HTTP
			jsonRpcOverHttp(w, r)
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
		//log.Fatal(http.ListenAndServe(":26657", serverMux))
		rpcServer = &http.Server{Addr: ":26657", Handler: serverMux}
		log.Fatal(rpcServer.ListenAndServe())

	}()
}

func ShutdownRpcServer() {
	if err := rpcServer.Shutdown(context.Background()); err != nil {
		log.Printf("rpcServer Shutdown: %v", err)
	}
}

func SendError(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte("Oops! Something was wrong"))
}

func SendResult(w http.ResponseWriter, body []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}
