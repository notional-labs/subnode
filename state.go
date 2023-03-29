package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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
	Pool []*BackendState
)

func InitPool() {
	Pool = Pool[:0] // Remove all elements

	cfg := GetConfig()
	for _, s := range cfg.Upstream {
		//target, err := url.Parse(s.Rpc)
		//if err != nil {
		//	panic(err)
		//}
		//hostProxy[s.Rpc] = httputil.NewSingleHostReverseProxy(target)

		backendState := BackendState{
			Name:      s.Rpc,
			NodeType:  GetBackendNodeType(&s),
			LastBlock: 0,
			Backend:   &s,
		}

		Pool = append(Pool, &backendState)
	}

	TaskUpdateState()
}

func SelectPrunedNode() *BackendState {
	for _, s := range Pool {
		if s.NodeType == BackendNodeTypePruned {
			return s
		}
	}

	return nil
}

func TaskUpdateState() {
	// call close(quit) to stop

	ticker := time.NewTicker(30 * time.Second)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				// do stuff
				for _, s := range Pool {
					if s.NodeType == BackendNodeTypePruned {
						fmt.Println(s.Name)

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

	fmt.Printf("body=%s\n", string(body))
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
