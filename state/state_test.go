package state

import (
	"github.com/notional-labs/subnode/config"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSelectPrunedNode(t *testing.T) {
	pruned := config.Backend{
		Rpc:    "http://pruned",
		Blocks: []int64{10},
	}
	archive := config.Backend{
		Rpc:    "http://archive",
		Blocks: []int64{},
	}

	cfg := &config.Config{
		Upstream: []config.Backend{pruned, archive},
	}

	config.SetConfig(cfg)

	InitPool()

	assert.Equal(t, SelectPrunedNode(config.ProtocolTypeRpc).Name, pruned.Rpc)
}

func TestReadHeightFromStatusJson(t *testing.T) {
	jsonText := []byte(`{"jsonrpc":"2.0","id":-1,"result":{"node_info":{"protocol_version":{"p2p":"8","block":"11","app":"15"},"id":"4d2605490fcd7369800cb0e1e27ef6d433c1cd96","listen_addr":"tcp://0.0.0.0:26656","network":"osmosis-1","version":"0.34.24","channels":"40202122233038606100","moniker":"test","other":{"tx_index":"on","rpc_address":"tcp://0.0.0.0:26657"}},"sync_info":{"latest_block_hash":"CD21FB291BBBA867D685187286F75BE5E0D8140B42E2CF6DC90C3B7AC40B1640","latest_app_hash":"52BA40F31319D1885CC2016FD106B2FFD301ECAB29876B7BD9107F9B710C82C0","latest_block_height":"8917933","latest_block_time":"2023-03-28T21:58:12.008170916Z","earliest_block_hash":"C8DC787FAAE0941EF05C75C3AECCF04B85DFB1D4A8D054A463F323B0D9459719","earliest_app_hash":"E3B0C44298FC1C149AFBF4C8996FB92427AE41E4649B934CA495991B7852B855","earliest_block_height":"1","earliest_block_time":"2021-06-18T17:00:00Z","catching_up":false},"validator_info":{"address":"75DC42DB1FF3C16952ACF482196791EA17D3745F","pub_key":{"type":"tendermint/PubKeyEd25519","value":"THfCGJaQ+DeUD8rI+TyFfs6ThE48RILOHI/j3dI2bZg="},"voting_power":"0"}}}`)

	height, err := ReadHeightFromStatusJson(jsonText)
	assert.NoError(t, err)
	assert.Equal(t, height, int64(8917933))
}
