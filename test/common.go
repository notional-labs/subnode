package test

import (
	"github.com/notional-labs/subnode/config"
	"github.com/notional-labs/subnode/server"
	"github.com/notional-labs/subnode/state"
	sn "github.com/notional-labs/subnode/utils"
	"github.com/thedevsaddam/gojsonq/v2"
)

var isServerStarted = false

func startServer() {
	if isServerStarted { // already started
		return
	}

	isServerStarted = true
	conf := "../subnode.yaml"
	c, err := config.LoadConfigFromFile(conf)
	if err != nil {
		panic(err)
	}
	//fmt.Printf("%+v\n", c)
	config.SetConfig(c)
	server.Start()
}

func getNetworkName() string {
	rpcUrl := state.SelectPrunedNode(config.ProtocolTypeRpc).Backend.Rpc + "/status?"
	body, err := sn.FetchUriOverHttp(rpcUrl)
	if err != nil {
		panic(err)
	}
	v_network := gojsonq.New().FromString(string(body)).Find("result.node_info.network")
	network := v_network.(string)
	if len(network) <= 0 {
		panic("invalid network")
	}
	return network
}
