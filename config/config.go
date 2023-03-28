package config

import (
	"gopkg.in/yaml.v3"
	"os"
)

/**

[
	{
		rpc: "http://1.2.3.4.5:26656",
		blocks: [-100000, 99999999]
	},
	{
		rpc: "http://1.2.3.4.5:26656",
		blocks: [1000000, 0]
	},
	{
		rpc: "http://1.2.3.4.5:26656",
		blocks: [0, 1000000]
	}
]
*/

type Backend struct {
	Rpc    string
	Blocks []int
}

type Config struct {
	Upstream []Backend `yaml:",flow"`
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
