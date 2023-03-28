package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLoadConfigFromBytes(t *testing.T) {
	var data = `
upstream:
  - rpc: "http://pruned"
    blocks: [50]
  - rpc: "http://subnode"
    blocks: [100, 200]
  - rpc: "http://archive"
    blocks: []
`
	cfg, err := LoadConfigFromBytes([]byte(data))
	assert.NoError(t, err)
	fmt.Printf("%+v\n", cfg)
}

func TestGetBackendNodeType(t *testing.T) {
	b := Backend{
		Rpc:    "http://pruned",
		Blocks: []int{10},
	}

	assert.Equal(t, GetBackendNodeType(&b), BackendNodeTypePruned)
}

func TestSelectPrunedNode(t *testing.T) {
	pruned := Backend{
		Rpc:    "http://pruned",
		Blocks: []int{10},
	}
	archive := Backend{
		Rpc:    "http://archive",
		Blocks: []int{},
	}

	cfg := Config{
		Upstream: []Backend{pruned, archive},
	}

	assert.Equal(t, SelectPrunedNode(&cfg).Rpc, pruned.Rpc)
}
