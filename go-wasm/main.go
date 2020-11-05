package main

import (
	"context"
	"syscall/js"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p/p2p/protocol/ping"
	ws "github.com/libp2p/go-ws-transport"
)

var (
	node        host.Host
	pingService *ping.PingService
)

func main() {
	ctx := context.Background()

	var err error
	node, err = libp2p.New(ctx,
		libp2p.Transport(ws.New),
		libp2p.Ping(false),
	)
	if err != nil {
		panic(err)
	}

	// configure our own ping protocol
	pingService = &ping.PingService{Host: node}
	node.SetStreamHandler(ping.ID, pingService.PingHandler)

	js.Global().Set("ping", js.FuncOf(pingNode))
	js.Global().Set("graphSyncFetch", js.FuncOf(graphSyncFetch))
	js.Global().Set("clientQueryAsk", js.FuncOf(clientQueryAsk))

	println("WASM Go Initialized 2")

	c := make(chan struct{}, 0)
	<-c // wait forever
}
