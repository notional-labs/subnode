package config

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLoadConfigFromBytes(t *testing.T) {
	var data = `
upstream:
  - rpc: "http://rpc"
    blocks: [200, -1]
  - rpc: "http://rpc2"
    blocks: [100, 200]
  - rpc: "http://rpc1"
    blocks: [0, 100]
`
	cfg, err := LoadConfigFromBytes([]byte(data))

	assert.NoError(t, err)

	fmt.Printf("%+v\n", cfg)
}
