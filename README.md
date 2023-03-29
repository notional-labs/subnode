# subnode
aggregated proxy for sub-archive-nodes of tendermint/cosmos chains

![Subnode Architecture](doc/architecture.png)


Archive node data is big, major chains could grow 5TB/year and will stop working at some point.


This project is to make archive node could scale forever by breaking data into multiple smaller nodes (called sub-node).
Each subnode stores data of 5 millions blocks or 5 TB. Old subnodes are read-only.


As data is spreaded over multiple sub-nodes, its required to have a proxy which aggreates data from sub-nodes and provides compatible rpc/api.


#### Supported Protocols
- [Tendermint RPC/JSONRPC](doc/rpc.md) on port 26657
- LCD/API on port 1337
- GRPC on port 9090



#### Configuration
See sample config [subnode.yaml](subnode.yaml).
This sample includes 1 pruned node and 1 archive node.

```yaml
upstream:
  - rpc: "https://rpc-osmosis-ia.cosmosia.notional.ventures:443"
    api: "https://api-osmosis-ia.cosmosia.notional.ventures:443"
    grpc: "grpc-osmosis-ia.cosmosia.notional.ventures:443"
    blocks: [362880]
  - rpc: "https://rpc-osmosis-archive-ia.cosmosia.notional.ventures:443"
    api: "https://api-osmosis-archive-ia.cosmosia.notional.ventures:443"
    grpc: "grpc-osmosis-archive-ia.cosmosia.notional.ventures:443"
    blocks: []
```

`blocks` config example:
- `[1, 100]` => from block 1 to block 100 (subnode)
- `[10]` => last 10 recent blocks (for pruned node)
- `[]` => for archive node

Node on the top of the list has higher priority when selecting.


### Usage
install:
```console
make install
```


start:
```console
subnode start --conf=/path/to/config/file
```