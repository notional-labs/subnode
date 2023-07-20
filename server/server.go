package server

import "github.com/notional-labs/subnode/state"

type Node struct {
	rpcServer   *RpcServer
	apiServer   *ApiServer
	grpcServer  *GrpcServer
	ethServer   *EthServer
	ethWsServer *EthWsServer
}

func NewNode() *Node {
	newItem := &Node{
		rpcServer:   NewRpcServer(),
		apiServer:   NewApiServer(),
		grpcServer:  NewGrpcServer(),
		ethServer:   NewEthServer(),
		ethWsServer: NewEthWsServer(),
	}
	return newItem
}

func (m *Node) Start() {
	state.Init()

	m.rpcServer.StartRpcServer()
	m.apiServer.StartApiServer()
	m.grpcServer.StartGrpcServer()
	m.ethServer.StartEthServer()
	m.ethWsServer.StartEthWsServer()

	select {}
}

func (m *Node) Shutdown() {
	m.rpcServer.ShutdownRpcServer()
	m.apiServer.ShutdownApiServer()
	m.grpcServer.ShutdownGrpcServer()
	m.ethServer.ShutdownEthServer()
	m.ethWsServer.ShutdownEthWsServer()
}
