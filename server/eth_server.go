package server

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/notional-labs/subnode/config"
	"github.com/notional-labs/subnode/state"
	"github.com/notional-labs/subnode/utils"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
)

var ethServer *http.Server

func ethJsonRpcOverHttp(w http.ResponseWriter, r *http.Request) {
	prunedNode := state.SelectPrunedNode(config.ProtocolTypeEth)
	selectedHost := prunedNode.Backend.Eth // default to pruned node

	body, err := io.ReadAll(r.Body)
	if err != nil {
		_ = utils.SendError(w)
		return
	}

	fmt.Printf("body=%s\n", string(body))

	var j0 interface{}
	err = json.Unmarshal(body, &j0)
	if err != nil {
		_ = utils.SendError(w)
		return
	}

	if m0, ok := j0.(map[string]interface{}); ok {
		if method, ok := m0["method"].(string); ok {
			//params := m0["params"].([]interface{})

			fmt.Printf("method=%s, params=%+v\n", method, m0["params"])

			//if method == "web3_clientVersion" ||
			//	method == "web3_sha3" ||
			//	method == "net_version" ||
			//	method == "net_peerCount" ||
			//	method == "net_listening" ||
			//	method == "eth_protocolVersion" ||
			//	method == "eth_syncing" ||
			//	method == "eth_gasPrice" ||
			//	method == "eth_accounts" ||
			//	method == "eth_blockNumber" {
			//	selectedHost = prunedNode.Backend.Eth
			//} else
			if method == "eth_getBalance" {
				height := int64(-1)

				if positionalParams, ok := m0["params"].([]interface{}); ok {
					// height is 2nd param
					if len(positionalParams) < 2 {
						_ = utils.SendError(w)
						return
					}

					if heightParam, ok := positionalParams[1].(string); ok {
						if strings.HasPrefix(heightParam, "0x") {
							heightParam = strings.TrimPrefix(heightParam, "0x")
							height, err = strconv.ParseInt(heightParam, 16, 64)
							if err != nil {
								_ = utils.SendError(w)
								return
							}
						}
					}
				}

				if height >= 0 {
					node, err := state.SelectMatchedBackend(height, config.ProtocolTypeEth)
					if err != nil {
						_ = utils.SendError(w)
						return
					}

					selectedHost = node.Backend.Eth
				}
			} else if method == "eth_getStorageAt" {
				height := int64(-1)

				if positionalParams, ok := m0["params"].([]interface{}); ok {
					// height is 3rd param
					if len(positionalParams) < 3 {
						_ = utils.SendError(w)
						return
					}

					if heightParam, ok := positionalParams[2].(string); ok {
						if strings.HasPrefix(heightParam, "0x") {
							heightParam = strings.TrimPrefix(heightParam, "0x")
							height, err = strconv.ParseInt(heightParam, 16, 64)
							if err != nil {
								_ = utils.SendError(w)
								return
							}
						}
					}
				}

				if height >= 0 {
					node, err := state.SelectMatchedBackend(height, config.ProtocolTypeEth)
					if err != nil {
						_ = utils.SendError(w)
						return
					}

					selectedHost = node.Backend.Eth
				}
			}
		}
	}

	r.Body = io.NopCloser(bytes.NewBuffer(body)) // assign a new body with previous byte slice
	r.Host = r.URL.Host
	state.ProxyMapEth[selectedHost].ServeHTTP(w, r)
}

func StartEthServer() {
	fmt.Println("StartEthServer...")
	handler := func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" { // JSONRPC over HTTP
			ethJsonRpcOverHttp(w, r)
		} else {
			_ = utils.SendError(w)
			return
		}
	}
	// handle all requests to your server using the proxy
	//http.HandleFunc("/", handler)
	serverMux := http.NewServeMux()
	serverMux.HandleFunc("/", handler)
	go func() {
		ethServer = &http.Server{Addr: ":8545", Handler: serverMux}
		log.Fatal(ethServer.ListenAndServe())
	}()
}

func ShutdownEthServer() {
	if err := ethServer.Shutdown(context.Background()); err != nil {
		log.Printf("ethServer Shutdown: %v", err)
	}
}
