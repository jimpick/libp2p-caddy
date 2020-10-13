module github.com/jimpick/libp2p-caddy/go-wasm

go 1.15

require (
	github.com/libp2p/go-libp2p v0.11.0
	github.com/libp2p/go-libp2p-core v0.6.1
	github.com/libp2p/go-libp2p-nat v0.0.6 // indirect
	github.com/libp2p/go-ws-transport v0.3.1
	github.com/multiformats/go-multiaddr v0.3.1
)

replace github.com/libp2p/go-libp2p => ../../go-libp2p

replace github.com/libp2p/go-ws-transport => ../../go-ws-transport

replace github.com/libp2p/go-reuseport-transport => ../../go-reuseport-transport
