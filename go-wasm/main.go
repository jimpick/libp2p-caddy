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

var (
	ch         <-chan ping.Result
	sendPingCh chan string
	pingCh     chan ping.Result
)

/*
func add(this js.Value, i []js.Value) interface{} {
	result := js.ValueOf(i[0].Int() + i[1].Int())
	println(result.String())
	return js.ValueOf(result)
}
*/

func pingNode(this js.Value, param []js.Value) interface{} {
	maddr := param[0].String()
	println("Go maddr: ", maddr)
	sendPingCh <- maddr
	res := <-pingCh

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
	fmt.Println("click pinged", maddr, "in", res.RTT)
	return js.ValueOf(res.RTT.Seconds())
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

	/*
		maddr := "/dns4/libp2p-caddy-ws.localhost/tcp/9056/wss/p2p/QmRvgKx6ffb9SNJBbVWTqhtZWVM36cf7zB65k8CdywPRXS"
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
		go func() {
			fmt.Println("sending ping messages to", addr)
			ch = pingService.Ping(ctx, peer.ID)
			for i := 0; i < 500; i++ {
				res := <-ch
				fmt.Println(i, "pinged", addr, "in", res.RTT)
				time.Sleep(10 * time.Second)
			}
		}()
	*/

	sendPingCh = make(chan string)
	pingCh = make(chan ping.Result)
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

	js.Global().Set("ping", js.FuncOf(pingNode))

	println("WASM Go Initialized")

	c := make(chan struct{}, 0)
	<-c // wait forever
}
