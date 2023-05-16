package internal

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
		Blocks: []int64{10},
	}

	assert.Equal(t, GetBackendNodeType(&b), BackendNodeTypePruned)
}
