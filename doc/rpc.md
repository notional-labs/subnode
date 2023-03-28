### Tendermint RPC

#### URI over HTTP

https://docs.tendermint.com/v0.34/rpc/#/


    [ ] /abci_info?
        route to pruned node
    [ ] /abci_query?path=_&data=_&height=_&prove=_
        base on height
    [ ] /block?height=_
        base on height
    [ ] /block_by_hash?hash=_
        not supported, should use indexer instead
    [ ] /block_results?height=_
        base on height
    [ ] /block_search?query=_&page=_&per_page=_&order_by=_
        not supported, should be used with indexer
    [ ] /blockchain?minHeight=_&maxHeight=_
        base on minHeight & maxHeight
    [ ] /broadcast_evidence?evidence=_
        route to pruned node
    [ ] /broadcast_tx_async?tx=_
        route to pruned node
    [ ] /broadcast_tx_commit?tx=_
        route to pruned node
    [ ] /broadcast_tx_sync?tx=_
        route to pruned node
    [ ] /check_tx?tx=_
        route to pruned node
    [ ] /commit?height=_
        base on height
    [ ] /consensus_params?height=_
        base on height
    [ ] /consensus_state?
         route to pruned node
    [ ] /dump_consensus_state?
        route to pruned node
    [ ] /genesis?
        route to pruned node
    [ ] /genesis_chunked?chunk=_
        route to pruned node
    [ ] /health?
        route to pruned node
    [ ] /net_info?
        route to pruned node
    [ ] /num_unconfirmed_txs?
         route to pruned node
    [ ] /status?
         route to pruned node
    [ ] /subscribe?query=_
        not supported, use pruned node directly
    [ ] /tx?hash=_&prove=_
        not supported, should use indexer instead
    [ ] /tx_search?query=_&prove=_&page=_&per_page=_&order_by=_
        not supported, should use indexer instead
    [ ] /unconfirmed_txs?limit=_
        route to pruned node
    [ ] /unsubscribe?query=_
        not supported, use pruned node directly
    [ ] /unsubscribe_all?
        not supported, use pruned node directly
    [ ] /validators?height=_&page=_&per_page=_
        base on height


