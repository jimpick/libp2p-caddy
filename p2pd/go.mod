module github.com/jimpick/libp2p-caddy/p2pd

go 1.15

require (
	github.com/libp2p/go-libp2p v0.11.0
	github.com/libp2p/go-libp2p-circuit v0.3.1
	github.com/libp2p/go-libp2p-connmgr v0.2.4
	github.com/libp2p/go-libp2p-daemon v0.3.0
	github.com/libp2p/go-libp2p-noise v0.1.1
	github.com/libp2p/go-libp2p-pubsub v0.3.5
	github.com/libp2p/go-libp2p-quic-transport v0.8.0
	github.com/libp2p/go-libp2p-secio v0.2.2
	github.com/libp2p/go-libp2p-tls v0.1.3
	github.com/multiformats/go-multiaddr v0.3.1
	github.com/prometheus/client_golang v1.7.1
)

replace github.com/libp2p/go-libp2p-daemon => ../../go-libp2p-daemon-ws

replace github.com/libp2p/go-libp2p-transport-upgrader => ../../go-libp2p-transport-upgrader

replace github.com/libp2p/go-ws-transport => ../../go-ws-transport-0xproject-feat-wss-dialing
