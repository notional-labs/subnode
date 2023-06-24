package aggerator

import (
	"encoding/json"
	"github.com/notional-labs/subnode/state"
	"github.com/notional-labs/subnode/utils"
	"log"
	"net/http"
)

func Eth_BathRequest(w http.ResponseWriter, jsonBody []byte) {
	log.Printf("Eth_BathRequest...")
	var arr []json.RawMessage

	err := json.Unmarshal(jsonBody, &arr)
	if err != nil {
		_ = utils.SendError(w)
		return
	}

	var arrSize = len(arr)
	var arrRes = make([]json.RawMessage, arrSize)

	for i, s := range arr {
		bodyChild, err := utils.FetchJsonRpcOverHttp("http://localhost:8545", s)
		if err != nil {
			_ = utils.SendError(w)
			return
		}

		arrRes[i] = bodyChild
	}

	jsonBytes, err := json.Marshal(arrRes)
	if err != nil {
		_ = utils.SendError(w)
	}

	_ = utils.SendResult(w, jsonBytes)
}

func Eth_getBlockTransactionCountByHash(w http.ResponseWriter, jsonBody []byte) {
	for i, s := range state.PoolEth {
		ethUrl := s.Backend.Eth
		body, err := utils.FetchJsonRpcOverHttp(ethUrl, jsonBody)
		if err != nil {
			continue
		}

		var j0 interface{}
		err = json.Unmarshal(body, &j0)
		if err != nil {
			_ = utils.SendError(w)
			return
		}

		if m0, ok := j0.(map[string]interface{}); ok {
			if (m0["result"] != nil) || (i >= len(state.PoolEth)-1) { // found result or last node, send it
				_ = utils.SendResult(w, body)
				return
			}
		}

	}

	_ = utils.SendError(w)
}

func Eth_getBlockByHash(w http.ResponseWriter, jsonBody []byte) {
	for i, s := range state.PoolEth {
		ethUrl := s.Backend.Eth
		body, err := utils.FetchJsonRpcOverHttp(ethUrl, jsonBody)
		if err != nil {
			continue
		}

		var j0 interface{}
		err = json.Unmarshal(body, &j0)
		if err != nil {
			_ = utils.SendError(w)
			return
		}

		if m0, ok := j0.(map[string]interface{}); ok {
			if (m0["result"] != nil) || (i >= len(state.PoolEth)-1) { // found result or last node, send it
				_ = utils.SendResult(w, body)
				return
			}
		}

	}

	_ = utils.SendError(w)
}

func Eth_getTransactionByHash(w http.ResponseWriter, jsonBody []byte) {
	for i, s := range state.PoolEth {
		ethUrl := s.Backend.Eth
		body, err := utils.FetchJsonRpcOverHttp(ethUrl, jsonBody)
		if err != nil {
			continue
		}

		var j0 interface{}
		err = json.Unmarshal(body, &j0)
		if err != nil {
			_ = utils.SendError(w)
			return
		}

		if m0, ok := j0.(map[string]interface{}); ok {
			if (m0["result"] != nil) || (i >= len(state.PoolEth)-1) { // found result or last node, send it
				_ = utils.SendResult(w, body)
				return
			}
		}

	}

	_ = utils.SendError(w)
}

func Eth_getTransactionByBlockHashAndIndex(w http.ResponseWriter, jsonBody []byte) {
	for i, s := range state.PoolEth {
		ethUrl := s.Backend.Eth
		body, err := utils.FetchJsonRpcOverHttp(ethUrl, jsonBody)
		if err != nil {
			continue
		}

		var j0 interface{}
		err = json.Unmarshal(body, &j0)
		if err != nil {
			_ = utils.SendError(w)
			return
		}

		if m0, ok := j0.(map[string]interface{}); ok {
			if (m0["result"] != nil) || (i >= len(state.PoolEth)-1) { // found result or last node, send it
				_ = utils.SendResult(w, body)
				return
			}
		}

	}

	_ = utils.SendError(w)
}

func Eth_getTransactionReceipt(w http.ResponseWriter, jsonBody []byte) {
	for i, s := range state.PoolEth {
		ethUrl := s.Backend.Eth
		body, err := utils.FetchJsonRpcOverHttp(ethUrl, jsonBody)
		if err != nil {
			continue
		}

		var j0 interface{}
		err = json.Unmarshal(body, &j0)
		if err != nil {
			_ = utils.SendError(w)
			return
		}

		if m0, ok := j0.(map[string]interface{}); ok {
			if (m0["result"] != nil) || (i >= len(state.PoolEth)-1) { // found result or last node, send it
				_ = utils.SendResult(w, body)
				return
			}
		}

	}

	_ = utils.SendError(w)
}
