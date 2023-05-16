package server

import "github.com/notional-labs/subnode/state"

func Start() error {
	state.Init()
	StartRpcServer()
	StartApiServer()
	StartGrpcServer()

	select {}

	return nil
}

func Shutdown() {
	ShutdownRpcServer()
	ShutdownApiServer()
	ShutdownGrpcServer()
}
