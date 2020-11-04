module github.com/jimpick/libp2p-caddy/go-wasm

go 1.15

require (
	github.com/filecoin-project/go-address v0.0.4
	github.com/filecoin-project/go-fil-markets v1.0.0
	github.com/filecoin-project/go-jsonrpc v0.1.2-0.20201008195726-68c6a2704e49
	github.com/filecoin-project/lotus v1.1.0
	github.com/githubnemo/CompileDaemon v1.2.1 // indirect
	github.com/ipfs/go-blockservice v0.1.4-0.20200624145336-a978cec6e834
	github.com/ipfs/go-cid v0.0.7
	github.com/ipfs/go-datastore v0.4.5
	github.com/ipfs/go-graphsync v0.3.1
	github.com/ipfs/go-ipfs-blockstore v1.0.1
	github.com/ipfs/go-ipfs-exchange-offline v0.0.1
	github.com/ipfs/go-merkledag v0.3.2
	github.com/ipfs/go-unixfs v0.2.4
	github.com/ipld/go-ipld-prime v0.5.1-0.20200828233916-988837377a7f
	github.com/libp2p/go-libp2p v0.11.0
	github.com/libp2p/go-libp2p-core v0.6.1
	github.com/libp2p/go-libp2p-nat v0.0.6 // indirect
	github.com/libp2p/go-ws-transport v0.3.1
	github.com/multiformats/go-multiaddr v0.3.1
)

replace github.com/filecoin-project/lotus => ../../lotus-wasm

replace github.com/filecoin-project/go-fil-markets => ../../go-fil-markets

replace github.com/libp2p/go-libp2p => github.com/jimpick/go-libp2p v0.3.2-0.20201019200821-ea23aad1363a

replace github.com/libp2p/go-ws-transport => github.com/jimpick/go-ws-transport v0.1.1-0.20201019202657-b30eedb00d94

replace github.com/libp2p/go-reuseport-transport => github.com/jimpick/go-reuseport-transport v0.0.5-0.20201019202422-85fd62f8a44c

replace github.com/filecoin-project/go-jsonrpc => ../../go-jsonrpc-wasm
