package state

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/notional-labs/subnode/config"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"time"
)

type BackendState struct {
	Name      string
	NodeType  config.BackendNodeType
	LastBlock int64

	Backend *config.Backend
}

var (
	PoolRpc     []*BackendState
	PoolApi     []*BackendState
	PoolGrpc    []*BackendState
	PoolEth     []*BackendState
	ProxyMapRpc = make(map[string]*httputil.ReverseProxy)
	ProxyMapApi = make(map[string]*httputil.ReverseProxy)
	ProxyMapEth = make(map[string]*httputil.ReverseProxy)
)

func Init() {
	InitPool()

	for i := 0; i < len(PoolRpc); i++ {
		rpcItem := PoolRpc[i]
		apiItem := PoolApi[i]
		ethItem := PoolEth[i]

		// rpc
		targetRpc, err := url.Parse(rpcItem.Name)
		if err != nil {
			panic(err)
		}
		ProxyMapRpc[rpcItem.Name] = httputil.NewSingleHostReverseProxy(targetRpc)

		// api
		targetApi, err := url.Parse(apiItem.Name)
		if err != nil {
			panic(err)
		}
		ProxyMapApi[apiItem.Name] = httputil.NewSingleHostReverseProxy(targetApi)

		// eth
		targetEth, err := url.Parse(ethItem.Name)
		if err != nil {
			panic(err)
		}
		ProxyMapEth[ethItem.Name] = httputil.NewSingleHostReverseProxy(targetEth)
	}
}

func InitPool() {
	PoolRpc = PoolRpc[:0] // Remove all elements
	PoolApi = PoolApi[:0] // Remove all elements
	PoolEth = PoolEth[:0] // Remove all elements

	for _, s := range config.GetConfig().Upstream {
		be := s // fix Copying the address of a loop variable in Go

		backendStateRpc := BackendState{
			Name:      s.Rpc,
			NodeType:  config.GetBackendNodeType(&be),
			LastBlock: 0,
			Backend:   &be,
		}
		fmt.Printf("debug: %+v\n", backendStateRpc)
		PoolRpc = append(PoolRpc, &backendStateRpc)

		//
		backendStateApi := BackendState{
			Name:      s.Api,
			NodeType:  config.GetBackendNodeType(&be),
			LastBlock: 0,
			Backend:   &be,
		}
		fmt.Printf("debug: %+v\n", backendStateApi)
		PoolApi = append(PoolApi, &backendStateApi)

		//
		backendStateGrpc := BackendState{
			Name:      s.Grpc,
			NodeType:  config.GetBackendNodeType(&be),
			LastBlock: 0,
			Backend:   &be,
		}
		fmt.Printf("debug: %+v\n", backendStateGrpc)
		PoolGrpc = append(PoolGrpc, &backendStateGrpc)

		//
		backendStateEth := BackendState{
			Name:      s.Eth,
			NodeType:  config.GetBackendNodeType(&be),
			LastBlock: 0,
			Backend:   &be,
		}
		fmt.Printf("debug: %+v\n", backendStateEth)
		PoolEth = append(PoolEth, &backendStateEth)
	}

	TaskUpdateState()
}

func SelectPrunedNode(t config.ProtocolType) *BackendState {
	pool := GetPool(t)

	for _, s := range pool {
		if s.NodeType == config.BackendNodeTypePruned {
			return s
		}
	}

	return nil
}

func GetPool(t config.ProtocolType) []*BackendState {
	if t == config.ProtocolTypeRpc {
		return PoolRpc
	} else if t == config.ProtocolTypeApi {
		return PoolApi
	} else if t == config.ProtocolTypeGrpc {
		return PoolGrpc
	} else if t == config.ProtocolTypeEth {
		return PoolEth
	}

	return nil
}

func SelectMatchedBackend(height int64, t config.ProtocolType) (*BackendState, error) {
	pool := GetPool(t)

	for _, s := range pool {
		//fmt.Printf("debug: %+v\n", s)
		if s.NodeType == config.BackendNodeTypePruned {
			if s.LastBlock > 0 { // fetched last block
				earliestHeight := s.LastBlock - s.Backend.Blocks[0]
				if height >= earliestHeight {
					return s, nil
				}
			}
		} else if s.NodeType == config.BackendNodeTypeSubNode {
			if height >= s.Backend.Blocks[0] {
				if s.Backend.Blocks[1] == 0 { // to the newest block
					return s, nil
				} else if height <= s.Backend.Blocks[1] {
					return s, nil
				}
			}
		} else if s.NodeType == config.BackendNodeTypeArchive {
			return s, nil
		}
	}

	return nil, errors.New("no node matched")
}

func IsNeededToFetchLastBlock(s *BackendState) bool {
	return s.NodeType == config.BackendNodeTypePruned
}

func TaskUpdateState() {
	// call close(quit) to stop

	doUpdateState() // call the 1st time immediately after starting

	ticker := time.NewTicker(30 * time.Second)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				doUpdateState()

			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
}

func doUpdateState() {
	for i := 0; i < len(PoolRpc); i++ {
		rpcItem := PoolRpc[i]
		apiItem := PoolApi[i]
		grpcItem := PoolGrpc[i]
		ethItem := PoolEth[i]

		if IsNeededToFetchLastBlock(rpcItem) == false {
			continue
		}

		height, err := FetchHeightFromStatus(rpcItem.Backend.Rpc)
		if err != nil {
			fmt.Println("Err FetchHeightFromStatus", err)
			continue
		}

		rpcItem.LastBlock = height
		apiItem.LastBlock = height
		grpcItem.LastBlock = height
		ethItem.LastBlock = height
	}
}

func FetchHeightFromStatus(rpcUrl string) (int64, error) {
	urlRequest := rpcUrl + "/status"
	spaceClient := http.Client{Timeout: time.Second * 30}

	req, err := http.NewRequest(http.MethodGet, urlRequest, nil)
	if err != nil {
		return 0, err
	}

	res, getErr := spaceClient.Do(req)
	if getErr != nil {
		return 0, err
	}

	if res.Body != nil {
		defer func() {
			_ = res.Body.Close()
		}()
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return 0, err
	}

	//fmt.Printf("body=%s\n", string(body))
	height, err := ReadHeightFromStatusJson(body)
	if err != nil {
		return 0, err
	}

	return height, nil
}

func ReadHeightFromStatusJson(jsonText []byte) (int64, error) {
	var j0 interface{}
	err := json.Unmarshal(jsonText, &j0)
	if err != nil {
		return 0, err
	}
	m0 := j0.(map[string]interface{})
	j1 := m0["result"]
	m1 := j1.(map[string]interface{})
	j2 := m1["sync_info"]
	m2 := j2.(map[string]interface{})
	j3 := m2["latest_block_height"]
	v := j3.(string)

	height, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		return 0, err
	}

	return height, nil
}
