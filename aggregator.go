package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"time"
)

// DoAggeratorJsonRpcOverHttp_tx_search
// the logic here is it will iterate all the subnodes until result is found
// https://github.com/notional-labs/subnode/issues/20
func DoAggeratorJsonRpcOverHttp_tx_search(w http.ResponseWriter, jsonBody []byte) {
	for i, s := range PoolRpc {
		rpcUrl := s.Backend.Rpc
		body, err := FetchJsonRpcOverHttp(rpcUrl, jsonBody)
		if err == nil {
			//fmt.Printf("DoAggeratorJsonRpcOverHttp_tx_search, body=%s\n", string(body))

			var j0 interface{}
			err = json.Unmarshal(body, &j0)
			if err != nil {
				SendError(w)
				return
			}

			m0 := j0.(map[string]interface{})
			result := m0["result"].(map[string]interface{})
			txs := result["txs"].([]interface{})

			if (len(txs) > 0) || (i >= len(PoolRpc)-1) { // found result or last node, send it
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write(body)
				return
			}
		}
	}
}

func FetchJsonRpcOverHttp(url string, jsonBody []byte) ([]byte, error) {
	hc := http.Client{Timeout: time.Second * 10}

	bodyReader := bytes.NewReader(jsonBody)
	req, err := http.NewRequest(http.MethodPost, url, bodyReader)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	res, getErr := hc.Do(req)
	if getErr != nil {
		return nil, err
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}