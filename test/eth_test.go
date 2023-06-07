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
