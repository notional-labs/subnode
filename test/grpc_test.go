package test

import (
	"context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
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
	time.Sleep(5 * time.Second)

	s.UrlEndpoint = "localhost:9090"

	const Bech32PrefixAccAddr = "osmo"
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
	time.Sleep(1 * time.Second)
}

func (s *GrpcTestSuite) TestGrpc_GetBalance() {
	myAddress, err := sdk.AccAddressFromBech32("osmo1083svrca4t350mphfv9x45wq9asrs60cq5yv9n")
	s.NoError(err)

	// Create a connection to the gRPC server.
	grpcConn, err := grpc.Dial(
		s.UrlEndpoint,
		grpc.WithInsecure(),
		grpc.WithDefaultCallOptions(grpc.ForceCodec(codec.NewProtoCodec(nil).GRPCCodec())),
	)
	s.NoError(err)
	defer grpcConn.Close()

	// This creates a gRPC client to query the x/bank service.
	bankClient := banktypes.NewQueryClient(grpcConn)
	bankRes, err := bankClient.Balance(
		context.Background(),
		&banktypes.QueryBalanceRequest{Address: myAddress.String(), Denom: "uosmo"},
	)
	s.NoError(err)
	//fmt.Println(bankRes.GetBalance()) // Prints the account balance
	s.True(bankRes.GetBalance().Amount.Int64() > 0)
}
