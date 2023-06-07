package test

import (
	"fmt"
	"github.com/notional-labs/subnode/config"
	"github.com/notional-labs/subnode/server"
	"time"
)

// Chain is set at compile time `-X github.com/notional-labs/subnode/test.Chain=osmosis`
// supported value: osmosis, evmos
// default is osmosis
var Chain = "osmosis"

const SleepBeforeEachTest = 2 * time.Second

var isServerStarted = false

func startServer() {
	if isServerStarted { // already started
		return
	}

	isServerStarted = true
	conf := fmt.Sprintf("./test.config.%s.yaml", Chain)
	fmt.Printf("config file= %s\n", conf)
	c, err := config.LoadConfigFromFile(conf)
	if err != nil {
		panic(err)
	}
	//fmt.Printf("%+v\n", c)
	config.SetConfig(c)
	server.Start()
}
