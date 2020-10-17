package main

import (
	"context"
	"fmt"
	"syscall/js"
	"time"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/host"
	peerstore "github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p/p2p/protocol/ping"
	ws "github.com/libp2p/go-ws-transport"
	multiaddr "github.com/multiformats/go-multiaddr"
)

var (
	node        host.Host
	pingService *ping.PingService
	// ch           <-chan ping.Result
	// pingResultCh chan ping.Result
	pingMaddr string
)

func pingNode(this js.Value, param []js.Value) interface{} {
	maddr := param[0].String()
	println("Go maddr: ", maddr)

	go func() {
		time.Sleep(1 * time.Second)
		ctx := context.Background()
		fmt.Printf("ping maddr found: %v\n", maddr)
		addr, err := multiaddr.NewMultiaddr(maddr)
		if err != nil {
			fmt.Printf("NewMultiaddr error %v\n", err)
			return
		}
		peer, err := peerstore.AddrInfoFromP2pAddr(addr)
		if err != nil {
			fmt.Printf("AddInfoFromP2pAddr error %v\n", err)
			return
		}
		if err := node.Connect(ctx, *peer); err != nil {
			fmt.Printf("Connect error %v\n", err)
			return
		}
		ch := pingService.Ping(ctx, peer.ID)
		res := <-ch
		fmt.Println("pinged", addr, "in", res.RTT)
		// pingResultCh <- res
		pingMaddr = ""
	}()

	return js.ValueOf(1234)

	// res := <-pingResultCh

	// create a background context (i.e. one that never cancels)
	// ctx := context.Background()

	/*
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
	*/

	// res := <-ch
	/*
		fmt.Println("click pinged", maddr, "in", res.RTT)
		return js.ValueOf(res.RTT.Milliseconds())
	*/
}

func main() {
	// create a background context (i.e. one that never cancels)
	ctx := context.Background()

	// start a libp2p node that listens on a random local TCP port,
	// but without running the built-in ping protocol
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

	println("WASM Go Initialized")

	c := make(chan struct{}, 0)
	<-c // wait forever
}
