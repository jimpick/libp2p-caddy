package main

import (
	"context"
	"fmt"
	"syscall/js"
	"time"

	"github.com/libp2p/go-libp2p"
	peerstore "github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p/p2p/protocol/ping"
	ws "github.com/libp2p/go-ws-transport"
	multiaddr "github.com/multiformats/go-multiaddr"
)

var (
	ch           <-chan ping.Result
	pingResultCh chan ping.Result
	pingMaddr    string
)

func pingNode(this js.Value, param []js.Value) interface{} {
	maddr := param[0].String()
	println("Go maddr: ", maddr)
	pingMaddr = maddr
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

	pingResultCh = make(chan ping.Result)

	go func() {
		for {
			if pingMaddr != "" {
				// maddr := "/dns4/libp2p-caddy-ws.localhost/tcp/9056/wss/p2p/QmRvgKx6ffb9SNJBbVWTqhtZWVM36cf7zB65k8CdywPRXS"
				fmt.Printf("Jim ping maddr found: %v\n", pingMaddr)
				addr, err := multiaddr.NewMultiaddr(pingMaddr)
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
				fmt.Println("sending ping messages to", addr)
				ch = pingService.Ping(ctx, peer.ID)
				res := <-ch
				fmt.Println("pinged", addr, "in", res.RTT)
				// pingResultCh <- res
				pingMaddr = ""
			}
			time.Sleep(5 * time.Second)
		}
	}()

	/*
		sendPingCh = make(chan string)
		go func() {
			for {
				maddr := <-sendPingCh
				fmt.Println("request to ping", maddr)
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
				fmt.Println("connected, pinging", peer.ID)
				newPingCh := pingService.Ping(ctx, peer.ID)
				res := <-newPingCh
				fmt.Println("pinged", addr, "in", res.RTT)
				pingCh <- res
			}
		}()
	*/

	js.Global().Set("ping", js.FuncOf(pingNode))

	println("WASM Go Initialized")

	c := make(chan struct{}, 0)
	<-c // wait forever
}
