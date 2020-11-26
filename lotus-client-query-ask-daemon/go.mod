module github.com/jimpick/libp2p-caddy/lotus-client-query-ask-daemon

go 1.15

require (
	github.com/filecoin-project/go-address v0.0.5-0.20201103152444-f2023ef3f5bb
	github.com/filecoin-project/go-fil-markets v1.0.5
	github.com/filecoin-project/go-jsonrpc v0.1.2-0.20201008195726-68c6a2704e49
	github.com/filecoin-project/lotus v1.1.2
	github.com/libp2p/go-libp2p-core v0.6.1
	github.com/libp2p/go-libp2p-kad-dht v0.8.3
	github.com/libp2p/go-libp2p-peer v0.2.0
	github.com/libp2p/go-libp2p-peerstore v0.2.6
	github.com/libp2p/go-libp2p-record v0.1.3
	github.com/urfave/cli/v2 v2.2.0
	github.com/wangjia184/sortedset v0.0.0-20160527075905-f5d03557ba30 // indirect
	go.uber.org/fx v1.13.1
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1
)

replace github.com/supranational/blst => ../extern/fil-blst/blst

replace github.com/filecoin-project/fil-blst => ../extern/fil-blst

replace github.com/filecoin-project/lotus => ../../lotus-client-query-ask-native

replace github.com/filecoin-project/go-fil-markets => ../../go-fil-markets-for-query-ask-api-native
