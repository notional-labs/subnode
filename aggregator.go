package main

import (
	"bytes"
	"encoding/json"
	"github.com/tidwall/sjson"
	"io"
	"net/http"
	"time"
)

func DoAggeratorUriOverHttp_status(w http.ResponseWriter, strQuery string) {
	prunedNode := SelectPrunedNode(ProtocolTypeRpc)
	rpcUrl := prunedNode.Backend.Rpc + "/status?" + strQuery
	body, err := FetchUriOverHttp(rpcUrl)
	if err != nil {
		SendError(w)
		return
	}

	var j0 interface{}
	err = json.Unmarshal(body, &j0)
	if err != nil {
		SendError(w)
		return
	}

	if m0, ok := j0.(map[string]interface{}); ok {
		if result, ok := m0["result"].(map[string]interface{}); ok {
			if node_info, ok := result["node_info"].(map[string]interface{}); ok {
				if node_info["network"] == "osmosis-1" {
					body, _ = sjson.SetBytes(body, "result.sync_info.earliest_block_hash", "C8DC787FAAE0941EF05C75C3AECCF04B85DFB1D4A8D054A463F323B0D9459719")
					body, _ = sjson.SetBytes(body, "result.sync_info.earliest_app_hash", "E3B0C44298FC1C149AFBF4C8996FB92427AE41E4649B934CA495991B7852B855")
					body, _ = sjson.SetBytes(body, "result.sync_info.earliest_block_height", "1")
					body, _ = sjson.SetBytes(body, "result.sync_info.earliest_block_time", "2021-06-18T17:00:00Z")
				}
			}
		}
	}

	SendResult(w, body)
}

// the logic here is it will iterate all the subnodes until result is found
// https://github.com/notional-labs/subnode/issues/20

func DoAggeratorJsonRpcOverHttp_block_search(w http.ResponseWriter, jsonBody []byte) {
	for i, s := range PoolRpc {
		rpcUrl := s.Backend.Rpc
		body, err := FetchJsonRpcOverHttp(rpcUrl, jsonBody)
		if err == nil {
			var j0 interface{}
			err = json.Unmarshal(body, &j0)
			if err != nil {
				SendError(w)
				return
			}

			if m0, ok := j0.(map[string]interface{}); ok {
				if result, ok := m0["result"].(map[string]interface{}); ok {
					if blocks, ok := result["blocks"].([]interface{}); ok {
						if (len(blocks) > 0) || (i >= len(PoolRpc)-1) { // found result or last node, send it
							SendResult(w, body)
							return
						}
					}
				}
			}
		}
	}

	SendError(w)
}

func DoAggeratorJsonRpcOverHttp_tx(w http.ResponseWriter, jsonBody []byte) {
	for i, s := range PoolRpc {
		rpcUrl := s.Backend.Rpc
		body, err := FetchJsonRpcOverHttp(rpcUrl, jsonBody)
		if err == nil {
			var j0 interface{}
			err = json.Unmarshal(body, &j0)
			if err != nil {
				SendError(w)
				return
			}

			if m0, ok := j0.(map[string]interface{}); ok {
				if result, ok := m0["result"].(map[string]interface{}); ok {
					if (result["hash"] != nil) || (i >= len(PoolRpc)-1) { // found result or last node, send it
						SendResult(w, body)
						return
					}
				}
			}
		}
	}

	SendError(w)
}

func DoAggeratorJsonRpcOverHttp_block_by_hash(w http.ResponseWriter, jsonBody []byte) {
	for i, s := range PoolRpc {
		rpcUrl := s.Backend.Rpc
		body, err := FetchJsonRpcOverHttp(rpcUrl, jsonBody)
		if err == nil {
			var j0 interface{}
			err = json.Unmarshal(body, &j0)
			if err != nil {
				SendError(w)
				return
			}

			if m0, ok := j0.(map[string]interface{}); ok {
				if result, ok := m0["result"].(map[string]interface{}); ok {
					if (result["block"] != nil) || (i >= len(PoolRpc)-1) { // found result or last node, send it
						SendResult(w, body)
						return
					}
				}
			}
		}
	}

	SendError(w)
}

func DoAggeratorJsonRpcOverHttp_tx_search(w http.ResponseWriter, jsonBody []byte) {
	for i, s := range PoolRpc {
		rpcUrl := s.Backend.Rpc
		body, err := FetchJsonRpcOverHttp(rpcUrl, jsonBody)
		if err == nil {
			var j0 interface{}
			err = json.Unmarshal(body, &j0)
			if err != nil {
				SendError(w)
				return
			}

			if m0, ok := j0.(map[string]interface{}); ok {
				if result, ok := m0["result"].(map[string]interface{}); ok {
					if txs, ok := result["txs"].([]interface{}); ok {
						if (len(txs) > 0) || (i >= len(PoolRpc)-1) { // found result or last node, send it
							SendResult(w, body)
							return
						}
					}
				}
			}
		}
	}

	SendError(w)
}

func DoAggeratorUriOverHttp_tx(w http.ResponseWriter, strQuery string) {
	for i, s := range PoolRpc {
		rpcUrl := s.Backend.Rpc + "/tx?" + strQuery
		body, err := FetchUriOverHttp(rpcUrl)
		if err == nil {
			var j0 interface{}
			err = json.Unmarshal(body, &j0)
			if err != nil {
				SendError(w)
				return
			}

			if m0, ok := j0.(map[string]interface{}); ok {
				if result, ok := m0["result"].(map[string]interface{}); ok {
					if (result["hash"] != nil) || (i >= len(PoolRpc)-1) { // found result or last node, send it
						SendResult(w, body)
						return
					}
				}
			}
		}
	}

	SendError(w)
}

func DoAggeratorUriOverHttp_block_by_hash(w http.ResponseWriter, strQuery string) {
	for i, s := range PoolRpc {
		rpcUrl := s.Backend.Rpc + "/block_by_hash?" + strQuery
		body, err := FetchUriOverHttp(rpcUrl)
		if err == nil {
			var j0 interface{}
			err = json.Unmarshal(body, &j0)
			if err != nil {
				SendError(w)
				return
			}

			if m0, ok := j0.(map[string]interface{}); ok {
				if result, ok := m0["result"].(map[string]interface{}); ok {
					if (result["block"] != nil) || (i >= len(PoolRpc)-1) { // found result or last node, send it
						SendResult(w, body)
						return
					}
				}
			}
		}
	}

	SendError(w)
}

func DoAggeratorUriOverHttp_block_search(w http.ResponseWriter, strQuery string) {
	for i, s := range PoolRpc {
		rpcUrl := s.Backend.Rpc + "/block_search?" + strQuery
		body, err := FetchUriOverHttp(rpcUrl)
		if err == nil {
			var j0 interface{}
			err = json.Unmarshal(body, &j0)
			if err != nil {
				SendError(w)
				return
			}

			if m0, ok := j0.(map[string]interface{}); ok {
				if result, ok := m0["result"].(map[string]interface{}); ok {
					if blocks, ok := result["blocks"].([]interface{}); ok {
						if (len(blocks) > 0) || (i >= len(PoolRpc)-1) { // found result or last node, send it
							SendResult(w, body)
							return
						}
					}
				}
			}
		}
	}

	SendError(w)
}

func DoAggeratorUriOverHttp_tx_search(w http.ResponseWriter, strQuery string) {
	for i, s := range PoolRpc {
		rpcUrl := s.Backend.Rpc + "/tx_search?" + strQuery
		body, err := FetchUriOverHttp(rpcUrl)
		if err == nil {
			var j0 interface{}
			err = json.Unmarshal(body, &j0)
			if err != nil {
				SendError(w)
				return
			}

			if m0, ok := j0.(map[string]interface{}); ok {
				if result, ok := m0["result"].(map[string]interface{}); ok {
					if txs, ok := result["txs"].([]interface{}); ok {
						if (len(txs) > 0) || (i >= len(PoolRpc)-1) { // found result or last node, send it
							SendResult(w, body)
							return
						}
					}
				}
			}
		}
	}

	SendError(w)
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
