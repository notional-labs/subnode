### Tendermint RPC

Support both `URI over HTTP` and `JSONRPC over HTTP`

https://docs.tendermint.com/v0.34/rpc/#/


    [x] /abci_info?
        route to pruned node
    [x] /abci_query?path=_&data=_&height=_&prove=_
        base on height
    [x] /block?height=_
        base on height
    [x] /block_by_hash?hash=_
        iterate all the subnodes
    [x] /block_results?height=_
        base on height
    [x] /block_search?query=_&page=_&per_page=_&order_by=_
        iterate all the subnodes, support partially
    [x] /blockchain?minHeight=_&maxHeight=_
        base on maxHeight
    [x] /broadcast_evidence?evidence=_
        route to pruned node
    [x] /broadcast_tx_async?tx=_
        route to pruned node
    [x] /broadcast_tx_commit?tx=_
        route to pruned node
    [x] /broadcast_tx_sync?tx=_
        route to pruned node
    [x] /check_tx?tx=_
        route to pruned node
    [x] /commit?height=_
        base on height
    [x] /consensus_params?height=_
        base on height
    [x] /consensus_state?
         route to pruned node
    [x] /dump_consensus_state?
        route to pruned node
    [x] /genesis?
        route to pruned node
    [x] /genesis_chunked?chunk=_
        route to pruned node
    [x] /health?
        route to pruned node
    [x] /net_info?
        route to pruned node
    [x] /num_unconfirmed_txs?
         route to pruned node
    [x] /status?
         route to pruned node
    [x] /subscribe?query=_
        not supported, use pruned node directly
    [x] /tx?hash=_&prove=_
        iterate all the subnodes
    [x] /tx_search?query=_&prove=_&page=_&per_page=_&order_by=_
        iterate all the subnodes, support partially
    [x] /unconfirmed_txs?limit=_
        route to pruned node
    [x] /unsubscribe?query=_
        not supported, use pruned node directly
    [x] /unsubscribe_all?
        not supported, use pruned node directly
    [x] /validators?height=_&page=_&per_page=_
        base on height


