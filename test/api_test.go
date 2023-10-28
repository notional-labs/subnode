package test

import (
	sn "github.com/notional-labs/subnode/utils"
	"github.com/stretchr/testify/suite"
	"github.com/thedevsaddam/gojsonq/v2"
	"testing"
	"time"
)

func TestApiTestSuite(t *testing.T) {
	suite.Run(t, new(ApiTestSuite))
}

type ApiTestSuite struct {
	suite.Suite
	UrlEndpoint string
}

func (s *ApiTestSuite) SetupSuite() {
	go startServer()

	// wait few secs for the server to init
	time.Sleep(15 * time.Second)

	s.UrlEndpoint = "http://localhost:1317"
}

func (s *ApiTestSuite) TearDownSuite() {
	//server.Shutdown()
}

func (s *ApiTestSuite) SetupTest() {
	time.Sleep(SleepBeforeEachTest)
}

func (s *ApiTestSuite) TestApi_cosmos_staking_v1beta1_params() {
	//{
	//  "params": {
	//    "unbonding_time": "1209600s",
	//    "max_validators": 150,
	//    "max_entries": 7,
	//    "historical_entries": 10000,
	//    "bond_denom": "uosmo",
	//    "min_commission_rate": "0.050000000000000000",
	//    "min_self_delegation": "0"
	//  }
	//}

	rpcUrl := s.UrlEndpoint + "/cosmos/staking/v1beta1/params"
	body, err := sn.FetchUriOverHttp(rpcUrl)
	s.NoError(err)

	v_unbonding_time := gojsonq.New().FromString(string(body)).Find("params.unbonding_time")
	s.NoError(err)
	s.True(len(v_unbonding_time.(string)) > 0)
}
