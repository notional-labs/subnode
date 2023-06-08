# Eth JsonRpc

List of Methods

ref: https://github.com/evmos/docs/blob/main/docs/develop/api/ethereum-json-rpc/methods.md

    [x] web3_clientVersion
        route to pruned node
    [x] web3_sha3
        route to pruned node
    [x] net_version
        route to pruned node
    [x] net_peerCount
        route to pruned node
    [x] net_listening
        route to pruned node
    [x] eth_protocolVersion
        route to pruned node
    [x] eth_syncing
        route to pruned node
    [x] eth_gasPrice
        route to pruned node
    [x] eth_accounts
        route to pruned node
    [x] eth_blockNumber
        route to pruned node
    [x] eth_getBalance
        base on Block Number only, otherwise route to pruned node
    [x] eth_getStorageAt
        base on Block Number only, otherwise route to pruned node
    [x] eth_getTransactionCount
        base on Block Number only, otherwise route to pruned node
    [x] eth_getBlockTransactionCountByNumber
        base on Block Number only, otherwise route to pruned node
    [x] eth_getBlockTransactionCountByHash
        iterate all the subnodes
    [x] eth_getCode
        base on Block Number only, otherwise route to pruned node
    [x] eth_sign
        route to pruned node
    [x] eth_sendTransaction
        route to pruned node
    [x] eth_sendRawTransaction
        route to pruned node
    [x] eth_call
        base on Block Number only, otherwise route to pruned node
    [x] eth_estimateGas
        route to pruned node
    [x] eth_getBlockByNumber
        base on Block Number only, otherwise route to pruned node
    [x] eth_getBlockByHash
        iterate all the subnodes
    [x] eth_getTransactionByHash
        iterate all the subnodes
    [x] eth_getTransactionByBlockHashAndIndex
        iterate all the subnodes
    [x] eth_getTransactionReceipt
        iterate all the subnodes
    [x] eth_newFilter
        route to pruned node
    [x] eth_newBlockFilter
        route to pruned node
    [x] eth_newPendingTransactionFilter
        route to pruned node
    [x] eth_uninstallFilter
        route to pruned node
    [x] eth_getFilterChanges
        route to pruned node
    [x] eth_getFilterLogs
        route to pruned node
    [ ] eth_getLogs
    [ ] eth_coinbase
        route to pruned node
    [x] eth_getProof
        base on Block Number only, otherwise route to pruned node
