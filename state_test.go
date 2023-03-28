package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSelectPrunedNode(t *testing.T) {
	pruned := Backend{
		Rpc:    "http://pruned",
		Blocks: []int{10},
	}
	archive := Backend{
		Rpc:    "http://archive",
		Blocks: []int{},
	}

	cfg = &Config{
		Upstream: []Backend{pruned, archive},
	}

	InitPool()

	assert.Equal(t, SelectPrunedNode().Name, pruned.Rpc)
}
