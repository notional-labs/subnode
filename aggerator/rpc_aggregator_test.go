package aggerator

import (
	"fmt"
	"github.com/notional-labs/subnode/utils"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFetchJsonRpcOverHttp(t *testing.T) {
	url := "https://rpc-osmosis-ia.cosmosia.notional.ventures/"
	jsonText := []byte(`{"jsonrpc": "2.0", "id": "0", "method": "validators", "params": { "height": "9045128", "page": "1", "per_page": "30" }}`)

	body, err := utils.FetchJsonRpcOverHttp(url, jsonText)
	assert.NoError(t, err)

	fmt.Printf("body=%s\n", string(body))
}
