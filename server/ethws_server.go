package server

import (
	"context"
	"fmt"
	"github.com/notional-labs/subnode/config"
	"github.com/notional-labs/subnode/state"
	"github.com/notional-labs/subnode/utils"
	"log"
	"net/http"
)

var ethWsServer *http.Server

func ethWsHandle(w http.ResponseWriter, r *http.Request) {
	prunedNode := state.SelectPrunedNode(config.ProtocolTypeEth)
	selectedHost := prunedNode.Backend.Eth // default to pruned node

	r.Host = r.URL.Host
	state.ProxyMapEth[selectedHost].ServeHTTP(w, r)
}

func StartEthWsServer() {
	fmt.Println("StartEthWsServer...")
	handler := func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" { // URI over HTTP
			ethWsHandle(w, r)
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
		ethWsServer = &http.Server{Addr: ":8546", Handler: serverMux}
		log.Fatal(ethWsServer.ListenAndServe())

	}()
}

func ShutdownEthWsServer() {
	if err := ethWsServer.Shutdown(context.Background()); err != nil {
		log.Printf("ethWsServer Shutdown: %v", err)
	}
}
