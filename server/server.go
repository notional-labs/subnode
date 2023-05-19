package server

import "github.com/notional-labs/subnode/state"

func Start() {
	state.Init()
	StartRpcServer()
	StartApiServer()
	StartGrpcServer()

	select {}
}

func Shutdown() {
	ShutdownRpcServer()
	ShutdownApiServer()
	ShutdownGrpcServer()
}
