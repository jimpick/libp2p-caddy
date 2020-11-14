module github.com/jimpick/libp2p-caddy/go-libp2p-tutorial-ws/cli-client-using-p2pclient

go 1.14

require (
	github.com/libp2p/go-libp2p v0.11.0
	github.com/libp2p/go-libp2p-core v0.6.1
	github.com/libp2p/go-libp2p-daemon v0.3.0
	github.com/libp2p/go-ws-transport v0.3.1
	github.com/multiformats/go-multiaddr v0.3.1
)

replace github.com/libp2p/go-libp2p => ../../../go-libp2p

replace github.com/libp2p/go-libp2p-core => ../../../go-libp2p-core

replace github.com/libp2p/go-ws-transport => ../../../go-ws-transport-0xproject-feat-wss-dialing

replace github.com/libp2p/go-reuseport-transport => ../../../go-reuseport-transport

replace github.com/libp2p/go-libp2p-swarm => ../../../go-libp2p-swarm
