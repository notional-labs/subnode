package config

import (
	"gopkg.in/yaml.v3"
	"os"
)

type Backend struct {
	Rpc string

	// examples:
	// 	[1, 100] => from block 1 to block 100
	// 	[10] => last 10 recent blocks
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
