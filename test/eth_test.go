package test

import (
	"github.com/notional-labs/subnode/utils"
	"github.com/stretchr/testify/suite"
	"github.com/thedevsaddam/gojsonq/v2"
	"testing"
	"time"
)

func TestEthTestSuite(t *testing.T) {
	if Chain != "evmos" {
		t.Log("Chain:", Chain, ". Ignore testing eth-jsonrpc for this chain. ")
		return
	}

	suite.Run(t, new(EthTestSuite))
}

type EthTestSuite struct {
	suite.Suite
	UrlEndpoint string
}

func (s *EthTestSuite) SetupSuite() {
	go startServer()

	// wait few secs for the server to init
	time.Sleep(15 * time.Second)

	s.UrlEndpoint = "http://localhost:8545"
}

func (s *EthTestSuite) TearDownSuite() {
	//server.Shutdown()
}

func (s *EthTestSuite) SetupTest() {
	time.Sleep(SleepBeforeEachTest)
}

func (s *EthTestSuite) TestEth_getBlockTransactionCountByHash() {
	jsonText := []byte(`{"jsonrpc":"2.0","method":"eth_getBlockTransactionCountByHash","params":["0x8f0708ce38bdb91b099555dc354716ff7af4bc85acdab7ba45638f8cfab3696a"],"id":1}`)
	body, err := utils.FetchJsonRpcOverHttp(s.UrlEndpoint, jsonText)
	s.NoError(err)

	v_result := gojsonq.New().FromString(string(body)).Find("result")
	s.True(v_result != nil)
}

func (s *EthTestSuite) TestEth_getBlockByHash() {
	jsonText := []byte(`{"jsonrpc":"2.0","method":"eth_getBlockByHash","params":["0x8f0708ce38bdb91b099555dc354716ff7af4bc85acdab7ba45638f8cfab3696a", false],"id":1}`)
	body, err := utils.FetchJsonRpcOverHttp(s.UrlEndpoint, jsonText)
	s.NoError(err)

	v_result := gojsonq.New().FromString(string(body)).Find("result")
	s.True(v_result != nil)
}

func (s *EthTestSuite) TestEth_getTransactionByHash() {
	jsonText := []byte(`{"jsonrpc":"2.0","method":"eth_getTransactionByHash","params":["0x992dcea82962e5d424b13357f5ffb69f317b072f1d435f0e0581e7876284782f"],"id":1}`)
	body, err := utils.FetchJsonRpcOverHttp(s.UrlEndpoint, jsonText)
	s.NoError(err)

	v_result := gojsonq.New().FromString(string(body)).Find("result")
	s.True(v_result != nil)
}

func (s *EthTestSuite) TestEth_getTransactionByBlockHashAndIndex() {
	jsonText := []byte(`{"jsonrpc":"2.0","method":"eth_getTransactionByBlockHashAndIndex","params":["0x514786aee553492748c63d7994654d0074e780ecc1d3ef5ffa68251b8810e2bd", "0x0"],"id":1}`)
	body, err := utils.FetchJsonRpcOverHttp(s.UrlEndpoint, jsonText)
	s.NoError(err)

	v_result := gojsonq.New().FromString(string(body)).Find("result")
	s.True(v_result != nil)
}

func (s *EthTestSuite) TestEth_getTransactionReceipt() {
	jsonText := []byte(`{"jsonrpc":"2.0","method":"eth_getTransactionReceipt","params":["0x992dcea82962e5d424b13357f5ffb69f317b072f1d435f0e0581e7876284782f"],"id":1}`)
	body, err := utils.FetchJsonRpcOverHttp(s.UrlEndpoint, jsonText)
	s.NoError(err)

	v_result := gojsonq.New().FromString(string(body)).Find("result")
	s.True(v_result != nil)
}

func (s *EthTestSuite) TestEth_BatchRequest() {
	jsonText := []byte(`[{"jsonrpc":"2.0","method":"web3_clientVersion","params":[],"id":1},{"jsonrpc":"2.0","method":"net_version","params":[],"id":2}]`)
	body, err := utils.FetchJsonRpcOverHttp(s.UrlEndpoint, jsonText)
	s.NoError(err)

	// [{"jsonrpc":"2.0","id":1,"result":"Version dev ()\nCompiled at  using Go go1.20.4 (amd64)"},{"jsonrpc":"2.0","id":2,"result":"9001"}]
	s.T().Log("res=:", string(body))

	jq := gojsonq.New().FromString(string(body))

	v_r0 := jq.Find("[0].result")
	s.True(v_r0 != nil)

	jq.Reset()
	v_r1 := jq.Find("[1].result")
	s.True(v_r1 != nil)
}
