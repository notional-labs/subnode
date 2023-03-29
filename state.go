package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"time"
)

type BackendState struct {
	Name      string
	NodeType  BackendNodeType
	LastBlock int64

	Backend *Backend
}

var (
	PoolRpc     []*BackendState
	PoolApi     []*BackendState
	ProxyMapRpc = make(map[string]*httputil.ReverseProxy)
	ProxyMapApi = make(map[string]*httputil.ReverseProxy)
)

func Init() {
	InitPool()

	for _, s := range PoolRpc {
		target, err := url.Parse(s.Name)
		if err != nil {
			panic(err)
		}
		ProxyMapRpc[s.Name] = httputil.NewSingleHostReverseProxy(target)
	}

	for _, s := range PoolApi {
		target, err := url.Parse(s.Name)
		if err != nil {
			panic(err)
		}
		ProxyMapApi[s.Name] = httputil.NewSingleHostReverseProxy(target)
	}
}

func InitPool() {
	cfg := GetConfig()

	PoolRpc = PoolRpc[:0] // Remove all elements
	PoolApi = PoolApi[:0] // Remove all elements
	for _, s := range cfg.Upstream {
		be := s // fix Copying the address of a loop variable in Go

		backendStateRpc := BackendState{
			Name:      s.Rpc,
			NodeType:  GetBackendNodeType(&be),
			LastBlock: 0,
			Backend:   &be,
		}
		fmt.Printf("debug: %+v\n", backendStateRpc)
		PoolRpc = append(PoolRpc, &backendStateRpc)

		//
		backendStateApi := BackendState{
			Name:      s.Api,
			NodeType:  GetBackendNodeType(&be),
			LastBlock: 0,
			Backend:   &be,
		}
		fmt.Printf("debug: %+v\n", backendStateApi)
		PoolApi = append(PoolApi, &backendStateApi)
	}

	TaskUpdateState()
}

func SelectPrunedNodeRpc() *BackendState {
	for _, s := range PoolRpc {
		if s.NodeType == BackendNodeTypePruned {
			return s
		}
	}

	return nil
}

func SelectPrunedNodeApi() *BackendState {
	for _, s := range PoolApi {
		if s.NodeType == BackendNodeTypePruned {
			return s
		}
	}

	return nil
}

func SelectMatchedNodeRpc(height int64) (*BackendState, error) {
	for _, s := range PoolRpc {
		fmt.Printf("debug: %+v\n", s)
		if s.NodeType == BackendNodeTypePruned {
			earliestHeight := s.LastBlock - s.Backend.Blocks[0]
			if height >= earliestHeight {
				return s, nil
			}
		} else if s.NodeType == BackendNodeTypeSubNode {
			if (height >= s.Backend.Blocks[0]) && (height <= s.Backend.Blocks[0]) {
				return s, nil
			}
		} else if s.NodeType == BackendNodeTypeArchive {
			return s, nil
		}
	}

	return nil, errors.New("no node matched")
}

func SelectMatchedNodeApi(height int64) (*BackendState, error) {
	for _, s := range PoolApi {
		fmt.Printf("debug: %+v\n", s)
		if s.NodeType == BackendNodeTypePruned {
			earliestHeight := s.LastBlock - s.Backend.Blocks[0]
			if height >= earliestHeight {
				return s, nil
			}
		} else if s.NodeType == BackendNodeTypeSubNode {
			if (height >= s.Backend.Blocks[0]) && (height <= s.Backend.Blocks[0]) {
				return s, nil
			}
		} else if s.NodeType == BackendNodeTypeArchive {
			return s, nil
		}
	}

	return nil, errors.New("no node matched")
}

func TaskUpdateState() {
	// call close(quit) to stop

	ticker := time.NewTicker(30 * time.Second)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				for _, s := range PoolRpc {
					if s.NodeType == BackendNodeTypePruned {
						height, err := FetchHeightFromStatus(s.Backend.Rpc)
						if err == nil {
							s.LastBlock = height
						} else {
							fmt.Println("Err FetchHeightFromStatus", err)
						}
					}
				}

				for _, s := range PoolApi {
					if s.NodeType == BackendNodeTypePruned {
						height, err := FetchHeightFromStatus(s.Backend.Rpc)
						if err == nil {
							s.LastBlock = height
						} else {
							fmt.Println("Err FetchHeightFromStatus", err)
						}
					}
				}

			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
}

func FetchHeightFromStatus(rpcUrl string) (int64, error) {
	url := rpcUrl + "/status"
	spaceClient := http.Client{Timeout: time.Second * 10}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return 0, err
	}

	res, getErr := spaceClient.Do(req)
	if getErr != nil {
		return 0, err
	}

	if res.Body != nil {
		defer res.Body.Close()
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
