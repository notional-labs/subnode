package server

import (
	"context"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/notional-labs/subnode/config"
	"github.com/notional-labs/subnode/state"
	"log"
	"net/http"
	"net/url"
	"sync"
)

var ethWsServer *http.Server
var upgrader = websocket.Upgrader{} // use default options

func createWSClient() (*websocket.Conn, error) {
	prunedNode := state.SelectPrunedNode(config.ProtocolTypeEthWs)
	selectedHost := prunedNode.Backend.EthWs // default to pruned node
	targetEthWs, err := url.Parse(selectedHost)
	if err != nil {
		return nil, err
	}
	log.Printf("connecting to %s", targetEthWs.String())

	c, _, err := websocket.DefaultDialer.Dial(targetEthWs.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}

	return c, nil
}

func ethWsHandle(w http.ResponseWriter, r *http.Request) {
	wsConServer, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer wsConServer.Close()

	var wg sync.WaitGroup

	//---------------------------------
	// ws-client
	wsConClient, err := createWSClient()
	if err != nil {
		log.Print("error:", err)
		return
	}
	defer wsConClient.Close()

	clientChannel := make(chan []byte) // struct{}
	serverChannel := make(chan []byte)

	wg.Add(1)
	go func() {
		defer wg.Done()
		//defer close(clientChannel)

		for {
			msg := <-clientChannel // receive msg from clientChannel

			// relay to server
			log.Println("relay to server")
			err = wsConServer.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				log.Println("relay to server:", err)
				break
			}
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		//defer close(serverChannel)

		for {
			msg := <-serverChannel // receive msg from serverChannel

			// relay to client
			log.Println("relay to client")
			err = wsConClient.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				log.Println("relay to client:", err)
				break
			}
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		for {
			_, msg, err := wsConClient.ReadMessage()
			if err != nil {
				log.Println("ws-client read:", err)
				break
			}
			log.Printf("ws-client recv: %s", msg)
			clientChannel <- msg // send msg to clientChannel
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		for {
			_, msg, err := wsConServer.ReadMessage()
			if err != nil {
				log.Println("ws-server read:", err)
				break
			}
			log.Printf("ws-server recv: %s", msg)
			serverChannel <- msg // send msg to serverChannel
		}
	}()

	wg.Wait()
}

func StartEthWsServer() {
	fmt.Println("StartEthWsServer...")

	handler := func(w http.ResponseWriter, r *http.Request) {
		ethWsHandle(w, r)
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
