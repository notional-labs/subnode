package config

import (
	"gopkg.in/yaml.v3"
	"os"
)

type Backend struct {
	Rpc   string // url to rpc, eg., https://rpc-osmosis-ia.cosmosia.notional.ventures:443
	Api   string // url to api, eg., https://api-osmosis-ia.cosmosia.notional.ventures:443
	Grpc  string // url to grpc, eg., grpc-osmosis-ia.cosmosia.notional.ventures:443
	Eth   string // url to eth-http, eg., https://jsonrpc-evmos-ia.cosmosia.notional.ventures:443
	EthWs string // url to eth-ws, eg., https://jsonrpc-evmos-ia.cosmosia.notional.ventures:443/websocket/

	// examples:
	// 	[1, 100] => from block 1 to block 100 (subnode)
	// 	[10] => last 10 recent blocks
	//	[] => archive node
	Blocks []int64
}

type Config struct {
	Upstream []Backend `yaml:",flow"`
}

type BackendNodeType uint8

const (
	BackendNodeTypePruned  BackendNodeType = 0
	BackendNodeTypeSubNode BackendNodeType = 1
	BackendNodeTypeArchive BackendNodeType = 2
)

type ProtocolType uint8

const (
	ProtocolTypeRpc   ProtocolType = 0
	ProtocolTypeApi   ProtocolType = 1
	ProtocolTypeGrpc  ProtocolType = 2
	ProtocolTypeEth   ProtocolType = 3
	ProtocolTypeEthWs ProtocolType = 4
)

var (
	cfg *Config
)

func GetBackendNodeType(b *Backend) BackendNodeType {
	switch c := len(b.Blocks); c {
	case 0:
		return BackendNodeTypeArchive
	case 1:
		return BackendNodeTypePruned
	case 2:
		return BackendNodeTypeSubNode
	default:
		panic("invalid blocks config")
	}
}

func LoadConfigFromFile(filename string) (*Config, error) {
	buf, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	c, err := LoadConfigFromBytes(buf)
	if err != nil {
		return nil, err
	}

	return c, err
}

func LoadConfigFromBytes(buf []byte) (*Config, error) {
	c := &Config{}
	err := yaml.Unmarshal(buf, c)
	if err != nil {
		return nil, err
	}

	return c, err
}

func GetConfig() *Config {
	return cfg
}

func SetConfig(newcfg *Config) {
	cfg = newcfg
}
