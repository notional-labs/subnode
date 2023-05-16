package test

import (
	sn "github.com/notional-labs/subnode/utils"
	"github.com/stretchr/testify/suite"
	"github.com/thedevsaddam/gojsonq/v2"
	"strconv"
	"testing"
	"time"
)

type RpcTestSuite struct {
	suite.Suite
	UrlEndpoint string
}

func (s *RpcTestSuite) SetupTest() {
	s.UrlEndpoint = "https://rpc-osmosis-sub.cosmosia.notional.ventures"
	time.Sleep(1 * time.Second)
}

func TestRpcTestSuite(t *testing.T) {
	suite.Run(t, new(RpcTestSuite))
}

func (s *RpcTestSuite) Test_abci_info() {
	// {"jsonrpc":"2.0","id":-1,"result":{"response":{"data":"OsmosisApp","app_version":"15","last_block_height":"9647581","last_block_app_hash":"dc6xiKez6O+kQ67w2Qh4/sR3PsbhDcrJScqtbSDQXR4="}}}
	rpcUrl := s.UrlEndpoint + "/abci_info?"

	body, err := sn.FetchUriOverHttp(rpcUrl)
	s.NoError(err)

	str_last_block_height := gojsonq.New().FromString(string(body)).Find("result.response.last_block_height")
	last_block_height, err := strconv.ParseInt(str_last_block_height.(string), 10, 64)
	s.NoError(err)
	s.True(last_block_height > 0)
}

func (s *RpcTestSuite) Test_abci_query() {
	// {"jsonrpc":"2.0","id":-1,"result":{"response":{"code":0,"log":"","info":"","index":"0","key":null,"value":"","proofOps":null,"height":"9650945","codespace":"sdk"}}}
	rpcUrl := s.UrlEndpoint + "/abci_query?path=\"/app/version\""

	body, err := sn.FetchUriOverHttp(rpcUrl)
	s.NoError(err)

	str_height := gojsonq.New().FromString(string(body)).Find("result.response.height")
	height, err := strconv.ParseInt(str_height.(string), 10, 64)
	s.NoError(err)
	s.True(height > 0)
}
