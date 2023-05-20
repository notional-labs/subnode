package server

import (
	"context"
	"fmt"
	"github.com/notional-labs/subnode/config"
	"github.com/notional-labs/subnode/state"
	"github.com/notional-labs/subnode/utils"
	"log"
	"net/http"
	"strconv"
)

var apiServer *http.Server

func StartApiServer() {
	fmt.Println("StartApiServer...")
	handler := func(w http.ResponseWriter, r *http.Request) {
		prunedNode := state.SelectPrunedNode(config.ProtocolTypeApi)

		if r.Method != "GET" && r.Method != "POST" {
			_ = utils.SendError(w)
		}

		selectedHost := prunedNode.Backend.Api // default to pruned node

		if r.Method == "GET" {
			fmt.Printf("r.RequestURI=%s\n", r.RequestURI)

			xCosmosBlockHeight := r.Header.Get("x-cosmos-block-height")
			if xCosmosBlockHeight != "" {
				height, err := strconv.ParseInt(xCosmosBlockHeight, 10, 64)
				if err != nil {
					_ = utils.SendError(w)
				}

				node, err := state.SelectMatchedBackend(height, config.ProtocolTypeApi)
				if err != nil {
					_ = utils.SendError(w)
				}

				selectedHost = node.Backend.Api
			}
		}

		r.Host = r.URL.Host
		state.ProxyMapApi[selectedHost].ServeHTTP(w, r)
	}
	// handle all requests to your server using the proxy
	//http.HandleFunc("/", handler)
	serverMux := http.NewServeMux()
	serverMux.HandleFunc("/", handler)
	go func() {
		//log.Fatal(http.ListenAndServe(":1337", serverMux))
		apiServer = &http.Server{Addr: ":1337", Handler: serverMux}
		log.Fatal(apiServer.ListenAndServe())
	}()
}

func ShutdownApiServer() {
	if err := apiServer.Shutdown(context.Background()); err != nil {
		log.Printf("apiServer Shutdown: %v", err)
	}
}
