package server

import (
	"context"
	"fmt"
	"github.com/notional-labs/subnode/config"
	"github.com/notional-labs/subnode/state"
	"github.com/notional-labs/subnode/utils"
	"log"
	"net/http"
	"path"
	"strconv"
	"strings"
)

type ApiServer struct {
	apiServer *http.Server
}

func NewApiServer() *ApiServer {
	newItem := &ApiServer{}
	return newItem
}

func (m *ApiServer) StartApiServer() {
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
			} else {
				// /cosmos/base/tendermint/v1beta1/blocks/{height}
				// /cosmos/base/tendermint/v1beta1/validatorsets/{height}

				urlPath := r.URL.Path
				//fmt.Printf("urlPath=%s\n", urlPath)

				if strings.HasPrefix(urlPath, "/cosmos/base/tendermint/v1beta1/blocks/") ||
					strings.HasPrefix(urlPath, "/cosmos/base/tendermint/v1beta1/validatorsets/") {
					pHeight := path.Base(r.URL.Path)
					fmt.Printf("pHeight=%s\n", pHeight)
					height, err := strconv.ParseInt(pHeight, 10, 64)
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
		}

		//fmt.Printf("selectedHost=%s\n", selectedHost)

		r.Host = r.URL.Host
		state.ProxyMapApi[selectedHost].ServeHTTP(w, r)
	}

	serverMux := http.NewServeMux()
	serverMux.HandleFunc("/", handler)
	go func() {
		m.apiServer = &http.Server{Addr: ":1317", Handler: serverMux}
		log.Fatal(m.apiServer.ListenAndServe())
	}()
}

func (m *ApiServer) ShutdownApiServer() {
	if err := m.apiServer.Shutdown(context.Background()); err != nil {
		log.Printf("apiServer Shutdown: %v", err)
	}
}
