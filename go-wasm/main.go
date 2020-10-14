package main

import (
	"context"
	"fmt"
	"syscall/js"

	"github.com/libp2p/go-libp2p"
	peerstore "github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p/p2p/protocol/ping"
	ws "github.com/libp2p/go-ws-transport"
	multiaddr "github.com/multiformats/go-multiaddr"
)

func add(this js.Value, i []js.Value) interface{} {
	result := js.ValueOf(i[0].Int() + i[1].Int())
	println(result.String())
	return js.ValueOf(result)
}

func main() {
	// create a background context (i.e. one that never cancels)
	ctx := context.Background()

	// start a libp2p node that listens on a random local TCP port,
	// but without running the built-in ping protocol
	node, err := libp2p.New(ctx,
		libp2p.Transport(ws.New),
		libp2p.Ping(false),
	)
	if err != nil {
		panic(err)
	}

	// configure our own ping protocol
	pingService := &ping.PingService{Host: node}
	node.SetStreamHandler(ping.ID, pingService.PingHandler)

	maddr := "/dns4/libp2p-caddy-ws.localhost/tcp/9056/wss/p2p/Qmf23ycNkZFqFKWTmrNmRWcTKviuVfnGnnmZLQHM3BQB9t"
	addr, err := multiaddr.NewMultiaddr(maddr)
	if err != nil {
		panic(err)
	}
	peer, err := peerstore.AddrInfoFromP2pAddr(addr)
	if err != nil {
		panic(err)
	}
	if err := node.Connect(ctx, *peer); err != nil {
		panic(err)
	}
	fmt.Println("sending 5 ping messages to", addr)
	ch := pingService.Ping(ctx, peer.ID)
	for i := 0; i < 5; i++ {
		res := <-ch
		fmt.Println("pinged", addr, "in", res.RTT)
	}

	js.Global().Set("add", js.FuncOf(add))

	println("WASM Go Initialized 2")

	c := make(chan struct{}, 0)
	<-c // wait forever
}
