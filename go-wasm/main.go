package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"syscall/js"

	"github.com/filecoin-project/go-jsonrpc"
	"github.com/ipfs/go-blockservice"
	"github.com/ipfs/go-cid"
	"github.com/ipfs/go-datastore"
	dssync "github.com/ipfs/go-datastore/sync"
	blockstore "github.com/ipfs/go-ipfs-blockstore"
	offline "github.com/ipfs/go-ipfs-exchange-offline"
	"github.com/ipfs/go-merkledag"
	uio "github.com/ipfs/go-unixfs/io"
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
)

// The following functions implement a window.ping() entrypoint callable from
// JS that returns a promise.

// See: https://withblue.ink/2020/10/03/go-webassembly-http-requests-and-promises.html

func pingNode(this js.Value, param []js.Value) interface{} {
	maddr := param[0].String()
	println("Go maddr: ", maddr)

	pingHandler := func(this js.Value, args []js.Value) interface{} {
		resolve := args[0]
		reject := args[1]

		go func() {
			addr, err := multiaddr.NewMultiaddr(maddr)
			if err != nil {
				fmt.Printf("NewMultiaddr error %v\n", err)
				reject.Invoke(js.ValueOf("NewMultiaddr error"))
				return
			}

			peer, err := peerstore.AddrInfoFromP2pAddr(addr)
			if err != nil {
				fmt.Printf("AddInfoFromP2pAddr error %v\n", err)
				reject.Invoke(js.ValueOf("AddInfoFromP2pAddr error"))
				return
			}

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			if err := node.Connect(ctx, *peer); err != nil {
				fmt.Printf("Connect error %v\n", err)
				reject.Invoke(js.ValueOf("Connect error"))
				return
			}

			ch := pingService.Ping(ctx, peer.ID)
			res := <-ch
			fmt.Println("pinged", addr, "in", res.RTT)
			resolve.Invoke(js.ValueOf(res.RTT.Milliseconds()))
		}()
		return nil
	}

	promiseConstructor := js.Global().Get("Promise")
	return promiseConstructor.New(js.FuncOf(pingHandler))
}

func graphSyncFetch(this js.Value, param []js.Value) interface{} {
	maddr := param[0].String()
	cidArg := param[1].String()
	fmt.Printf("Go maddr %v\n", maddr)
	fmt.Printf("Go cid: %v\n", cidArg)

	graphSyncFetchHandler := func(this js.Value, args []js.Value) interface{} {
		resolve := args[0]
		reject := args[1]

		go func() {
			addr, err := multiaddr.NewMultiaddr(maddr)
			if err != nil {
				fmt.Printf("NewMultiaddr error %v\n", err)
				reject.Invoke(js.ValueOf("NewMultiaddr error"))
				return
			}

			peer, err := peerstore.AddrInfoFromP2pAddr(addr)
			if err != nil {
				fmt.Printf("AddInfoFromP2pAddr error %v\n", err)
				reject.Invoke(js.ValueOf("AddInfoFromP2pAddr error"))
				return
			}

			target, err := cid.Decode(cidArg)
			if err != nil {
				fmt.Printf("failed to decode CID '%q': %s\n", cidArg, err)
				reject.Invoke(js.ValueOf("CID Decode error"))
			}
			fmt.Printf("Go target: %v\n", target)

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			if err := node.Connect(ctx, *peer); err != nil {
				fmt.Printf("Connect error %v\n", err)
				reject.Invoke(js.ValueOf("Connect error"))
				return
			}

			bs := blockstore.NewBlockstore(dssync.MutexWrap(datastore.NewMapDatastore()))
			gs, err := newGraphsync(ctx, node, bs)
			if err != nil {
				fmt.Printf("newGraphsync error %v\n", err)
				reject.Invoke(js.ValueOf("newGraphsync error"))
				return
			}

			err = fetch(ctx, gs, peer.ID, target)
			if err != nil {
				fmt.Printf("fetch error %v\n", err)
				reject.Invoke(js.ValueOf("fetch error"))
				return
			}

			dag := merkledag.NewDAGService(blockservice.New(bs, offline.Exchange(bs)))
			root, err := dag.Get(ctx, target)
			if err != nil {
				fmt.Printf("get.Get error %v\n", err)
				reject.Invoke(js.ValueOf("get.Get error"))
				return
			}

			reader, err := uio.NewDagReader(ctx, root, dag)
			if err != nil {
				fmt.Printf("NewDagReader error %v\n", err)
				reject.Invoke(js.ValueOf("NewDagReader error"))
				return
			}

			buf := new(bytes.Buffer)
			_, err = io.Copy(buf, reader)
			if err != nil {
				fmt.Printf("io.Copy error %v\n", err)
				reject.Invoke(js.ValueOf("io.Copy error"))
				return
			}

			resolve.Invoke(js.ValueOf(buf.String()))
		}()
		return nil
	}

	promiseConstructor := js.Global().Get("Promise")
	return promiseConstructor.New(js.FuncOf(graphSyncFetchHandler))
}

func chainHead(this js.Value, param []js.Value) interface{} {
	jsHandler := param[0]
	println("Go chainHead")

	go func() {
		var api FullNodeStruct
		closer, err := jsonrpc.NewJsMergeClient(context.Background(), jsHandler, "Filecoin", []interface{}{&api})
		if err != nil {
			fmt.Printf("connecting with lotus failed: %s\n", err)
			return
		}
		defer closer()

		// Now you can call any API you're interested in.
		tipset, err := api.ChainHead(context.Background())
		if err != nil {
			fmt.Printf("calling chain head: %s\n", err)
			return
		}
		fmt.Printf("Current chain head is: %s\n", tipset.String())
	}()

	return js.ValueOf(1234)
}

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
	js.Global().Set("chainHead", js.FuncOf(chainHead))

	println("WASM Go Initialized")

	c := make(chan struct{}, 0)
	<-c // wait forever
}
