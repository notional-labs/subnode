package aggerator

import (
	"encoding/json"
	"github.com/notional-labs/subnode/state"
	"github.com/notional-labs/subnode/utils"
	"net/http"
)

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
