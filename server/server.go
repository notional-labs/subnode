package server

import "github.com/notional-labs/subnode/state"

func Start() {
	state.Init()
	StartRpcServer()
	StartApiServer()
	StartGrpcServer()
	StartEthServer()

	select {}
}

func Shutdown() {
	ShutdownRpcServer()
	ShutdownApiServer()
	ShutdownGrpcServer()
	ShutdownEthServer()
}
