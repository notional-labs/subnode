package server

func Start() error {
	Init()
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
