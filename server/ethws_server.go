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
		log.Println("dial err:", err)
		return nil, err
	}

	return c, nil
}

func isClosed(ch <-chan []byte) bool {
	select {
	case <-ch:
		return true
	default:
	}

	return false
}

func closeAll(wsConServer *websocket.Conn, wsConClient *websocket.Conn, clientChannel chan []byte, serverChannel chan []byte) {
	if !isClosed(clientChannel) {
		close(clientChannel)
	}
	if !isClosed(serverChannel) {
		close(serverChannel)
	}
	if wsConServer != nil {
		wsConServer.Close()
	}
	if wsConClient != nil {
		wsConClient.Close()
	}
}

func wsClientConRelay(wsConServer *websocket.Conn, wsConClient *websocket.Conn, clientChannel chan []byte, serverChannel chan []byte, wg *sync.WaitGroup) {
	defer wg.Done()
	//defer close(clientChannel)

	for {
		msg := <-clientChannel // receive msg from clientChannel

		// relay to server
		log.Println("relay to server")
		err := wsConServer.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			log.Println("relay to server err:", err)
			closeAll(wsConServer, wsConClient, clientChannel, serverChannel)
			break
		}
	}

	log.Println("exit processing clientChannel")
}

func wsServerConRelay(wsConServer *websocket.Conn, wsConClient *websocket.Conn, clientChannel chan []byte, serverChannel chan []byte, wg *sync.WaitGroup) {
	defer wg.Done()
	//defer close(serverChannel)

	for {
		msg := <-serverChannel // receive msg from serverChannel

		// relay to client
		log.Println("relay to client")
		err := wsConClient.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			log.Println("relay to client err:", err)
			closeAll(wsConServer, wsConClient, clientChannel, serverChannel)
			break
		}
	}

	log.Println("exit processing serverChannel")
}

func wsClientHandle(wsConServer *websocket.Conn, wsConClient *websocket.Conn, clientChannel chan []byte, serverChannel chan []byte, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		_, msg, err := wsConClient.ReadMessage()
		if err != nil {
			log.Println("ws-client read err:", err)
			closeAll(wsConServer, wsConClient, clientChannel, serverChannel)
			break
		}
		log.Printf("ws-client recv: %s", msg)
		clientChannel <- msg // send msg to clientChannel
	}

	log.Println("exit ws-client")
}

func wsServerHandle(wsConServer *websocket.Conn, wsConClient *websocket.Conn, clientChannel chan []byte, serverChannel chan []byte, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		_, msg, err := wsConServer.ReadMessage()
		if err != nil {
			log.Println("ws-server read err:", err)
			closeAll(wsConServer, wsConClient, clientChannel, serverChannel)
			break
		}
		log.Printf("ws-server recv: %s", msg)
		serverChannel <- msg // send msg to serverChannel
	}

	log.Println("exit ws-server")
}

func ethWsHandle(w http.ResponseWriter, r *http.Request) {
	var wg sync.WaitGroup
	var wsConServer *websocket.Conn
	var wsConClient *websocket.Conn
	clientChannel := make(chan []byte)
	serverChannel := make(chan []byte)

	var err error
	wsConServer, err = upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade err:", err)
		return
	}

	defer closeAll(wsConServer, wsConClient, clientChannel, serverChannel)

	//---------------------------------
	// ws-client
	wsConClient, err = createWSClient()
	if err != nil {
		log.Print("error with createWSClient:", err)
		return
	}

	wg.Add(1)
	go wsClientConRelay(wsConServer, wsConClient, clientChannel, serverChannel, &wg)
	wg.Add(1)
	go wsServerConRelay(wsConServer, wsConClient, clientChannel, serverChannel, &wg)
	wg.Add(1)
	go wsClientHandle(wsConServer, wsConClient, clientChannel, serverChannel, &wg)
	wg.Add(1)
	go wsServerHandle(wsConServer, wsConClient, clientChannel, serverChannel, &wg)

	wg.Wait()
	log.Printf("WaitGroup counter is zero")
}

func StartEthWsServer() {
	fmt.Println("StartEthWsServer...")

	handler := func(w http.ResponseWriter, r *http.Request) {
		ethWsHandle(w, r)
	}

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
