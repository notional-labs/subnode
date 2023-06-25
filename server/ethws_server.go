package server

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/notional-labs/subnode/config"
	"github.com/notional-labs/subnode/state"
	"github.com/notional-labs/subnode/utils"
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

		//--------------------------------------------
		// process batch request
		if utils.IsBatch(msg) {
			log.Printf("ethws: process batch request")
			var arr []json.RawMessage

			err := json.Unmarshal(msg, &arr)
			if err != nil {
				errMsg := fmt.Sprintf("{\"jsonrpc\":\"2.0\",\"error\":{\"code\":500,\"message\":\"%s\"},\"id\":null}", err)
				_ = wsConServer.WriteMessage(websocket.TextMessage, []byte(errMsg))
				continue
			}

			var arrSize = len(arr)
			var arrRes = make([]json.RawMessage, arrSize)

			for i, s := range arr {
				bodyChild, err := utils.FetchJsonRpcOverHttp("http://localhost:8545", s)
				if err != nil {
					errMsg := fmt.Sprintf("{\"jsonrpc\":\"2.0\",\"error\":{\"code\":500,\"message\":\"%s\"},\"id\":null}", err)
					arrRes[i] = []byte(errMsg)
					continue
				}

				arrRes[i] = bodyChild
			}

			jsonBytes, err := json.Marshal(arrRes)
			if err != nil {
				errMsg := fmt.Sprintf("{\"jsonrpc\":\"2.0\",\"error\":{\"code\":500,\"message\":\"%s\"},\"id\":null}", err)
				_ = wsConServer.WriteMessage(websocket.TextMessage, []byte(errMsg))
				continue
			}

			_ = wsConServer.WriteMessage(websocket.TextMessage, jsonBytes)

			continue
		}

		//--------------------------------------------
		// process single request
		err = processSingleMsg(wsConServer, serverChannel, msg)
		if err != nil {
			errMsg := fmt.Sprintf("{\"jsonrpc\":\"2.0\",\"error\":{\"code\":500,\"message\":\"%s\"},\"id\":null}", err)
			_ = wsConServer.WriteMessage(websocket.TextMessage, []byte(errMsg))
		}
	}

	log.Println("exit ws-server")
}

func processSingleMsg(wsConServer *websocket.Conn, serverChannel chan []byte, msg []byte) error {
	var j0 interface{}
	err := json.Unmarshal(msg, &j0)
	if err == nil {
		if m0, ok := j0.(map[string]interface{}); ok {
			if method, ok := m0["method"].(string); ok {
				fmt.Printf("method=%s, params=%+v\n", method, m0["params"])
				if method != "eth_subscribe" && method != "eth_unsubscribe" {
					res, err := utils.FetchJsonRpcOverHttp("http://localhost:8545", msg)
					if err != nil {
						return err
					}

					_ = wsConServer.WriteMessage(websocket.TextMessage, res)
					return nil
				}
			}
		}
	}

	serverChannel <- msg // send msg to serverChannel
	return nil
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
