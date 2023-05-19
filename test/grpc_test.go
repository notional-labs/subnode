package test

import (
	"context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	grpctypes "github.com/cosmos/cosmos-sdk/types/grpc"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"testing"
	"time"
)

func TestGrpcTestSuite(t *testing.T) {
	suite.Run(t, new(GrpcTestSuite))
}

type GrpcTestSuite struct {
	suite.Suite
	UrlEndpoint string
}

func (s *GrpcTestSuite) SetupSuite() {
	go startServer()

	// wait few secs for the server to init
	time.Sleep(15 * time.Second)

	s.UrlEndpoint = "localhost:9090"

	var Bech32PrefixAccAddr = "osmo"
	if Chain == "evmos" {
		Bech32PrefixAccAddr = "evmos"
	} else {
		panic("not supported chain " + Chain)
	}

	var (
		// Bech32PrefixAccPub defines the Bech32 prefix of an account's public key.
		Bech32PrefixAccPub = Bech32PrefixAccAddr + "pub"
		//// Bech32PrefixValAddr defines the Bech32 prefix of a validator's operator address.
		//Bech32PrefixValAddr = Bech32PrefixAccAddr + "valoper"
		//// Bech32PrefixValPub defines the Bech32 prefix of a validator's operator public key.
		//Bech32PrefixValPub = Bech32PrefixAccAddr + "valoperpub"
		//// Bech32PrefixConsAddr defines the Bech32 prefix of a consensus node address.
		//Bech32PrefixConsAddr = Bech32PrefixAccAddr + "valcons"
		//// Bech32PrefixConsPub defines the Bech32 prefix of a consensus node public key.
		//Bech32PrefixConsPub = Bech32PrefixAccAddr + "valconspub"
	)
	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount(Bech32PrefixAccAddr, Bech32PrefixAccPub)
	//config.SetBech32PrefixForValidator(Bech32PrefixValAddr, Bech32PrefixValPub)
	//config.SetBech32PrefixForConsensusNode(Bech32PrefixConsAddr, Bech32PrefixConsPub)
}

func (s *GrpcTestSuite) TearDownSuite() {
	//server.Shutdown()
}

func (s *GrpcTestSuite) SetupTest() {
	time.Sleep(SleepBeforeEachTest)
}

func (s *GrpcTestSuite) TestGrpc_GetBalance() {
	addrstr := "osmo1083svrca4t350mphfv9x45wq9asrs60cq5yv9n"
	denomstr := "uosmo"
	if Chain == "evmos" {
		addrstr = "evmos1rv94jqhlhx6makfwl6qs390e4shg32m6r5zkre"
		denomstr = "aevmos"
	} else {
		panic("not supported chain " + Chain)
	}

	////////////////
	myAddress, err := sdk.AccAddressFromBech32(addrstr)
	s.NoError(err)

	// Create a connection to the gRPC server.
	grpcConn, err := grpc.Dial(
		s.UrlEndpoint,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(grpc.ForceCodec(codec.NewProtoCodec(nil).GRPCCodec())),
	)
	s.NoError(err)
	defer func() {
		_ = grpcConn.Close()
	}()

	// This creates a gRPC client to query the x/bank service.
	bankClient := banktypes.NewQueryClient(grpcConn)
	bankRes, err := bankClient.Balance(
		context.Background(),
		&banktypes.QueryBalanceRequest{Address: myAddress.String(), Denom: denomstr},
	)
	s.NoError(err)
	s.T().Log("balance=", bankRes.GetBalance())
	s.True(len(bankRes.GetBalance().Amount.String()) > 0)

	////////////////////////////////////////////////////////////////////////////////////////
	// Query for historical state
	var header metadata.MD
	bankRes, err = bankClient.Balance(
		metadata.AppendToOutgoingContext(context.Background(), grpctypes.GRPCBlockHeightHeader, "8000000"), // Add metadata to request
		&banktypes.QueryBalanceRequest{Address: myAddress.String(), Denom: denomstr},
		grpc.Header(&header), // Retrieve header from response
	)
	s.NoError(err)
	blockHeight := header.Get(grpctypes.GRPCBlockHeightHeader)
	s.T().Log("balance=", bankRes.GetBalance(), "at height=", blockHeight)
	s.True(len(bankRes.GetBalance().Amount.String()) > 0)
}
