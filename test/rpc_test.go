package test

import (
	"fmt"
	sn "github.com/notional-labs/subnode/utils"
	"github.com/stretchr/testify/suite"
	"github.com/thedevsaddam/gojsonq/v2"
	"strconv"
	"testing"
	"time"
)

func TestRpcTestSuite(t *testing.T) {
	suite.Run(t, new(RpcTestSuite))
}

type RpcTestSuite struct {
	suite.Suite
	UrlEndpoint string
}

func (s *RpcTestSuite) SetupSuite() {
	go startServer()

	// wait few secs for the server to init
	time.Sleep(15 * time.Second)

	s.UrlEndpoint = "http://localhost:26657"
}

func (s *RpcTestSuite) TearDownSuite() {
	//server.Shutdown()
}

func (s *RpcTestSuite) SetupTest() {
	time.Sleep(SleepBeforeEachTest)
}

func (s *RpcTestSuite) TestRpc_abci_info() {
	// {"jsonrpc":"2.0","id":-1,"result":{"response":{"data":"OsmosisApp","app_version":"15","last_block_height":"9647581","last_block_app_hash":"dc6xiKez6O+kQ67w2Qh4/sR3PsbhDcrJScqtbSDQXR4="}}}
	rpcUrl := s.UrlEndpoint + "/abci_info?"
	body, err := sn.FetchUriOverHttp(rpcUrl)
	s.NoError(err)

	str_last_block_height := gojsonq.New().FromString(string(body)).Find("result.response.last_block_height")
	last_block_height, err := strconv.ParseInt(str_last_block_height.(string), 10, 64)
	s.NoError(err)
	s.True(last_block_height > 0)
}

func (s *RpcTestSuite) TestRpc_abci_query() {
	// {"jsonrpc":"2.0","id":-1,"result":{"response":{"code":0,"log":"","info":"","index":"0","key":null,"value":"","proofOps":null,"height":"9650945","codespace":"sdk"}}}
	rpcUrl := s.UrlEndpoint + "/abci_query?path=\"/app/version\""
	body, err := sn.FetchUriOverHttp(rpcUrl)
	s.NoError(err)

	str_height := gojsonq.New().FromString(string(body)).Find("result.response.height")
	height, err := strconv.ParseInt(str_height.(string), 10, 64)
	s.NoError(err)
	s.True(height > 0)
}

func (s *RpcTestSuite) TestRpc_block() {
	// {"jsonrpc":"2.0","id":-1,"result":{"block_id":{"hash":"1FD08E9E72D3A19FA6A4F48F88026D8B74D594C3C7EE10B26A1E268E93043BA4","...
	rpcUrl := s.UrlEndpoint + "/block?"
	body, err := sn.FetchUriOverHttp(rpcUrl)
	s.NoError(err)

	hash := gojsonq.New().FromString(string(body)).Find("result.block_id.hash")
	s.True(len(hash.(string)) == 64)
}

func (s *RpcTestSuite) TestRpc_block_by_hash() {
	// {"jsonrpc":"2.0","id":-1,"result":{"block_id":{"hash":"1FD08E9E72D3A19FA6A4F48F88026D8B74D594C3C7EE10B26A1E268E93043BA4","...

	// get hash from last block first
	rpcUrl := s.UrlEndpoint + "/block?"
	body, err := sn.FetchUriOverHttp(rpcUrl)
	s.NoError(err)
	hash := gojsonq.New().FromString(string(body)).Find("result.block_id.hash")

	rpcUrl = s.UrlEndpoint + "/block_by_hash?hash=0x" + hash.(string)
	body, err = sn.FetchUriOverHttp(rpcUrl)
	s.NoError(err)

	hash2 := gojsonq.New().FromString(string(body)).Find("result.block_id.hash")
	s.Equal(hash, hash2)
}

func (s *RpcTestSuite) TestRpc_block_results() {
	// {"jsonrpc":"2.0","id":-1,"result":{"height":"9651394","txs_results":[{"code":0,"data":"CiUKIy9pYmMuY29yZS5jbGllbnQudjEuTXNnVXBkYXR...
	rpcUrl := s.UrlEndpoint + "/block_results?"
	body, err := sn.FetchUriOverHttp(rpcUrl)
	s.NoError(err)

	v_height := gojsonq.New().FromString(string(body)).Find("result.height")
	height, err := strconv.ParseInt(v_height.(string), 10, 64)
	s.NoError(err)
	s.True(height > 0)
}

func (s *RpcTestSuite) TestRpc_block_search() {
	// {"jsonrpc":"2.0","id":-1,"result":{"blocks":[{"block_id":{"hash":"D9CE09E9B332C4374FD03EAE5211AA306A87A14BD74E99785515A79B3C5057F7"...

	// get last_block_height first
	rpcUrl := s.UrlEndpoint + "/block?"
	body, err := sn.FetchUriOverHttp(rpcUrl)
	s.NoError(err)
	v_height := gojsonq.New().FromString(string(body)).Find("result.block.header.height")
	height, err := strconv.ParseInt(v_height.(string), 10, 64)
	s.NoError(err)

	///////////
	rpcUrl = fmt.Sprint(s.UrlEndpoint, "/block_search?query=\"block.height%20=%20", height, "\"&page=1&per_page=1&order_by=\"asc\"")
	s.T().Log("rpcUrl=", rpcUrl)
	body, err = sn.FetchUriOverHttp(rpcUrl)
	s.NoError(err)

	v_hash := gojsonq.New().FromString(string(body)).Find("result.blocks.[0].block_id.hash")
	s.True(len(v_hash.(string)) == 64)
}

func (s *RpcTestSuite) TestRpc_blockchain() {
	// {"jsonrpc":"2.0","id":-1,"result":{"last_height":"9652346","block_metas":[{"block_id":{"hash":"334962A99991EF83EFFBBD066A91CE5A285C45BE7714C862B0476F72BD826DBA",...

	// get last_block_height first
	rpcUrl := s.UrlEndpoint + "/block?"
	body, err := sn.FetchUriOverHttp(rpcUrl)
	s.NoError(err)
	v_height := gojsonq.New().FromString(string(body)).Find("result.block.header.height")
	height, err := strconv.ParseInt(v_height.(string), 10, 64)
	s.NoError(err)

	rpcUrl = fmt.Sprint(s.UrlEndpoint, "/blockchain?minHeight=", height-1, "&maxHeight=", height)
	body, err = sn.FetchUriOverHttp(rpcUrl)
	s.NoError(err)
	v_height = gojsonq.New().FromString(string(body)).Find("result.last_height")
	height, err = strconv.ParseInt(v_height.(string), 10, 64)
	s.NoError(err)
	s.True(height > 0)
}

func (s *RpcTestSuite) TestRpc_commit() {
	//{
	//  "jsonrpc": "2.0",
	//  "id": -1,
	//  "result": {
	//    "signed_header": {
	//      "header": {
	//        "version": {
	//          "block": "11",
	//          "app": "15"
	//        },
	//        "chain_id": "osmosis-1",
	//        "height": "9653379",
	//        "time": "2023-05-16T13:21:27.960697653Z",
	//        "last_block_id": {
	//          "hash": "3EAC0A0A5B9FBF49E933314EF712C036C80379B31773438C72C4A4489DDE1B57",
	//          "parts": {
	//            "total": 1,
	//            "hash": "3E76071B95FBA4F33E91905EBBDE8E7F7F88ABAA27EFCE7102552352508FC7C9"
	//          }
	//        },
	//        "last_commit_hash": "6031EBA4FE965DC2BE7032599D9BD80DCC00F6BB60CC6237462AE9294DABF144",
	//        "data_hash": "CEB7DDFCC16941B34ECFE5F006DB4396D87554087F3DA8F9D7E73D00F65D2214",
	//        "validators_hash": "0DB80E04299EE9375DC265A22F65471EBC650D51AFE811E7824F87215F16CA50",
	//        "next_validators_hash": "0DB80E04299EE9375DC265A22F65471EBC650D51AFE811E7824F87215F16CA50",
	//        "consensus_hash": "A967D55FACBBA19AB96149048F2476C4657EC03D25B78A81AF5B8F0A08F61DFF",
	//        "app_hash": "0A975888D47943643531180D8EA035D2EC5EF65A2ADDA76D86ECFEBC745EE11A",
	//        "last_results_hash": "4C3D461A79BC0CFF5936558E95819266CCEB3DF556D26507E4A01E413A3DBA48",
	//        "evidence_hash": "E3B0C44298FC1C149AFBF4C8996FB92427AE41E4649B934CA495991B7852B855",
	//        "proposer_address": "E5CBA199E045E7036711D814E57E2B66C3CC0391"
	//      },
	//      "commit": {
	//        "height": "9653379",
	//        "round": 0,
	//        "block_id": {
	//          "hash": "22D9FEFE09E1316DDACB9B6401FD90C56EAF7821C1E32DC5E9ABB1B918FAAD84",
	//          "parts": {
	//            "total": 1,
	//            "hash": "59C7A2FDF47D3FB96DED16AA3FA2DB016E6A8292A4FBDCBC37477DE68DAFEC67"
	//          }
	//        },
	// ...
	rpcUrl := s.UrlEndpoint + "/commit?"
	body, err := sn.FetchUriOverHttp(rpcUrl)
	s.NoError(err)

	v_hash := gojsonq.New().FromString(string(body)).Find("result.signed_header.header.last_commit_hash")
	s.True(len(v_hash.(string)) == 64)
}

func (s *RpcTestSuite) TestRpc_consensus_params() {
	//{
	//  "jsonrpc": "2.0",
	//  "id": -1,
	//  "result": {
	//    "block_height": "9654243",
	//    "consensus_params": {
	//      "block": {
	//        "max_bytes": "10485760",
	//        "max_gas": "120000000",
	//        "time_iota_ms": "1000"
	//      },
	//      "evidence": {
	//        "max_age_num_blocks": "403200",
	//        "max_age_duration": "1209600000000000",
	//        "max_bytes": "1048576"
	//      },
	//      "validator": {
	//        "pub_key_types": [
	//          "ed25519"
	//        ]
	//      },
	//      "version": {
	//        "app_version": "15"
	//      }
	//    }
	//  }
	//}

	rpcUrl := s.UrlEndpoint + "/consensus_params?"
	body, err := sn.FetchUriOverHttp(rpcUrl)
	s.NoError(err)

	v_height := gojsonq.New().FromString(string(body)).Find("result.block_height")
	height, err := strconv.ParseInt(v_height.(string), 10, 64)
	s.NoError(err)
	s.True(height > 0)
}

func (s *RpcTestSuite) TestRpc_consensus_state() {
	//{
	//  "jsonrpc": "2.0",
	//  "id": -1,
	//  "result": {
	//    "round_state": {
	//      "height/round/step": "9654436/0/3",
	//      "start_time": "2023-05-16T15:02:41.933248936Z",
	//      "proposal_block_hash": "",
	//      "locked_block_hash": "",
	//      "valid_block_hash": "",
	//      "height_vote_set": [

	rpcUrl := s.UrlEndpoint + "/consensus_state?"
	body, err := sn.FetchUriOverHttp(rpcUrl)
	s.NoError(err)

	v_addr := gojsonq.New().FromString(string(body)).Find("result.round_state.proposer.address")
	s.NoError(err)
	s.True(len(v_addr.(string)) == 40)
}

func (s *RpcTestSuite) TestRpc_dump_consensus_state() {
	// {"jsonrpc":"2.0","id":-1,"result":{"round_state":{"height":"9655627","round":0,"st...

	rpcUrl := s.UrlEndpoint + "/dump_consensus_state?"

	body, err := sn.FetchUriOverHttp(rpcUrl)
	s.NoError(err)

	v_height := gojsonq.New().FromString(string(body)).Find("result.round_state.height")
	height, err := strconv.ParseInt(v_height.(string), 10, 64)
	s.NoError(err)
	s.True(height > 0)
}

// skip as /genesis is blocked
func (s *RpcTestSuite) TestRpc_genesis() {
}

// skip as /genesis_chunked is blocked
func (s *RpcTestSuite) TestRpc_genesis_chunked() {
}

func (s *RpcTestSuite) TestRpc_health() {
	// {"jsonrpc":"2.0","id":-1,"result":{}}

	rpcUrl := s.UrlEndpoint + "/health?"

	body, err := sn.FetchUriOverHttp(rpcUrl)
	s.NoError(err)

	v_jsonrpc := gojsonq.New().FromString(string(body)).Find("jsonrpc")
	s.NoError(err)
	s.True(v_jsonrpc.(string) == "2.0")
}

func (s *RpcTestSuite) TestRpc_net_info() {
	// {"jsonrpc":"2.0","id":-1,"result":{"listening":true,"listeners":["Listener(@)"],"n_peers":"133",...

	rpcUrl := s.UrlEndpoint + "/net_info?"

	body, err := sn.FetchUriOverHttp(rpcUrl)
	s.NoError(err)

	v_n_peers := gojsonq.New().FromString(string(body)).Find("result.n_peers")
	n_peers, err := strconv.ParseInt(v_n_peers.(string), 10, 64)
	s.NoError(err)
	s.True(n_peers >= 0)
}

func (s *RpcTestSuite) TestRpc_num_unconfirmed_txs() {
	// {"jsonrpc":"2.0","id":-1,"result":{"n_txs":"6","total":"6","total_bytes":"2378","txs":null}}

	rpcUrl := s.UrlEndpoint + "/num_unconfirmed_txs?"

	body, err := sn.FetchUriOverHttp(rpcUrl)
	s.NoError(err)

	v_n_txs := gojsonq.New().FromString(string(body)).Find("result.n_txs")
	n_txs, err := strconv.ParseInt(v_n_txs.(string), 10, 64)
	s.NoError(err)
	s.True(n_txs >= 0)
}

func (s *RpcTestSuite) TestRpc_status() {
	//{
	//  "jsonrpc": "2.0",
	//  "id": -1,
	//  "result": {
	//    "node_info": {
	//      "protocol_version": {
	//        "p2p": "8",
	//        "block": "11",
	//        "app": "15"
	//      },
	//      "id": "e306116770450276fb9c6c4e54fdc1f4a62f0c64",
	//      "listen_addr": "tcp://0.0.0.0:26656",
	//      "network": "osmosis-1",
	//      "version": "0.34.24",
	//      "channels": "40202122233038606100",
	//      "moniker": "test",
	//      "other": {
	//        "tx_index": "on",
	//        "rpc_address": "tcp://0.0.0.0:26657"
	//      }
	//    },
	//    "sync_info": {
	//      "latest_block_hash": "0F1B3BF38FFCF292FEFE43AB33B295322E1FBC6469D69FECD204248EE231BB72",
	//      "latest_app_hash": "46CB155A9928BF0987305187222F81E316E79947C8F73FF8D0362F4B60014D55",
	//      "latest_block_height": "9657269",

	rpcUrl := s.UrlEndpoint + "/status?"

	body, err := sn.FetchUriOverHttp(rpcUrl)
	s.NoError(err)

	v_latest_block_height := gojsonq.New().FromString(string(body)).Find("result.sync_info.latest_block_height")
	latest_block_height, err := strconv.ParseInt(v_latest_block_height.(string), 10, 64)
	s.NoError(err)
	s.True(latest_block_height > 0)
}

func (s *RpcTestSuite) TestRpc_tx() {
	//{
	//  "jsonrpc": "2.0",
	//  "id": -1,
	//  "result": {
	//    "hash": "7E651387114BCFAAC7AA9A49489C39D6D7D3EB7272025D973EC6E58C02A6B849",
	//    "height": "9657343",

	// figure out which chain running test on
	tx_hash := "0x7E651387114BCFAAC7AA9A49489C39D6D7D3EB7272025D973EC6E58C02A6B849"
	if Chain == "evmos" {
		tx_hash = "0xC3A6D0ED36D543BB3A587F3DED7C4B8478E6CD59C494EB202B8B5FB37E5DE879"
	}

	rpcUrl := fmt.Sprint(s.UrlEndpoint, "/tx?hash=", tx_hash, "&prove=true")
	body, err := sn.FetchUriOverHttp(rpcUrl)
	s.NoError(err)

	v_hash := gojsonq.New().FromString(string(body)).Find("result.hash")
	s.True(len(v_hash.(string)) == 64)
}

func (s *RpcTestSuite) TestRpc_tx_search() {
	//{
	//  "jsonrpc": "2.0",
	//  "id": -1,
	//  "result": {
	//    "txs": [
	//      {
	//        "hash": "474882D59D192FB7825868511E3478197FD18C45EC2002CC75169D04B8CDE1D6",
	//        "height": "9657343",

	// figure out which chain running test on
	block_num_test := 9657343 // default for osmosis
	if Chain == "evmos" {
		block_num_test = 13393844
	}

	rpcUrl := fmt.Sprint(s.UrlEndpoint, "/tx_search?query=\"tx.height=", block_num_test, "\"&prove=false&page=1&per_page=1&order_by=\"asc\"")
	body, err := sn.FetchUriOverHttp(rpcUrl)
	s.NoError(err)

	v_height := gojsonq.New().FromString(string(body)).Find("result.txs.[0].height")
	height, err := strconv.ParseInt(v_height.(string), 10, 64)
	s.NoError(err)
	s.True(height > 0)
}

func (s *RpcTestSuite) TestRpc_validators() {
	// {"jsonrpc":"2.0","id":-1,"result":{"block_height":"9657919",...

	rpcUrl := s.UrlEndpoint + "/validators?"
	body, err := sn.FetchUriOverHttp(rpcUrl)
	s.NoError(err)

	v_block_height := gojsonq.New().FromString(string(body)).Find("result.block_height")
	block_height, err := strconv.ParseInt(v_block_height.(string), 10, 64)
	s.NoError(err)
	s.True(block_height > 0)
}
