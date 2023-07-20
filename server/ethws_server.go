package server

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/cometbft/cometbft/rpc/jsonrpc/types"
	"github.com/gorilla/websocket"
	"github.com/notional-labs/subnode/config"
	"github.com/notional-labs/subnode/state"
	"github.com/notional-labs/subnode/utils"
	"log"
	"net/http"
	"net/url"
	"sync"
)

var upgrader = websocket.Upgrader{} // use default options

type EthWsServer struct {
	ethWsServer *http.Server
}

func NewEthWsServer() *EthWsServer {
	newItem := &EthWsServer{}
	return newItem
}

func (m *EthWsServer) createWSClient() (*websocket.Conn, error) {
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

func (m *EthWsServer) wsClientConRelay(wsConServer *websocket.Conn, wsConClient *websocket.Conn, clientChannel chan []byte, serverChannel chan []byte, wg *sync.WaitGroup) {
	defer wg.Done()
	//defer close(clientChannel)

	for {
		msg := <-clientChannel // receive msg from clientChannel

		// relay to server
		log.Println("relay to server")
		err := wsConServer.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			log.Println("relay to server err:", err)
			m.closeAll(wsConServer, wsConClient, clientChannel, serverChannel)
			break
		}
	}

	log.Println("exit processing clientChannel")
}

func (m *EthWsServer) wsServerConRelay(wsConServer *websocket.Conn, wsConClient *websocket.Conn, clientChannel chan []byte, serverChannel chan []byte, wg *sync.WaitGroup) {
	defer wg.Done()
	//defer close(serverChannel)

	for {
		msg := <-serverChannel // receive msg from serverChannel

		// relay to client
		log.Println("relay to client")
		err := wsConClient.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			log.Println("relay to client err:", err)
			m.closeAll(wsConServer, wsConClient, clientChannel, serverChannel)
			break
		}
	}

	log.Println("exit processing serverChannel")
}

func (m *EthWsServer) wsClientHandle(wsConServer *websocket.Conn, wsConClient *websocket.Conn, clientChannel chan []byte, serverChannel chan []byte, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		_, msg, err := wsConClient.ReadMessage()
		if err != nil {
			log.Println("ws-client read err:", err)
			m.closeAll(wsConServer, wsConClient, clientChannel, serverChannel)
			break
		}
		log.Printf("ws-client recv: %s", msg)
		clientChannel <- msg // send msg to clientChannel
	}

	log.Println("exit ws-client")
}

func (m *EthWsServer) wsServerHandle(wsConServer *websocket.Conn, wsConClient *websocket.Conn, clientChannel chan []byte, serverChannel chan []byte, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		_, msg, err := wsConServer.ReadMessage()
		if err != nil {
			log.Println("ws-server read err:", err)
			m.closeAll(wsConServer, wsConClient, clientChannel, serverChannel)
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
				_ = m.sendWsRpcResponseErr(wsConServer, err)
				continue
			}

			var arrSize = len(arr)
			var arrRes = make([]json.RawMessage, arrSize)

			for i, s := range arr {
				bodyChild, err := utils.FetchJsonRpcOverHttp("http://localhost:8545", s)
				if err != nil {
					errMsg := m.getWsJsonRpcErr(s, err)
					arrRes[i] = errMsg
					continue
				}

				arrRes[i] = bodyChild
			}

			jsonBytes, err := json.Marshal(arrRes)
			if err != nil {
				_ = m.sendWsRpcResponseErr(wsConServer, err)
				continue
			}

			_ = wsConServer.WriteMessage(websocket.TextMessage, jsonBytes)
			continue
		}

		//--------------------------------------------
		// process single request
		var rpcReq = types.RPCRequest{}
		err = rpcReq.UnmarshalJSON(msg)
		if err != nil {
			_ = m.sendWsRpcResponseErr(wsConServer, err)
			continue
		}

		m.processSingleMsg(wsConServer, serverChannel, &rpcReq, msg)
	}

	log.Println("exit ws-server")
}

func (m *EthWsServer) ethWsHandle(w http.ResponseWriter, r *http.Request) {
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

	defer m.closeAll(wsConServer, wsConClient, clientChannel, serverChannel)

	//---------------------------------
	// ws-client
	wsConClient, err = m.createWSClient()
	if err != nil {
		log.Print("error with createWSClient:", err)
		return
	}

	wg.Add(1)
	go m.wsClientConRelay(wsConServer, wsConClient, clientChannel, serverChannel, &wg)
	wg.Add(1)
	go m.wsServerConRelay(wsConServer, wsConClient, clientChannel, serverChannel, &wg)
	wg.Add(1)
	go m.wsClientHandle(wsConServer, wsConClient, clientChannel, serverChannel, &wg)
	wg.Add(1)
	go m.wsServerHandle(wsConServer, wsConClient, clientChannel, serverChannel, &wg)

	wg.Wait()
	log.Printf("WaitGroup counter is zero")
}

func (m *EthWsServer) StartEthWsServer() {
	fmt.Println("StartEthWsServer...")

	handler := func(w http.ResponseWriter, r *http.Request) {
		m.ethWsHandle(w, r)
	}

	serverMux := http.NewServeMux()
	serverMux.HandleFunc("/", handler)
	go func() {
		m.ethWsServer = &http.Server{Addr: ":8546", Handler: serverMux}
		log.Fatal(m.ethWsServer.ListenAndServe())

	}()
}

func (m *EthWsServer) ShutdownEthWsServer() {
	if err := m.ethWsServer.Shutdown(context.Background()); err != nil {
		log.Printf("ethWsServer Shutdown: %v", err)
	}
}

// ------------------
// helper

func (m *EthWsServer) getWsJsonRpcErr(jsonRawReq []byte, err error) []byte {
	var rpcReq = types.RPCRequest{}
	errJson := rpcReq.UnmarshalJSON(jsonRawReq)
	if errJson != nil {
		errMsg := "{\"jsonrpc\":\"2.0\",\"error\":{\"code\":-32600,\"message\":\"invalid request\"},\"id\":null}"
		return []byte(errMsg)
	}

	rpcRes := types.NewRPCErrorResponse(rpcReq.ID, 500, err.Error(), "")
	jsonStrRes, _ := json.Marshal(rpcRes)
	return jsonStrRes
}

func (m *EthWsServer) sendWsRpcResponseErr(wsConServer *websocket.Conn, err error) error {
	errMsg := fmt.Sprintf("{\"jsonrpc\":\"2.0\",\"error\":{\"code\":-32600,\"message\":\"%s\"},\"id\":null}", err)
	errWrite := wsConServer.WriteMessage(websocket.TextMessage, []byte(errMsg))
	return errWrite
}

func (m *EthWsServer) processSingleMsg(wsConServer *websocket.Conn, serverChannel chan []byte, rpcReq *types.RPCRequest, jsonRaw []byte) {
	if rpcReq.Method != "eth_subscribe" && rpcReq.Method != "eth_unsubscribe" {
		res, err := utils.FetchJsonRpcOverHttp("http://localhost:8545", jsonRaw)
		if err != nil {
			_ = m.sendWsRpcResponseErr(wsConServer, err)
			return
		}

		_ = wsConServer.WriteMessage(websocket.TextMessage, res)
		return
	}

	serverChannel <- jsonRaw // send msg to serverChannel
}

func (m *EthWsServer) isClosed(ch <-chan []byte) bool {
	select {
	case <-ch:
		return true
	default:
	}

	return false
}

func (m *EthWsServer) closeAll(wsConServer *websocket.Conn, wsConClient *websocket.Conn, clientChannel chan []byte, serverChannel chan []byte) {
	if !m.isClosed(clientChannel) {
		close(clientChannel)
	}
	if !m.isClosed(serverChannel) {
		close(serverChannel)
	}
	if wsConServer != nil {
		wsConServer.Close()
	}
	if wsConClient != nil {
		wsConClient.Close()
	}
}
