package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
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
						url := s.Backend.Rpc + "/status"
						spaceClient := http.Client{Timeout: time.Second * 10}

						req, err := http.NewRequest(http.MethodGet, url, nil)
						if err != nil {
							log.Fatal(err)
						}

						//req.Header.Set("User-Agent", "subnode")

						res, getErr := spaceClient.Do(req)
						if getErr != nil {
							log.Fatal(getErr)
						}

						if res.Body != nil {
							defer res.Body.Close()
						}

						body, readErr := io.ReadAll(res.Body)
						if readErr != nil {
							log.Fatal(readErr)
						}

						fmt.Printf("body=%s\n", string(body))

					}
				}

			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
}
