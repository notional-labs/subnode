package main

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestJsonRpcParams(t *testing.T) {
	jsonText := []byte(`{"jsonrpc": "2.0", "id": "0", "method": "validators", "params": { "height": "9045128", "page": "1", "per_page": "30" }}`)
	//jsonText := []byte(`{"jsonrpc": "2.0", "id": "0", "method": "validators", "params": ["9045128", "1", "30"]}`)

	var j0 interface{}
	err := json.Unmarshal(jsonText, &j0)
	assert.NoError(t, err)

	m0 := j0.(map[string]interface{})
	//method := m0["method"].(string)
	//params := m0["params"].([]interface{})
	//
	//fmt.Printf("method=%s, params=%+v\n", method, params)

	fmt.Println(reflect.TypeOf(m0["params"]))
}

func TestDummy(t *testing.T) {
	jsonText := []byte(`{"a": "aaa", "b": null}`)

	var j0 interface{}
	err := json.Unmarshal(jsonText, &j0)
	assert.NoError(t, err)
	m0, ok := j0.(map[string]interface{})
	assert.True(t, ok)

	fmt.Printf("b=%+v\n", m0["b"])
}
