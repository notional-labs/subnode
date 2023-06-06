package test

import (
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

func TestEthTestSuite(t *testing.T) {
	if Chain != "evmos" {
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

func (s *EthTestSuite) TestEth_dummy() {

}
