package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"time"
)

// DoAggeratorJsonRpcOverHttp_block_by_hash
// the logic here is it will iterate all the subnodes until result is found
// https://github.com/notional-labs/subnode/issues/20
func DoAggeratorJsonRpcOverHttp_block_by_hash(w http.ResponseWriter, jsonBody []byte) {
	for i, s := range PoolRpc {
		rpcUrl := s.Backend.Rpc
		body, err := FetchJsonRpcOverHttp(rpcUrl, jsonBody)
		if err == nil {
			//fmt.Printf("DoAggeratorJsonRpcOverHttp_block_by_hash, body=%s\n", string(body))

			var j0 interface{}
			err = json.Unmarshal(body, &j0)
			if err != nil {
				SendError(w)
				return
			}

			m0 := j0.(map[string]interface{})
			result := m0["result"].(map[string]interface{})
			block := result["block"]

			if (block != nil) || (i >= len(PoolRpc)-1) { // found result or last node, send it
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write(body)
				return
			}
		}
	}
}

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

// DoAggeratorUriOverHttp_block_by_hash
// the logic here is it will iterate all the subnodes until result is found
// https://github.com/notional-labs/subnode/issues/20
func DoAggeratorUriOverHttp_block_by_hash(w http.ResponseWriter, strQuery string) {
	for i, s := range PoolRpc {
		rpcUrl := s.Backend.Rpc + "/block_by_hash?" + strQuery

		//fmt.Printf("DoAggeratorUriOverHttp_block_by_hash, rpcUrl=%s\n", rpcUrl)
		body, err := FetchUriOverHttp(rpcUrl)
		if err == nil {
			//fmt.Printf("DoAggeratorUriOverHttp_block_by_hash, body=%s\n", string(body))

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

// DoAggeratorUriOverHttp_tx_search
// the logic here is it will iterate all the subnodes until result is found
// https://github.com/notional-labs/subnode/issues/20
func DoAggeratorUriOverHttp_tx_search(w http.ResponseWriter, strQuery string) {
	for i, s := range PoolRpc {
		rpcUrl := s.Backend.Rpc + "/tx_search?" + strQuery

		//fmt.Printf("DoAggeratorUriOverHttp_tx_search, rpcUrl=%s\n", rpcUrl)
		body, err := FetchUriOverHttp(rpcUrl)
		if err == nil {
			//fmt.Printf("DoAggeratorUriOverHttp_tx_search, body=%s\n", string(body))

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

func FetchUriOverHttp(url string) ([]byte, error) {
	hc := http.Client{Timeout: time.Second * 10}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

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
