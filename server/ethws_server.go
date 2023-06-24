package server

import (
	"context"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var ethWsServer *http.Server

//func ethWsHandle(w http.ResponseWriter, r *http.Request) {
//	prunedNode := state.SelectPrunedNode(config.ProtocolTypeEthWs)
//	selectedHost := prunedNode.Backend.Eth // default to pruned node
//
//	r.Host = r.URL.Host
//	state.ProxyMapEthWs[selectedHost].ServeHTTP(w, r)
//}

func StartEthWsServer() {
	fmt.Println("StartEthWsServer...")

	var upgrader = websocket.Upgrader{} // use default options

	handler := func(w http.ResponseWriter, r *http.Request) {
		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Print("upgrade:", err)
			return
		}
		defer c.Close()
		for {
			mt, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				break
			}
			log.Printf("recv: %s", message)
			err = c.WriteMessage(mt, message)
			if err != nil {
				log.Println("write:", err)
				break
			}
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
