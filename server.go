package main

func Start() error {
	Init()
	StartRpcServer()
	StartApiServer()
	StartGrpcServer()

	select {}

	return nil
}
